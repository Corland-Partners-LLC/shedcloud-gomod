package partnerapi_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi"
)

func TestSiteEventsTrack(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/partner/v1/site-events" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(body, &got)
		if got["session_id"] != "sess-1" || got["visitor_id"] != "visitor-1" || got["site_host"] != "lelandssheds.com" {
			t.Errorf("envelope = %s", body)
		}
		// Forwarded end-client context must go out snake_case on the envelope.
		if got["client_ip"] != "203.0.113.9" || got["client_user_agent"] != "Mozilla/5.0 Safari/604.1" {
			t.Errorf("client context = %s", body)
		}
		events, _ := got["events"].([]any)
		if len(events) != 3 {
			t.Fatalf("events = %v", events)
		}
		first, _ := events[0].(map[string]any)
		if first["event_type"] != "page.view" || first["page"] != "/sheds/lofted-barn" {
			t.Errorf("first event = %v", first)
		}
		third, _ := events[2].(map[string]any)
		customer, _ := third["customer"].(map[string]any)
		if third["event_type"] != "customer.profile" || customer["email"] != "jess@example.com" || customer["first_name"] != "Jess" {
			t.Errorf("customer event = %v", third)
		}
		payment, _ := third["payment"].(map[string]any)
		if payment["payment_type"] != "CASH" {
			t.Errorf("payment snapshot = %v", third)
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "accepted": 3})
	}))
	t.Cleanup(srv.Close)

	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_testkey"},
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.SiteEvents.Track(context.Background(), partnerapi.SiteEventsTrackRequest{
		SessionID:       "sess-1",
		VisitorID:       "visitor-1",
		SiteHost:        "lelandssheds.com",
		ClientIP:        "203.0.113.9",
		ClientUserAgent: "Mozilla/5.0 Safari/604.1",
		Events: []partnerapi.SiteEventInput{
			{EventType: "page.view", Page: "/sheds/lofted-barn"},
			{EventType: "cta.click", Payload: map[string]string{"cta": "design-your-own"}},
			{
				EventType: "customer.profile",
				Customer:  &partnerapi.SiteEventCustomer{FirstName: "Jess", LastName: "Miller", Email: "jess@example.com", Phone: "+18175550123"},
				Payment:   &partnerapi.SiteEventPayment{PaymentType: "CASH"},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Success || res.Accepted != 3 {
		t.Fatalf("res = %+v", res)
	}
}

func TestSiteEventsListAndEach(t *testing.T) {
	t.Parallel()
	call := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/partner/v1/site-events" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("types") != "page.view,cta.click" {
			t.Errorf("types = %q", r.URL.Query().Get("types"))
		}
		if r.URL.Query().Get("sessionId") != "sess-1" {
			t.Errorf("sessionId = %q", r.URL.Query().Get("sessionId"))
		}
		call++
		if call == 1 {
			if r.URL.Query().Has("cursor") {
				t.Errorf("first call must not send a cursor, got %q", r.URL.Query().Get("cursor"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{{
					"eventId": "ev-1", "eventType": "page.view",
					"occurredAt": "2026-07-01T12:00:00Z", "sessionId": "sess-1",
					"source": "marketing",
				}},
				"nextCursor": "cur-1",
				"hasMore":    true,
			})
			return
		}
		if r.URL.Query().Get("cursor") != "cur-1" {
			t.Errorf("cursor = %q", r.URL.Query().Get("cursor"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{
				"eventId": "ev-2", "eventType": "cta.click",
				"occurredAt": "2026-07-01T12:01:00Z", "sessionId": "sess-1",
			}},
			"hasMore": false,
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

	var seen []string
	err = client.SiteEvents.Each(context.Background(), partnerapi.SiteEventListParams{
		SessionID: "sess-1",
		Types:     []string{"page.view", "cta.click"},
	}, func(ev partnerapi.SiteEventItem) error {
		seen = append(seen, ev.EventID)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(seen) != 2 || seen[0] != "ev-1" || seen[1] != "ev-2" {
		t.Fatalf("seen = %v", seen)
	}
}
