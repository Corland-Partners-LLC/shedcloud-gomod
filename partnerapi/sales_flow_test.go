package partnerapi

import (
	"context"
	"testing"
)

func TestStockUnits(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"stock_units":[{"id":"u1","product_id":"p1","serial_number":"SN-1","status":"available","origin":"built","unit_cost":"4000.00","asking_price":"6500.00","partner_keep_visible":false,"hide_from_embed":false,"created_at":"2026-08-01T00:00:00Z","updated_at":"2026-08-02T00:00:00Z","width_ft":10,"length_ft":12,"product_name":"Barn 10x12"}],"pagination":{"page":2,"page_size":25,"total":31}}`)
	page, err := c.StockUnits.List(ctx, StockUnitListParams{
		ListParams:      ListParams{Page: 2, PageSize: 25, Q: "barn"},
		Statuses:        []string{"available", "reserved"},
		IncludeReceived: true,
		Origin:          "built",
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/stock-units")
	if q(rec, "statuses") != "available,reserved" || q(rec, "include_received") != "true" || q(rec, "page_size") != "25" || q(rec, "q") != "barn" || q(rec, "origin") != "built" {
		t.Fatalf("query = %v", rec.Query)
	}
	if q(rec, "include_expired") != "" {
		t.Fatal("false flags must be omitted")
	}
	u := page.Items[0]
	if len(page.Items) != 1 || page.Pagination.Total != 31 || u.AskingPrice != "6500.00" || *u.WidthFt != 10 || u.ProductName != "Barn 10x12" {
		t.Fatalf("page = %+v", page)
	}

	if _, err := c.StockUnits.GetWithParams(ctx, "u1", StockUnitGetParams{IncludeReceived: true}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/stock-units/u1")
	if q(rec, "include_received") != "true" {
		t.Fatal("include_received missing")
	}

	if _, err := c.StockUnits.Create(ctx, StockUnitCreateRequest{ProductID: "p1", UnitCost: "4000.00", AskingPrice: "6500.00"}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/stock-units")
	if b := rec.json(t); b["product_id"] != "p1" || b["unit_cost"] != "4000.00" || b["serial_number"] != nil {
		t.Fatalf("body = %v", b)
	}

	if _, err := c.StockUnits.RegisterReturn(ctx, StockUnitRegisterReturnRequest{ProductID: "p1", ConditionGrade: "B"}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/stock-units/register-return")

	price := "5900.00"
	keep := true
	if _, err := c.StockUnits.Update(ctx, "u1", StockUnitUpdateRequest{AskingPrice: &price, PartnerKeepVisible: &keep}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/stock-units/u1")
	if b := rec.json(t); b["asking_price"] != "5900.00" || b["partner_keep_visible"] != true || len(b) != 2 {
		t.Fatalf("body = %v", b)
	}

	if _, err := c.StockUnits.Release(ctx, "u1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/stock-units/u1/release")
}

func TestStockUnitPDF(t *testing.T) {
	c, rec := fakeERP(t, 200, "%PDF-1.4 fake")
	pdf, err := c.StockUnits.PDF(context.Background(), "u1", StockUnitGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/stock-units/u1/pdf")
	if string(pdf) != "%PDF-1.4 fake" || rec.Header.Get("Accept") != "application/pdf" {
		t.Fatalf("pdf=%q accept=%q", pdf, rec.Header.Get("Accept"))
	}
}

func TestLeads(t *testing.T) {
	ctx := context.Background()
	lead := `{"id":"l1","name":"Sam River","status":"new","lifecycle_stage":"lead","lead_score":0,"priority":"normal","currency":"USD","estimated_value":"8000.00","created_at":"2026-08-01T00:00:00Z","updated_at":"2026-08-01T00:00:00Z","revision":3}`
	c, rec := fakeERP(t, 200, `{"leads":[`+lead+`],"pagination":{"page":1,"page_size":50,"total":1,"total_pages":1}}`)
	page, err := c.Leads.List(ctx, LeadListParams{ListParams: ListParams{Q: "sam"}, Status: "new"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/leads")
	if q(rec, "status") != "new" || q(rec, "q") != "sam" || page.Items[0].EstimatedValue != "8000.00" || page.Pagination.TotalPages != 1 {
		t.Fatalf("page = %+v query = %v", page, rec.Query)
	}

	c, rec = fakeERP(t, 201, lead)
	l, err := c.Leads.Create(ctx, LeadCreateRequest{Name: "Sam River", Email: "sam@example.com", EstimatedValue: "8000.00", SourceSlug: "website"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/leads")
	if b := rec.json(t); b["estimated_value"] != "8000.00" || b["source_slug"] != "website" || b["lead_score"] != nil {
		t.Fatalf("body = %v", b)
	}
	if l.ID != "l1" || l.Revision != 3 {
		t.Fatalf("lead = %+v", l)
	}

	if _, err := c.Leads.UpdateStatus(ctx, "l1", LeadStatusRequest{Status: "contacted", Note: "called"}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/leads/l1/status")
	if b := rec.json(t); b["status"] != "contacted" || b["note"] != "called" {
		t.Fatalf("body = %v", b)
	}

	if _, err := c.Leads.Assign(ctx, "l1", "u9"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/leads/l1/assign")
	if rec.json(t)["assigned_to"] != "u9" {
		t.Fatal("assigned_to missing")
	}

	phone := "+15550100"
	if _, err := c.Leads.Update(ctx, "l1", LeadUpdateRequest{Phone: &phone}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/leads/l1")

	c, rec = fakeERP(t, 200, `{"Lead":`+lead+`,"Customer":{"id":"c1","name":"Sam River","preferred_language":"en","status":"active","created_at":"2026-08-01T00:00:00Z","updated_at":"2026-08-01T00:00:00Z"},"Created":true}`)
	res, err := c.Leads.Convert(ctx, "l1", LeadConvertRequest{})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/leads/l1/convert")
	if res.Lead.ID != "l1" || res.Customer.ID != "c1" || !res.Created {
		t.Fatalf("convert = %+v (Go-cased keys must decode)", res)
	}

	c, rec = fakeERP(t, 200, `{"history":[{"id":"h1","lead_id":"l1","from_status":"new","to_status":"contacted","changed_at":"2026-08-02T00:00:00Z"}]}`)
	hist, err := c.Leads.StatusHistory(ctx, "l1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/leads/l1/status-history")
	if len(hist) != 1 || hist[0].ToStatus != "contacted" {
		t.Fatalf("history = %+v", hist)
	}
}

const quoteJSON = `{"id":"q1","number":"Q-100","status":"open","deal_type":"cash","fulfillment_method":"delivery","currency":"USD","base_price":"6500.00","total_upgrades":"0.00","subtotal":"6500.00","discount_amount":"0.00","delivery_fee":"0.00","building_anchored_fee":"0.00","material_surcharge_fee":"0.00","fuel_surcharge_fee":"0.00","total_fees":"0.00","total_pretax":"6500.00","tax_rate":"0.0825","tax_amount":"536.25","tax_exempt":false,"total":"7036.25","mileage_total":"0","mileage_free":"0","mileage_extra":"0","mileage_rate":"0","mileage_charge":"0.00","has_share_token":true,"created_at":"2026-08-01T00:00:00Z","updated_at":"2026-08-01T00:00:00Z","converted_order_id":"o1","lines":[{"id":"ln1","quote_id":"q1","line_type":"building","line_source":"stock","stock_unit_id":"u1","quantity":"1","base_price":"6500.00","upgrades_total":"0.00","removed_upgrades_total":"0.00","discount":"0.00","line_total":"6500.00","unavailable":false,"sort_order":0,"is_custom":false,"created_at":"2026-08-01T00:00:00Z","updated_at":"2026-08-01T00:00:00Z"}]}`

func TestQuotes(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 201, quoteJSON[:len(quoteJSON)-1]+`,"view_url":"https://app.example/quotes/view/tok","share_token":"tok"}`)
	created, err := c.Quotes.Create(ctx, QuoteCreateRequest{StoreLocationID: "loc1", DealType: "cash", BillToName: "Sam"}, WithIdempotencyKey("idem-q"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/quotes")
	if rec.Header.Get("Idempotency-Key") != "idem-q" {
		t.Fatal("idempotency key not sent")
	}
	if b := rec.json(t); b["store_location_id"] != "loc1" || b["bill_to_name"] != "Sam" {
		t.Fatalf("body = %v", b)
	}
	if created.ViewURL != "https://app.example/quotes/view/tok" || created.ShareToken != "tok" || created.ID != "q1" || created.Total != "7036.25" || !created.HasShareToken {
		t.Fatalf("created = %+v", created)
	}
	if len(created.Lines) != 1 || created.Lines[0].StockUnitID != "u1" || created.Lines[0].LineTotal != "6500.00" {
		t.Fatalf("lines = %+v", created.Lines)
	}

	c, rec = fakeERP(t, 200, `{"quotes":[`+quoteJSON+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	page, err := c.Quotes.List(ctx, QuoteListParams{ListParams: ListParams{PageSize: 10}, Status: "open", CustomerID: "c1", StoreID: "s1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/quotes")
	if q(rec, "status") != "open" || q(rec, "customer_id") != "c1" || q(rec, "store_id") != "s1" || len(page.Items) != 1 || page.Items[0].Number != "Q-100" {
		t.Fatalf("query=%v page=%+v", rec.Query, page)
	}

	c, rec = fakeERP(t, 200, quoteJSON)
	if _, err := c.Quotes.ReplaceLines(ctx, "q1", []LineInput{{StockUnitID: "u1", Quantity: "1"}}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PUT", "/partner/v1/quotes/q1/lines")
	lines, _ := rec.json(t)["lines"].([]any)
	if len(lines) != 1 || lines[0].(map[string]any)["stock_unit_id"] != "u1" {
		t.Fatalf("body = %s", rec.Body)
	}

	disc := "100.00"
	if _, err := c.Quotes.Update(ctx, "q1", QuoteUpdateRequest{DiscountAmount: &disc}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/quotes/q1")

	if _, err := c.Quotes.SetCustomer(ctx, "q1", CustomerLinkRequest{Name: "Sam", Email: "s@e.com"}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PUT", "/partner/v1/quotes/q1/customer")

	if _, err := c.Quotes.UpdateStatus(ctx, "q1", "lost", "price"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/quotes/q1/status")
	if b := rec.json(t); b["status"] != "lost" || b["reason"] != "price" {
		t.Fatalf("body = %v", b)
	}

	if _, err := c.Quotes.ComputeCashTerms(ctx, "q1", CashTermsComputeRequest{PromoCode: "SPRING"}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/quotes/q1/cash-terms/compute")
	if _, err := c.Quotes.SetRTOTerms(ctx, "q1", RTOTermsRequest{TermMonths: 36, MonthlyPayment: "120.00"}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PUT", "/partner/v1/quotes/q1/rto-terms")
	if _, err := c.Quotes.ComputeMileage(ctx, "q1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/quotes/q1/mileage")
	if _, err := c.Quotes.ResolveTax(ctx, "q1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/quotes/q1/tax/resolve")
	if _, err := c.Quotes.UpdateLineBuild(ctx, "q1", "ln1", LineBuildRequest{ConfigurationID: "cfg", DeriveUpgrades: true}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/quotes/q1/lines/ln1/build")
	if _, err := c.Quotes.Reprice(ctx, "q1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/quotes/q1/reprice")

	// Sell replaces Convert.
	sold, err := c.Quotes.Sell(ctx, "q1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/quotes/q1/sell")
	if sold.ConvertedOrderID != "o1" {
		t.Fatalf("sold = %+v", sold)
	}
}

func TestQuotePreviewsAndEvents(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"cash_purchase_total":"7036.25","suggested_down_payment":"700.00"}`)
	p, err := c.Quotes.PreviewCashTerms(ctx, "q1", CashTermsComputeRequest{})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/quotes/q1/cash-terms/preview")
	if p["suggested_down_payment"] != "700.00" {
		t.Fatalf("preview = %v", p)
	}

	c, rec = fakeERP(t, 200, `{"applicable":true,"stale":true,"can_apply":true,"reason":"price list changed","currency":"USD","document_base_price":"6500.00","book_base_price":"6800.00","document_total":"7036.25","book_total":"7361.00"}`)
	rp, err := c.Quotes.PreviewReprice(ctx, "q1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/quotes/q1/reprice/preview")
	if !rp.CanApply || rp.BookBasePrice != "6800.00" {
		t.Fatalf("reprice = %+v", rp)
	}

	c, rec = fakeERP(t, 200, `{"events":[{"id":"e1","quote_id":"q1","event_type":"created","created_at":"2026-08-01T00:00:00Z","detail":{"a":1}}]}`)
	ev, err := c.Quotes.Events(ctx, "q1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/quotes/q1/events")
	if len(ev) != 1 || ev[0].EventType != "created" {
		t.Fatalf("events = %+v", ev)
	}
}

func TestOrders(t *testing.T) {
	ctx := context.Background()
	orderJSON := `{"id":"o1","number":"SO-9","status":"open","deal_type":"cash","fulfillment_method":"delivery","currency":"USD","base_price":"6500.00","total":"7036.25","tax_exempt":false,"source_quote_id":"q1","placed_at":"2026-08-03T00:00:00Z","created_at":"2026-08-03T00:00:00Z","updated_at":"2026-08-03T00:00:00Z","sale_date":"2026-08-03T00:00:00Z","ledger_booked":true,"ledger_cancelled":false,"ledger_amount":"6500.00","signature_process_status":"none"}`
	c, rec := fakeERP(t, 200, `{"orders":[`+orderJSON+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	page, err := c.Orders.List(ctx, OrderListParams{Status: "open"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/orders")
	o := page.Items[0]
	if q(rec, "status") != "open" || o.SourceQuoteID != "q1" || !o.LedgerBooked || o.LedgerAmount != "6500.00" || o.SaleDate == nil {
		t.Fatalf("order = %+v", o)
	}

	c, rec = fakeERP(t, 200, orderJSON)
	if _, err := c.Orders.Get(ctx, "o1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/orders/o1")
	if _, err := c.Orders.UpdateStatus(ctx, "o1", "unprocessed"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/orders/o1/status")
	if rec.json(t)["status"] != "unprocessed" {
		t.Fatal("status missing")
	}
	if _, err := c.Orders.Create(ctx, OrderCreateRequest{StoreLocationID: "loc1"}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/orders")
	if _, err := c.Orders.SetCashTerms(ctx, "o1", CashTermsRequest{DownPayment: "700.00"}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PUT", "/partner/v1/orders/o1/cash-terms")
	if _, err := c.Orders.ReplaceLines(ctx, "o1", []LineInput{{ProductID: "p1"}}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PUT", "/partner/v1/orders/o1/lines")

	c, rec = fakeERP(t, 200, `{"events":[{"id":"e","order_id":"o1","event_type":"status","created_at":"2026-08-03T00:00:00Z"}]}`)
	if ev, err := c.Orders.Events(ctx, "o1"); err != nil || len(ev) != 1 || ev[0].OrderID != "o1" {
		t.Fatalf("events = %+v err=%v", ev, err)
	}
	expect(t, rec, "GET", "/partner/v1/orders/o1/events")

	c, rec = fakeERP(t, 200, `{"order_id":"o1","number":"SO-9","status":"open","signature_process_status":"completed","contract_pdf_status":"ready","contract_pdf_error":"","signature_process_started_at":null,"signature_process_completed_at":"2026-08-04T00:00:00Z"}`)
	ct, err := c.Orders.Contract(ctx, "o1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/orders/o1/contract")
	if ct.SignatureProcessStatus != "completed" || ct.SignatureProcessCompletedAt == nil || ct.SignatureProcessStartedAt != nil {
		t.Fatalf("contract = %+v", ct)
	}
}
