package partnerapi_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi"
)

func TestResolveBaseURL(t *testing.T) {
	t.Parallel()
	if got := partnerapi.ResolveBaseURL("", ""); got != partnerapi.HostProduction {
		t.Fatalf("default = %q, want production", got)
	}
	if got := partnerapi.ResolveBaseURL("", partnerapi.EnvironmentSandbox); got != partnerapi.HostSandbox {
		t.Fatalf("sandbox = %q", got)
	}
	if got := partnerapi.ResolveBaseURL("https://localhost:8080/", partnerapi.EnvironmentSandbox); got != "https://localhost:8080" {
		t.Fatalf("baseUrl override = %q", got)
	}
}

func TestClientDefaultsToProduction(t *testing.T) {
	t.Parallel()
	var hit string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = r.URL.Path
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []any{}, "page": 1, "limit": 50, "total": 0})
	}))
	t.Cleanup(srv.Close)

	// Override via BaseURL so we don't hit the real production host.
	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.BaseURL != srv.URL {
		t.Fatalf("BaseURL = %q", client.BaseURL)
	}
	if _, err := client.LotStock.List(context.Background(), partnerapi.LotStockListParams{}); err != nil {
		t.Fatal(err)
	}
	if hit != "/partner/v1/lot-stock" {
		t.Fatalf("path = %q", hit)
	}
}

func TestLotStockListSendsAPIKeyAndQuery(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sc_live_testkey" {
			t.Errorf("Authorization = %q", r.Header.Get("Authorization"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("limit = %q", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("purchaseType") != "Lot Stock" {
			t.Errorf("purchaseType = %q", r.URL.Query().Get("purchaseType"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{
				"id":    "abc",
				"title": "10x16 Lofted Barn",
				"price": 8995,
				"sold":  false,
			}},
			"page":  1,
			"limit": 10,
			"total": 1,
		})
	}))
	t.Cleanup(srv.Close)

	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_testkey"},
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.LotStock.List(context.Background(), partnerapi.LotStockListParams{
		PaginationParams: partnerapi.PaginationParams{Limit: 10},
		PurchaseType:     "Lot Stock",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 1 || len(res.Data) != 1 || res.Data[0].Title != "10x16 Lofted Barn" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestOAuthTokenExchangeAndCache(t *testing.T) {
	t.Parallel()
	var tokenCalls atomic.Int32
	var resourceCalls atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth/token":
			tokenCalls.Add(1)
			if r.Method != http.MethodPost {
				t.Errorf("method = %s", r.Method)
			}
			user, pass, ok := r.BasicAuth()
			if !ok || user != "sc_client_abc" || pass != "sc_secret_xyz" {
				t.Errorf("basic auth = %q %q ok=%v", user, pass, ok)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "sc_at_cached",
				"token_type":   "Bearer",
				"expires_in":   3600,
				"scope":        "partner-api.lot-stock.read",
			})
		default:
			resourceCalls.Add(1)
			if r.Header.Get("Authorization") != "Bearer sc_at_cached" {
				t.Errorf("Authorization = %q", r.Header.Get("Authorization"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []any{}, "page": 1, "limit": 50, "total": 0})
		}
	}))
	t.Cleanup(srv.Close)

	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth: partnerapi.Auth{
			ClientID:     "sc_client_abc",
			ClientSecret: "sc_secret_xyz",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.LotStock.List(context.Background(), partnerapi.LotStockListParams{}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Leads.List(context.Background(), partnerapi.SalesListParams{PaginationParams: partnerapi.PaginationParams{Page: 1}}); err != nil {
		t.Fatal(err)
	}

	if tokenCalls.Load() != 1 {
		t.Fatalf("token exchanges = %d, want 1", tokenCalls.Load())
	}
	if resourceCalls.Load() != 2 {
		t.Fatalf("resource calls = %d, want 2", resourceCalls.Load())
	}
}

func TestPartnerAPIErrorOnForbidden(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"credential is not authorized for this scope"}`))
	}))
	t.Cleanup(srv.Close)

	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_testkey"},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Orders.Get(context.Background(), "missing")
	apiErr, ok := err.(*partnerapi.Error)
	if !ok {
		t.Fatalf("err type = %T, want *partnerapi.Error", err)
	}
	if !apiErr.IsForbidden() || apiErr.Status != 403 {
		t.Fatalf("unexpected error: %+v", apiErr)
	}
	if !strings.Contains(apiErr.Message, "not authorized") {
		t.Fatalf("message = %q", apiErr.Message)
	}
}

func TestLeadPatch(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/partner/v1/leads/lead-1" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]string
		_ = json.Unmarshal(body, &got)
		if got["salespersonName"] != "Alex Rep" || got["salesLocation"] != "loc-1" {
			t.Errorf("body = %s", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":          "lead-1",
			"salesperson": map[string]string{"name": "Alex Rep"},
		})
	}))
	t.Cleanup(srv.Close)

	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_testkey"},
	})
	if err != nil {
		t.Fatal(err)
	}

	lead, err := client.Leads.Update(context.Background(), "lead-1", partnerapi.LeadPatchRequest{
		SalespersonName: "Alex Rep",
		SalesLocation:   "loc-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if lead.ID != "lead-1" {
		t.Fatalf("id = %q", lead.ID)
	}
}

func TestQuoteGetIncludesSerialAndWorkOrder(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/partner/v1/quotes/quote-1" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":           "quote-1",
			"orderNumber":  13847,
			"serialNumber": "SC-2024-00123",
			"workOrderId":  "wo-1",
			"converted":    false,
			"pricing":      map[string]any{"subtotal": 8200, "total": 8995, "paymentType": "rto"},
		})
	}))
	t.Cleanup(srv.Close)

	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_testkey"},
	})
	if err != nil {
		t.Fatal(err)
	}

	quote, err := client.Quotes.Get(context.Background(), "quote-1")
	if err != nil {
		t.Fatal(err)
	}
	if quote.SerialNumber != "SC-2024-00123" || quote.WorkOrderID != "wo-1" || quote.Converted {
		t.Fatalf("unexpected quote: %+v", quote)
	}
}

func TestQuoteListConvertedQuery(t *testing.T) {
	t.Parallel()
	converted := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("converted") != "false" {
			t.Errorf("converted = %q", r.URL.Query().Get("converted"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []any{}, "page": 1, "limit": 50, "total": 0})
	}))
	t.Cleanup(srv.Close)

	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_testkey"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Quotes.List(context.Background(), partnerapi.QuoteListParams{Converted: &converted}); err != nil {
		t.Fatal(err)
	}
}

func TestNewRequiresAuth(t *testing.T) {
	t.Parallel()
	if _, err := partnerapi.New(partnerapi.Options{}); err == nil {
		t.Fatal("expected error for empty auth")
	}
	if _, err := partnerapi.New(partnerapi.Options{
		Auth: partnerapi.Auth{APIKey: "k", ClientID: "c", ClientSecret: "s"},
	}); err == nil {
		t.Fatal("expected error when both auth modes set")
	}
}
