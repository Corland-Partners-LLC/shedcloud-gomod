package partnerapi

import (
	"context"
	"net/http"
)

// DocumentUploadURLRequest is POST /partner/v1/documents/upload-url.
type DocumentUploadURLRequest struct {
	EntityType  string `json:"entity_type,omitempty"`
	EntityID    string `json:"entity_id,omitempty"`
	Role        string `json:"role,omitempty"`
	FileName    string `json:"file_name,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}

// DocumentUploadTicket is the reply of UploadURL: PUT the bytes to Upload.URL
// (sending Upload.Headers), then register the file with Attach using
// StorageKey.
type DocumentUploadTicket struct {
	StorageKey string       `json:"storage_key"`
	Upload     PresignedURL `json:"upload"`
}

// DocumentAttachRequest is POST /partner/v1/documents: it registers a file
// already uploaded to StorageKey and links it to an entity.
type DocumentAttachRequest struct {
	EntityType string         `json:"entity_type,omitempty"`
	EntityID   string         `json:"entity_id,omitempty"`
	Role       string         `json:"role,omitempty"`
	FileName   string         `json:"file_name,omitempty"`
	MimeType   string         `json:"mime_type,omitempty"`
	SizeBytes  int64          `json:"size_bytes,omitempty"`
	StorageKey string         `json:"storage_key,omitempty"`
	Checksum   string         `json:"checksum,omitempty"`
	Caption    string         `json:"caption,omitempty"`
	Position   int            `json:"position,omitempty"`
	IsPrimary  bool           `json:"is_primary,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

// DocumentAttachment is the reply of Attach: the Document plus the storage
// key, a URL and the access level when the ERP sends them.
type DocumentAttachment struct {
	Document
	StorageKey string `json:"storage_key,omitempty"`
	URL        string `json:"url,omitempty"`
	Access     string `json:"access,omitempty"`
}

// UploadURL requests a presigned PUT URL for a new file. Scope
// partner-api.documents.write.
func (s *DocumentsService) UploadURL(ctx context.Context, body DocumentUploadURLRequest) (*DocumentUploadTicket, error) {
	return call[DocumentUploadTicket](ctx, s.c, http.MethodPost, "/partner/v1/documents/upload-url", nil, body)
}

// Attach registers an uploaded file and links it to an entity (POST
// /partner/v1/documents). Scope partner-api.documents.write. Pass
// WithIdempotencyKey to make retries safe.
func (s *DocumentsService) Attach(ctx context.Context, body DocumentAttachRequest, opts ...RequestOption) (*DocumentAttachment, error) {
	return call[DocumentAttachment](ctx, s.c, http.MethodPost, "/partner/v1/documents", nil, body, opts...)
}
