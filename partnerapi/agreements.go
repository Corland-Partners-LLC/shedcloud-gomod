package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// AgreementsService calls /partner/v1/agreements (read-only). B2B dealer ↔
// RTO / ERP / DMAN partnerships — invitations and acceptance stay portal-only.
type AgreementsService struct {
	c *Client
}

// List returns the company's B2B agreements.
func (s *AgreementsService) List(ctx context.Context, params AgreementListParams) (*PaginatedResponse[AgreementItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[AgreementItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/agreements", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one agreement when the caller's company is either endpoint.
func (s *AgreementsService) Get(ctx context.Context, id string) (*AgreementItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out AgreementItem
	path := fmt.Sprintf("/partner/v1/agreements/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Active returns the most recently updated active rate-program agreement.
func (s *AgreementsService) Active(ctx context.Context) (*AgreementItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out AgreementItem
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/agreements/active", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListStateLegal returns per-state RTO legal configuration for an agreement.
func (s *AgreementsService) ListStateLegal(ctx context.Context, agreementID string) (*AgreementStateLegalListResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out AgreementStateLegalListResponse
	path := fmt.Sprintf("/partner/v1/agreements/%s/state-legal", url.PathEscape(agreementID))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetStateLegal returns the full legal appendix for one agreement and US state.
func (s *AgreementsService) GetStateLegal(ctx context.Context, agreementID, state string) (*AgreementStateLegalItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out AgreementStateLegalItem
	path := fmt.Sprintf("/partner/v1/agreements/%s/state-legal/%s", url.PathEscape(agreementID), url.PathEscape(state))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
