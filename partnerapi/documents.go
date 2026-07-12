package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DocumentsService calls /partner/v1/documents — file metadata plus
// short-lived presigned downloads for documents attached to orders, quotes,
// and work orders (contract PDFs, photos, …). Requires the
// partner-api.documents.read scope.
type DocumentsService struct {
	c *Client
}

// List returns documents attached to one entity. Params EntityType and
// EntityID are required.
func (s *DocumentsService) List(ctx context.Context, params DocumentListParams) (*PaginatedResponse[DocumentItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[DocumentItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/documents", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Download returns a presigned download URL for one document. The URL
// expires quickly (~10 minutes) — fetch it right before downloading, don't
// store it.
func (s *DocumentsService) Download(ctx context.Context, id string) (*DocumentDownload, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out DocumentDownload
	path := fmt.Sprintf("/partner/v1/documents/%s/download", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
