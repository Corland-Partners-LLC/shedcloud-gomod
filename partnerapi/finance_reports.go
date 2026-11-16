package partnerapi

import (
	"context"
	"net/http"
)

// FinancialStatementRow is one GL account rollup on a P&L or balance sheet.
// Money is a decimal string.
type FinancialStatementRow struct {
	GLAccountID   string `json:"gl_account_id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	AccountType   string `json:"account_type"`
	NormalBalance string `json:"normal_balance"`
	TotalDebit    string `json:"total_debit"`
	TotalCredit   string `json:"total_credit"`
	Balance       string `json:"balance"`
}

// FinancialStatementGroup is one section of a financial statement.
type FinancialStatementGroup struct {
	Key   string                  `json:"key"`
	Label string                  `json:"label"`
	Rows  []FinancialStatementRow `json:"rows"`
	Total string                  `json:"total"`
}

// ProfitAndLoss is a dated income-statement rollup of posted journal lines
// (GET /partner/v1/financial-statements/profit-and-loss). From/To are
// YYYY-MM-DD.
type ProfitAndLoss struct {
	From        string                    `json:"from"`
	To          string                    `json:"to"`
	Rows        []FinancialStatementRow   `json:"rows"`
	Groups      []FinancialStatementGroup `json:"groups"`
	TotalDebit  string                    `json:"total_debit"`
	TotalCredit string                    `json:"total_credit"`
	NetIncome   string                    `json:"net_income"`
}

// BalanceSheet is an as-of rollup of posted balance-sheet accounts
// (GET /partner/v1/financial-statements/balance-sheet).
type BalanceSheet struct {
	AsOf        string                    `json:"as_of"`
	Rows        []FinancialStatementRow   `json:"rows"`
	Groups      []FinancialStatementGroup `json:"groups"`
	TotalDebit  string                    `json:"total_debit"`
	TotalCredit string                    `json:"total_credit"`
	NetIncome   string                    `json:"net_income"`
	InBalance   bool                      `json:"in_balance"`
}

// ProfitAndLossParams are the query params of the P&L (dates YYYY-MM-DD).
type ProfitAndLossParams struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// BalanceSheetParams are the query params of the balance sheet.
type BalanceSheetParams struct {
	AsOf string `json:"as_of,omitempty"`
}

// FinancialStatementsService calls /partner/v1/financial-statements.
type FinancialStatementsService struct{ c *Client }

// ProfitAndLoss returns the income statement for From..To. Scope
// partner-api.financial-statements.read.
func (s *FinancialStatementsService) ProfitAndLoss(ctx context.Context, params ProfitAndLossParams) (*ProfitAndLoss, error) {
	return call[ProfitAndLoss](ctx, s.c, http.MethodGet, "/partner/v1/financial-statements/profit-and-loss", params, nil)
}

// BalanceSheet returns the balance sheet as of a date. Scope
// partner-api.financial-statements.read.
func (s *FinancialStatementsService) BalanceSheet(ctx context.Context, params BalanceSheetParams) (*BalanceSheet, error) {
	return call[BalanceSheet](ctx, s.c, http.MethodGet, "/partner/v1/financial-statements/balance-sheet", params, nil)
}
