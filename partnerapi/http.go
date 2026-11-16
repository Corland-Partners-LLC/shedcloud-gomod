package partnerapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type httpClient struct {
	baseURL   string
	getToken  func(context.Context) (string, error)
	httpDo    func(*http.Request) (*http.Response, error)
	userAgent string
	timeout   time.Duration
}

// RequestOption customizes one request (e.g. WithIdempotencyKey). Accepted by
// the create/convert methods as trailing variadic arguments.
type RequestOption func(*http.Request)

// WithIdempotencyKey sets the Idempotency-Key header on a create request. A
// retried request with the same key (within 24 hours) replays the stored
// response instead of creating a duplicate resource.
func WithIdempotencyKey(key string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Idempotency-Key", key)
	}
}

// withAccept overrides the Accept header (used for the PDF download).
func withAccept(v string) RequestOption {
	return func(req *http.Request) { req.Header.Set("Accept", v) }
}

// request sends one call and decodes a JSON response into out (when non-nil).
func (c *httpClient) request(ctx context.Context, method, path string, query any, body any, out any, opts ...RequestOption) error {
	raw, err := c.do(ctx, method, path, query, body, opts...)
	if err != nil {
		return err
	}
	if out == nil || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("partnerapi: decode response: %w", err)
	}
	return nil
}

// do sends one call and returns the raw 2xx response body.
func (c *httpClient) do(ctx context.Context, method, path string, query any, body any, opts ...RequestOption) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}
	if query != nil {
		q := encodeQuery(query)
		if q != "" {
			u.RawQuery = q
		}
	}

	var bodyReader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(raw)
	}

	reqCtx := ctx
	var cancel context.CancelFunc
	if c.timeout > 0 {
		reqCtx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(reqCtx, method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	token, err := c.getToken(reqCtx)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, newError("no access token available", http.StatusUnauthorized, nil, "")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.httpDo(req)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			return nil, newError(fmt.Sprintf("request timed out after %s", c.timeout), http.StatusRequestTimeout, nil, "")
		}
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	var parsed any
	if len(raw) > 0 {
		if json.Unmarshal(raw, &parsed) != nil {
			parsed = string(raw)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, code := extractError(parsed, resp.Status)
		return nil, newError(msg, resp.StatusCode, parsed, code)
	}

	return raw, nil
}

// encodeQuery turns a params struct into a query string using json tags as
// keys. Zero values and nil pointers are omitted; string slices are joined
// with commas (the ERP reads statuses=a,b).
func encodeQuery(v any) string {
	if v == nil {
		return ""
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return ""
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return ""
	}

	params := url.Values{}
	encodeStructQuery(rv, params)
	return params.Encode()
}

func encodeStructQuery(rv reflect.Value, params url.Values) {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Anonymous {
			encodeStructQuery(rv.Field(i), params)
			continue
		}
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		if name == "" {
			continue
		}
		fv := rv.Field(i)
		if !fv.IsValid() || !fv.CanInterface() {
			continue
		}
		fromPtr := false
		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				continue
			}
			fromPtr = true
			fv = fv.Elem()
		}
		// Non-nil pointers are always encoded (so converted=false is sent).
		// Plain zero values (empty string, 0, false) are omitted.
		if !fromPtr && isZeroValue(fv) {
			continue
		}
		params.Set(name, formatQueryValue(fv))
	}
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	default:
		return false
	}
}

func formatQueryValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Slice:
		parts := make([]string, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			parts = append(parts, formatQueryValue(v.Index(i)))
		}
		return strings.Join(parts, ",")
	default:
		return fmt.Sprint(v.Interface())
	}
}
