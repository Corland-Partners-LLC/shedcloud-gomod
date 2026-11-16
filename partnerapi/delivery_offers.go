package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// DeliveryOfferCharges is the carrier pay on an offer (delivery + fuel +
// additional lines). Money values are decimal strings. The ERP only sets it
// for carrier callers that may see pay.
type DeliveryOfferCharges struct {
	DeliveryCharge    string                     `json:"delivery_charge"`
	FuelSurcharge     string                     `json:"fuel_surcharge"`
	AdditionalCharges []DeliveryOfferExtraCharge `json:"additional_charges"`
	Total             string                     `json:"total"`
}

// DeliveryOfferExtraCharge is one named additional charge line.
type DeliveryOfferExtraCharge struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

// DeliveryOfferDisplay is the work-order / carrier join the ERP attaches to
// manufacturer list rows.
type DeliveryOfferDisplay struct {
	WONumber         string   `json:"wo_number,omitempty"`
	CustomerName     string   `json:"customer_name,omitempty"`
	ProductName      string   `json:"product_name,omitempty"`
	Address          string   `json:"address,omitempty"`
	Lat              *float64 `json:"lat,omitempty"`
	Lng              *float64 `json:"lng,omitempty"`
	CarrierName      string   `json:"carrier_name,omitempty"`
	SerialNumber     string   `json:"serial_number,omitempty"`
	WOStatus         string   `json:"wo_status,omitempty"`
	InvoiceNumber    string   `json:"invoice_number,omitempty"`
	InvoiceAt        string   `json:"invoice_at,omitempty"`
	InvoiceEmailedTo string   `json:"invoice_emailed_to,omitempty"`
	InvoiceEmailedAt string   `json:"invoice_emailed_at,omitempty"`
}

// DeliveryOffer is a manufacturer to carrier share of one work order
// (ShedRun). Status is offered, accepted, declined, cancelled, in_transit,
// delivered, failed or rescheduled. Money is a decimal string.
type DeliveryOffer struct {
	ID                      string                `json:"id"`
	ManufacturerTenantID    string                `json:"manufacturer_tenant_id"`
	CarrierTenantID         string                `json:"carrier_tenant_id"`
	WorkOrderID             string                `json:"work_order_id"`
	Kind                    string                `json:"kind"`
	RentalReturnID          string                `json:"rental_return_id,omitempty"`
	ReturnLocationID        string                `json:"return_location_id,omitempty"`
	TransportationCompanyID string                `json:"transportation_company_id,omitempty"`
	Status                  string                `json:"status"`
	OfferedBy               string                `json:"offered_by,omitempty"`
	OfferedAt               time.Time             `json:"offered_at"`
	RespondedBy             string                `json:"responded_by,omitempty"`
	RespondedAt             *time.Time            `json:"responded_at,omitempty"`
	DeclineReason           string                `json:"decline_reason,omitempty"`
	CancelledReason         string                `json:"cancelled_reason,omitempty"`
	FailedReason            string                `json:"failed_reason,omitempty"`
	RescheduleTo            string                `json:"reschedule_to,omitempty"`
	PlannedPickupAt         string                `json:"planned_pickup_at,omitempty"`
	PlannedDeliveryAt       string                `json:"planned_delivery_at,omitempty"`
	AssignedDriverID        string                `json:"assigned_driver_id,omitempty"`
	AssignedDriverName      string                `json:"assigned_driver_name,omitempty"`
	AssignedVehicleID       string                `json:"assigned_vehicle_id,omitempty"`
	AssignedVehicleName     string                `json:"assigned_vehicle_name,omitempty"`
	HaulFee                 string                `json:"haul_fee,omitempty"`
	RunID                   string                `json:"run_id,omitempty"`
	TransportOrderID        string                `json:"transport_order_id,omitempty"`
	DeliveredAt             *time.Time            `json:"delivered_at,omitempty"`
	PublicTrackToken        string                `json:"public_track_token,omitempty"`
	Display                 *DeliveryOfferDisplay `json:"display,omitempty"`
	Charges                 *DeliveryOfferCharges `json:"charges,omitempty"`
	CreatedAt               time.Time             `json:"created_at"`
	UpdatedAt               time.Time             `json:"updated_at"`
	HandoffID               string                `json:"handoff_id,omitempty"`
	Executor                string                `json:"executor,omitempty"`
	ContractorCompanyName   string                `json:"contractor_company_name,omitempty"`
	FromCompanyName         string                `json:"from_company_name,omitempty"`
	ScheduledPickupAt       string                `json:"scheduled_pickup_at,omitempty"`
	ScheduledDeliveryAt     string                `json:"scheduled_delivery_at,omitempty"`
}

// DeliveryOfferListParams filters GET /partner/v1/delivery-offers. The list
// is not paginated.
type DeliveryOfferListParams struct {
	WorkOrderID string `json:"work_order_id,omitempty"`
}

// DeliveryOfferCreateRequest is POST /partner/v1/work-orders/{id}/delivery-offers
// (share a work order with a ShedRun carrier). HaulFee is a decimal string.
type DeliveryOfferCreateRequest struct {
	TransportationCompanyID string `json:"transportation_company_id,omitempty"`
	CarrierTenantID         string `json:"carrier_tenant_id,omitempty"`
	AgreementCode           string `json:"agreement_code,omitempty"`
	Message                 string `json:"message,omitempty"`
	HaulFee                 string `json:"haul_fee,omitempty"`
}

// DeliveryOfferAcceptResult is the accept reply. PublicTrackToken is only
// returned once, on this call.
type DeliveryOfferAcceptResult struct {
	Offer            DeliveryOffer `json:"offer"`
	PublicTrackToken string        `json:"public_track_token"`
}

// DeliveryOfferTransitionRequest is POST /partner/v1/delivery-offers/{id}/transition.
type DeliveryOfferTransitionRequest struct {
	Status string `json:"status,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// DeliveryOfferPODRequest is POST /partner/v1/delivery-offers/{id}/pod
// (proof of delivery). Checklist and Photos are free-form arrays.
type DeliveryOfferPODRequest struct {
	Signature string `json:"signature,omitempty"`
	Checklist []any  `json:"checklist,omitempty"`
	Photos    []any  `json:"photos,omitempty"`
	GeoLat    string `json:"geo_lat,omitempty"`
	GeoLng    string `json:"geo_lng,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

// DeliveryOfferPOD is a stored proof of delivery.
type DeliveryOfferPOD struct {
	ID             string     `json:"id"`
	Signature      string     `json:"signature,omitempty"`
	Checklist      []any      `json:"checklist,omitempty"`
	Photos         []any      `json:"photos,omitempty"`
	GeoLat         string     `json:"geo_lat,omitempty"`
	GeoLng         string     `json:"geo_lng,omitempty"`
	Notes          string     `json:"notes,omitempty"`
	ProofNoteID    string     `json:"proof_note_id,omitempty"`
	ProofEmailedAt *time.Time `json:"proof_emailed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// DeliveryOffersService calls the ShedRun delivery-offer routes under
// /partner/v1.
type DeliveryOffersService struct{ c *Client }

// List returns the delivery offers of the company (optionally one work
// order's). Scope partner-api.transport.read.
func (s *DeliveryOffersService) List(ctx context.Context, params DeliveryOfferListParams) ([]DeliveryOffer, error) {
	return listKey[DeliveryOffer](ctx, s.c, http.MethodGet, "/partner/v1/delivery-offers", "offers", params)
}

// Create shares a work order with a ShedRun carrier. Scope
// partner-api.transport.write. Pass WithIdempotencyKey to make retries safe.
func (s *DeliveryOffersService) Create(ctx context.Context, workOrderID string, body DeliveryOfferCreateRequest, opts ...RequestOption) (*DeliveryOffer, error) {
	return call[DeliveryOffer](ctx, s.c, http.MethodPost, "/partner/v1/work-orders/"+seg(workOrderID)+"/delivery-offers", nil, body, opts...)
}

// Accept accepts a share and returns the public track token (only once).
// Scope partner-api.transport.write.
func (s *DeliveryOffersService) Accept(ctx context.Context, offerID string) (*DeliveryOfferAcceptResult, error) {
	return call[DeliveryOfferAcceptResult](ctx, s.c, http.MethodPost, "/partner/v1/delivery-offers/"+seg(offerID)+"/accept", nil, nil)
}

// Transition moves the offer to another status (carrier status move). Scope
// partner-api.transport.write.
func (s *DeliveryOffersService) Transition(ctx context.Context, offerID string, body DeliveryOfferTransitionRequest) (*DeliveryOffer, error) {
	return call[DeliveryOffer](ctx, s.c, http.MethodPost, "/partner/v1/delivery-offers/"+seg(offerID)+"/transition", nil, body)
}

// SubmitPOD records the proof of delivery for an offer. Scope
// partner-api.transport.write.
func (s *DeliveryOffersService) SubmitPOD(ctx context.Context, offerID string, body DeliveryOfferPODRequest) (*DeliveryOfferPOD, error) {
	return call[DeliveryOfferPOD](ctx, s.c, http.MethodPost, "/partner/v1/delivery-offers/"+seg(offerID)+"/pod", nil, body)
}
