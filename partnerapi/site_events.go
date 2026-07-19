package partnerapi

import (
	"context"
	"net/http"
	"strings"
)

// SiteEventsService calls /partner/v1/site-events — visitor behavioral
// tracking for partner marketing sites. Events flow into the same analytics
// pipeline as the 3D configurator tracker, so a shopper's journey stitches
// across both properties via a shared visitor id.
//
// Track requires the partner-api.site-events.write scope; List/Each require
// partner-api.site-events.read.
//
// Partner keys are server-side secrets: proxy events through your backend,
// never call this from the browser. Batch events client-side before
// forwarding — the endpoint is rate limited per company.
type SiteEventsService struct {
	c *Client
}

// Track ingests a batch of behavioral events (batch-only; max 25 per call).
func (s *SiteEventsService) Track(ctx context.Context, req SiteEventsTrackRequest) (*SiteEventsTrackResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out SiteEventsTrackResponse
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/site-events", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// siteEventListQuery flattens SiteEventListParams for query encoding.
type siteEventListQuery struct {
	Cursor    string `json:"cursor,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	SessionID string `json:"sessionId,omitempty"`
	Types     string `json:"types,omitempty"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
}

// List returns one page of tracked events, newest first (~90-day retention).
func (s *SiteEventsService) List(ctx context.Context, params SiteEventListParams) (*SiteEventListResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := siteEventListQuery{
		Cursor:    params.Cursor,
		Limit:     params.Limit,
		SessionID: params.SessionID,
		Types:     strings.Join(params.Types, ","),
		From:      params.From,
		To:        params.To,
	}
	var out SiteEventListResponse
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/site-events", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Each calls fn for every matching event, transparently following
// pagination. Iteration stops when fn returns a non-nil error (returned
// as-is) or the feed is exhausted.
func (s *SiteEventsService) Each(ctx context.Context, params SiteEventListParams, fn func(SiteEventItem) error) error {
	if ctx == nil {
		ctx = context.Background()
	}
	cursor := params.Cursor
	for {
		page := params
		page.Cursor = cursor
		resp, err := s.List(ctx, page)
		if err != nil {
			return err
		}
		for _, ev := range resp.Data {
			if err := fn(ev); err != nil {
				return err
			}
		}
		if !resp.HasMore || resp.NextCursor == "" {
			return nil
		}
		cursor = resp.NextCursor
	}
}
