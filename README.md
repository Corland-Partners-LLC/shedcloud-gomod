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
| `client.Leads` | `GET/PATCH /partner/v1/leads`, `POST .../status` |
| `client.Quotes` | `GET/PATCH /partner/v1/quotes`, `POST .../status` |
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
