package partnerapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// UsersService calls /partner/v1/users — salespeople/users of the company.
// Resolve PartnerSalesperson.ID from sales DTOs into a full profile.
// Requires the partner-api.users.read scope.
type UsersService struct {
	c *Client
}

// List returns users for the authenticated company.
func (s *UsersService) List(ctx context.Context, params UserListParams) (*PaginatedResponse[UserItem], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out PaginatedResponse[UserItem]
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/users", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one user by id.
func (s *UsersService) Get(ctx context.Context, id string) (*UserItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out UserItem
	path := fmt.Sprintf("/partner/v1/users/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
