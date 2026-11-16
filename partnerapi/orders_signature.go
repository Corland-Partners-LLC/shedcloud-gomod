package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// SignatureProcessReference is one personal reference on the contract sheet.
type SignatureProcessReference struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Relationship string `json:"relationship"`
}

// SignatureProcessSheet is the primary renter's contract data sheet. The SSN
// is never returned in full (SSNMasked is ***-**-1234 or fully masked).
type SignatureProcessSheet struct {
	OrderID            string                      `json:"order_id"`
	TenantID           string                      `json:"tenant_id"`
	FullName           string                      `json:"full_name"`
	PrimaryPhone       string                      `json:"primary_phone"`
	SecondaryPhone     string                      `json:"secondary_phone"`
	Email              string                      `json:"email,omitempty"`
	DateOfBirth        string                      `json:"date_of_birth"`
	SSNMasked          string                      `json:"ssn_masked,omitempty"`
	HasSSN             bool                        `json:"has_ssn"`
	DriversLicense     string                      `json:"drivers_license"`
	Employer           string                      `json:"employer"`
	Position           string                      `json:"position"`
	WorkPhone          string                      `json:"work_phone"`
	LengthOnJob        string                      `json:"length_on_job"`
	YearsAtAddress     string                      `json:"years_at_address"`
	MailingLine1       string                      `json:"mailing_line1"`
	MailingLine2       string                      `json:"mailing_line2"`
	MailingCity        string                      `json:"mailing_city"`
	MailingState       string                      `json:"mailing_state"`
	MailingPostal      string                      `json:"mailing_postal"`
	AuthorizeContactAt string                      `json:"authorize_contact_at"`
	References         []SignatureProcessReference `json:"references"`
	ConditionNew       *bool                       `json:"condition_new,omitempty"`
	ResidenceOwn       *bool                       `json:"residence_own,omitempty"`
	PaperlessBilling   bool                        `json:"paperless_billing"`
	DueDateFirst       bool                        `json:"due_date_first"`
	DueDateFifteenth   bool                        `json:"due_date_fifteenth"`
	LDWAccepted        *bool                       `json:"ldw_accepted,omitempty"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
}

// SignatureProcessCorenter is the optional co-renter party.
type SignatureProcessCorenter struct {
	OrderID            string    `json:"order_id"`
	TenantID           string    `json:"tenant_id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	MiddleName         string    `json:"middle_name"`
	Phone              string    `json:"phone"`
	Email              string    `json:"email,omitempty"`
	DateOfBirth        string    `json:"date_of_birth"`
	SSNMasked          string    `json:"ssn_masked,omitempty"`
	HasSSN             bool      `json:"has_ssn"`
	DriversLicense     string    `json:"drivers_license"`
	AuthorizeContactAt string    `json:"authorize_contact_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// SignatureProcessSignature is one signature artifact collected for the
// order. Role is customer, corenter or salesperson.
type SignatureProcessSignature struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	OrderID     string    `json:"order_id"`
	Role        string    `json:"role"`
	StorageKey  string    `json:"storage_key"`
	URL         string    `json:"url,omitempty"`
	SignerLabel string    `json:"signer_label,omitempty"`
	SignedAt    time.Time `json:"signed_at"`
	IPAddress   string    `json:"ip_address,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// SignatureProcessMissingField is one completeness gap that blocks signing.
type SignatureProcessMissingField struct {
	Party     string `json:"party"`
	SectionID string `json:"section_id"`
	Field     string `json:"field"`
	Label     string `json:"label"`
}

// SignatureProcessFollow is the signature-process state of an order.
// ProcessStatus is none, out_for_signature, generating_pdf, complete or
// cancelled; ContractPDFStatus is none, pending, generating, ready or failed.
type SignatureProcessFollow struct {
	ProcessStatus     string                         `json:"signature_process_status"`
	ContractPDFStatus string                         `json:"contract_pdf_status"`
	ContractPDFError  string                         `json:"contract_pdf_error,omitempty"`
	StartedAt         *time.Time                     `json:"started_at,omitempty"`
	CompletedAt       *time.Time                     `json:"completed_at,omitempty"`
	IsRTO             bool                           `json:"is_rto"`
	CorenterRequired  bool                           `json:"corenter_required"`
	Sheet             *SignatureProcessSheet         `json:"sheet,omitempty"`
	Corenter          *SignatureProcessCorenter      `json:"corenter,omitempty"`
	Signatures        []SignatureProcessSignature    `json:"signatures"`
	Missing           []SignatureProcessMissingField `json:"missing"`
	PartiesComplete   bool                           `json:"parties_complete"`
	NextParty         string                         `json:"next_party,omitempty"`
	HasCustomerShare  bool                           `json:"has_customer_share"`
	HasCorenterShare  bool                           `json:"has_corenter_share"`
	CustomerSigned    bool                           `json:"customer_signed"`
	CorenterSigned    bool                           `json:"corenter_signed"`
	SalespersonSigned bool                           `json:"salesperson_signed"`

	SignatureProvider  string   `json:"signature_provider"`
	BoldSignDocumentID string   `json:"boldsign_document_id,omitempty"`
	BoldSignStatus     string   `json:"boldsign_status,omitempty"`
	BoldSignConnected  bool     `json:"boldsign_connected"`
	DocuSignEnvelopeID string   `json:"docusign_envelope_id,omitempty"`
	DocuSignStatus     string   `json:"docusign_status,omitempty"`
	DocuSignConnected  bool     `json:"docusign_connected"`
	AvailableProviders []string `json:"available_providers"`
}

// SignatureStartRequest is POST /partner/v1/orders/{id}/signature-process/start.
// Send true also emails the signing link. Provider may only be "native" (or
// empty); the ERP rejects any other value.
type SignatureStartRequest struct {
	Send     bool   `json:"send,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// SignatureSendRequest is POST /partner/v1/orders/{id}/signature-process/send.
// To is "customer" or "corenter".
type SignatureSendRequest struct {
	To string `json:"to,omitempty"`
}

// SignatureRestartRequest is POST /partner/v1/orders/{id}/signature-process/restart.
type SignatureRestartRequest struct {
	ClearForms bool `json:"clear_forms,omitempty"`
}

// SignatureProcess returns the order's signature process state. Scope
// partner-api.contracts.read.
func (s *OrdersService) SignatureProcess(ctx context.Context, id string) (*SignatureProcessFollow, error) {
	return call[SignatureProcessFollow](ctx, s.c, http.MethodGet, s.path(id, "/signature-process"), nil, nil)
}

// StartSignature starts the native signature process. Scope
// partner-api.orders.write.
func (s *OrdersService) StartSignature(ctx context.Context, id string, body SignatureStartRequest) (*SignatureProcessFollow, error) {
	return call[SignatureProcessFollow](ctx, s.c, http.MethodPost, s.path(id, "/signature-process/start"), nil, body)
}

// SendSignature emails the signing link to the customer or co-renter and
// returns the order. Scope partner-api.orders.write.
func (s *OrdersService) SendSignature(ctx context.Context, id string, body SignatureSendRequest) (*Order, error) {
	return call[Order](ctx, s.c, http.MethodPost, s.path(id, "/signature-process/send"), nil, body)
}

// CancelSignature cancels the running signature process. Scope
// partner-api.orders.write.
func (s *OrdersService) CancelSignature(ctx context.Context, id string) (*SignatureProcessFollow, error) {
	return call[SignatureProcessFollow](ctx, s.c, http.MethodPost, s.path(id, "/signature-process/cancel"), nil, nil)
}

// RestartSignature restarts the signature process (ClearForms also clears
// the filled forms). Scope partner-api.orders.write.
func (s *OrdersService) RestartSignature(ctx context.Context, id string, body SignatureRestartRequest) (*SignatureProcessFollow, error) {
	return call[SignatureProcessFollow](ctx, s.c, http.MethodPost, s.path(id, "/signature-process/restart"), nil, body)
}

// RetryContractPDF retries a failed contract PDF generation. Scope
// partner-api.orders.write.
func (s *OrdersService) RetryContractPDF(ctx context.Context, id string) (*SignatureProcessFollow, error) {
	return call[SignatureProcessFollow](ctx, s.c, http.MethodPost, s.path(id, "/signature-process/retry-pdf"), nil, nil)
}
