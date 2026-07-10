# Changelog

## 0.1.0

- Initial `partnerapi` Go client
- Built-in hosts: production `https://go.shedcloud.com` (default), sandbox `https://api.shedcloudtest.com`
- Auth: API key bearer + OAuth2 client-credentials with token cache
- Resources: lot-stock, leads, quotes, orders, work-orders (list/get/update/status)
- Typed errors (`Error`, `AuthError`) and scope constants
- Parity with `@shedcloud/partner-api` (including `QuoteItem.serialNumber` / `workOrderId`)
