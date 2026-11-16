package partnerapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRequiresBaseURLAndAuth(t *testing.T) {
	if _, err := New(Options{Auth: Auth{APIKey: "sc_live_x"}}); !errors.Is(err, ErrBaseURLRequired) {
		t.Fatalf("err = %v, want ErrBaseURLRequired", err)
	}
	if _, err := New(Options{BaseURL: "https://erp.example"}); err == nil {
		t.Fatal("expected an auth error")
	}
	c, err := New(Options{BaseURL: "https://erp.example///", Auth: Auth{APIKey: "sc_live_x"}})
	if err != nil {
		t.Fatal(err)
	}
	if c.BaseURL != "https://erp.example" {
		t.Fatalf("BaseURL = %q", c.BaseURL)
	}
}

func TestErrorEnvelope(t *testing.T) {
	c, _ := fakeERP(t, 422, `{"error":"invalid transition","code":"configuration_rule_violation","message":"This build breaks a rule.","violations":[{"code":"x"}]}`)
	_, err := c.Quotes.Get(context.Background(), "q1")
	var pe *Error
	if !errors.As(err, &pe) {
		t.Fatalf("err = %v", err)
	}
	if pe.Status != 422 || pe.Code != "configuration_rule_violation" || pe.Message != "This build breaks a rule." {
		t.Fatalf("error = %+v", pe)
	}
	if !pe.IsUnprocessable() || pe.IsForbidden() {
		t.Fatal("status helpers wrong")
	}
	if body, _ := pe.Body.(map[string]any); body["violations"] == nil {
		t.Fatal("violations should stay in Body")
	}

	c, _ = fakeERP(t, 403, `{"error":"forbidden","code":"forbidden"}`)
	_, err = c.Me(context.Background())
	if !errors.As(err, &pe) || !pe.IsForbidden() || pe.Message != "forbidden" {
		t.Fatalf("403 = %v", err)
	}
}

func TestMeAndCompany(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"tenant_id":"t1","credential_id":"c1","kind":"api_key","scopes":["partner-api.leads.read"]}`)
	me, err := c.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/me")
	if me.TenantID != "t1" || me.Kind != "api_key" || !HasScope(me.Scopes, ScopeLeadsRead) || HasScope(me.Scopes, ScopeLeadsWrite) {
		t.Fatalf("me = %+v", me)
	}

	c, rec = fakeERP(t, 200, `{"id":"a","name":"Austin Sheds","slug":"austin-sheds","type":"erp","locale":"en","currency":"USD"}`)
	co, err := c.Company.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/company")
	if co.Name != "Austin Sheds" || co.Currency != "USD" {
		t.Fatalf("company = %+v", co)
	}
}

func TestOAuthClientCredentials(t *testing.T) {
	var tokenCalls, apiCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
			atomic.AddInt32(&tokenCalls, 1)
			raw, _ := io.ReadAll(r.Body)
			form := string(raw)
			for _, want := range []string{"grant_type=client_credentials", "client_id=sc_client_1", "client_secret=sc_secret_1"} {
				if !strings.Contains(form, want) {
					t.Errorf("token form %q lacks %q", form, want)
				}
			}
			if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
				t.Errorf("content-type = %q", ct)
			}
			_, _ = w.Write([]byte(`{"access_token":"sc_at_abc","token_type":"Bearer","expires_in":3600,"scope":"partner-api.leads.read"}`))
		case "/partner/v1/me":
			atomic.AddInt32(&apiCalls, 1)
			if got := r.Header.Get("Authorization"); got != "Bearer sc_at_abc" {
				t.Errorf("Authorization = %q", got)
			}
			_, _ = w.Write([]byte(`{"tenant_id":"t","credential_id":"c","kind":"oauth","scopes":[]}`))
		}
	}))
	defer srv.Close()

	c, err := New(Options{BaseURL: srv.URL, Auth: Auth{ClientID: "sc_client_1", ClientSecret: "sc_secret_1"}})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err := c.Me(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	if tokenCalls != 1 || apiCalls != 2 {
		t.Fatalf("tokenCalls=%d apiCalls=%d, want 1 and 2 (token cached)", tokenCalls, apiCalls)
	}
	if _, err := c.RefreshAccessToken(context.Background()); err != nil || tokenCalls != 2 {
		t.Fatalf("refresh err=%v tokenCalls=%d", err, tokenCalls)
	}
}

func TestOAuthFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = w.Write([]byte(`{"error":"invalid client credentials","code":"unauthorized"}`))
	}))
	defer srv.Close()
	c, _ := New(Options{BaseURL: srv.URL, Auth: Auth{ClientID: "a", ClientSecret: "b"}})
	_, err := c.Me(context.Background())
	var ae *AuthError
	if !errors.As(err, &ae) || !ae.IsUnauthorized() || ae.Code != "unauthorized" {
		t.Fatalf("err = %#v", err)
	}
}

func TestScopeCatalog(t *testing.T) {
	if len(AllScopes) != 91 {
		t.Fatalf("AllScopes has %d entries, want the ERP catalog size (91)", len(AllScopes))
	}
	seenScope := map[string]bool{}
	for _, s := range AllScopes {
		if !strings.HasPrefix(s, "partner-api.") || seenScope[s] {
			t.Fatalf("bad or duplicate scope %q", s)
		}
		seenScope[s] = true
	}
	for _, s := range []string{ScopeLotStockWrite, ScopePricingRead, ScopeSalesLedgerRead, ScopeRevenueForecastRead, ScopeSiteEventsWrite} {
		if !seenScope[s] {
			t.Fatalf("scope %q missing", s)
		}
	}
	if ScopeLotStockRead != "partner-api.lot-stock.read" {
		t.Fatal("lot-stock scope renamed?")
	}
}

func TestWebhookSignature(t *testing.T) {
	body := []byte(`{"id":"evt_1","entity_type":"quote","entity_id":"q","action":"sold","occurred_at":"2026-08-14T12:00:00Z"}`)
	now := time.Unix(1_800_000_000, 0)
	sig := ComputeWebhookSignature("whsec", now, body)
	header := "t=1800000000,v1=" + sig
	if err := verifyWebhookSignatureAt("whsec", header, body, now, 0); err != nil {
		t.Fatal(err)
	}
	if err := verifyWebhookSignatureAt("whsec", header, []byte("tampered"), now, 0); !errors.Is(err, ErrWebhookSignature) {
		t.Fatalf("tampered body should fail, got %v", err)
	}
	if err := verifyWebhookSignatureAt("whsec", header, body, now.Add(time.Hour), 0); !errors.Is(err, ErrWebhookSignature) {
		t.Fatalf("old timestamp should fail, got %v", err)
	}
	if WebhookEventHeader != "X-ShedCloud-Event" || WebhookDeliveryHeader != "X-ShedCloud-Delivery" {
		t.Fatal("ERP webhook header names changed")
	}
}

func TestIdempotencyKeyOnCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"id":"c1","name":"Sam River","email":"sam@example.com","preferred_language":"en","status":"active","created_at":"2026-08-19T00:00:00Z","updated_at":"2026-08-19T00:00:00Z"}`)
	cu, err := c.Customers.Create(context.Background(), CustomerCreateRequest{Name: "Sam River", Email: "sam@example.com"}, WithIdempotencyKey("k-1"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/customers")
	if rec.Header.Get("Idempotency-Key") != "k-1" || cu.ID != "c1" {
		t.Fatalf("idem=%q cu=%+v", rec.Header.Get("Idempotency-Key"), cu)
	}
	body := rec.json(t)
	if body["name"] != "Sam River" || body["phone"] != nil {
		t.Fatalf("body = %v (empty fields must be omitted: the ERP rejects unknown fields only, but keeps omitempty)", body)
	}
}
