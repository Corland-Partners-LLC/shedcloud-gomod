package partnerapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *partnerapi.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_testkey"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestUsersListAndGet(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/partner/v1/users":
			if r.URL.Query().Get("search") != "alex" || r.URL.Query().Get("active") != "true" {
				t.Errorf("query = %q", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{{
					"id": "user-1", "name": "Alex Rep", "active": true,
					"locationIds": []string{"loc-1"}, "allLocations": false, "inLeadRoutingPool": true,
				}},
				"page": 1, "limit": 50, "total": 1,
			})
		case "/partner/v1/users/user-1":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "user-1", "active": true, "locationIds": []string{}, "allLocations": true,
			})
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
		}
	})

	active := true
	list, err := client.Users.List(context.Background(), partnerapi.UserListParams{Search: "alex", Active: &active})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 1 || !list.Data[0].InLeadRoutingPool {
		t.Fatalf("unexpected list: %+v", list)
	}

	user, err := client.Users.Get(context.Background(), "user-1")
	if err != nil {
		t.Fatal(err)
	}
	if !user.AllLocations {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestPaymentsListAndOrderPayments(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/partner/v1/payments":
			if r.URL.Query().Get("status") != "paid" || r.URL.Query().Get("createdFrom") != "2026-01-01" {
				t.Errorf("query = %q", r.URL.RawQuery)
			}
		case "/partner/v1/orders/order-1/payments":
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{
				"id": "pay-1", "orderId": "order-1", "amount": 500.0,
				"method": "card", "status": "paid", "providerReference": "cs_test_abc",
			}},
			"page": 1, "limit": 50, "total": 1,
		})
	})

	list, err := client.Payments.List(context.Background(), partnerapi.PaymentListParams{Status: "paid", CreatedFrom: "2026-01-01"})
	if err != nil {
		t.Fatal(err)
	}
	if list.Data[0].ProviderReference != "cs_test_abc" {
		t.Fatalf("unexpected payment: %+v", list.Data[0])
	}

	forOrder, err := client.Orders.Payments(context.Background(), "order-1", partnerapi.PaginationParams{})
	if err != nil {
		t.Fatal(err)
	}
	if forOrder.Data[0].Amount != 500 {
		t.Fatalf("unexpected order payment: %+v", forOrder.Data[0])
	}
}

func TestDocumentsListAndDownload(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/partner/v1/documents":
			q := r.URL.Query()
			if q.Get("entityType") != "order" || q.Get("entityId") != "order-1" || q.Get("type") != "Contract" {
				t.Errorf("query = %q", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{{
					"id": "doc-1", "fileName": "contract.pdf", "type": "Contract", "sizeBytes": 12345,
				}},
				"page": 1, "limit": 50, "total": 1,
			})
		case "/partner/v1/documents/doc-1/download":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"downloadUrl": "https://s3.example.com/signed?sig=abc",
				"fileName":    "contract.pdf",
				"expiresAt":   "2026-07-11T22:30:00Z",
			})
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
		}
	})

	list, err := client.Documents.List(context.Background(), partnerapi.DocumentListParams{
		EntityType: "order", EntityID: "order-1", Type: "Contract",
	})
	if err != nil {
		t.Fatal(err)
	}
	if list.Data[0].SizeBytes != 12345 {
		t.Fatalf("unexpected document: %+v", list.Data[0])
	}

	dl, err := client.Documents.Download(context.Background(), "doc-1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(dl.DownloadURL, "sig=abc") {
		t.Fatalf("unexpected download: %+v", dl)
	}
}

func TestEventsListEachAndRedeliver(t *testing.T) {
	t.Parallel()
	var listCalls int
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/partner/v1/events" && r.Method == http.MethodGet:
			listCalls++
			if listCalls == 1 {
				if got := r.URL.Query().Get("types"); got != "order.created,payment.paid" {
					t.Errorf("types = %q", got)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"data": []map[string]any{
						{"id": "ev-1", "type": "order.created", "occurredAt": "2026-07-11T20:00:00Z", "resourceType": "order", "resourceId": "order-1"},
						{"id": "ev-2", "type": "payment.paid", "occurredAt": "2026-07-11T20:05:00Z", "resourceType": "payment", "resourceId": "pay-1"},
					},
					"nextCursor": "ev-2",
					"hasMore":    true,
				})
				return
			}
			if got := r.URL.Query().Get("cursor"); got != "ev-2" {
				t.Errorf("cursor = %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "ev-3", "type": "order.created", "occurredAt": "2026-07-11T20:10:00Z", "resourceType": "order", "resourceId": "order-2"},
				},
				"nextCursor": "ev-3",
				"hasMore":    false,
			})
		case r.URL.Path == "/partner/v1/events/ev-1/redeliver" && r.Method == http.MethodPost:
			_ = json.NewEncoder(w).Encode(map[string]any{"eventId": "ev-1", "enqueued": 2})
		case r.URL.Path == "/partner/v1/webhook-deliveries":
			if r.URL.Query().Get("eventId") != "ev-1" {
				t.Errorf("eventId = %q", r.URL.Query().Get("eventId"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{{
					"id": "del-1", "subscriptionId": "sub-1", "eventId": "ev-1",
					"eventType": "order.created", "url": "https://partner.example.com/hook",
					"attempt": 1, "statusCode": 200, "ok": true, "deliveredAt": "2026-07-11T20:00:01Z",
				}},
				"page": 1, "limit": 50, "total": 1,
			})
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	var seen []string
	err := client.Events.Each(context.Background(),
		partnerapi.EventListParams{Types: []string{"order.created", "payment.paid"}},
		func(ev partnerapi.EventItem) error {
			seen = append(seen, ev.ID)
			return nil
		})
	if err != nil {
		t.Fatal(err)
	}
	if len(seen) != 3 || seen[0] != "ev-1" || seen[2] != "ev-3" {
		t.Fatalf("seen = %v", seen)
	}
	if listCalls != 2 {
		t.Fatalf("listCalls = %d", listCalls)
	}

	red, err := client.Events.Redeliver(context.Background(), "ev-1")
	if err != nil {
		t.Fatal(err)
	}
	if red.Enqueued != 2 {
		t.Fatalf("enqueued = %d", red.Enqueued)
	}

	deliveries, err := client.Events.Deliveries(context.Background(), partnerapi.WebhookDeliveryListParams{EventID: "ev-1"})
	if err != nil {
		t.Fatal(err)
	}
	if !deliveries.Data[0].OK {
		t.Fatalf("unexpected delivery: %+v", deliveries.Data[0])
	}
}

func TestEventsEachStopsOnCallbackError(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "ev-1", "type": "order.created", "occurredAt": "x", "resourceType": "order", "resourceId": "o"},
			},
			"nextCursor": "ev-1",
			"hasMore":    true,
		})
	})

	wantErr := errors.New("stop here")
	err := client.Events.Each(context.Background(), partnerapi.EventListParams{}, func(partnerapi.EventItem) error {
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v", err)
	}
}

func TestConfiguratorSessionCreate(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/partner/v1/configurator-sessions" || r.Method != http.MethodPost {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Idempotency-Key") != "sess-idem-1" {
			t.Errorf("Idempotency-Key = %q", r.Header.Get("Idempotency-Key"))
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		if body["customerId"] != "cust-1" || body["serialNumber"] != "SC-2024-00123" {
			t.Errorf("body = %s", raw)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"sessionId":  "ps-1",
			"launchUrl":  "https://go.shedcloud.com/api/external/partner/launch?token=tok",
			"expiresAt":  "2026-07-11T20:15:00Z",
			"customerId": "cust-1",
		})
	})

	session, err := client.ConfiguratorSessions.Create(context.Background(), partnerapi.ConfiguratorSessionCreateRequest{
		CustomerID:   "cust-1",
		LocationID:   "loc-1",
		SerialNumber: "SC-2024-00123",
		ReturnURL:    "https://crm.example.com/return",
	}, partnerapi.WithIdempotencyKey("sess-idem-1"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(session.LaunchURL, "token=tok") {
		t.Fatalf("unexpected session: %+v", session)
	}
}

func TestUpdateSendsIfMatchAndSubResources(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPatch && r.URL.Path == "/partner/v1/orders/order-1":
			if r.Header.Get("If-Match") != "7" {
				t.Errorf("If-Match = %q", r.Header.Get("If-Match"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "order-1", "version": 8})
		case r.URL.Path == "/partner/v1/orders/order-1/status-history":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{{
					"status": "Submitted", "previousStatus": "Unsubmitted",
					"changedAt": "2026-07-10T12:00:00Z", "actor": map[string]any{"name": "Alex Rep"},
				}},
				"page": 1, "limit": 50, "total": 1,
			})
		case r.URL.Path == "/partner/v1/orders/order-1/line-items":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{{
					"id": "li-1", "name": "Extra Window", "quantity": 2.0, "amount": 150.0,
					"status": "added", "isStandardFeature": false,
				}},
				"totals":        map[string]any{"included": 4, "added": 1, "removed": 0},
				"configuration": map[string]any{"model": "Lofted Barn", "siding": "LP Smart Panel"},
			})
		case r.URL.Path == "/partner/v1/orders/order-1/contract":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"orderId": "order-1", "status": "partially_signed",
				"customerSigned": true, "salespersonSigned": false,
			})
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	order, err := client.Orders.Update(context.Background(), "order-1",
		partnerapi.OrderPatchRequest{CustomerPhone: "555-0100"},
		partnerapi.WithIfMatch(7))
	if err != nil {
		t.Fatal(err)
	}
	if order.Version != 8 {
		t.Fatalf("version = %d", order.Version)
	}

	history, err := client.Orders.StatusHistory(context.Background(), "order-1", partnerapi.PaginationParams{})
	if err != nil {
		t.Fatal(err)
	}
	if history.Data[0].PreviousStatus != "Unsubmitted" {
		t.Fatalf("unexpected history: %+v", history.Data[0])
	}

	lineItems, err := client.Orders.LineItems(context.Background(), "order-1")
	if err != nil {
		t.Fatal(err)
	}
	if lineItems.Totals.Added != 1 || lineItems.Configuration == nil || lineItems.Configuration.Model != "Lofted Barn" {
		t.Fatalf("unexpected line items: %+v", lineItems)
	}

	contract, err := client.Orders.Contract(context.Background(), "order-1")
	if err != nil {
		t.Fatal(err)
	}
	if contract.Status != "partially_signed" || !contract.CustomerSigned {
		t.Fatalf("unexpected contract: %+v", contract)
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	t.Parallel()
	secret := "whsec_test_secret"
	body := []byte(`{"id":"ev-1","type":"order.created"}`)

	now := time.Now()
	header := fmt.Sprintf("t=%d,v1=%s", now.Unix(), partnerapi.ComputeWebhookSignature(secret, now, body))

	if err := partnerapi.VerifyWebhookSignature(secret, header, body, 0); err != nil {
		t.Fatalf("valid signature rejected: %v", err)
	}

	if err := partnerapi.VerifyWebhookSignature(secret, header, append(body, 'x'), 0); !errors.Is(err, partnerapi.ErrWebhookSignature) {
		t.Fatalf("tampered body: err = %v", err)
	}

	if err := partnerapi.VerifyWebhookSignature("whsec_other", header, body, 0); !errors.Is(err, partnerapi.ErrWebhookSignature) {
		t.Fatalf("wrong secret: err = %v", err)
	}

	stale := now.Add(-10 * time.Minute)
	staleHeader := fmt.Sprintf("t=%d,v1=%s", stale.Unix(), partnerapi.ComputeWebhookSignature(secret, stale, body))
	if err := partnerapi.VerifyWebhookSignature(secret, staleHeader, body, 0); !errors.Is(err, partnerapi.ErrWebhookSignature) {
		t.Fatalf("stale timestamp: err = %v", err)
	}
	// The same stale header verifies with a wider tolerance.
	if err := partnerapi.VerifyWebhookSignature(secret, staleHeader, body, time.Hour); err != nil {
		t.Fatalf("stale within tolerance rejected: %v", err)
	}

	for _, bad := range []string{"", "nonsense", "t=abc,v1=00", "v1=deadbeef"} {
		if err := partnerapi.VerifyWebhookSignature(secret, bad, body, 0); !errors.Is(err, partnerapi.ErrWebhookSignature) {
			t.Fatalf("malformed header %q: err = %v", bad, err)
		}
	}
}
