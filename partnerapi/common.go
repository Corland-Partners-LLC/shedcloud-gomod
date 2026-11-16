package partnerapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ListParams are the paging/search query params shared by every ERP list:
// page (1-based), page_size (default 50, max 100) and free-text q. Embed it in
// resource-specific params.
type ListParams struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	Q        string `json:"q,omitempty"`
}

// Pagination is the {page, page_size, total} block of every list response.
// TotalPages is only sent by leads and users.
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// Page is one page of a list. The ERP wraps each list as
// {"<resource>": [...], "pagination": {...}}; Items holds the array whatever
// the resource key is.
type Page[T any] struct {
	Items      []T
	Pagination Pagination
}

// OAuthTokenResponse is the POST /oauth/token reply (snake_case).
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

// listPage runs a GET list call and unpacks {key: [...], pagination: {...}}.
func listPage[T any](ctx context.Context, c *Client, path, key string, params any) (*Page[T], error) {
	var env map[string]json.RawMessage
	if err := c.http.request(ctx, http.MethodGet, path, params, nil, &env); err != nil {
		return nil, err
	}
	out := &Page[T]{Items: []T{}}
	if raw, ok := env[key]; ok && string(raw) != "null" {
		if err := json.Unmarshal(raw, &out.Items); err != nil {
			return nil, fmt.Errorf("partnerapi: decode %s: %w", key, err)
		}
	}
	if raw, ok := env["pagination"]; ok {
		_ = json.Unmarshal(raw, &out.Pagination)
	}
	return out, nil
}

// listKey unpacks {key: [...]} (lists without pagination: events, history, roles).
func listKey[T any](ctx context.Context, c *Client, method, path, key string, params any) ([]T, error) {
	var env map[string]json.RawMessage
	if err := c.http.request(ctx, method, path, params, nil, &env); err != nil {
		return nil, err
	}
	out := []T{}
	if raw, ok := env[key]; ok && string(raw) != "null" {
		if err := json.Unmarshal(raw, &out); err != nil {
			return nil, fmt.Errorf("partnerapi: decode %s: %w", key, err)
		}
	}
	return out, nil
}

// call is the one-shot helper: request + decode into a fresh *T.
func call[T any](ctx context.Context, c *Client, method, path string, query, body any, opts ...RequestOption) (*T, error) {
	var out T
	if err := c.http.request(ctx, method, path, query, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func seg(id string) string { return url.PathEscape(id) }
