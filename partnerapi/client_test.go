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

func TestLocationListBooleanFilters(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/partner/v1/locations" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("salesLot") != "true" {
			t.Errorf("salesLot = %q", r.URL.Query().Get("salesLot"))
		}
		// Non-nil false pointers must still be sent.
		if r.URL.Query().Get("plant") != "false" {
			t.Errorf("plant = %q", r.URL.Query().Get("plant"))
		}
		if r.URL.Query().Has("active") {
			t.Errorf("nil active must not be sent, got %q", r.URL.Query().Get("active"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{
				"id": "loc-1", "name": "Dallas Lot", "code": "DAL01",
				"active": true, "salesLot": true, "plant": false,
			}},
			"page": 1, "limit": 50, "total": 1,
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

	salesLot, plant := true, false
	res, err := client.Locations.List(context.Background(), partnerapi.LocationListParams{
		SalesLot: &salesLot,
		Plant:    &plant,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Data) != 1 || res.Data[0].Code != "DAL01" || !res.Data[0].SalesLot {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestLocationCreate(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/partner/v1/locations" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(body, &got)
		if got["name"] != "Fort Worth Lot" || got["code"] != "FTW01" || got["salesLot"] != true {
			t.Errorf("body = %s", body)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "loc-2", "name": "Fort Worth Lot", "code": "FTW01",
			"active": true, "salesLot": true, "plant": false,
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

	salesLot := true
	loc, err := client.Locations.Create(context.Background(), partnerapi.LocationCreateRequest{
		Name:     "Fort Worth Lot",
		Code:     "FTW01",
		Address:  "500 Elm St",
		SalesLot: &salesLot,
	})
	if err != nil {
		t.Fatal(err)
	}
	if loc.ID != "loc-2" {
		t.Fatalf("id = %q", loc.ID)
	}
}

func TestQuoteCreate(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/partner/v1/quotes" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(body, &got)
		if got["serialNumber"] != "SC-2024-00123" {
			t.Errorf("body = %s", body)
		}
		customer, _ := got["customer"].(map[string]any)
		if customer["email"] != "jane@example.com" || customer["name"] != "Jane Doe" {
			t.Errorf("customer = %v", customer)
		}
		delivery, _ := got["deliveryAddress"].(map[string]any)
		if delivery["city"] != "Dallas" || delivery["zipCode"] != "75201" {
			t.Errorf("deliveryAddress = %v", delivery)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "quote-1", "orderNumber": 13848, "status": "Open",
			"serialNumber": "SC-2024-00123", "workOrderId": "wo-1", "converted": false,
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

	quote, err := client.Quotes.Create(context.Background(), partnerapi.QuoteCreateRequest{
		SerialNumber: "SC-2024-00123",
		Customer: partnerapi.QuoteCreateCustomer{
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Phone: "555-0100",
		},
		DeliveryAddress: &partnerapi.QuoteCreateDeliveryAddress{
			Address: "42 Oak Ave",
			City:    "Dallas",
			State:   "TX",
			ZipCode: "75201",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if quote.ID != "quote-1" || quote.WorkOrderID != "wo-1" {
		t.Fatalf("quote = %+v", quote)
	}
}

func TestLeadCreateWithIdempotencyKey(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/partner/v1/leads" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Idempotency-Key"); got != "idem-123" {
			t.Errorf("Idempotency-Key = %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(body, &got)
		if got["locationId"] != "66c00443c2d8aa83c5757dcf" {
			t.Errorf("body = %s", body)
		}
		customer, _ := got["customer"].(map[string]any)
		if customer["email"] != "jane@example.com" {
			t.Errorf("customer = %v", customer)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "lead-1", "orderNumber": 13902, "status": "Open",
			"salesperson": map[string]any{"name": "John Rep", "email": "john@dealer.com"},
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

	lead, err := client.Leads.Create(context.Background(), partnerapi.LeadCreateRequest{
		LocationID: "66c00443c2d8aa83c5757dcf",
		Customer: partnerapi.LeadCreateCustomer{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		},
	}, partnerapi.WithIdempotencyKey("idem-123"))
	if err != nil {
		t.Fatal(err)
	}
	if lead.ID != "lead-1" || lead.Salesperson.Email != "john@dealer.com" {
		t.Fatalf("lead = %+v", lead)
	}
}

func TestQuoteConvert(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/partner/v1/quotes/quote-1/convert" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(body, &got)
		if got["salespersonEmail"] != "newrep@dealer.com" {
			t.Errorf("body = %s", body)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "order-1", "orderNumber": 13849, "status": "Unsubmitted",
			"sourceQuoteId": "quote-1", "sourceQuoteNumber": 13848,
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

	order, err := client.Quotes.Convert(context.Background(), "quote-1", partnerapi.QuoteConvertRequest{
		SalespersonEmail: "newrep@dealer.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if order.ID != "order-1" || order.Status != "Unsubmitted" || order.SourceQuoteID != "quote-1" {
		t.Fatalf("order = %+v", order)
	}
}

func TestCustomerCreateAndPatch(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			if r.URL.Path != "/partner/v1/customers" {
				t.Errorf("path = %s", r.URL.Path)
			}
			body, _ := io.ReadAll(r.Body)
			var got map[string]string
			_ = json.Unmarshal(body, &got)
			if got["email"] != "jane@example.com" || got["name"] != "Jane Smith" {
				t.Errorf("body = %s", body)
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "cust-1", "name": "Jane Smith", "email": "jane@example.com", "active": true,
			})
		case http.MethodPatch:
			if r.URL.Path != "/partner/v1/customers/cust-1" {
				t.Errorf("path = %s", r.URL.Path)
			}
			body, _ := io.ReadAll(r.Body)
			var got map[string]any
			_ = json.Unmarshal(body, &got)
			if got["phone"] != "555-0100" || len(got) != 1 {
				t.Errorf("body = %s", body)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "cust-1", "phone": "555-0100", "active": true,
			})
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_testkey"},
	})
	if err != nil {
		t.Fatal(err)
	}

	created, err := client.Customers.Create(context.Background(), partnerapi.CustomerCreateRequest{
		Email: "jane@example.com",
		Name:  "Jane Smith",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != "cust-1" {
		t.Fatalf("id = %q", created.ID)
	}

	updated, err := client.Customers.Update(context.Background(), "cust-1", partnerapi.CustomerPatchRequest{
		Phone: "555-0100",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Phone != "555-0100" {
		t.Fatalf("phone = %q", updated.Phone)
	}
}

func TestProductList(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/partner/v1/products" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("search") != "barn" || r.URL.Query().Get("active") != "true" {
			t.Errorf("query = %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{
				"id": "prod-1", "name": "12x24 Lofted Barn", "sku": "LB-1224",
				"price": 8995.5, "width": 12, "length": 24, "active": true,
			}},
			"page": 1, "limit": 50, "total": 1,
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

	active := true
	res, err := client.Products.List(context.Background(), partnerapi.ProductListParams{
		Search: "barn",
		Active: &active,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Data) != 1 || res.Data[0].Price != 8995.5 || res.Data[0].SKU != "LB-1224" {
		t.Fatalf("unexpected response: %+v", res)
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
