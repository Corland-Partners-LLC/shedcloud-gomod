# Changelog

## 1.0.0 — unreleased

**Major release: the SDK now talks to the ShedCloud ERP Partner API.** It is
not source-compatible with 0.x. Anyone who stays on `v0.12.0` (or any 0.x) keeps
calling the Legacy Partner API unchanged; pin with
`go get github.com/Corland-Partners-LLC/shedcloud-gomod@v0.12.0`. The import
path is unchanged (no `/v2`: the module had no `v1+` release before). See the
migration guide in `README.md`.

Source of truth: `shedcloud-erp-api` `docs/partner/openapi.yaml` (167 paths),
`PARTNER_API.md`, `PARTNER_AUTH.md`, `internal/partnerauth/scopes.go` and the
handlers under `internal/handler/partnerapi`. The OpenAPI file has no response
schemas, so response types were modelled from the handler/domain structs.

### Breaking changes

- **Host.** No built-in hosts (the ERP docs publish none).
  `Options.BaseURL` is required (`ErrBaseURLRequired`). Removed `Environment`,
  `EnvironmentProduction`, `EnvironmentSandbox`, `HostProduction`,
  `HostSandbox`, `DefaultBaseURL`; `ResolveBaseURL` takes one argument.
- **JSON is snake_case and money is a decimal string** (0.x: camelCase,
  `float64`). Every request/response type was re-modelled.
- **Lists.** `PaginatedResponse[T]` / `PaginationParams` / `SortOrder` replaced
  by `Page[T]{Items, Pagination}` and `ListParams{Page, PageSize, Q}`.
- **Errors** follow `{error, code, message}`: `Error.Message` is the
  plain-language `message` (else `error`), `Error.Code` is `code`. Added
  `IsConflict`, `IsUnprocessable`.
- **OAuth** sends `client_id` / `client_secret` as form fields (0.x used HTTP
  Basic).
- **`WithIfMatch` removed** (the ERP does not document `If-Match`).
- **Scopes** are now the ERP catalog (91, all 0.x scopes kept). Added
  `HasScope`.
- **Webhooks.** Header constants are now `WebhookEventHeader`
  (`X-ShedCloud-Event`) and `WebhookDeliveryHeader` (`X-ShedCloud-Delivery`);
  `WebhookEventIDHeader` / `WebhookEventTypeHeader` removed; body type is
  `WebhookEvent`. Signature scheme and `VerifyWebhookSignature` unchanged.
- **`LotStock` → `StockUnits`** (`/stock-units`): List, Get, Create,
  RegisterReturn, Update, Release, PDF. Typed photo galleries are not on the
  unit.
- **`Quotes.Convert` → `Quotes.Sell`** (`POST /quotes/{id}/sell`).
- **Lines.** `AddLineItem` / `DeleteLineItem` / `LineItems` →
  `ReplaceLines` (`PUT …/lines`, one item in v1) + `UpdateLineBuild`.
- **Status history.** Quotes/orders: `Events`; work orders/leads:
  `StatusHistory` returns a slice. `UpdateStatus` is `PATCH` (was `POST`).
- **Quote create** returns `QuoteCreated` with `ViewURL` and a one-time
  `ShareToken`; `ViewURL` is no longer on get/list (only `HasShareToken`).
- **Orders.** `SaleDate` / `SoldPricing` patching removed (`sale_date`
  read-only); `Cancelled` → `LedgerCancelled` (+ `LedgerBooked`,
  `LedgerAmount`). Order payments → `Customers.CreatePayment` /
  `Payments.List` with `ReferenceType: "order"`. `CreatePaymentLink` removed
  (no ERP equivalent).
- **`SalesLedger.List` → `SalesLedger.Report`** (`POST`, new filters
  `ledger_from/to`, `sale_from/to`, `store_id`, `movement_mode`,
  `include_sums`; rows are ledger movements).
- **`RevenueForecast.Get`** returns the ERP forecast strip
  `[{label, amount}]` (camelCase query params, `amount` is a JSON number); the
  0.x KPI blocks, week buckets and detail rows do not exist in the ERP.
- **`SiteEvents`.** `Delivery` / `Payment` blocks removed (put detail in
  `Payload`); `Each` and cursors removed (the ERP returns one batch, `hasMore`
  is always false). `Track` stays snake_case.
- **`Me`** returns `{tenant_id, credential_id, kind, scopes}` (no company or
  rate-limit block).
- **Removed, no ERP equivalent:** `Customers.Merge`, `Orders.CreatePaymentLink`,
  `StockTemplates`, BOM `dimension-variables`, `Products.CreateSize`.
- Renames: `Users.Create` → `Users.Invite`; `Domains.ForLocation` →
  `Locations.Domains`; `Agreements.Active()` now returns a page; `Events.Each`
  removed.

### Added

- `Leads.Assign`, `Leads.Convert`; `Quotes`/`Orders` terms (cash and RTO
  preview/compute/set), mileage, tax, reprice preview/apply; `Orders.Contract`;
  `Customers.Ledger`; `Payments`; `Products.ResolvePrice`
  (`/pricing/resolve`); `Locations.Domains`; and every resource listed under
  Full ERP coverage.

### Full ERP coverage

All 214 operations (166 paths) of the ERP `/partner/v1` `openapi.yaml` have a method (see the
Resources table in `README.md`), including products write and pricing matrix,
catalog, promo codes, delivery offers, RTO agreements, order signature process,
sold upgrades, inventory settings / adjustments / transfers, Request for Sale,
sales returns, stock reconciliations, cycle counts, warehouse tasks, service
tickets, notes, audit, communications, opportunities, document uploads,
banking, financial statements, commissions, freight bills, fixed assets, FX,
consolidations, subscriptions, MRP, quality, analytics exports, EDI maps and
engineering changes. A route test pins method and path for each one.

### Notes

- Where `PARTNER_API.md` examples disagree with the handler code, the SDK
  follows the code (customer ledger `net_balance`, sales ledger `rows`).
- Some list responses have no `pagination` block in the ERP (service tickets,
  promo codes, sale requests, catalog, banks…); those return slices or a zero
  `Pagination`. Scope and filter details are in each method's doc comment.
- Cash/RTO previews (`TermsPreview`) and `ResolvePrice` return generic maps;
  nested blobs (addresses, snapshots, work-order components, agreement
  schedules) stay `json.RawMessage`.

## 0.12.0 — 2026-09-09

- **Quote `viewUrl`.** `QuoteItem` adds `ViewURL` — the full customer Share
  page (`{PORTAL_BASE_URL}/quotes/{state}/view/{id}`), the same link the
  portal Share Customer View dialog copies. Present on create, get, and list
  when `PORTAL_BASE_URL` is configured.
- Parity with the `shedcloud-api-go` Partner API changelog entry (2026-09-09).

## 0.11.0 — 2026-08-25

- **Revenue / forecast.** `client.RevenueForecast.Get` for
  `GET /partner/v1/revenue-forecast` (scope `ScopeRevenueForecastRead`).
  New types `RevenueForecastParams`, `RevenueForecastResponse`, KPI blocks,
  week buckets, and detail rows. Qualification is delivery / promised dates;
  money is the Sales Ledger order amount so a related order matches
  `SalesLedger.List` `newOrderAmount`. Period totals do not equal the Sales
  Ledger report for the same calendar window.
- Parity with the `shedcloud-api-go` Partner API changelog entry (2026-08-25).

## 0.10.0 — 2026-07-30

- **Orders: sticky `Cancelled`.** `OrderItem` adds `Cancelled` (same meaning as
  sales-ledger / portal Sales Ledger sticky cancelled). Present on list and get.
  Store reporting should exclude `Cancelled` from net sold totals.
- Parity with `@shedcloud/partner-api` 0.10.0 and the `shedcloud-api-go`
  Partner API changelog entry (2026-07-30).

## 0.9.0 — 2026-07-30

- **Sales ledger.** `client.SalesLedger.List` for `GET /partner/v1/sales-ledger`
  (scope `ScopeSalesLedgerRead`). New types `SalesLedgerItem` and
  `SalesLedgerListParams`. Rows include `SaleDate` and sticky `Cancelled`; list
  accepts `SaleDateFrom` / `SaleDateTo`, `Cancelled`, status, and related filters.
- Parity with `@shedcloud/partner-api` 0.9.0 and the `shedcloud-api-go`
  Partner API sales-ledger surface.

## 0.8.1 — 2026-07-28

- **Lot-stock typed photo galleries.** `LotStockItem` adds `FrontImages`,
  `BackImages`, `LeftImages`, `RightImages`, `ExteriorImages`, and
  `InteriorImages` (public CDN URL slices; empty/omitted when none).
  Existing `Images` and `HeroImages` are unchanged.
- Parity with the `shedcloud-api-go` Partner API changelog entry (2026-07-28).

## 0.8.0 — 2026-07-27

- **Orders: sale date + sold pricing.** `OrderItem` adds `SaleDate` and
  `SoldPricing` (`PartnerSoldPricing`). List with `SaleDateFrom` /
  `SaleDateTo`. `OrderPatchRequest` accepts `SaleDate` and `SoldPricing`.
- Parity with `@shedcloud/partner-api` 0.8.0 and the `shedcloud-api-go`
  Partner API changelog entry (2026-07-27).

## 0.7.0 — 2026-07-26

- **Location sales budgets** — `client.LocationBudgets.List/Get/Create/Update` for
  `GET/POST/PUT /partner/v1/location-budgets` (scopes
  `ScopeLocationBudgetsRead` / `ScopeLocationBudgetsWrite`). New types:
  `LocationBudgetItem`, `LocationBudgetUpsertRequest`, `LocationBudgetListParams`.
- Parity with `@shedcloud/partner-api` 0.7.0 and the `shedcloud-api-go` Partner
  API changelog entry (2026-07-25).

## 0.3.0 — 2026-07-22

- **Site events** — identity snapshots on `SiteEventInput`: `Customer`
  (`SiteEventCustomer`: first/last name, email, phone, zip), `Delivery`
  (`SiteEventDelivery`: address + lat/lng), and `Payment`
  (`SiteEventPayment`: payment_type `RENT_TO_OWN`/`CASH`, terms, totals).
  Attach to `customer.profile` / `identity.resolved` / `delivery.address` /
  `payment.select` events so visitor profiles stop showing as anonymous.
- **Site events** — `SiteEventsTrackRequest.ClientIP` / `ClientUserAgent`:
  backend proxies forward the end shopper's IP and User-Agent so events
  aren't stamped with the server's own context (geo accuracy + per-IP rate
  fairness upstream).
- Parity with `@shedcloud/partner-api` 0.6.0, `shedcloud-partner-api`
  (PyPI) 0.6.0, and the `shedcloud-api-go` partner ingest change.

## 0.2.0 — 2026-07-19

- **Site events** — `client.SiteEvents.Track()` and `List()` / `Each()` for
  `POST/GET /partner/v1/site-events` (visitor behavioral tracking from partner
  marketing sites). New scopes `ScopeSiteEventsRead` / `ScopeSiteEventsWrite`
  (`partner-api.site-events.read` / `.write`). Ingest body is **snake_case**
  on the wire (`session_id`, `events[].event_type`, …) — pass the struct/map
  through as-is. Reads use camelCase query params (`sessionId`, `types`,
  `from`, `to`). Batch-only ingest (max 25 events per call); server-side use
  only.
- Parity with `@shedcloud/partner-api` site-events resource and
  `shedcloud-api-go` Partner API changelog entry (2026-07-19).

## 0.1.0

- Initial `partnerapi` Go client
- Built-in hosts: production `https://go.shedcloud.com` (default), sandbox `https://api.shedcloudtest.com`
- Auth: API key bearer + OAuth2 client-credentials with token cache
- Resources: lot-stock, leads, quotes, orders, work-orders (list/get/update/status)
- Typed errors (`Error`, `AuthError`) and scope constants
- Parity with `@shedcloud/partner-api` (including `QuoteItem.serialNumber` / `workOrderId`)
