package partnerapi

// Assignable Partner API scopes (credential permissions). Keep in sync with
// shedcloud-api-go/internal/partnerauth/scopes.go.
const (
	ScopeLotStockRead    = "partner-api.lot-stock.read"
	ScopeLeadsRead       = "partner-api.leads.read"
	ScopeLeadsWrite      = "partner-api.leads.write"
	ScopeQuotesRead      = "partner-api.quotes.read"
	ScopeQuotesWrite     = "partner-api.quotes.write"
	ScopeOrdersRead      = "partner-api.orders.read"
	ScopeOrdersWrite     = "partner-api.orders.write"
	ScopeWorkOrdersRead  = "partner-api.work-orders.read"
	ScopeWorkOrdersWrite = "partner-api.work-orders.write"
	ScopeLocationsRead   = "partner-api.locations.read"
	ScopeLocationsWrite  = "partner-api.locations.write"
	ScopeCustomersRead   = "partner-api.customers.read"
	ScopeCustomersWrite  = "partner-api.customers.write"
	ScopeProductsRead    = "partner-api.products.read"
)

// AllScopes is the full catalog in display order.
var AllScopes = []string{
	ScopeLotStockRead,
	ScopeLeadsRead,
	ScopeLeadsWrite,
	ScopeQuotesRead,
	ScopeQuotesWrite,
	ScopeOrdersRead,
	ScopeOrdersWrite,
	ScopeWorkOrdersRead,
	ScopeWorkOrdersWrite,
	ScopeLocationsRead,
	ScopeLocationsWrite,
	ScopeCustomersRead,
	ScopeCustomersWrite,
	ScopeProductsRead,
}
