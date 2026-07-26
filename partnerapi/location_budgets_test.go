package partnerapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi"
)

func TestLocationBudgetsList(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/partner/v1/location-budgets" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("year") != "2026" {
			t.Fatalf("year = %q", r.URL.Query().Get("year"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":  []any{map[string]any{"id": "bud-1", "locationId": "loc-1", "year": 2026}},
			"page":  1,
			"limit": 25,
			"total": 1,
		})
	}))
	t.Cleanup(srv.Close)

	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: srv.URL,
		Auth:    partnerapi.Auth{APIKey: "sc_live_test"},
	})
	if err != nil {
		t.Fatal(err)
	}

	out, err := client.LocationBudgets.List(context.Background(), partnerapi.LocationBudgetListParams{
		Year: 2026,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Data) != 1 || out.Data[0].ID != "bud-1" {
		t.Fatalf("unexpected list: %+v", out.Data)
	}
}
