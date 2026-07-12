package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// EventsService calls /partner/v1/events — the cursor-based change event
// feed (lossless reconciliation) plus webhook delivery-log access and event
// redelivery. Requires the partner-api.events.read scope.
type EventsService struct {
	c *Client
}

// eventListQuery flattens EventListParams for query encoding (Types joined).
type eventListQuery struct {
	Cursor string `json:"cursor,omitempty"`
	Types  string `json:"types,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// List returns one page of events after params.Cursor.
func (s *EventsService) List(ctx context.Context, params EventListParams) (*EventListResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	q := eventListQuery{
		Cursor: params.Cursor,
		Types:  strings.Join(params.Types, ","),
		Limit:  params.Limit,
	}
	var out EventListResponse
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/events", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Each calls fn for every event after params.Cursor, transparently following
// pagination. Iteration stops when fn returns a non-nil error (returned
// as-is) or the feed is exhausted.
//
//	err := client.Events.Each(ctx, partnerapi.EventListParams{Cursor: lastSeen},
//		func(ev partnerapi.EventItem) error {
//			lastSeen = ev.ID
//			return handle(ev)
//		})
func (s *EventsService) Each(ctx context.Context, params EventListParams, fn func(EventItem) error) error {
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

// Redeliver re-enqueues one event to every active webhook subscription.
func (s *EventsService) Redeliver(ctx context.Context, id string) (*EventRedeliverResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out EventRedeliverResponse
	path := fmt.Sprintf("/partner/v1/events/%s/redeliver", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPost, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Deliveries returns the webhook delivery attempt log (newest first).
func (s *EventsService) Deliveries(ctx context.Context, params WebhookDeliveryListParams) (*PaginatedResponse[WebhookDeliveryItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[WebhookDeliveryItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/webhook-deliveries", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
