package partnerapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Webhook signature verification.
//
// Every webhook delivery is signed Stripe-style with the subscription secret
// (shown once at create/rotate time):
//
//	X-ShedCloud-Signature: t=<unix seconds>,v1=<hex hmac-sha256 of "<t>.<body>">
//
// The timestamp is bound into the signed string, so an old signed body
// cannot be replayed with a fresh timestamp.

// WebhookSignatureHeader carries the signature on every webhook delivery.
const WebhookSignatureHeader = "X-ShedCloud-Signature"

// WebhookEventIDHeader / WebhookEventTypeHeader accompany every delivery for
// cheap dedupe before parsing the body.
const (
	WebhookEventIDHeader   = "X-ShedCloud-Event-Id"
	WebhookEventTypeHeader = "X-ShedCloud-Event-Type"
)

// DefaultWebhookTolerance is the maximum accepted signature age when
// VerifyWebhookSignature is called with tolerance <= 0.
const DefaultWebhookTolerance = 5 * time.Minute

// ErrWebhookSignature is wrapped by every VerifyWebhookSignature failure —
// check with errors.Is.
var ErrWebhookSignature = errors.New("partnerapi: webhook signature verification failed")

// ComputeWebhookSignature returns the hex HMAC-SHA256 of
// "<unix seconds>.<body>" — the value carried in the v1= part of the
// signature header. Exposed mainly for constructing test fixtures.
func ComputeWebhookSignature(secret string, t time.Time, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(t.Unix(), 10)))
	mac.Write([]byte("."))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyWebhookSignature validates a received X-ShedCloud-Signature header
// against the raw request body (exactly as received — do not re-serialize
// parsed JSON) and the subscription secret. tolerance <= 0 uses
// DefaultWebhookTolerance. Returns nil when the signature is valid; every
// failure wraps ErrWebhookSignature.
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//		body, _ := io.ReadAll(r.Body)
//		err := partnerapi.VerifyWebhookSignature(secret,
//			r.Header.Get(partnerapi.WebhookSignatureHeader), body, 0)
//		if err != nil {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//		w.WriteHeader(http.StatusOK) // ack fast; process async, dedupe by event id
//	}
func VerifyWebhookSignature(secret, header string, body []byte, tolerance time.Duration) error {
	return verifyWebhookSignatureAt(secret, header, body, time.Now(), tolerance)
}

func verifyWebhookSignatureAt(secret, header string, body []byte, now time.Time, tolerance time.Duration) error {
	if secret == "" {
		return errors.Join(ErrWebhookSignature, errors.New("secret is required"))
	}
	if strings.TrimSpace(header) == "" {
		return errors.Join(ErrWebhookSignature, errors.New("signature header is missing"))
	}
	if tolerance <= 0 {
		tolerance = DefaultWebhookTolerance
	}

	var ts int64
	var sigs []string
	for _, part := range strings.Split(header, ",") {
		k, v, found := strings.Cut(strings.TrimSpace(part), "=")
		if !found {
			continue
		}
		switch k {
		case "t":
			parsed, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return errors.Join(ErrWebhookSignature, errors.New("invalid timestamp in signature header"))
			}
			ts = parsed
		case "v1":
			sigs = append(sigs, v)
		}
	}
	if ts == 0 || len(sigs) == 0 {
		return errors.Join(ErrWebhookSignature, errors.New("signature header missing t= or v1="))
	}

	at := time.Unix(ts, 0)
	if d := now.Sub(at); d > tolerance || d < -tolerance {
		return errors.Join(ErrWebhookSignature, errors.New("signature timestamp outside tolerance"))
	}

	expected := ComputeWebhookSignature(secret, at, body)
	for _, sig := range sigs {
		if hmac.Equal([]byte(strings.ToLower(sig)), []byte(expected)) {
			return nil
		}
	}
	return errors.Join(ErrWebhookSignature, errors.New("signature mismatch"))
}
