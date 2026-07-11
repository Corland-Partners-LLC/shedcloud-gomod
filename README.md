# shedcloud-gomod

Official Go client for the ShedCloud Partner API (`/partner/v1/*`).

Use this module from Go 1.22+ to call company-scoped Partner API endpoints with an API key (`sc_live_…`) or OAuth2 client credentials.

## Install

```bash
go get github.com/Corland-Partners-LLC/shedcloud-gomod@latest
```

## Hosts

| Environment | Host |
|-------------|------|
| `production` (default) | `https://go.shedcloud.com` |
| `sandbox` | `https://api.shedcloudtest.com` |

Pass `Environment` for sandbox, or `BaseURL` for a custom/local override.

## Quick start

### API key (production)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi"
)

func main() {
	client, err := partnerapi.New(partnerapi.Options{
		Auth: partnerapi.Auth{APIKey: os.Getenv("SHEDCLOUD_API_KEY")},
	})
	if err != nil {
		log.Fatal(err)
	}

	stock, err := client.LotStock.List(context.Background(), partnerapi.LotStockListParams{
		PaginationParams: partnerapi.PaginationParams{Limit: 50},
		PurchaseType:     "Lot Stock",
		Sort:             "price",
		Order:            partnerapi.SortAsc,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stock.Total, stock.Data[0].Title)
}
```

### Sandbox

```go
client, err := partnerapi.New(partnerapi.Options{
	Environment: partnerapi.EnvironmentSandbox,
	Auth:        partnerapi.Auth{APIKey: os.Getenv("SHEDCLOUD_API_KEY")},
})
```

### OAuth2 client credentials

```go
client, err := partnerapi.New(partnerapi.Options{
	Auth: partnerapi.Auth{
		ClientID:     os.Getenv("SHEDCLOUD_CLIENT_ID"),
		ClientSecret: os.Getenv("SHEDCLOUD_CLIENT_SECRET"),
	},
})
// Access tokens are fetched from POST /oauth/token and cached until near expiry.
orders, err := client.Orders.List(ctx, partnerapi.OrderListParams{
	SalesListParams: partnerapi.SalesListParams{Status: "Unprocessed", PaginationParams: partnerapi.PaginationParams{Limit: 25}},
})
```

Create credentials in the ShedCloud portal under **Settings → Developer API**.

## Resources

| Client field | Endpoints |
|--------------|-----------|
| `client.LotStock` | `GET /partner/v1/lot-stock` |
| `client.Leads` | `GET/POST/PATCH /partner/v1/leads`, `POST .../status` — `Create(...)` makes a lead with location lead-routing |
| `client.Quotes` | `GET/POST/PATCH /partner/v1/quotes`, `POST .../status` — `Create(...)` makes an in-stock quote from a serial number, `Convert(id, ...)` places it as a sales order |
| `client.Orders` | `GET/PATCH /partner/v1/orders`, `POST .../status` |
| `client.WorkOrders` | `GET/PATCH /partner/v1/work-orders`, `POST .../status` |
| `client.Locations` | `GET/POST/PATCH /partner/v1/locations` |
| `client.Customers` | `GET/POST/PATCH /partner/v1/customers` |
| `client.Products` | `GET /partner/v1/products` (read-only catalog) |

### Examples

```go
lead, err := client.Leads.Get(ctx, "665f0a1b2c3d4e5f60718293")

_, err = client.Leads.Update(ctx, lead.ID, partnerapi.LeadPatchRequest{
	SalespersonName: "Alex Rep",
	SalesLocation:   "66c00443c2d8aa83c5757dcf",
})

_, err = client.Orders.UpdateStatus(ctx, orderID, partnerapi.StatusUpdateRequest{
	Status:            "On hold",
	ActionDescription: "Waiting on customer financing",
})

// Create a quote from an in-stock unit: the serial number's work order is
// linked to the new quote, the sales location is assigned, and lead routing
// auto-assigns a salesperson when the location has routing configured.
quote, err := client.Quotes.Create(ctx, partnerapi.QuoteCreateRequest{
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

// Create a lead: when no salesperson is given, the location's lead-routing
// strategy (round-robin, availability, skill-based) auto-assigns one.
newLead, err := client.Leads.Create(ctx, partnerapi.LeadCreateRequest{
	LocationID: "66c00443c2d8aa83c5757dcf",
	Customer:   partnerapi.LeadCreateCustomer{Name: "Jane Doe", Email: "jane@example.com"},
})

// Convert a quote to a sales order (requires partner-api.orders.write).
// The new order starts in "Unsubmitted" — submit with UpdateStatus.
order, err := client.Quotes.Convert(ctx, quote.ID, partnerapi.QuoteConvertRequest{})

// All create/convert methods accept an idempotency key: a retried request
// with the same key replays the stored response instead of duplicating.
_, err = client.Quotes.Create(ctx, req, partnerapi.WithIdempotencyKey(uuid.NewString()))
```

## Errors

Failed responses return `*partnerapi.Error` (or `*partnerapi.AuthError` for OAuth token failures):

```go
_, err := client.Orders.Get(ctx, id)
if apiErr, ok := err.(*partnerapi.Error); ok {
	fmt.Println(apiErr.Status, apiErr.Message, apiErr.Body)
	// apiErr.IsUnauthorized() / IsForbidden() / IsNotFound() / IsRateLimited()
}
```

## Scopes

```go
partnerapi.ScopeLotStockRead    // partner-api.lot-stock.read
partnerapi.ScopeOrdersWrite     // partner-api.orders.write
```

## Development

```bash
go test ./...
go vet ./...
```

## Docs

- Partner API reference: https://go.shedcloud.com/partner/reference
- Backend source of truth: `shedcloud-api-go/docs/PARTNER_API.md`
- TypeScript twin: [`@shedcloud/partner-api`](https://github.com/Corland-Partners-LLC/shedcloud-npm)
