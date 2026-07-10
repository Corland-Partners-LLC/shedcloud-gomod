# AGENTS.md — shedcloud-gomod

Orientation for coding agents working in this repository.

## What this is

`github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi` is the official Go SDK for ShedCloud's public **Partner API** (`/partner/v1/*`). Partners authenticate with a company-scoped API key (`sc_live_…`) or OAuth2 client credentials (`POST /oauth/token`).

This repo is the **consumer library**. The API itself lives in `shedcloud-api-go` (`internal/handler/partnerapi`, `docs/PARTNER_API.md`, `docs/partner-api/swagger.yaml`). The TypeScript twin is `shedcloud-npm` (`@shedcloud/partner-api`).

## Layout

```
shedcloud-gomod/
├── partnerapi/           # importable package
│   ├── client.go         # New / Client
│   ├── hosts.go          # production / sandbox base URLs
│   ├── auth.go           # API key + OAuth token cache
│   ├── http.go           # request helper + query encoding
│   ├── errors.go
│   ├── scopes.go
│   ├── types.go          # request/response types
│   ├── lot_stock.go
│   ├── leads.go
│   ├── quotes.go
│   ├── orders.go
│   ├── work_orders.go
│   └── client_test.go
├── examples/
├── go.mod
└── README.md
```

## Rules of the road

1. **Never invent endpoints or field names.** Mirror `shedcloud-api-go/docs/PARTNER_API.md` and keep parity with `shedcloud-npm`.
2. **Hosts** are fixed in `hosts.go`: production `https://go.shedcloud.com`, sandbox `https://api.shedcloudtest.com`.
3. **Scopes** live in `scopes.go` and must stay in sync with `shedcloud-api-go/internal/partnerauth/scopes.go`.
4. **Auth stays outside resource methods.** Resources only call `httpClient.request`. Token exchange/caching belongs in `auth.go`.
5. **Keep the surface small.** No portal admin endpoints (`/v1/settings/api-keys`, etc.).
6. After changing code, run `go test ./... && go vet ./...`.

## Adding a resource

1. Confirm the route + scopes in `shedcloud-api-go/internal/router/partner_api_routes.go`.
2. Add types to `types.go`.
3. Add `<name>.go` with `List` / `Get` / `Update` / `UpdateStatus` as applicable.
4. Wire it on `Client` in `client.go`.
5. Add a unit test with `httptest`.
6. Mirror the same change in `shedcloud-npm` when practical.
