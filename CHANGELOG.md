# Changelog

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
