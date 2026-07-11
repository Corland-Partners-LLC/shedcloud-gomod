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

// LotStockAttributes is the exterior configuration of a lot-stock unit — the
// same values ShedCloud's work-order detail page shows. Custom colors surface
// as "Custom Color - #HEX".
type LotStockAttributes struct {
	Siding      string `json:"siding,omitempty"`
	SidingColor string `json:"sidingColor,omitempty"`
	TrimColor   string `json:"trimColor,omitempty"`
	// RoofMaterial is "Metal Roof" or "Shingle Roof".
	RoofMaterial string `json:"roofMaterial,omitempty"`
	RoofColor    string `json:"roofColor,omitempty"`
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
	// Attributes is nil when the unit has no configurator.
	Attributes *LotStockAttributes `json:"attributes,omitempty"`
	Sold       bool                `json:"sold"`
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

// LocationItem is one full location record (sales lot, plant, warehouse).
type LocationItem struct {
	ID            string   `json:"id"`
	Name          string   `json:"name,omitempty"`
	Slug          string   `json:"slug,omitempty"`
	Code          string   `json:"code,omitempty"`
	Address       string   `json:"address,omitempty"`
	City          string   `json:"city,omitempty"`
	State         string   `json:"state,omitempty"`
	ZipCode       string   `json:"zipCode,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	ContactPerson string   `json:"contactPerson,omitempty"`
	ContactEmail  string   `json:"contactEmail,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	Active        bool     `json:"active"`
	SalesLot      bool     `json:"salesLot"`
	Plant         bool     `json:"plant"`
	// StoreHours is the weekly store schedule keyed by day ("mon".."sun").
	// Empty when the location has no schedule configured.
	StoreHours map[string]StoreHoursDay `json:"storeHours,omitempty"`
	CreatedAt  string                   `json:"createdAt,omitempty"`
	UpdatedAt  string                   `json:"updatedAt,omitempty"`
}

// StoreHoursDay is one day inside LocationItem.StoreHours. Times are 24-hour
// "HH:MM" strings in the location's local time. Disabled days may still carry
// their last from/to values.
type StoreHoursDay struct {
	Enabled bool   `json:"enabled"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
}

// CustomerItem is one full customer record.
type CustomerItem struct {
	ID            string `json:"id"`
	Name          string `json:"name,omitempty"`
	ContactName   string `json:"contactName,omitempty"`
	ContactPerson string `json:"contactPerson,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	ZipCode       string `json:"zipCode,omitempty"`
	Code          string `json:"code,omitempty"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"createdAt,omitempty"`
	UpdatedAt     string `json:"updatedAt,omitempty"`
}

// ProductItem is one finished catalog product (read-only). Raw materials and
// kits are never returned by the Partner API.
type ProductItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	SKU         string  `json:"sku,omitempty"`
	Price       float64 `json:"price"`
	Width       float64 `json:"width,omitempty"`
	Length      float64 `json:"length,omitempty"`
	Active      bool    `json:"active"`
	// Images are ready-to-use public URLs: the uploaded gallery (upload
	// order) first, then legacy image slots on the product record.
	Images    []string `json:"images,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	UpdatedAt string   `json:"updatedAt,omitempty"`
}

// PaginationParams are shared by every list endpoint.
type PaginationParams struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// LotStockListParams are query params for GET /partner/v1/lot-stock.
type LotStockListParams struct {
	PaginationParams
	// PurchaseType is ALL (the default — no category filter), Lot Stock,
	// Rental Return, or Immediate Sale. Case-insensitive; hyphens and
	// underscores are accepted (e.g. "lot-stock").
	PurchaseType string    `json:"purchaseType,omitempty"`
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

// LocationListParams are query params for GET /partner/v1/locations.
type LocationListParams struct {
	PaginationParams
	Search string    `json:"search,omitempty"`
	Sort   string    `json:"sort,omitempty"` // name | code | city | createdAt | updatedAt
	Order  SortOrder `json:"order,omitempty"`
	// Boolean flag filters apply only when non-nil (so false is still sent).
	Active   *bool `json:"active,omitempty"`
	SalesLot *bool `json:"salesLot,omitempty"`
	Plant    *bool `json:"plant,omitempty"`
}

// CustomerListParams are query params for GET /partner/v1/customers.
type CustomerListParams struct {
	PaginationParams
	Search      string    `json:"search,omitempty"`
	Sort        string    `json:"sort,omitempty"` // name | email | createdAt | updatedAt
	Order       SortOrder `json:"order,omitempty"`
	Email       string    `json:"email,omitempty"` // substring match
	Phone       string    `json:"phone,omitempty"` // substring match
	CreatedFrom string    `json:"createdFrom,omitempty"`
	CreatedTo   string    `json:"createdTo,omitempty"`
	UpdatedFrom string    `json:"updatedFrom,omitempty"`
	UpdatedTo   string    `json:"updatedTo,omitempty"`
}

// ProductListParams are query params for GET /partner/v1/products.
type ProductListParams struct {
	PaginationParams
	Search string    `json:"search,omitempty"`
	Sort   string    `json:"sort,omitempty"` // name | sku | price | createdAt | updatedAt
	Order  SortOrder `json:"order,omitempty"`
	SKU    string    `json:"sku,omitempty"` // substring match
	// Active filters by the active flag when non-nil.
	Active      *bool  `json:"active,omitempty"`
	CreatedFrom string `json:"createdFrom,omitempty"`
	CreatedTo   string `json:"createdTo,omitempty"`
	UpdatedFrom string `json:"updatedFrom,omitempty"`
	UpdatedTo   string `json:"updatedTo,omitempty"`
}

// LeadPatchRequest is the body for PATCH /partner/v1/leads/{id}.
type LeadPatchRequest struct {
	SalesLocation    string `json:"salesLocation,omitempty"`
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
}

// LeadCreateCustomer is the customer block of LeadCreateRequest. At least one
// of name, email, or phone is required.
type LeadCreateCustomer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// LeadCreateRequest is the body for POST /partner/v1/leads. When the
// salesperson fields are omitted, the location's lead-routing strategy
// auto-assigns one.
type LeadCreateRequest struct {
	LocationID       string             `json:"locationId"`
	Customer         LeadCreateCustomer `json:"customer"`
	SalespersonName  string             `json:"salespersonName,omitempty"`
	SalespersonEmail string             `json:"salespersonEmail,omitempty"`
}

// QuoteConvertRequest is the optional body for
// POST /partner/v1/quotes/{id}/convert — overrides the salesperson copied
// from the quote.
type QuoteConvertRequest struct {
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
}

// QuoteCreateCustomer is the customer block of QuoteCreateRequest. An email,
// or a name + phone, is required; the customer is matched by email when one
// already exists, otherwise created.
type QuoteCreateCustomer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// QuoteCreateDeliveryAddress is the optional delivery address stamped onto
// the created quote (and cascaded to the linked work order).
type QuoteCreateDeliveryAddress struct {
	Address   string   `json:"address,omitempty"`
	City      string   `json:"city,omitempty"`
	State     string   `json:"state,omitempty"`
	ZipCode   string   `json:"zipCode,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// QuoteCreateRequest is the body for POST /partner/v1/quotes — create a quote
// from an in-stock (on-lot) unit. Identify the unit with SerialNumber or
// WorkOrderID (WorkOrderID wins when both are sent).
type QuoteCreateRequest struct {
	SerialNumber string `json:"serialNumber,omitempty"`
	WorkOrderID  string `json:"workOrderId,omitempty"`
	LocationID   string `json:"locationId,omitempty"`
	// PurchaseType is "Lot Stock", "Rental Return", or "Immediate Sale";
	// inferred from the unit when omitted.
	PurchaseType    string                      `json:"purchaseType,omitempty"`
	Price           *float64                    `json:"price,omitempty"`
	Note            string                      `json:"note,omitempty"`
	Customer        QuoteCreateCustomer         `json:"customer"`
	DeliveryAddress *QuoteCreateDeliveryAddress `json:"deliveryAddress,omitempty"`
}

// QuotePatchRequest is the body for PATCH /partner/v1/quotes/{id}.
type QuotePatchRequest struct {
	CustomerName     string `json:"customerName,omitempty"`
	CustomerEmail    string `json:"customerEmail,omitempty"`
	CustomerPhone    string `json:"customerPhone,omitempty"`
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
	SalesLocation    string `json:"salesLocation,omitempty"`
	DeliveryAddress  string `json:"deliveryAddress,omitempty"`
	DeliveryCity     string `json:"deliveryCity,omitempty"`
	DeliveryState    string `json:"deliveryState,omitempty"`
	DeliveryZipCode  string `json:"deliveryZipCode,omitempty"`
}

// OrderPatchRequest is the body for PATCH /partner/v1/orders/{id}.
type OrderPatchRequest struct {
	CustomerName     string `json:"customerName,omitempty"`
	CustomerEmail    string `json:"customerEmail,omitempty"`
	CustomerPhone    string `json:"customerPhone,omitempty"`
	SalespersonName  string `json:"salespersonName,omitempty"`
	SalespersonEmail string `json:"salespersonEmail,omitempty"`
	SalesLocation    string `json:"salesLocation,omitempty"`
	DeliveryAddress  string `json:"deliveryAddress,omitempty"`
	DeliveryCity     string `json:"deliveryCity,omitempty"`
	DeliveryState    string `json:"deliveryState,omitempty"`
	DeliveryZipCode  string `json:"deliveryZipCode,omitempty"`
}

// WorkOrderPatchRequest is the body for PATCH /partner/v1/work-orders/{id}.
type WorkOrderPatchRequest struct {
	Title            string `json:"title,omitempty"`
	Description      string `json:"description,omitempty"`
	BuildingLocation string `json:"buildingLocation,omitempty"`
	PromisedDate     string `json:"promisedDate,omitempty"`
}

// LocationCreateRequest is the body for POST /partner/v1/locations.
// Name is required, plus either an Address or a Latitude/Longitude pair
// (provided together).
type LocationCreateRequest struct {
	Name          string   `json:"name"`
	Slug          string   `json:"slug,omitempty"`
	Code          string   `json:"code,omitempty"`
	Address       string   `json:"address,omitempty"`
	City          string   `json:"city,omitempty"`
	State         string   `json:"state,omitempty"`
	ZipCode       string   `json:"zipCode,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	ContactPerson string   `json:"contactPerson,omitempty"`
	ContactEmail  string   `json:"contactEmail,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	Active        *bool    `json:"active,omitempty"`
	SalesLot      *bool    `json:"salesLot,omitempty"`
	Plant         *bool    `json:"plant,omitempty"`
}

// LocationPatchRequest is the body for PATCH /partner/v1/locations/{id}.
// Same key set as create, all optional.
type LocationPatchRequest struct {
	Name          string   `json:"name,omitempty"`
	Slug          string   `json:"slug,omitempty"`
	Code          string   `json:"code,omitempty"`
	Address       string   `json:"address,omitempty"`
	City          string   `json:"city,omitempty"`
	State         string   `json:"state,omitempty"`
	ZipCode       string   `json:"zipCode,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	ContactPerson string   `json:"contactPerson,omitempty"`
	ContactEmail  string   `json:"contactEmail,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	Active        *bool    `json:"active,omitempty"`
	SalesLot      *bool    `json:"salesLot,omitempty"`
	Plant         *bool    `json:"plant,omitempty"`
}

// CustomerCreateRequest is the body for POST /partner/v1/customers.
// Email is required and must be unique within the company.
type CustomerCreateRequest struct {
	Email         string `json:"email"`
	Name          string `json:"name,omitempty"`
	ContactName   string `json:"contactName,omitempty"`
	ContactPerson string `json:"contactPerson,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	ZipCode       string `json:"zipCode,omitempty"`
	Code          string `json:"code,omitempty"`
}

// CustomerPatchRequest is the body for PATCH /partner/v1/customers/{id}.
type CustomerPatchRequest struct {
	Name          string `json:"name,omitempty"`
	ContactName   string `json:"contactName,omitempty"`
	ContactPerson string `json:"contactPerson,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	ZipCode       string `json:"zipCode,omitempty"`
	Code          string `json:"code,omitempty"`
	Active        *bool  `json:"active,omitempty"`
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
