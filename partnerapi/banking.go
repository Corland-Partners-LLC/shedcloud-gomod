package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// BankAccount is a bank account tied to a GL account
// (/partner/v1/bank-accounts). Money is a decimal string.
type BankAccount struct {
	ID                    string     `json:"id"`
	TenantID              string     `json:"tenant_id,omitempty"`
	Name                  string     `json:"name"`
	GLAccountID           string     `json:"gl_account_id"`
	LastReconciledBalance string     `json:"last_reconciled_balance"`
	LastReconciledAt      *time.Time `json:"last_reconciled_at,omitempty"`
	IsActive              bool       `json:"is_active"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// BankAccountsService calls /partner/v1/bank-accounts.
type BankAccountsService struct{ c *Client }

// List returns every bank account (no pagination). Scope
// partner-api.bank-accounts.read.
func (s *BankAccountsService) List(ctx context.Context) ([]BankAccount, error) {
	return listKey[BankAccount](ctx, s.c, http.MethodGet, "/partner/v1/bank-accounts", "bank_accounts", nil)
}

// BankReconciliation is one bank reconciliation period
// (/partner/v1/bank-reconciliations). Dates are YYYY-MM-DD; money is a decimal
// string. Status is "open" or "locked".
type BankReconciliation struct {
	ID                  string     `json:"id"`
	TenantID            string     `json:"tenant_id,omitempty"`
	BankAccountID       string     `json:"bank_account_id"`
	PeriodStart         string     `json:"period_start"`
	PeriodEnd           string     `json:"period_end"`
	StatementEndBalance string     `json:"statement_end_balance"`
	BookEndBalance      string     `json:"book_end_balance"`
	Status              string     `json:"status"`
	LockedAt            *time.Time `json:"locked_at,omitempty"`
	LockedBy            string     `json:"locked_by,omitempty"`
	CreatedBy           string     `json:"created_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// BankReconciliationListParams filters GET /partner/v1/bank-reconciliations.
type BankReconciliationListParams struct {
	BankAccountID string `json:"bank_account_id,omitempty"`
}

// BankReconciliationsService calls /partner/v1/bank-reconciliations.
type BankReconciliationsService struct{ c *Client }

// List returns reconciliations (no pagination), optionally for one bank
// account. Scope partner-api.bank-reconciliations.read.
func (s *BankReconciliationsService) List(ctx context.Context, params BankReconciliationListParams) ([]BankReconciliation, error) {
	return listKey[BankReconciliation](ctx, s.c, http.MethodGet, "/partner/v1/bank-reconciliations", "bank_reconciliations", params)
}
