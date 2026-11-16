package partnerapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// seen is what the fake ERP recorded for the last request.
type seen struct {
	Method string
	Path   string
	Query  map[string][]string
	Header http.Header
	Body   []byte
}

func (s *seen) json(t *testing.T) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(s.Body, &m); err != nil {
		t.Fatalf("request body is not a JSON object: %v (%q)", err, s.Body)
	}
	return m
}

// fakeERP serves status+body for every request and records the last one.
// Use apiKey auth ("sc_live_test").
func fakeERP(t *testing.T, status int, body string) (*Client, *seen) {
	t.Helper()
	rec := &seen{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		rec.Method, rec.Path, rec.Query, rec.Header, rec.Body = r.Method, r.URL.Path, r.URL.Query(), r.Header.Clone(), raw
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	c, err := New(Options{BaseURL: srv.URL + "/", Auth: Auth{APIKey: "sc_live_test"}})
	if err != nil {
		t.Fatal(err)
	}
	return c, rec
}

func expect(t *testing.T, rec *seen, method, path string) {
	t.Helper()
	if rec.Method != method || rec.Path != path {
		t.Fatalf("request = %s %s, want %s %s", rec.Method, rec.Path, method, path)
	}
	if got := rec.Header.Get("Authorization"); got != "Bearer sc_live_test" {
		t.Fatalf("Authorization = %q", got)
	}
}

func q(rec *seen, key string) string {
	if v := rec.Query[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}
