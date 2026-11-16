package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// SoldUpgradesReportRequest is the POST /partner/v1/sold-upgrades body. Action
// is "added" (priced extras kept on the sold building, the ERP default) or
// "removed" (standard features taken off). SaleFrom and SaleTo are dates or
// timestamps; Page/PageSize default to 1/50 (max 100).
type SoldUpgradesReportRequest struct {
	Action          string `json:"action,omitempty"`
	SaleFrom        string `json:"sale_from,omitempty"`
	SaleTo          string `json:"sale_to,omitempty"`
	StoreLocationID string `json:"store_location_id,omitempty"`
	Q               string `json:"q,omitempty"`
	Page            int    `json:"page,omitempty"`
	PageSize        int    `json:"page_size,omitempty"`
	IncludeSums     bool   `json:"include_sums,omitempty"`
}

// SoldUpgradeRow is one exploded upgrade line on a booked order. Amount is a
// decimal string.
type SoldUpgradeRow struct {
	SoldAt            *time.Time `json:"sold_at,omitempty"`
	SaleDate          string     `json:"sale_date,omitempty"`
	OrderID           string     `json:"order_id"`
	OrderNumber       string     `json:"order_number"`
	CustomerName      string     `json:"customer_name"`
	StoreLocationID   string     `json:"store_location_id,omitempty"`
	StoreLocationName string     `json:"store_location_name"`
	SalespersonName   string     `json:"salesperson_name"`
	ProductName       string     `json:"product_name"`
	UpgradeLabel      string     `json:"upgrade_label"`
	Amount            string     `json:"amount"`
	Action            string     `json:"action"`
}

// SoldUpgradeLabelRollup is the count and sum for one upgrade label.
type SoldUpgradeLabelRollup struct {
	Label     string `json:"label"`
	Count     int64  `json:"count"`
	SumAmount string `json:"sum_amount"`
}

// SoldUpgradesPagination is the paging block of the sold-upgrades report
// (TotalItems, not Total).
type SoldUpgradesPagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
}

// SoldUpgradesReport is the sold-upgrades reply. SumAmount and ByLabel are
// only filled when the request set IncludeSums.
type SoldUpgradesReport struct {
	Rows       []SoldUpgradeRow         `json:"rows"`
	Pagination SoldUpgradesPagination   `json:"pagination"`
	SumAmount  string                   `json:"sum_amount,omitempty"`
	ByLabel    []SoldUpgradeLabelRollup `json:"by_label,omitempty"`
}

// SoldUpgradesService calls /partner/v1/sold-upgrades.
type SoldUpgradesService struct{ c *Client }

// Report returns the booked-order sold-upgrades report (a POST read). Scope
// partner-api.sold-upgrades.read.
func (s *SoldUpgradesService) Report(ctx context.Context, body SoldUpgradesReportRequest) (*SoldUpgradesReport, error) {
	return call[SoldUpgradesReport](ctx, s.c, http.MethodPost, "/partner/v1/sold-upgrades", nil, body)
}
