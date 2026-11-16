package partnerapi

import (
	"context"
	"testing"
)

const pageBlock = `"pagination":{"page":2,"page_size":10,"total":25}`

func TestFinEntBankAccountsList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"bank_accounts":[{"id":"b1","name":"Main","gl_account_id":"g1","last_reconciled_balance":"100.50","last_reconciled_at":"2026-01-31T00:00:00Z","is_active":true}]}`)
	got, err := (&BankAccountsService{c: c}).List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/bank-accounts")
	if len(got) != 1 || got[0].GLAccountID != "g1" || got[0].LastReconciledBalance != "100.50" || !got[0].IsActive || got[0].LastReconciledAt == nil {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntBankReconciliationsList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"bank_reconciliations":[{"id":"r1","bank_account_id":"b1","period_start":"2026-01-01","period_end":"2026-01-31","statement_end_balance":"10.00","book_end_balance":"9.99","status":"locked","locked_by":"u1"}]}`)
	got, err := (&BankReconciliationsService{c: c}).List(context.Background(), BankReconciliationListParams{BankAccountID: "b1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/bank-reconciliations")
	if q(rec, "bank_account_id") != "b1" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(got) != 1 || got[0].PeriodEnd != "2026-01-31" || got[0].BookEndBalance != "9.99" || got[0].Status != "locked" {
		t.Fatalf("got %+v", got)
	}
}

const stmtRows = `"rows":[{"gl_account_id":"g1","code":"4000","name":"Sales","account_type":"revenue","normal_balance":"credit","total_debit":"0.00","total_credit":"500.00","balance":"500.00"}],"groups":[{"key":"revenue","label":"Revenue","rows":[],"total":"500.00"}]`

func TestFinEntProfitAndLoss(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"from":"2026-01-01","to":"2026-01-31",`+stmtRows+`,"total_debit":"0.00","total_credit":"500.00","net_income":"500.00"}`)
	got, err := (&FinancialStatementsService{c: c}).ProfitAndLoss(context.Background(), ProfitAndLossParams{From: "2026-01-01", To: "2026-01-31"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/financial-statements/profit-and-loss")
	if q(rec, "from") != "2026-01-01" || q(rec, "to") != "2026-01-31" {
		t.Fatalf("query = %v", rec.Query)
	}
	if got.NetIncome != "500.00" || len(got.Rows) != 1 || got.Rows[0].NormalBalance != "credit" || got.Groups[0].Total != "500.00" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntBalanceSheet(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"as_of":"2026-01-31",`+stmtRows+`,"total_debit":"1.00","total_credit":"1.00","net_income":"0.00","in_balance":true}`)
	got, err := (&FinancialStatementsService{c: c}).BalanceSheet(context.Background(), BalanceSheetParams{AsOf: "2026-01-31"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/financial-statements/balance-sheet")
	if q(rec, "as_of") != "2026-01-31" {
		t.Fatalf("query = %v", rec.Query)
	}
	if got.AsOf != "2026-01-31" || !got.InBalance || got.TotalDebit != "1.00" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntCommissionsPlans(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"plans":[{"id":"p1","name":"Std","is_active":true,"rules":[{"id":"r1","plan_id":"p1","applies_to":"order","rate_pct":"2.5"}]}],`+pageBlock+`}`)
	got, err := (&CommissionsService{c: c}).Plans(context.Background(), CommissionListParams{ListParams: ListParams{Page: 2, PageSize: 10, Q: "std"}})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/commissions/plans")
	if q(rec, "page") != "2" || q(rec, "page_size") != "10" || q(rec, "q") != "std" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(got.Items) != 1 || got.Items[0].Rules[0].RatePct != "2.5" || got.Pagination.Total != 25 {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntCommissionsAccruals(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"accruals":[{"id":"a1","plan_id":"p1","order_id":"o1","salesperson_id":"s1","status":"accrued","amount":"12.34"}],`+pageBlock+`}`)
	got, err := (&CommissionsService{c: c}).Accruals(context.Background(), CommissionListParams{})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/commissions/accruals")
	if len(got.Items) != 1 || got.Items[0].Amount != "12.34" || got.Items[0].PlanID == nil || *got.Items[0].PlanID != "p1" || got.Items[0].SalespersonID != "s1" {
		t.Fatalf("got %+v", got)
	}
}

const billJSON = `{"id":"fb1","number":"FB-1","carrier_id":"c1","run_id":"ru1","invoice_amount":"110.00","rated_amount":"100.00","variance_amount":"10.00","status":"pending"}`

func TestFinEntFreightBillsList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"freight_bills":[`+billJSON+`],`+pageBlock+`}`)
	got, err := (&FreightBillsService{c: c}).List(context.Background(), FreightBillListParams{ListParams: ListParams{Page: 2}, Status: "pending"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/freight-bills")
	if q(rec, "status") != "pending" || q(rec, "page") != "2" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(got.Items) != 1 || got.Items[0].VarianceAmount != "10.00" || got.Items[0].CarrierID == nil || *got.Items[0].CarrierID != "c1" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntFreightBillsGet(t *testing.T) {
	c, rec := fakeERP(t, 200, billJSON)
	got, err := (&FreightBillsService{c: c}).Get(context.Background(), "fb1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/freight-bills/fb1")
	if got.Number != "FB-1" || got.InvoiceAmount != "110.00" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntFreightBillsCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, billJSON)
	got, err := (&FreightBillsService{c: c}).Create(context.Background(), FreightBillRequest{
		CarrierID: "c1", RunID: "ru1", TransportOrderID: "t1", InvoiceAmount: "110.00", RatedAmount: "100.00", Status: "pending",
	}, WithIdempotencyKey("k1"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/freight-bills")
	if rec.Header.Get("Idempotency-Key") != "k1" {
		t.Fatalf("idempotency header = %q", rec.Header.Get("Idempotency-Key"))
	}
	b := rec.json(t)
	if b["carrier_id"] != "c1" || b["invoice_amount"] != "110.00" || b["transport_order_id"] != "t1" {
		t.Fatalf("body = %v", b)
	}
	if _, ok := b["offer_id"]; ok {
		t.Fatalf("offer_id should be omitted: %v", b)
	}
	if got.ID != "fb1" {
		t.Fatalf("got %+v", got)
	}
}

const assetJSON = `{"id":"fa1","name":"Forklift","asset_account_id":"g1","placed_in_service":"2026-01-01","cost":"10000.00","salvage":"500.00","life_months":60,"method":"straight_line","status":"active"}`

func TestFinEntFixedAssetsList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"fixed_assets":[`+assetJSON+`],`+pageBlock+`}`)
	got, err := (&FixedAssetsService{c: c}).List(context.Background(), FixedAssetListParams{Status: "active"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/fixed-assets")
	if q(rec, "status") != "active" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(got.Items) != 1 || got.Items[0].Cost != "10000.00" || got.Items[0].LifeMonths != 60 || got.Items[0].PlacedInService == nil || got.Pagination.Total != 25 {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntFixedAssetsGet(t *testing.T) {
	c, rec := fakeERP(t, 200, assetJSON)
	got, err := (&FixedAssetsService{c: c}).Get(context.Background(), "fa1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/fixed-assets/fa1")
	if got.Name != "Forklift" || got.Salvage != "500.00" || got.Method != "straight_line" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntFixedAssetsCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, assetJSON)
	got, err := (&FixedAssetsService{c: c}).Create(context.Background(), FixedAssetRequest{
		Name: "Forklift", AssetAccountID: "g1", PlacedInService: "2026-01-01", Cost: "10000.00", Salvage: "500.00", LifeMonths: 60,
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/fixed-assets")
	b := rec.json(t)
	if b["name"] != "Forklift" || b["cost"] != "10000.00" || b["life_months"] != float64(60) || b["placed_in_service"] != "2026-01-01" {
		t.Fatalf("body = %v", b)
	}
	if _, ok := b["expense_account_id"]; ok {
		t.Fatalf("expense_account_id should be omitted: %v", b)
	}
	if got.ID != "fa1" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntFXRates(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"currency_rates":[{"id":"x1","from_currency":"USD","to_currency":"CAD","rate":"1.350000","as_of":"2026-01-01"}],`+pageBlock+`}`)
	got, err := (&FXService{c: c}).Rates(context.Background(), FXRateListParams{ListParams: ListParams{PageSize: 10}})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/fx/rates")
	if q(rec, "page_size") != "10" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(got.Items) != 1 || got.Items[0].Rate != "1.350000" || got.Items[0].ToCurrency != "CAD" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntFXUpsertRate(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"id":"x1","from_currency":"USD","to_currency":"CAD","rate":"1.35","as_of":"2026-01-01"}`)
	got, err := (&FXService{c: c}).UpsertRate(context.Background(), FXRateRequest{FromCurrency: "USD", ToCurrency: "CAD", Rate: "1.35", AsOf: "2026-01-01"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/fx/rates")
	b := rec.json(t)
	if b["from_currency"] != "USD" || b["to_currency"] != "CAD" || b["rate"] != "1.35" || b["as_of"] != "2026-01-01" {
		t.Fatalf("body = %v", b)
	}
	if got.Rate != "1.35" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntFXConvert(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"amount":"135.00","rate":{"id":"x1","from_currency":"USD","to_currency":"CAD","rate":"1.35","as_of":"2026-01-01"}}`)
	got, err := (&FXService{c: c}).Convert(context.Background(), FXConvertParams{Amount: "100.00", From: "USD", To: "CAD", AsOf: "2026-01-01"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/fx/convert")
	if q(rec, "amount") != "100.00" || q(rec, "from") != "USD" || q(rec, "to") != "CAD" || q(rec, "as_of") != "2026-01-01" {
		t.Fatalf("query = %v", rec.Query)
	}
	if got.Amount != "135.00" || got.Rate.Rate != "1.35" || got.Rate.FromCurrency != "USD" {
		t.Fatalf("got %+v", got)
	}
}

const runJSON = `{"id":"cr1","as_of":"2026-01-31","status":"posted","journal_entry_id":"je1","created_by":"u1"}`

func TestFinEntConsolidationsList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"consolidations":[`+runJSON+`],`+pageBlock+`}`)
	got, err := (&ConsolidationsService{c: c}).List(context.Background(), ConsolidationListParams{Status: "posted"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/consolidations")
	if q(rec, "status") != "posted" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(got.Items) != 1 || got.Items[0].AsOf != "2026-01-31" || got.Items[0].JournalEntryID == nil || *got.Items[0].JournalEntryID != "je1" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntConsolidationsGet(t *testing.T) {
	c, rec := fakeERP(t, 200, runJSON)
	got, err := (&ConsolidationsService{c: c}).Get(context.Background(), "cr1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/consolidations/cr1")
	if got.Status != "posted" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntConsolidationsCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, runJSON)
	got, err := (&ConsolidationsService{c: c}).Create(context.Background(), ConsolidationRequest{AsOf: "2026-01-31"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/consolidations")
	if b := rec.json(t); b["as_of"] != "2026-01-31" || len(b) != 1 {
		t.Fatalf("body = %v", b)
	}
	if got.ID != "cr1" {
		t.Fatalf("got %+v", got)
	}
}

const subJSON = `{"id":"s1","customer_id":"cu1","amount":"29.99","cadence":"monthly","next_charge_on":"2026-02-01","status":"active"}`

func TestFinEntSubscriptionsList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"subscriptions":[`+subJSON+`],`+pageBlock+`}`)
	got, err := (&SubscriptionsService{c: c}).List(context.Background(), SubscriptionListParams{ListParams: ListParams{Q: "acme"}, Status: "active"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/subscriptions")
	if q(rec, "status") != "active" || q(rec, "q") != "acme" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(got.Items) != 1 || got.Items[0].Amount != "29.99" || got.Items[0].NextChargeOn != "2026-02-01" || got.Pagination.PageSize != 10 {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntSubscriptionsGet(t *testing.T) {
	c, rec := fakeERP(t, 200, subJSON)
	got, err := (&SubscriptionsService{c: c}).Get(context.Background(), "s1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/subscriptions/s1")
	if got.Cadence != "monthly" || got.CustomerID != "cu1" {
		t.Fatalf("got %+v", got)
	}
}

func TestFinEntSubscriptionsCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, subJSON)
	got, err := (&SubscriptionsService{c: c}).Create(context.Background(), SubscriptionRequest{
		CustomerID: "cu1", Amount: "29.99", Cadence: "monthly", NextChargeOn: "2026-02-01",
	}, WithIdempotencyKey("k2"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/subscriptions")
	if rec.Header.Get("Idempotency-Key") != "k2" {
		t.Fatalf("idempotency header = %q", rec.Header.Get("Idempotency-Key"))
	}
	b := rec.json(t)
	if b["customer_id"] != "cu1" || b["amount"] != "29.99" || b["cadence"] != "monthly" || b["next_charge_on"] != "2026-02-01" {
		t.Fatalf("body = %v", b)
	}
	if got.ID != "s1" {
		t.Fatalf("got %+v", got)
	}
}
