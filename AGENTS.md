# AGENTS.md — shedcloud-gomod

Orientation for coding agents working in this repository.

## What this is

`github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi` is the official Go SDK for the ShedCloud **ERP** Partner API (`/partner/v1/*`), from **v1.0.0**. Partners authenticate with a company-scoped API key (`sc_live_…`) or OAuth2 client credentials (`POST /oauth/token`, form fields → `sc_at_…`).

Versions: **v1.x → ERP** (snake_case JSON, money as decimal strings). **v0.x → Legacy** Partner API (camelCase). Do not mix them; the `v0.12.0` tag is the last Legacy release.

This repo is the **consumer library**. The API lives in `shedcloud-erp-api` (read-only for us): `docs/partner/openapi.yaml`, `PARTNER_API.md`, `PARTNER_AUTH.md`, `internal/partnerauth/scopes.go`, `internal/handler/partnerapi/*`. The openapi file has **no response schemas**; response shapes come from the handler code and `internal/domain` structs. The npm/PyPI/PHP/Ruby SDKs still target Legacy and are aligned separately.

## Layout

```
shedcloud-gomod/
├── partnerapi/
│   ├── client.go          # New / Client / Options (BaseURL required)
│   ├── hosts.go           # package doc, ErrBaseURLRequired, ResolveBaseURL
│   ├── auth.go            # API key + OAuth token cache (form-field credentials)
│   ├── http.go            # request/do helper, query encoding, RequestOption
│   ├── common.go          # ListParams, Page[T], Pagination, listPage/listKey/call
│   ├── errors.go          # Error / AuthError, ERP {error, code, message} envelope
│   ├── scopes.go          # generated from the ERP AllScopes (91)
│   ├── webhooks.go        # signature verify + headers + WebhookEvent
│   ├── company.go         # Me, Company
│   ├── stock_units.go     # /stock-units (was lot-stock)
│   ├── leads.go
│   ├── commercial.go      # Quotes + Orders (shared generic service)
│   ├── work_orders.go
│   ├── location_budgets.go
│   ├── reports.go         # sales ledger (POST), revenue forecast
│   ├── site_events.go     # Track uses snake_case, List read side is camelCase
│   ├── customers.go, payments.go, locations.go
│   ├── catalog.go         # products (read), pricing resolve, agreements
│   ├── admin.go           # users/roles, domains, documents, events, configurator sessions
│   └── *_test.go          # httptest, one fake ERP per call (helpers_test.go)
├── examples/list-stock-units/
├── scripts/test-partner-api/   # smoke test against a real host
├── go.mod
└── README.md, CHANGELOG.md
```

## Rules of the road

1. **Never invent endpoints or field names.** Mirror the ERP `openapi.yaml` + handler/domain structs. The ERP decodes request bodies with `DisallowUnknownFields` (400 on unknown keys), so request types must contain only accepted fields, all `omitempty`.
2. **No built-in hosts.** The ERP docs publish no production/sandbox host; `Options.BaseURL` is required. If hosts become official, add them in a minor release.
3. **Scopes** in `scopes.go` are generated from `internal/partnerauth/scopes.go` (`AllScopes`, catalog order). Regenerate, don't hand-edit; `TestScopeCatalog` pins the count.
4. **Auth stays outside resource methods.** Resources only call `httpClient.request`/`do` (through `call`, `listPage`, `listKey`). Token exchange/caching belongs in `auth.go`.
5. **Money is a decimal string, never a float.** Exception: `RevenueForecastItem.Amount` (the ERP sends a JSON number there).
6. **PATCH request types use pointer fields**; only non-nil fields are sent. Create calls take `...RequestOption` (`WithIdempotencyKey`).
7. **Lists:** `ListParams{Page, PageSize, Q}` embedded in resource params → `*Page[T]`; the resource key differs per route (`stock_units`, `work_orders`, `location_budgets`, …) and is passed to `listPage`.
8. **`SiteEvents.Track` sends snake_case**, `SiteEvents.List` reads camelCase (`eventId`, `sessionId`) and takes camelCase query params; `RevenueForecast` query params are camelCase too. Do not "normalize" them.
9. Where `PARTNER_API.md` examples disagree with handler code, follow the code.
10. Keep the surface small. No Cognito admin routes (`/v1/companies/{id}/partner-api/*`).
11. Run `gofmt -l .` (the committed content is LF; the working tree may be CRLF on Windows: check `git show HEAD:<file> | gofmt -l`), then `go test ./... && go vet ./...`.

## Adding a resource

1. Find the route, scope and request/response structs in `shedcloud-erp-api` (`internal/handler/partnerapi/*.go`, `internal/domain/*.go`).
2. Add a `<name>.go` with types + a `<Name>Service`; wire it in `client.go`.
3. Use `listPage` / `listKey` / `call`; keep nested blobs as `json.RawMessage` unless the shape is stable.
4. Add an `httptest` test (`fakeERP` in `helpers_test.go`): method, path, query/body keys, decoded fields.
5. Document it in `README.md` (Resources table) and `CHANGELOG.md`.

## Pending ERP resources (not wrapped)

See the list in `README.md` ("ERP resources the SDK does not cover yet").
