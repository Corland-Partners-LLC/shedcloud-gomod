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

	ScopeUsersRead                 = "partner-api.users.read"
	ScopeUsersWrite                = "partner-api.users.write"
	ScopeContractsRead             = "partner-api.contracts.read"
	ScopePaymentsRead              = "partner-api.payments.read"
	ScopePaymentsWrite             = "partner-api.payments.write"
	ScopeDocumentsRead             = "partner-api.documents.read"
	ScopeEventsRead                = "partner-api.events.read"
	ScopeConfiguratorSessionsWrite = "partner-api.configurator-sessions.write"
	ScopeDomainsRead               = "partner-api.domains.read"
	ScopeAgreementsRead            = "partner-api.agreements.read"
	ScopeProductsWrite             = "partner-api.products.write"
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
	ScopeUsersRead,
	ScopeUsersWrite,
	ScopeContractsRead,
	ScopePaymentsRead,
	ScopePaymentsWrite,
	ScopeDocumentsRead,
	ScopeEventsRead,
	ScopeConfiguratorSessionsWrite,
	ScopeDomainsRead,
	ScopeAgreementsRead,
	ScopeProductsWrite,
}
