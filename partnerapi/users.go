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

// Create makes a company user through the portal's own pipeline (profile,
// RBAC role, locations, invitation email). When the email belongs to an
// existing ShedCloud account it is linked to the company instead of
// duplicated (409 when already a member). The Cognito login is created at
// the user's first sign-in via the invite. Requires the
// partner-api.users.write scope.
func (s *UsersService) Create(ctx context.Context, body UserCreateRequest, opts ...RequestOption) (*UserItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out UserItem
	if err := s.c.http.request(ctx, http.MethodPost, "/partner/v1/users", nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches a company user: profile fields, role, locations
// (LocationIDs replaces), and Active (enable/disable — the owner is
// protected). Requires the partner-api.users.write scope.
func (s *UsersService) Update(ctx context.Context, id string, body UserPatchRequest) (*UserItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out UserItem
	path := fmt.Sprintf("/partner/v1/users/%s", url.PathEscape(id))
	if err := s.c.http.request(ctx, http.MethodPatch, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Roles lists the company's assignable RBAC roles, for discovering valid
// RoleID values. Requires the partner-api.users.read scope.
func (s *UsersService) Roles(ctx context.Context) (*RolesResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out RolesResponse
	if err := s.c.http.request(ctx, http.MethodGet, "/partner/v1/roles", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
