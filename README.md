# shedcloud-gomod

Official Go client for the **ShedCloud ERP Partner API** (`/partner/v1/*`).

> **Which version do I use?**
>
> | Module version | Talks to | JSON |
> |---|---|---|
> | **v1.x** (this README) | the new ShedCloud **ERP** Partner API | snake_case, money as decimal strings |
> | **v0.x** (up to 0.12.0) | the **Legacy** Partner API | camelCase, money as numbers |
>
> The two APIs are different services. If you stay on `v0.x` you keep calling Legacy and nothing changes for you; pin it with `go get github.com/Corland-Partners-LLC/shedcloud-gomod@v0.12.0`. Move to `v1.x` only when your company runs on the ERP. See [Migrating from 0.x](#migrating-from-0x).

Go 1.22+. API key (`sc_live_…`) or OAuth2 client credentials.

## Install

```bash
go get github.com/Corland-Partners-LLC/shedcloud-gomod@v1.0.0
```

The import path stays `github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi` (no `/v2` suffix: the module never had a `v1+` release before).

## Host (you must set it)

The ERP Partner API docs publish no fixed production or sandbox host, so **this SDK has none built in**. `Options.BaseURL` is required; `New` returns `ErrBaseURLRequired` without it. Use the host your ShedCloud contact gave you (it is also shown in the portal, Settings → Developer API).

## Quick start

```go
client, err := partnerapi.New(partnerapi.Options{
	BaseURL: os.Getenv("SHEDCLOUD_BASE_URL"), // required, e.g. https://api.example.com
	Auth:    partnerapi.Auth{APIKey: os.Getenv("SHEDCLOUD_API_KEY")},
})
if err != nil {
	log.Fatal(err)
}

me, _ := client.Me(ctx) // credential id, kind, granted scopes
fmt.Println(me.Scopes)

page, err := client.StockUnits.List(ctx, partnerapi.StockUnitListParams{
	ListParams: partnerapi.ListParams{PageSize: 20},
	Statuses:   []string{"available"},
})
for _, u := range page.Items {
	fmt.Println(u.SerialNumber, u.ProductName, u.AskingPrice) // money is a string: "6500.00"
}
fmt.Println(page.Pagination.Total)
```

### OAuth2 client credentials

```go
client, err := partnerapi.New(partnerapi.Options{
	BaseURL: baseURL,
	Auth: partnerapi.Auth{
		ClientID:     os.Getenv("SHEDCLOUD_CLIENT_ID"),     // sc_client_…
		ClientSecret: os.Getenv("SHEDCLOUD_CLIENT_SECRET"), // sc_secret_…
	},
})
```

The token (`sc_at_…`) is fetched from `POST /oauth/token` (form fields), cached, and refreshed shortly before it expires. `client.RefreshAccessToken(ctx)` forces a refresh.

### Lists, money, ids

- Every list takes `ListParams{Page, PageSize, Q}` (1-based, `PageSize` max 100) and returns `*Page[T]{Items, Pagination{Page, PageSize, Total}}`.
- Money and quantities are **decimal strings** (`"6500.00"`). Use a decimal library (for example `shopspring/decimal`) to compute; never `float64`.
- Update (`PATCH`) request types use pointer fields: only non-nil fields are sent.
- The ERP rejects unknown JSON fields with 400, so request types contain exactly the fields the ERP accepts.
- Create calls accept `partnerapi.WithIdempotencyKey("key")` (24h replay; same key with a different body is a 409).

### Place to Order

```go
q, _ := client.Quotes.Create(ctx, partnerapi.QuoteCreateRequest{StoreLocationID: loc, DealType: "cash"})
fmt.Println(q.ViewURL) // customer page; q.ShareToken is returned only here, once

_, _ = client.Quotes.ReplaceLines(ctx, q.ID, []partnerapi.LineInput{{StockUnitID: unitID}})
_, _ = client.Quotes.SetCustomer(ctx, q.ID, partnerapi.CustomerLinkRequest{Name: "Sam River", Email: "sam@example.com"})
sold, _ := client.Quotes.Sell(ctx, q.ID) // replaces Convert
order, _ := client.Orders.Get(ctx, sold.ConvertedOrderID)
```

## Resources

| Client field | ERP routes | Notes |
|---|---|---|
| `Me(ctx)` | `GET /me` | credential, kind, scopes |
| `Company` | `/company` | |
| `StockUnits` | `/stock-units` | replaces `LotStock`; List, Get, Create, RegisterReturn, Update, Release, PDF |
| `Leads` | `/leads` | List, Get, Create, Update, UpdateStatus, Assign, Convert, StatusHistory |
| `Quotes` | `/quotes` | List, Get, Create, Update, UpdateStatus, SetCustomer, ReplaceLines, UpdateLineBuild, cash/RTO terms (Preview/Compute/Set), ComputeMileage, ResolveTax, PreviewReprice, Reprice, **Sell**, Events |
| `Orders` | `/orders` | same shared calls as Quotes (no Sell), plus UpdateStatus, Contract |
| `WorkOrders` | `/work-orders` | List, Get, Create, Update, UpdateStatus, StatusHistory |
| `LocationBudgets` | `/location-budgets` | List, Get, Create, Update (PUT) |
| `SalesLedger` | `POST /sales-ledger` | `Report` |
| `RevenueForecast` | `/revenue-forecast` | `Get` |
| `SiteEvents` | `/site-events` | `Track` (snake_case body), `List` |
| `Customers` | `/customers` | List, Get, Create, Update, Payments, CreatePayment, Ledger |
| `Payments` | `/payments` | List, Get |
| `Locations` | `/locations` | List, Get, Create, Update, Domains |
| `Products` | `/products`, `/pricing/resolve` | List, Get, ResolvePrice (read only) |
| `Agreements` | `/agreements` | List, Active, Get, ListStateLegal, GetStateLegal |
| `Users` | `/users`, `/roles` | List, Get, Invite, Update, Roles |
| `Domains` | `/domains` | List |
| `Documents` | `/documents` | List (by entity), Download |
| `Events` | `/events`, `/webhook-deliveries` | List, Deliveries, Redeliver |
| `ConfiguratorSessions` | `/configurator-sessions` | Create |

### ERP resources the SDK does not cover yet

Available in the ERP Partner API (see its `openapi.yaml`) but **not wrapped here** — pending, not part of 1.0.0. Call them with your own `net/http` for now.

`products` writes and `pricing-matrix`; `delivery-offers` (transport); `sold-upgrades`; `rto-agreements`; order `signature-process`; `notes`; `audit`; `promo-codes`; `inventory-settings`, `inventory/adjustments`, `inventory/transfers`; `service-tickets`; `stock-sale-requests` (Request for Sale); `catalog/templates`, `catalog/kits`; `communications` and customer conversations; `documents` uploads; `sales-returns`; `stock-reconciliations`; `cycle-counts`; `bank-accounts`, `bank-reconciliations`, `financial-statements`; `commissions`; `mrp`; `quality/ncrs`; `opportunities`; `warehouse-tasks`; `freight-bills`; `fixed-assets`; `analytics/exports`; `edi-maps`; `subscriptions`; `fx`; `consolidations`; `engineering-changes`.

## Migrating from 0.x

v1.0.0 is a **major** release: it is not source-compatible with 0.x because the API behind it changed. Plan on reading every call site.

### Global changes

| Area | 0.x (Legacy) | 1.0 (ERP) |
|---|---|---|
| Host | built in (`Environment`, `HostProduction`, `HostSandbox`) | **`BaseURL` required**; `Environment*`, `Host*`, `DefaultBaseURL` removed; `ResolveBaseURL(baseURL)` has one argument |
| JSON | camelCase | **snake_case** (`SiteEvents.Track` already was) |
| Money | `float64` | **decimal `string`** |
| Lists | `PaginatedResponse[T]{Data, Page, Limit, Total}`, `PaginationParams{Limit, Offset…}`, `Sort/Order` | `Page[T]{Items, Pagination}`, `ListParams{Page, PageSize, Q}`; no generic sort |
| Errors | legacy body | `{error, code, message}`: `Error.Message` = plain `message` (else `error`), `Error.Code` = `code`; new `IsConflict`, `IsUnprocessable` |
| OAuth | HTTP Basic client auth | `client_id`/`client_secret` as **form fields** |
| Concurrency | `WithIfMatch(version)` | **removed** (the ERP does not document `If-Match`) |
| Webhooks | `X-ShedCloud-Event-Id` / `-Event-Type`, `EventItem` body | `WebhookEventHeader` (`X-ShedCloud-Event`), `WebhookDeliveryHeader` (`X-ShedCloud-Delivery`), `WebhookEvent` body; signature scheme unchanged (`VerifyWebhookSignature`) |
| Scopes | 32 constants | the ERP catalog (91). Every 0.x scope still exists; new ones include `ScopeLotStockWrite`, `ScopePricingRead` |
| Types | `LotStockItem`, `QuoteItem`, `OrderItem`, `WorkOrderItem`, `LeadItem`, `CustomerItem`, … | `StockUnit`, `Quote`, `Order`, `WorkOrder`, `Lead`, `Customer`, … (all fields re-modelled from the ERP) |

### Per resource

| 0.x | 1.0 |
|---|---|
| `LotStock.List` | `StockUnits.List` (`/stock-units`). Filters: `Origin`, `Statuses`, `SourceWorkOrderID`, `IncludeExpired`, `IncludeReceived`. Sold units past their window and `received` units are hidden by default. Typed photo galleries (`FrontImages`…) are gone from the unit; use `Documents` |
| `Quotes.Convert` | `Quotes.Sell` (`POST /quotes/{id}/sell`, scope `orders.write`); the order id is `Quote.ConvertedOrderID` |
| `Quotes/Orders.AddLineItem`, `DeleteLineItem`, `LineItems` | `ReplaceLines` (`PUT …/lines`, one item in v1) and `UpdateLineBuild`; lines come inside `Get` |
| `Quotes/Orders/WorkOrders/Leads.StatusHistory` (paginated) | quotes/orders: `Events` (not paginated); work orders and leads: `StatusHistory` returns a slice |
| `UpdateStatus` (`POST`) | `PATCH …/status` (quotes also take a reason: `UpdateStatus(ctx, id, status, reason)`) |
| `Quotes.Create` body (nested customer, delivery address, line items) | `QuoteCreateRequest{StoreLocationID, DealType, …}` then `ReplaceLines` and `SetCustomer`; the reply `QuoteCreated` carries `ViewURL` and the one-time `ShareToken`. `ViewURL` is **not** on get/list any more (`HasShareToken` only) |
| `OrderPatchRequest{SaleDate, SoldPricing}`, `OrderItem.Cancelled` | `sale_date` is read-only; ledger state is `Order.LedgerBooked`, `LedgerAmount`, `LedgerCancelled`. Sold pricing detail lives in `SalesLedger.Report` |
| `Orders.Payments`, `Orders.CreatePayment` | `Customers.CreatePayment` / `Customers.Payments` and `Payments.List` with `ReferenceType: "order"`, `ReferenceID` |
| `Orders.CreatePaymentLink` | **removed** (no ERP equivalent) |
| `Orders.Contract` (summary) | `Orders.Contract` returns signature / contract-PDF status; PDF bytes are staff-only |
| `SalesLedger.List` (`GET`, `SaleDateFrom`…, `Cancelled`, status) | `SalesLedger.Report` (`POST`) with `SalesLedgerRequest{LedgerFrom, LedgerTo, SaleFrom, SaleTo, StoreID, MovementMode, Q, IncludeSums}`; rows are ledger movements (`Delta*`, `New*`) |
| `RevenueForecast.Get` (KPI blocks, week buckets, detail rows) | `RevenueForecast.Get` returns the ERP forecast strip `[{label, amount}]`; query params are camelCase (`locationIds`, `weekStart`…); `amount` is a JSON number |
| `LocationBudgets` (`LocationBudgetUpsertRequest`) | `LocationBudgetRequest` with `Months []LocationBudgetMonth{Month, WeightPct, Amount, WeeklyAmount, LeadTarget}` |
| `SiteEvents.Track` | unchanged envelope; `Delivery` / `Payment` blocks removed (use `Payload`); new `EventID`, `UserID`, `LocationID`. `List` has no cursor; `Each` removed |
| `Leads.Create` (location routing) | `LeadCreateRequest` is flat; new `Assign`, `Convert` |
| `Customers.Merge` | **removed** (no ERP equivalent) |
| `StockTemplates`, BOM `dimension-variables`, `Products.CreateSize` | **removed**: no confirmed ERP equivalent (`catalog/templates` browses ShedCloud catalog templates, which is a different thing) |
| `Products.Create/Update` | **removed for now**: the ERP has product writes (`partner-api.products.write`) but they are not wrapped yet (pending list) |
| `Client.Me` (company + rate-limit blocks) | `MeResponse{TenantID, CredentialID, Kind, Scopes}`; use `Company.Get` for the company |
| `Domains.List`, `Domains.ForLocation` | `Domains.List` and `Locations.Domains` return `[]Storefront` (not paginated) |
| `Agreements.Active()` (one item) | `Agreements.Active(params)` returns a `Page[Agreement]` |
| `Users.Create` | `Users.Invite` (`POST /users` invites by email) |
| `Events.List`, `Events.Redeliver`, `Events.Deliveries` | same routes, ERP shapes (`Event`, `WebhookDelivery`, `EventRedeliverResponse`); `List` returns a slice (`After` + `Limit`), `Each` removed |

## Webhooks

Verify deliveries with the subscription secret (Settings → Developer API). Always verify the **raw** request body. Deliveries are at-least-once: dedupe on the `X-ShedCloud-Delivery` header and the event `id`.

```go
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	if err := partnerapi.VerifyWebhookSignature(
		os.Getenv("SHEDCLOUD_WEBHOOK_SECRET"),
		r.Header.Get(partnerapi.WebhookSignatureHeader), body, 0); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ev partnerapi.WebhookEvent // id, entity_type, entity_id, action, occurred_at, payload
	_ = json.Unmarshal(body, &ev)
	w.WriteHeader(http.StatusOK) // ack fast; process async
}
```

## Errors

```go
_, err := client.Orders.Get(ctx, id)
var apiErr *partnerapi.Error
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.Status, apiErr.Code, apiErr.Message)
	// IsUnauthorized / IsForbidden / IsNotFound / IsConflict / IsUnprocessable / IsRateLimited
	// apiErr.Body keeps the full response (for example violations[] on 422 configuration_rule_violation)
}
```

OAuth token failures return `*partnerapi.AuthError`.

## Scopes

`partnerapi.AllScopes` is the full ERP catalog and each scope has a constant (`ScopeLotStockRead`, `ScopeLotStockWrite`, `ScopeQuotesWrite`, `ScopeOrdersWrite` — needed by `Quotes.Sell` —, `ScopeSalesLedgerRead`, …). `client.Me(ctx)` lists what your credential holds; `partnerapi.HasScope(me.Scopes, scope)` checks one.

## Development

```bash
go test ./...
go vet ./...
```

## Versioning

v1.x follows the ERP Partner API, which is additive within `/partner/v1`. Breaking SDK changes need a new major version. See [CHANGELOG.md](CHANGELOG.md).

## Related SDKs

The npm, PyPI, PHP and Ruby SDKs still target the **Legacy** API at the time of v1.0.0 and will be aligned to the ERP separately.

## Docs

- ERP Partner API: `docs/partner/openapi.yaml`, `PARTNER_API.md`, `PARTNER_AUTH.md` in `shedcloud-erp-api` (also served by the API under `/partner/reference`)
- This module on pkg.go.dev: [`partnerapi`](https://pkg.go.dev/github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi)
