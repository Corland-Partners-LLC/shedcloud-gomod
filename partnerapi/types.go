package partnerapi

// Shared Partner API types. Shapes mirror shedcloud-api-go/docs/PARTNER_API.md.

// SortOrder is asc or desc.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// PaginatedResponse is the list envelope for every Partner API list endpoint.
type PaginatedResponse[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// PartnerCustomer is the nested customer object on sales entities.
type PartnerCustomer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// PartnerSalesperson is the nested salesperson object on sales entities.
type PartnerSalesperson struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// PartnerLocation is the nested location object.
type PartnerLocation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

// PartnerPricing is the nested pricing object on quotes and orders.
type PartnerPricing struct {
	Subtotal       float64 `json:"subtotal"`
	Total          float64 `json:"total"`
	MonthlyPayment float64 `json:"monthlyPayment,omitempty"`
	// PaymentType is "rto" or "cash".
	PaymentType string `json:"paymentType,omitempty"`
}

// LotStockItem is one on-lot inventory row.
type LotStockItem struct {
	ID           string   `json:"id"`
	WorkOrderID  string   `json:"workOrderId,omitempty"`
	SerialNumber string   `json:"serialNumber,omitempty"`
	Title        string   `json:"title,omitempty"`
	PurchaseType string   `json:"purchaseType,omitempty"`
	BasePrice    float64  `json:"basePrice,omitempty"`
	Price        float64  `json:"price,omitempty"`
	LocationID   string   `json:"locationId,omitempty"`
	LocationName string   `json:"locationName,omitempty"`
	LocationSlug string   `json:"locationSlug,omitempty"`
	Images       []string `json:"images,omitempty"`
	Sold         bool     `json:"sold"`
}

// LeadItem is one lead.
type LeadItem struct {
	ID              string             `json:"id"`
	OrderNumber     int                `json:"orderNumber,omitempty"`
	Status          string             `json:"status,omitempty"`
	StatusChangedAt string             `json:"statusChangedAt,omitempty"`
	Customer        PartnerCustomer    `json:"customer"`
	Salesperson     PartnerSalesperson `json:"salesperson"`
	Location        PartnerLocation    `json:"location"`
	CreatedAt       string             `json:"createdAt,omitempty"`
	UpdatedAt       string             `json:"updatedAt,omitempty"`
}

// QuoteItem is one quote.
type QuoteItem struct {
	ID                   string             `json:"id"`
	OrderNumber          int                `json:"orderNumber,omitempty"`
	Status               string             `json:"status,omitempty"`
	StatusChangedAt      string             `json:"statusChangedAt,omitempty"`
	Customer             PartnerCustomer    `json:"customer"`
	Salesperson          PartnerSalesperson `json:"salesperson"`
	Location             PartnerLocation    `json:"location"`
	Pricing              PartnerPricing     `json:"pricing"`
	SerialNumber         string             `json:"serialNumber,omitempty"`
	WorkOrderID          string             `json:"workOrderId,omitempty"`
	Converted            bool               `json:"converted"`
	ConvertedOrderID     string             `json:"convertedOrderId,omitempty"`
	ConvertedOrderNumber int                `json:"convertedOrderNumber,omitempty"`
	CreatedAt            string             `json:"createdAt,omitempty"`
	UpdatedAt            string             `json:"updatedAt,omitempty"`
}

// OrderItem is one sales order.
type OrderItem struct {
	ID                string             `json:"id"`
	OrderNumber       int                `json:"orderNumber,omitempty"`
	Status            string             `json:"status,omitempty"`
	StatusChangedAt   string             `json:"statusChangedAt,omitempty"`
	Customer          PartnerCustomer    `json:"customer"`
	Salesperson       PartnerSalesperson `json:"salesperson"`
	Location          PartnerLocation    `json:"location"`
	Pricing           PartnerPricing     `json:"pricing"`
	SerialNumber      string             `json:"serialNumber,omitempty"`
	WorkOrderID       string             `json:"workOrderId,omitempty"`
	SourceQuoteID     string             `json:"sourceQuoteId,omitempty"`
	SourceQuoteNumber int                `json:"sourceQuoteNumber,omitempty"`
	CreatedAt         string             `json:"createdAt,omitempty"`
	UpdatedAt         string             `json:"updatedAt,omitempty"`
}

// WorkOrderItem is one work order.
type WorkOrderItem struct {
	ID              string          `json:"id"`
	WorkOrderNumber int             `json:"workOrderNumber,omitempty"`
	SerialNumber    string          `json:"serialNumber,omitempty"`
	Title           string          `json:"title,omitempty"`
	Status          string          `json:"status,omitempty"`
	StatusChangedAt string          `json:"statusChangedAt,omitempty"`
	OrderID         string          `json:"orderId,omitempty"`
	OrderNumber     int             `json:"orderNumber,omitempty"`
	Location        PartnerLocation `json:"location"`
	BasePrice       float64         `json:"basePrice,omitempty"`
	PromisedDate    string          `json:"promisedDate,omitempty"`
	CreatedAt       string          `json:"createdAt,omitempty"`
	UpdatedAt       string          `json:"updatedAt,omitempty"`
}

// PaginationParams are shared by every list endpoint.
type PaginationParams struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// LotStockListParams are query params for GET /partner/v1/lot-stock.
type LotStockListParams struct {
	PaginationParams
	PurchaseType string    `json:"purchaseType,omitempty"` // ALL | Lot Stock | Rental Return | Immediate Sale
	Location     string    `json:"location,omitempty"`
	Search       string    `json:"search,omitempty"`
	Sort         string    `json:"sort,omitempty"` // serialNumber | title | price | createdAt
	Order        SortOrder `json:"order,omitempty"`
}

// SalesListParams are shared list filters for leads, quotes, and orders.
type SalesListParams struct {
	PaginationParams
	Search        string    `json:"search,omitempty"`
	Sort          string    `json:"sort,omitempty"` // orderNumber | customerName | status | total | createdAt | updatedAt
	Order         SortOrder `json:"order,omitempty"`
	Status        string    `json:"status,omitempty"`
	Location      string    `json:"location,omitempty"`
	CustomerEmail string    `json:"customerEmail,omitempty"`
	CustomerPhone string    `json:"customerPhone,omitempty"`
	OrderNumber   string    `json:"orderNumber,omitempty"`
	Salesperson   string    `json:"salesperson,omitempty"`
	CreatedFrom   string    `json:"createdFrom,omitempty"`
	CreatedTo     string    `json:"createdTo,omitempty"`
	UpdatedFrom   string    `json:"updatedFrom,omitempty"`
	UpdatedTo     string    `json:"updatedTo,omitempty"`
}

// QuoteListParams extends SalesListParams with converted.
type QuoteListParams struct {
	SalesListParams
	// Converted filters converted quotes when non-nil.
	Converted *bool `json:"converted,omitempty"`
}

// OrderListParams extends SalesListParams with payment/serial filters.
type OrderListParams struct {
	SalesListParams
	PaymentType  string `json:"paymentType,omitempty"` // rto | cash
	SerialNumber string `json:"serialNumber,omitempty"`
}

// WorkOrderListParams are query params for GET /partner/v1/work-orders.
type WorkOrderListParams struct {
	PaginationParams
	Search        string    `json:"search,omitempty"`
	Sort          string    `json:"sort,omitempty"` // workOrderNumber | serialNumber | status | createdAt | updatedAt
	Order         SortOrder `json:"order,omitempty"`
	Status        string    `json:"status,omitempty"`
	SerialNumber  string    `json:"serialNumber,omitempty"`
	OrderNumber   string    `json:"orderNumber,omitempty"`
	LinkedOrderID string    `json:"linkedOrderId,omitempty"`
	Location      string    `json:"location,omitempty"`
	CreatedFrom   string    `json:"createdFrom,omitempty"`
	CreatedTo     string    `json:"createdTo,omitempty"`
	UpdatedFrom   string    `json:"updatedFrom,omitempty"`
	UpdatedTo     string    `json:"updatedTo,omitempty"`
}

// LeadPatchRequest is the body for PATCH /partner/v1/leads/{id}.
type LeadPatchRequest struct {
	SalesLocation    string `json:"salesLocation,omitempty"`
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
}

// QuotePatchRequest is the body for PATCH /partner/v1/quotes/{id}.
type QuotePatchRequest struct {
	CustomerName     string `json:"customerName,omitempty"`
	CustomerEmail    string `json:"customerEmail,omitempty"`
	CustomerPhone    string `json:"customerPhone,omitempty"`
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
	SalesLocation    string `json:"salesLocation,omitempty"`
}

// OrderPatchRequest is the body for PATCH /partner/v1/orders/{id}.
type OrderPatchRequest struct {
	CustomerName     string `json:"customerName,omitempty"`
	CustomerEmail    string `json:"customerEmail,omitempty"`
	CustomerPhone    string `json:"customerPhone,omitempty"`
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
	SalesLocation    string `json:"salesLocation,omitempty"`
}

// WorkOrderPatchRequest is the body for PATCH /partner/v1/work-orders/{id}.
type WorkOrderPatchRequest struct {
	Title            string `json:"title,omitempty"`
	Description      string `json:"description,omitempty"`
	BuildingLocation string `json:"buildingLocation,omitempty"`
	PromisedDate     string `json:"promisedDate,omitempty"`
}

// StatusUpdateRequest is the body for POST /partner/v1/{resource}/{id}/status.
type StatusUpdateRequest struct {
	Status            string `json:"status"`
	ActionDescription string `json:"actionDescription,omitempty"`
}

// OAuthTokenResponse is the JSON body from POST /oauth/token.
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}
