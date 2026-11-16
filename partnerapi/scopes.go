package partnerapi

// Assignable Partner API scopes (credential permissions). Mirrors
// shedcloud-erp-api/internal/partnerauth/scopes.go (AllScopes).
const (
	ScopeLotStockRead              = "partner-api.lot-stock.read"
	ScopeLotStockWrite             = "partner-api.lot-stock.write"
	ScopeLeadsRead                 = "partner-api.leads.read"
	ScopeLeadsWrite                = "partner-api.leads.write"
	ScopeQuotesRead                = "partner-api.quotes.read"
	ScopeQuotesWrite               = "partner-api.quotes.write"
	ScopeOrdersRead                = "partner-api.orders.read"
	ScopeOrdersWrite               = "partner-api.orders.write"
	ScopeWorkOrdersRead            = "partner-api.work-orders.read"
	ScopeWorkOrdersWrite           = "partner-api.work-orders.write"
	ScopeLocationsRead             = "partner-api.locations.read"
	ScopeLocationsWrite            = "partner-api.locations.write"
	ScopeLocationBudgetsRead       = "partner-api.location-budgets.read"
	ScopeLocationBudgetsWrite      = "partner-api.location-budgets.write"
	ScopeCustomersRead             = "partner-api.customers.read"
	ScopeCustomersWrite            = "partner-api.customers.write"
	ScopeProductsRead              = "partner-api.products.read"
	ScopeProductsWrite             = "partner-api.products.write"
	ScopePricingRead               = "partner-api.pricing.read"
	ScopeUsersRead                 = "partner-api.users.read"
	ScopeUsersWrite                = "partner-api.users.write"
	ScopePaymentsRead              = "partner-api.payments.read"
	ScopePaymentsWrite             = "partner-api.payments.write"
	ScopeDocumentsRead             = "partner-api.documents.read"
	ScopeDocumentsWrite            = "partner-api.documents.write"
	ScopeEventsRead                = "partner-api.events.read"
	ScopeConfiguratorSessionsWrite = "partner-api.configurator-sessions.write"
	ScopeAgreementsRead            = "partner-api.agreements.read"
	ScopeRTOAgreementsRead         = "partner-api.rto-agreements.read"
	ScopeRTOAgreementsWrite        = "partner-api.rto-agreements.write"
	ScopeCompanyRead               = "partner-api.company.read"
	ScopeDomainsRead               = "partner-api.domains.read"
	ScopeContractsRead             = "partner-api.contracts.read"
	ScopeSalesLedgerRead           = "partner-api.sales-ledger.read"
	ScopeSoldUpgradesRead          = "partner-api.sold-upgrades.read"
	ScopeRevenueForecastRead       = "partner-api.revenue-forecast.read"
	ScopeSiteEventsWrite           = "partner-api.site-events.write"
	ScopeSiteEventsRead            = "partner-api.site-events.read"
	ScopeSalesReturnsRead          = "partner-api.sales-returns.read"
	ScopeSalesReturnsWrite         = "partner-api.sales-returns.write"
	ScopeStockReconciliationRead   = "partner-api.stock-reconciliation.read"
	ScopeStockReconciliationWrite  = "partner-api.stock-reconciliation.write"
	ScopeCycleCountsRead           = "partner-api.cycle-counts.read"
	ScopeCycleCountsWrite          = "partner-api.cycle-counts.write"
	ScopeBankAccountsRead          = "partner-api.bank-accounts.read"
	ScopeBankReconciliationsRead   = "partner-api.bank-reconciliations.read"
	ScopeFinancialStatementsRead   = "partner-api.financial-statements.read"
	ScopeCommissionsRead           = "partner-api.commissions.read"
	ScopeMRPRead                   = "partner-api.mrp.read"
	ScopeMRPWrite                  = "partner-api.mrp.write"
	ScopeQualityRead               = "partner-api.quality.read"
	ScopeQualityWrite              = "partner-api.quality.write"
	ScopeOpportunitiesRead         = "partner-api.opportunities.read"
	ScopeOpportunitiesWrite        = "partner-api.opportunities.write"
	ScopeWarehouseTasksRead        = "partner-api.warehouse-tasks.read"
	ScopeWarehouseTasksWrite       = "partner-api.warehouse-tasks.write"
	ScopeFreightBillsRead          = "partner-api.freight-bills.read"
	ScopeFreightBillsWrite         = "partner-api.freight-bills.write"
	ScopeFixedAssetsRead           = "partner-api.fixed-assets.read"
	ScopeFixedAssetsWrite          = "partner-api.fixed-assets.write"
	ScopeAnalyticsRead             = "partner-api.analytics.read"
	ScopeEDIMapsRead               = "partner-api.edi-maps.read"
	ScopeEDIMapsWrite              = "partner-api.edi-maps.write"
	ScopeSubscriptionsRead         = "partner-api.subscriptions.read"
	ScopeSubscriptionsWrite        = "partner-api.subscriptions.write"
	ScopeFXRead                    = "partner-api.fx.read"
	ScopeFXWrite                   = "partner-api.fx.write"
	ScopeConsolidationsRead        = "partner-api.consolidations.read"
	ScopeConsolidationsWrite       = "partner-api.consolidations.write"
	ScopeEngineeringChangesRead    = "partner-api.engineering-changes.read"
	ScopeEngineeringChangesWrite   = "partner-api.engineering-changes.write"
	ScopeTransportRead             = "partner-api.transport.read"
	ScopeTransportWrite            = "partner-api.transport.write"
	ScopeNotesRead                 = "partner-api.notes.read"
	ScopeNotesWrite                = "partner-api.notes.write"
	ScopeAuditRead                 = "partner-api.audit.read"
	ScopePromoCodesRead            = "partner-api.promo-codes.read"
	ScopePromoCodesWrite           = "partner-api.promo-codes.write"
	ScopeInventorySettingsRead     = "partner-api.inventory-settings.read"
	ScopeInventorySettingsWrite    = "partner-api.inventory-settings.write"
	ScopeInventoryAdjustmentsRead  = "partner-api.inventory-adjustments.read"
	ScopeInventoryAdjustmentsWrite = "partner-api.inventory-adjustments.write"
	ScopeInventoryTransfersRead    = "partner-api.inventory-transfers.read"
	ScopeInventoryTransfersWrite   = "partner-api.inventory-transfers.write"
	ScopeServiceTicketsRead        = "partner-api.service-tickets.read"
	ScopeServiceTicketsWrite       = "partner-api.service-tickets.write"
	ScopeStockSaleRequestsRead     = "partner-api.stock-sale-requests.read"
	ScopeStockSaleRequestsWrite    = "partner-api.stock-sale-requests.write"
	ScopeCatalogRead               = "partner-api.catalog.read"
	ScopeCommunicationsRead        = "partner-api.communications.read"
	ScopeCommunicationsWrite       = "partner-api.communications.write"
)

// AllScopes is the full catalog in the ERP display order.
var AllScopes = []string{
	ScopeLotStockRead,
	ScopeLotStockWrite,
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
	ScopeLocationBudgetsRead,
	ScopeLocationBudgetsWrite,
	ScopeCustomersRead,
	ScopeCustomersWrite,
	ScopeProductsRead,
	ScopeProductsWrite,
	ScopePricingRead,
	ScopeUsersRead,
	ScopeUsersWrite,
	ScopePaymentsRead,
	ScopePaymentsWrite,
	ScopeDocumentsRead,
	ScopeDocumentsWrite,
	ScopeEventsRead,
	ScopeConfiguratorSessionsWrite,
	ScopeAgreementsRead,
	ScopeRTOAgreementsRead,
	ScopeRTOAgreementsWrite,
	ScopeCompanyRead,
	ScopeDomainsRead,
	ScopeContractsRead,
	ScopeSalesLedgerRead,
	ScopeSoldUpgradesRead,
	ScopeRevenueForecastRead,
	ScopeSiteEventsWrite,
	ScopeSiteEventsRead,
	ScopeSalesReturnsRead,
	ScopeSalesReturnsWrite,
	ScopeStockReconciliationRead,
	ScopeStockReconciliationWrite,
	ScopeCycleCountsRead,
	ScopeCycleCountsWrite,
	ScopeBankAccountsRead,
	ScopeBankReconciliationsRead,
	ScopeFinancialStatementsRead,
	ScopeCommissionsRead,
	ScopeMRPRead,
	ScopeMRPWrite,
	ScopeQualityRead,
	ScopeQualityWrite,
	ScopeOpportunitiesRead,
	ScopeOpportunitiesWrite,
	ScopeWarehouseTasksRead,
	ScopeWarehouseTasksWrite,
	ScopeFreightBillsRead,
	ScopeFreightBillsWrite,
	ScopeFixedAssetsRead,
	ScopeFixedAssetsWrite,
	ScopeAnalyticsRead,
	ScopeEDIMapsRead,
	ScopeEDIMapsWrite,
	ScopeSubscriptionsRead,
	ScopeSubscriptionsWrite,
	ScopeFXRead,
	ScopeFXWrite,
	ScopeConsolidationsRead,
	ScopeConsolidationsWrite,
	ScopeEngineeringChangesRead,
	ScopeEngineeringChangesWrite,
	ScopeTransportRead,
	ScopeTransportWrite,
	ScopeNotesRead,
	ScopeNotesWrite,
	ScopeAuditRead,
	ScopePromoCodesRead,
	ScopePromoCodesWrite,
	ScopeInventorySettingsRead,
	ScopeInventorySettingsWrite,
	ScopeInventoryAdjustmentsRead,
	ScopeInventoryAdjustmentsWrite,
	ScopeInventoryTransfersRead,
	ScopeInventoryTransfersWrite,
	ScopeServiceTicketsRead,
	ScopeServiceTicketsWrite,
	ScopeStockSaleRequestsRead,
	ScopeStockSaleRequestsWrite,
	ScopeCatalogRead,
	ScopeCommunicationsRead,
	ScopeCommunicationsWrite,
}

// HasScope reports whether granted contains required.
func HasScope(granted []string, required string) bool {
	for _, s := range granted {
		if s == required {
			return true
		}
	}
	return false
}
