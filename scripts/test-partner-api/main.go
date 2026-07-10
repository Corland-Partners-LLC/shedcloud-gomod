// Smoke-test the Partner API Go client with your API key.
//
// Usage:
//
//	go run ./scripts/test-partner-api -api-key sc_live_...
//	# or
//	$env:SHEDCLOUD_API_KEY="sc_live_..."; go run ./scripts/test-partner-api
//
// Optional flags:
//
//	-base-url   override host (default: production https://go.shedcloud.com)
//	-sandbox    use sandbox host
//	-resource   lot-stock | leads | quotes | orders | work-orders | locations | customers | products (default: lot-stock)
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
	apiKey := flag.String("api-key", os.Getenv("SHEDCLOUD_API_KEY"), "Partner API key (sc_live_…); or set SHEDCLOUD_API_KEY")
	baseURL := flag.String("base-url", os.Getenv("SHEDCLOUD_BASE_URL"), "API host override (optional)")
	sandbox := flag.Bool("sandbox", false, "Use sandbox host (api.shedcloudtest.com)")
	resource := flag.String("resource", "lot-stock", "Resource to call: lot-stock | leads | quotes | orders | work-orders | locations | customers | products")
	limit := flag.Int("limit", 5, "Page size")
	flag.Parse()

	key := strings.TrimSpace(*apiKey)
	if key == "" {
		fmt.Fprintln(os.Stderr, "ERROR: pass -api-key or set SHEDCLOUD_API_KEY")
		os.Exit(1)
	}
	if !strings.HasPrefix(key, "sc_live_") {
		fmt.Fprintln(os.Stderr, "WARNING: API keys normally start with sc_live_ — continuing anyway")
	}

	opts := partnerapi.Options{
		Auth: partnerapi.Auth{APIKey: key},
	}
	if u := strings.TrimSpace(*baseURL); u != "" {
		opts.BaseURL = u
	} else if *sandbox {
		opts.Environment = partnerapi.EnvironmentSandbox
	}

	client, err := partnerapi.New(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client init failed: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("Partner API Go client smoke test")
	fmt.Printf("  Host     : %s\n", client.BaseURL)
	fmt.Printf("  Resource : %s\n", *resource)
	fmt.Printf("  Key      : %s...\n", truncate(key, 16))
	fmt.Println(strings.Repeat("-", 60))

	start := time.Now()
	switch strings.ToLower(strings.TrimSpace(*resource)) {
	case "lot-stock", "lotstock":
		res, err := client.LotStock.List(ctx, partnerapi.LotStockListParams{
			PaginationParams: partnerapi.PaginationParams{Limit: *limit},
			PurchaseType:     "Lot Stock",
		})
		must(err)
		printList("lot-stock", res.Page, res.Limit, res.Total, len(res.Data), start)
		if len(res.Data) > 0 {
			item := res.Data[0]
			fmt.Printf("  first: %s | %s | %s | $%.2f | %s\n",
				item.SerialNumber, item.Title, item.PurchaseType, item.Price, item.LocationName)
		}
	case "leads":
		res, err := client.Leads.List(ctx, partnerapi.SalesListParams{
			PaginationParams: partnerapi.PaginationParams{Limit: *limit},
		})
		must(err)
		printList("leads", res.Page, res.Limit, res.Total, len(res.Data), start)
		if len(res.Data) > 0 {
			item := res.Data[0]
			fmt.Printf("  first: #%d | %s | %s | %s\n",
				item.OrderNumber, item.Status, item.Customer.Name, item.Customer.Email)
		}
	case "quotes":
		res, err := client.Quotes.List(ctx, partnerapi.QuoteListParams{
			SalesListParams: partnerapi.SalesListParams{
				PaginationParams: partnerapi.PaginationParams{Limit: *limit},
			},
		})
		must(err)
		printList("quotes", res.Page, res.Limit, res.Total, len(res.Data), start)
		if len(res.Data) > 0 {
			item := res.Data[0]
			fmt.Printf("  first: #%d | %s | total=$%.2f | serial=%s | converted=%v\n",
				item.OrderNumber, item.Status, item.Pricing.Total, item.SerialNumber, item.Converted)
		}
	case "orders":
		res, err := client.Orders.List(ctx, partnerapi.OrderListParams{
			SalesListParams: partnerapi.SalesListParams{
				PaginationParams: partnerapi.PaginationParams{Limit: *limit},
			},
		})
		must(err)
		printList("orders", res.Page, res.Limit, res.Total, len(res.Data), start)
		if len(res.Data) > 0 {
			item := res.Data[0]
			fmt.Printf("  first: #%d | %s | total=$%.2f | serial=%s\n",
				item.OrderNumber, item.Status, item.Pricing.Total, item.SerialNumber)
		}
	case "work-orders", "workorders", "wo":
		res, err := client.WorkOrders.List(ctx, partnerapi.WorkOrderListParams{
			PaginationParams: partnerapi.PaginationParams{Limit: *limit},
		})
		must(err)
		printList("work-orders", res.Page, res.Limit, res.Total, len(res.Data), start)
		if len(res.Data) > 0 {
			item := res.Data[0]
			fmt.Printf("  first: WO#%d | %s | %s | order=#%d\n",
				item.WorkOrderNumber, item.SerialNumber, item.Status, item.OrderNumber)
		}
	case "locations":
		res, err := client.Locations.List(ctx, partnerapi.LocationListParams{
			PaginationParams: partnerapi.PaginationParams{Limit: *limit},
		})
		must(err)
		printList("locations", res.Page, res.Limit, res.Total, len(res.Data), start)
		if len(res.Data) > 0 {
			item := res.Data[0]
			fmt.Printf("  first: %s | code=%s | %s, %s | active=%v salesLot=%v plant=%v\n",
				item.Name, item.Code, item.City, item.State, item.Active, item.SalesLot, item.Plant)
		}
	case "customers":
		res, err := client.Customers.List(ctx, partnerapi.CustomerListParams{
			PaginationParams: partnerapi.PaginationParams{Limit: *limit},
		})
		must(err)
		printList("customers", res.Page, res.Limit, res.Total, len(res.Data), start)
		if len(res.Data) > 0 {
			item := res.Data[0]
			fmt.Printf("  first: %s | %s | %s | active=%v\n",
				item.Name, item.Email, item.Phone, item.Active)
		}
	case "products":
		res, err := client.Products.List(ctx, partnerapi.ProductListParams{
			PaginationParams: partnerapi.PaginationParams{Limit: *limit},
		})
		must(err)
		printList("products", res.Page, res.Limit, res.Total, len(res.Data), start)
		if len(res.Data) > 0 {
			item := res.Data[0]
			fmt.Printf("  first: %s | sku=%s | $%.2f | %gx%g | active=%v\n",
				item.Name, item.SKU, item.Price, item.Width, item.Length, item.Active)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown -resource %q (use lot-stock|leads|quotes|orders|work-orders|locations|customers|products)\n", *resource)
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("Done.")
}

func printList(name string, page, limit int, total int64, returned int, start time.Time) {
	fmt.Printf(">> %s list OK  (%d ms)\n", name, time.Since(start).Milliseconds())
	fmt.Printf("   total=%d  page=%d  limit=%d  returned=%d\n", total, page, limit, returned)
}

func must(err error) {
	if err == nil {
		return
	}
	if apiErr, ok := err.(*partnerapi.Error); ok {
		fmt.Fprintf(os.Stderr, "FAIL HTTP %d: %s\n", apiErr.Status, apiErr.Message)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
	os.Exit(1)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
