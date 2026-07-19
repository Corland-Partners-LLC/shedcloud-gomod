package partnerapi

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// Client is the official Go client for the ShedCloud Partner API.
//
//	client, err := partnerapi.New(partnerapi.Options{
//		Auth: partnerapi.Auth{APIKey: os.Getenv("SHEDCLOUD_API_KEY")},
//	})
//	stock, err := client.LotStock.List(ctx, partnerapi.LotStockListParams{Limit: 50})
type Client struct {
	LotStock             *LotStockService
	StockTemplates       *StockTemplatesService
	Leads                *LeadsService
	Quotes               *QuotesService
	Orders               *OrdersService
	WorkOrders           *WorkOrdersService
	Locations            *LocationsService
	Customers            *CustomersService
	Company              *CompanyService
	Domains              *DomainsService
	Agreements           *AgreementsService
	Products             *ProductsService
	Users                *UsersService
	Payments             *PaymentsService
	Documents            *DocumentsService
	Events               *EventsService
	SiteEvents           *SiteEventsService
	ConfiguratorSessions *ConfiguratorSessionsService

	// BaseURL is the resolved API host used for all requests.
	BaseURL string

	auth *authProvider
	http *httpClient
}

// Options configures a Partner API Client.
type Options struct {
	// Environment selects a built-in host. Defaults to production.
	// Ignored when BaseURL is set.
	Environment Environment
	// BaseURL overrides the host (trailing slashes stripped). Prefer Environment
	// unless you need a custom/local host.
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
	if err := opts.Auth.validate(); err != nil {
		return nil, err
	}

	baseURL := ResolveBaseURL(opts.BaseURL, opts.Environment)
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
	c.LotStock = &LotStockService{c: c}
	c.StockTemplates = &StockTemplatesService{c: c}
	c.Leads = &LeadsService{c: c}
	c.Quotes = &QuotesService{c: c}
	c.Orders = &OrdersService{c: c}
	c.WorkOrders = &WorkOrdersService{c: c}
	c.Locations = &LocationsService{c: c}
	c.Customers = &CustomersService{c: c}
	c.Company = &CompanyService{c: c}
	c.Domains = &DomainsService{c: c}
	c.Agreements = &AgreementsService{c: c}
	c.Products = &ProductsService{c: c}
	c.Users = &UsersService{c: c}
	c.Payments = &PaymentsService{c: c}
	c.Documents = &DocumentsService{c: c}
	c.Events = &EventsService{c: c}
	c.SiteEvents = &SiteEventsService{c: c}
	c.ConfiguratorSessions = &ConfiguratorSessionsService{c: c}
	return c, nil
}

// RefreshAccessToken force-refreshes the OAuth access token (no-op for API key auth).
func (c *Client) RefreshAccessToken(ctx context.Context) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return c.auth.refresh(ctx)
}
