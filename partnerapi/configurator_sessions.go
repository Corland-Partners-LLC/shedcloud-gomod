package partnerapi

import (
	"context"
	"net/http"
)

// ConfiguratorSessionsService calls /partner/v1/configurator-sessions —
// mint single-use launch URLs that open the ShedCloud 3D configurator for a
// customer (blank, from a saved quote, or from an in-stock unit). Requires
// the partner-api.configurator-sessions.write scope.
type ConfiguratorSessionsService struct {
	c *Client
}

// Create mints a launch session. Supports WithIdempotencyKey.
func (s *ConfiguratorSessionsService) Create(ctx context.Context, body ConfiguratorSessionCreateRequest, opts ...RequestOption) (*ConfiguratorSessionCreateResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out ConfiguratorSessionCreateResponse
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/configurator-sessions", nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
