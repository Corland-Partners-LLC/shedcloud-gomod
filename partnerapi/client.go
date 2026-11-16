package partnerapi

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// Client is the official Go client for the ShedCloud ERP Partner API.
//
//	client, err := partnerapi.New(partnerapi.Options{
//		BaseURL: os.Getenv("SHEDCLOUD_BASE_URL"),
//		Auth:    partnerapi.Auth{APIKey: os.Getenv("SHEDCLOUD_API_KEY")},
//	})
//	page, err := client.StockUnits.List(ctx, partnerapi.StockUnitListParams{})
type Client struct {
	StockUnits           *StockUnitsService
	Leads                *LeadsService
	Quotes               *QuotesService
	Orders               *OrdersService
	WorkOrders           *WorkOrdersService
	Locations            *LocationsService
	LocationBudgets      *LocationBudgetsService
	Customers            *CustomersService
	Company              *CompanyService
	Domains              *DomainsService
	Agreements           *AgreementsService
	Products             *ProductsService
	Users                *UsersService
	Payments             *PaymentsService
	SalesLedger          *SalesLedgerService
	RevenueForecast      *RevenueForecastService
	Documents            *DocumentsService
	Events               *EventsService
	SiteEvents           *SiteEventsService
	ConfiguratorSessions *ConfiguratorSessionsService

	Catalog              *CatalogService
	PromoCodes           *PromoCodesService
	DeliveryOffers       *DeliveryOffersService
	RTOAgreements        *RTOAgreementsService
	SoldUpgrades         *SoldUpgradesService
	Inventory            *InventoryService
	StockSaleRequests    *StockSaleRequestsService
	SalesReturns         *SalesReturnsService
	StockReconciliations *StockReconciliationsService
	CycleCounts          *CycleCountsService
	WarehouseTasks       *WarehouseTasksService
	ServiceTickets       *ServiceTicketsService
	Notes                *NotesService
	Audit                *AuditService
	Communications       *CommunicationsService
	Opportunities        *OpportunitiesService
	BankAccounts         *BankAccountsService
	BankReconciliations  *BankReconciliationsService
	FinancialStatements  *FinancialStatementsService
	Commissions          *CommissionsService
	FreightBills         *FreightBillsService
	FixedAssets          *FixedAssetsService
	FX                   *FXService
	Consolidations       *ConsolidationsService
	Subscriptions        *SubscriptionsService
	MRP                  *MRPService
	Quality              *QualityService
	Analytics            *AnalyticsService
	EDIMaps              *EDIMapsService
	EngineeringChanges   *EngineeringChangesService

	// BaseURL is the resolved API host used for all requests.
	BaseURL string

	auth *authProvider
	http *httpClient
}

// Options configures a Partner API Client.
type Options struct {
	// BaseURL is the ERP Partner API host, required (trailing slashes are
	// stripped). The ERP docs publish no fixed production or sandbox host, so
	// the SDK has none built in.
	BaseURL string
	// Auth is required. Set APIKey or ClientID+ClientSecret.
	Auth Auth
	// HTTPClient overrides the default http.Client. Useful in tests.
	HTTPClient *http.Client
	// Timeout is the per-request timeout. Default: 30s. Set to a negative
	// duration to disable (not recommended).
	Timeout time.Duration
	// UserAgent overrides the User-Agent header. Default: shedcloud-gomod/partnerapi.
	UserAgent string
	// TokenSkew is how early an OAuth token is refreshed before expiry. Default: 60s.
	TokenSkew time.Duration
}

// New builds a Partner API client.
func New(opts Options) (*Client, error) {
	baseURL := ResolveBaseURL(opts.BaseURL)
	if baseURL == "" {
		return nil, ErrBaseURLRequired
	}
	if err := opts.Auth.validate(); err != nil {
		return nil, err
	}

	httpClientImpl := opts.HTTPClient
	if httpClientImpl == nil {
		httpClientImpl = &http.Client{}
	}
	httpDo := httpClientImpl.Do

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	if timeout < 0 {
		timeout = 0
	}

	userAgent := strings.TrimSpace(opts.UserAgent)
	if userAgent == "" {
		userAgent = "shedcloud-gomod/partnerapi"
	}

	auth := newAuthProvider(baseURL, opts.Auth, httpDo, opts.TokenSkew)
	hc := &httpClient{
		baseURL:   baseURL,
		getToken:  auth.getAccessToken,
		httpDo:    httpDo,
		userAgent: userAgent,
		timeout:   timeout,
	}

	c := &Client{
		BaseURL: baseURL,
		auth:    auth,
		http:    hc,
	}
	c.StockUnits = &StockUnitsService{c: c}
	c.Leads = &LeadsService{c: c}
	c.Quotes = newQuotesService(c)
	c.Orders = newOrdersService(c)
	c.WorkOrders = &WorkOrdersService{c: c}
	c.Locations = &LocationsService{c: c}
	c.LocationBudgets = &LocationBudgetsService{c: c}
	c.Customers = &CustomersService{c: c}
	c.Company = &CompanyService{c: c}
	c.Domains = &DomainsService{c: c}
	c.Agreements = &AgreementsService{c: c}
	c.Products = &ProductsService{c: c}
	c.Users = &UsersService{c: c}
	c.Payments = &PaymentsService{c: c}
	c.SalesLedger = &SalesLedgerService{c: c}
	c.RevenueForecast = &RevenueForecastService{c: c}
	c.Documents = &DocumentsService{c: c}
	c.Events = &EventsService{c: c}
	c.SiteEvents = &SiteEventsService{c: c}
	c.ConfiguratorSessions = &ConfiguratorSessionsService{c: c}
	c.Catalog = &CatalogService{c: c}
	c.PromoCodes = &PromoCodesService{c: c}
	c.DeliveryOffers = &DeliveryOffersService{c: c}
	c.RTOAgreements = &RTOAgreementsService{c: c}
	c.SoldUpgrades = &SoldUpgradesService{c: c}
	c.Inventory = &InventoryService{c: c}
	c.StockSaleRequests = &StockSaleRequestsService{c: c}
	c.SalesReturns = &SalesReturnsService{c: c}
	c.StockReconciliations = &StockReconciliationsService{c: c}
	c.CycleCounts = &CycleCountsService{c: c}
	c.WarehouseTasks = &WarehouseTasksService{c: c}
	c.ServiceTickets = &ServiceTicketsService{c: c}
	c.Notes = &NotesService{c: c}
	c.Audit = &AuditService{c: c}
	c.Communications = &CommunicationsService{c: c}
	c.Opportunities = &OpportunitiesService{c: c}
	c.BankAccounts = &BankAccountsService{c: c}
	c.BankReconciliations = &BankReconciliationsService{c: c}
	c.FinancialStatements = &FinancialStatementsService{c: c}
	c.Commissions = &CommissionsService{c: c}
	c.FreightBills = &FreightBillsService{c: c}
	c.FixedAssets = &FixedAssetsService{c: c}
	c.FX = &FXService{c: c}
	c.Consolidations = &ConsolidationsService{c: c}
	c.Subscriptions = &SubscriptionsService{c: c}
	c.MRP = &MRPService{c: c}
	c.Quality = &QualityService{c: c}
	c.Analytics = &AnalyticsService{c: c}
	c.EDIMaps = &EDIMapsService{c: c}
	c.EngineeringChanges = &EngineeringChangesService{c: c}
	return c, nil
}

// RefreshAccessToken force-refreshes the OAuth access token (no-op for API key auth).
func (c *Client) RefreshAccessToken(ctx context.Context) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return c.auth.refresh(ctx)
}
