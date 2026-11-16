// Smoke-test the Partner API Go client against an ERP host.
//
// Usage:
//
//	go run ./scripts/test-partner-api -base-url https://<erp-host> -api-key sc_live_...
//
// The ERP host is required (it is not built in). Flags:
//
//	-base-url   ERP Partner API host (or SHEDCLOUD_BASE_URL)
//	-api-key    API key (or SHEDCLOUD_API_KEY)
//	-resource   stock-units | leads | quotes | orders | work-orders | locations | customers | products (default: stock-units)
//	-limit      page size (default: 5)
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi"
)

func main() {
	apiKey := flag.String("api-key", os.Getenv("SHEDCLOUD_API_KEY"), "Partner API key (sc_live_…)")
	baseURL := flag.String("base-url", os.Getenv("SHEDCLOUD_BASE_URL"), "ERP Partner API host (required)")
	resource := flag.String("resource", "stock-units", "stock-units | leads | quotes | orders | work-orders | locations | customers | products")
	limit := flag.Int("limit", 5, "Page size")
	flag.Parse()

	if strings.TrimSpace(*apiKey) == "" || strings.TrimSpace(*baseURL) == "" {
		fmt.Fprintln(os.Stderr, "ERROR: pass -base-url and -api-key (or SHEDCLOUD_BASE_URL / SHEDCLOUD_API_KEY)")
		os.Exit(1)
	}
	client, err := partnerapi.New(partnerapi.Options{
		BaseURL: *baseURL,
		Auth:    partnerapi.Auth{APIKey: strings.TrimSpace(*apiKey)},
	})
	if err != nil {
		fail(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	me, err := client.Me(ctx)
	if err != nil {
		fail(err)
	}
	fmt.Printf("credential %s (%s), %d scopes\n", me.CredentialID, me.Kind, len(me.Scopes))

	lp := partnerapi.ListParams{PageSize: *limit}
	var n int
	var total int64
	switch *resource {
	case "stock-units":
		p, e := client.StockUnits.List(ctx, partnerapi.StockUnitListParams{ListParams: lp})
		err, n, total = e, countOf(p), totalOf(p)
	case "leads":
		p, e := client.Leads.List(ctx, partnerapi.LeadListParams{ListParams: lp})
		err, n, total = e, countOf(p), totalOf(p)
	case "quotes":
		p, e := client.Quotes.List(ctx, partnerapi.QuoteListParams{ListParams: lp})
		err, n, total = e, countOf(p), totalOf(p)
	case "orders":
		p, e := client.Orders.List(ctx, partnerapi.OrderListParams{ListParams: lp})
		err, n, total = e, countOf(p), totalOf(p)
	case "work-orders":
		p, e := client.WorkOrders.List(ctx, partnerapi.WorkOrderListParams{ListParams: lp})
		err, n, total = e, countOf(p), totalOf(p)
	case "locations":
		p, e := client.Locations.List(ctx, partnerapi.LocationListParams{ListParams: lp})
		err, n, total = e, countOf(p), totalOf(p)
	case "customers":
		p, e := client.Customers.List(ctx, partnerapi.CustomerListParams{ListParams: lp})
		err, n, total = e, countOf(p), totalOf(p)
	case "products":
		p, e := client.Products.List(ctx, partnerapi.ProductListParams{ListParams: lp})
		err, n, total = e, countOf(p), totalOf(p)
	default:
		fail(fmt.Errorf("unknown -resource %q", *resource))
	}
	if err != nil {
		fail(err)
	}
	fmt.Printf("%s: %d returned, %d total\n", *resource, n, total)
}

func countOf[T any](p *partnerapi.Page[T]) int {
	if p == nil {
		return 0
	}
	return len(p.Items)
}

func totalOf[T any](p *partnerapi.Page[T]) int64 {
	if p == nil {
		return 0
	}
	return p.Pagination.Total
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "ERROR:", err)
	os.Exit(1)
}
