package partnerapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi"
)

func TestRevenueForecastGet(t *testing.T) {
	t.Parallel()
	details := true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/partner/v1/revenue-forecast" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("period") != partnerapi.PeriodLastMonth {
			t.Fatalf("period = %q", q.Get("period"))
		}
		if q.Get("locationId") != "66c00443c2d8aa83c5757dcf" {
			t.Fatalf("locationId = %q", q.Get("locationId"))
		}
		if q.Get("details") != "true" {
			t.Fatalf("details = %q", q.Get("details"))
		}
		if q.Get("detailsFor") != "actual" {
			t.Fatalf("detailsFor = %q", q.Get("detailsFor"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"period": map[string]any{
				"key":      "last_month",
				"start":    "2026-07-01",
				"end":      "2026-07-31",
				"timezone": "America/Chicago",
			},
			"revenueForecast": map[string]any{
				"actual": map[string]any{
					"amount": 9200, "basePrice": 8000, "upgrades": 1200,
					"removedUpgrades": 0, "units": 1,
					"details": []any{map[string]any{
						"workOrderId": "wo-1", "orderId": "so-1",
						"qualifiedBy": "wo_delivered", "bucketDate": "2026-07-22",
						"dateSource": "wo_actual_delivered", "amount": 9200,
						"amountSource": "sales_ledger",
					}},
				},
				"afdTransferred":     map[string]any{"amount": 0, "units": 0},
				"forecastNextPeriod": map[string]any{"period": map[string]any{"key": "this_month"}, "amount": 0, "units": 0},
				"forecastPlus2":      map[string]any{"period": map[string]any{"key": "next_month"}, "amount": 0, "units": 0},
			},
			"authorizedDelivery": map[string]any{
				"authorized": map[string]any{"amount": 0, "units": 0},
				"weeks":      []any{},
			},
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

	out, err := client.RevenueForecast.Get(context.Background(), partnerapi.RevenueForecastParams{
		Period:     partnerapi.PeriodLastMonth,
		LocationID: "66c00443c2d8aa83c5757dcf",
		Details:    &details,
		DetailsFor: "actual",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Period.Key != partnerapi.PeriodLastMonth || out.Period.Timezone != "America/Chicago" {
		t.Fatalf("period = %+v", out.Period)
	}
	if out.RevenueForecast.Actual.Amount != 9200 || out.RevenueForecast.Actual.Units != 1 {
		t.Fatalf("actual = %+v", out.RevenueForecast.Actual)
	}
	if len(out.RevenueForecast.Actual.Details) != 1 {
		t.Fatalf("details = %+v", out.RevenueForecast.Actual.Details)
	}
	d := out.RevenueForecast.Actual.Details[0]
	if d.WorkOrderID != "wo-1" || d.OrderID != "so-1" || d.Amount != 9200 || d.AmountSource != "sales_ledger" {
		t.Fatalf("detail = %+v", d)
	}
}
