# Changelog

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
