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

Each resource maps to a section of the [hosted reference](https://go.shedcloud.com/partner/reference).

| Client field | Endpoints | Reference |
|--------------|-----------|-----------|
| `client.LotStock` | `GET /partner/v1/lot-stock` | [#lot-stock](https://go.shedcloud.com/partner/reference#lot-stock) |
| `client.StockTemplates` | `GET /partner/v1/stock-templates` (buildable catalog designs) | [#stock-templates](https://go.shedcloud.com/partner/reference#stock-templates) |
| `client.Leads` | `GET/POST/PATCH /partner/v1/leads`, `POST .../status`, `GET .../status-history` — `Create(...)` makes a lead with location lead-routing | [#leads](https://go.shedcloud.com/partner/reference#leads) |
| `client.Quotes` | `GET/POST/PATCH /partner/v1/quotes`, `POST .../status`, `GET .../status-history`, `GET/POST/DELETE .../line-items` — `Create(...)` makes an in-stock quote from a serial number, `Convert(id, ...)` places it as a sales order | [#quotes](https://go.shedcloud.com/partner/reference#quotes) |
| `client.Orders` | `GET/POST/PATCH /partner/v1/orders`, `POST .../status`, `GET .../status-history`, `GET/POST/DELETE .../line-items`, `GET .../contract`, `GET/POST .../payments`, `POST .../payment-links` — `Create(...)` makes a full order (customer, base product + size, upgrades, configurator); `CreatePayment(...)` records manual payments, `CreatePaymentLink(...)` returns a Stripe Checkout URL | [#orders](https://go.shedcloud.com/partner/reference#orders) |
| `client.WorkOrders` | `GET/POST/PATCH /partner/v1/work-orders`, `POST .../status`, `GET .../status-history` — `Create(...)` makes a work order with friendly type enums + optional `SizeID` | [#work-orders](https://go.shedcloud.com/partner/reference#work-orders) |
| `client.Locations` | `GET/POST/PATCH /partner/v1/locations` | [#locations](https://go.shedcloud.com/partner/reference#locations) |
| `client.Customers` | `GET/POST/PATCH /partner/v1/customers`, `POST .../{id}/merge` | [#customers](https://go.shedcloud.com/partner/reference#customers) |
| `client.Products` | `GET/POST/PATCH /partner/v1/products`, `CreateSize(id, ...)` for `POST .../{id}/sizes` (finished catalog products with gallery `Images`) | [#products](https://go.shedcloud.com/partner/reference#products) |
| `client.Domains` | `GET /partner/v1/domains`, `ForLocation(id)` for `GET /partner/v1/locations/{id}/domains` (white-label storefront domains, `defaultForStore` filter) | [#domains](https://go.shedcloud.com/partner/reference#domains) |
| `client.Users` | `GET/POST/PATCH /partner/v1/users`, `Roles(ctx)` for `GET /partner/v1/roles` — `Create(...)` makes a company user (role, locations, invite email); `Update(...)` patches profile/role/locations/`Active` | [#users](https://go.shedcloud.com/partner/reference#users) |
| `client.Payments` | `GET /partner/v1/payments[/{id}]` (read-only) | [#payments](https://go.shedcloud.com/partner/reference#payments) |
| `client.Documents` | `GET /partner/v1/documents`, `GET .../{id}/download` (short-lived presigned URL) | [#documents](https://go.shedcloud.com/partner/reference#documents) |
| `client.Events` | `GET /partner/v1/events` cursor feed, `Each(...)` iterator, `Redeliver(id)`, `Deliveries(...)` webhook delivery log | [#events](https://go.shedcloud.com/partner/reference#events) |
| `client.SiteEvents` | `POST /partner/v1/site-events` batch ingest (visitor behavioral tracking, snake_case body), `List(...)` / `Each(...)` read-back | [#site-events](https://go.shedcloud.com/partner/reference#site-events) |
| `client.ConfiguratorSessions` | `POST /partner/v1/configurator-sessions` (single-use 3D configurator launch URLs) | [#configurator-sessions](https://go.shedcloud.com/partner/reference#configurator-sessions) |

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

// Or create a full order directly: customer + location + base product
// (model + size), upgrades, and an optional configurator payload.
basePrice := 8995.0
fullOrder, err := client.Orders.Create(ctx, partnerapi.OrderCreateRequest{
	Customer:   &partnerapi.QuoteCreateCustomer{Name: "Jane Doe", Email: "jane@example.com"},
	LocationID: "66c00443c2d8aa83c5757dcf",
	ProductID:  "6659f3ab8e5a2c001f9b1c11",
	SizeID:     "6659f3ab8e5a2c001f9b1c22",
	BasePrice:  &basePrice,
	Upgrades:   []partnerapi.OrderCreateLineItem{{ProductID: "6659f3ab8e5a2c001f9b1c33", Quantity: 2}},
	Configuration: &partnerapi.OrderCreateConfiguration{
		SidingColor: "Barn Red", RoofMaterial: "Metal Roof", RoofColor: "Charcoal",
	},
}, partnerapi.WithIdempotencyKey(uuid.NewString()))

// Manage upgrade lines afterward (idempotent on LineKey); the base product
// line is protected. Update pricing fields yourself after changing lines.
line, err := client.Orders.AddLineItem(ctx, fullOrder.ID, partnerapi.LineItemCreateRequest{
	ProductID: "6659f3ab8e5a2c001f9b1c44",
	LineKey:   "crm-line-42",
})
err = client.Orders.DeleteLineItem(ctx, fullOrder.ID, line.LineID)

// Create a manufacturing work order (portal building-wizard parity): status
// starts in "Customer Care", the number is allocated automatically, and an
// optional SizeID attaches the product.
workOrder, err := client.WorkOrders.Create(ctx, partnerapi.WorkOrderCreateRequest{
	LocationID:    "dallas-lot",
	PurchaseType:  "new-build",
	WorkOrderType: "made-to-order",
	DeliveryType:  "delivery-from-factory",
	SerialNumber:  "SC-2026-00999",
	SizeID:        "6659f3ab8e5a2c001f9b1c22",
}, partnerapi.WithIdempotencyKey(uuid.NewString()))

// Record a manual payment (cash | check | financed | manual) — runs the
// portal's own pipeline: payment record + order-balance recalc + audit.
// Card/ACH are rejected; use payment links for those.
payment, err := client.Orders.CreatePayment(ctx, fullOrder.ID, partnerapi.PaymentCreateRequest{
	Method: "check", Amount: 500, CheckNumber: "1042",
}, partnerapi.WithIdempotencyKey(uuid.NewString()))

// Or send the customer a Stripe Checkout link (requires the company's
// Stripe integration): the webhook records the payment when they pay.
sendEmail := true
link, err := client.Orders.CreatePaymentLink(ctx, fullOrder.ID, partnerapi.PaymentLinkCreateRequest{
	Amount: 250, SendEmail: &sendEmail,
})
fmt.Println(link.URL, link.ExpiresAt)

// Create a company user: discover roles first, then create with the role,
// locations, and invite email. The user's login is created at their first
// sign-in via the invite (no Cognito account is made here).
roles, err := client.Users.Roles(ctx)
newUser, err := client.Users.Create(ctx, partnerapi.UserCreateRequest{
	FirstName:   "Alex",
	LastName:    "Rep",
	Email:       "alex@example.com",
	RoleID:      roles.Data[0].ID,
	LocationIDs: []string{"66c00443c2d8aa83c5757dcf"},
}, partnerapi.WithIdempotencyKey(uuid.NewString()))
// Deactivate later (the company owner is protected).
inactive := false
_, err = client.Users.Update(ctx, newUser.ID, partnerapi.UserPatchRequest{Active: &inactive})

// All create/convert methods accept an idempotency key: a retried request
// with the same key replays the stored response instead of duplicating.
_, err = client.Quotes.Create(ctx, req, partnerapi.WithIdempotencyKey(uuid.NewString()))

// Stamp your own correlation ids on records; filter lists by them later.
crmID := "deal-42"
_, err = client.Orders.Update(ctx, order.ID, partnerapi.OrderPatchRequest{
	ExternalRefs: partnerapi.ExternalReferencesPatch{"crmDealId": &crmID},
})
matches, err := client.Orders.List(ctx, partnerapi.OrderListParams{
	SalesListParams: partnerapi.SalesListParams{ExternalRef: "crmDealId:deal-42"},
})

// Optimistic concurrency: send the version you read as If-Match — the server
// answers 409 Conflict if someone else wrote in between.
_, err = client.Orders.Update(ctx, order.ID,
	partnerapi.OrderPatchRequest{CustomerPhone: "555-0100"},
	partnerapi.WithIfMatch(order.Version))

// Consume the change feed losslessly with a stored cursor.
err = client.Events.Each(ctx, partnerapi.EventListParams{Cursor: lastSeen},
	func(ev partnerapi.EventItem) error {
		lastSeen = ev.ID
		return handle(ev)
	})

// Catalog products: create a model, then give it real pricing via sizes.
product, err := client.Products.Create(ctx, partnerapi.ProductCreateRequest{
	Name: "10x16 Lofted Barn", SKU: "LB-1016",
}, partnerapi.WithIdempotencyKey(uuid.NewString()))
_, err = client.Products.CreateSize(ctx, product.ID, partnerapi.ProductSizeCreateRequest{
	Width: 10, Length: 16, Price: 8200,
})

// White-label storefront domains: find each location's primary storefront.
defaultOnly := true
domains, err := client.Domains.List(ctx, partnerapi.DomainListParams{
	DefaultForStore: &defaultOnly,
})
storeDomains, err := client.Domains.ForLocation(ctx, "66c00443c2d8aa83c5757dcf",
	partnerapi.LocationDomainListParams{})
```

## Webhooks

Verify webhook deliveries with the subscription secret (shown once when the webhook is created in Settings → Developer API). Always verify against the **raw** request body:

```go
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	err := partnerapi.VerifyWebhookSignature(
		os.Getenv("SHEDCLOUD_WEBHOOK_SECRET"),
		r.Header.Get(partnerapi.WebhookSignatureHeader),
		body, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var event partnerapi.EventItem
	_ = json.Unmarshal(body, &event)
	w.WriteHeader(http.StatusOK) // ack fast; process async and dedupe by event.ID
}
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

## Versioning & changelog

The Partner API is **additive-only within `/partner/v1`** — new endpoints and
fields may appear, but existing response fields are never renamed, retyped, or
removed. This module is tagged with semver in lockstep with API additions: new
API capabilities arrive as **minor** releases, fixes as **patches**. Pinning
an older version is always safe; it simply won't surface newer fields.

- API changes: [hosted changelog](https://go.shedcloud.com/partner/reference#changelog)

## Docs

- Partner API reference: https://go.shedcloud.com/partner/reference
- API changelog: https://go.shedcloud.com/partner/reference#changelog
- Backend source of truth: `shedcloud-api-go/docs/PARTNER_API.md`
- TypeScript/JavaScript: [`@shedcloud/partner-api`](https://github.com/Corland-Partners-LLC/shedcloud-npm)
- Python: [`shedcloud-partner-api`](https://github.com/Corland-Partners-LLC/shedcloud-pypi)
- PHP: [`shedcloud/partner-api`](https://github.com/Corland-Partners-LLC/shedcloud-php)
- Ruby: [`shedcloud-partner_api`](https://github.com/Corland-Partners-LLC/shedcloud-gem)
