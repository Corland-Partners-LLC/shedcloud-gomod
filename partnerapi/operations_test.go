package partnerapi

import (
	"context"
	"strings"
	"testing"
)

func TestWorkOrders(t *testing.T) {
	ctx := context.Background()
	wo := `{"id":"w1","wo_number":"WO-5","config_category":"building","priority":1,"rush":false,"destination_is_customer":true,"qty_ordered":"1","qty_completed":"0","qty_scrapped":"0","status":"draft","structure_frozen":false,"has_shortage":false,"items_set_at":null,"bom_mode":"snapshot","is_custom":false,"created_at":"2026-08-01T00:00:00Z","updated_at":"2026-08-01T00:00:00Z","estimated_delivery_date":"2026-09-01"}`
	c, rec := fakeERP(t, 200, `{"work_orders":[`+wo+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	page, err := c.WorkOrders.List(ctx, WorkOrderListParams{Status: "draft", ProductID: "p1", CustomerID: "c1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/work-orders")
	if q(rec, "status") != "draft" || q(rec, "product_id") != "p1" || q(rec, "customer_id") != "c1" || page.Items[0].WONumber != "WO-5" || page.Items[0].EstimatedDeliveryDate != "2026-09-01" {
		t.Fatalf("page = %+v query = %v", page, rec.Query)
	}

	c, rec = fakeERP(t, 201, wo)
	if _, err := c.WorkOrders.Create(ctx, WorkOrderCreateRequest{ProductID: "p1", QtyOrdered: "1", Rush: true, SnapshotBOM: true}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/work-orders")
	if b := rec.json(t); b["qty_ordered"] != "1" || b["rush"] != true || b["snapshot_bom"] != true {
		t.Fatalf("body = %v", b)
	}

	cm := "rush job"
	if _, err := c.WorkOrders.Update(ctx, "w1", WorkOrderUpdateRequest{Comments: &cm}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/work-orders/w1")

	if _, err := c.WorkOrders.UpdateStatus(ctx, "w1", WorkOrderStatusRequest{Status: "released", Note: "go", CapturedFields: map[string]any{"x": 1}}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/work-orders/w1/status")
	if b := rec.json(t); b["status"] != "released" || b["captured_fields"] == nil {
		t.Fatalf("body = %v", b)
	}

	c, rec = fakeERP(t, 200, `{"history":[{"id":"h","work_order_id":"w1","to_status":"released","changed_at":"2026-08-02T00:00:00Z"}]}`)
	h, err := c.WorkOrders.StatusHistory(ctx, "w1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/work-orders/w1/status-history")
	if len(h) != 1 || h[0].ToStatus != "released" {
		t.Fatalf("history = %+v", h)
	}
}

func TestLocationBudgets(t *testing.T) {
	ctx := context.Background()
	b := `{"id":"b1","location_id":"loc1","location_name":"Austin","year":2026,"annual_budget":"1200000.00","bonus_threshold_1":"0","bonus_threshold_2":"0","inventory_cap":"0","close_rate_target":"0.25","months":[{"month":1,"weight_pct":"8.33","amount":"100000.00","weekly_amount":"23000.00","lead_target":40}],"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`
	c, rec := fakeERP(t, 200, `{"location_budgets":[`+b+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	page, err := c.LocationBudgets.List(ctx, LocationBudgetListParams{Year: 2026, LocationID: "loc1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/location-budgets")
	if q(rec, "year") != "2026" || q(rec, "location_id") != "loc1" || page.Items[0].Months[0].LeadTarget != 40 || page.Items[0].AnnualBudget != "1200000.00" {
		t.Fatalf("page = %+v", page)
	}

	c, rec = fakeERP(t, 201, b)
	req := LocationBudgetRequest{LocationID: "loc1", Year: 2026, AnnualBudget: "1200000.00", Months: []LocationBudgetMonth{{Month: 1, Amount: "100000.00", LeadTarget: 40}}}
	if _, err := c.LocationBudgets.Create(ctx, req); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/location-budgets")
	if _, err := c.LocationBudgets.Update(ctx, "b1", req); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PUT", "/partner/v1/location-budgets/b1")
	if _, err := c.LocationBudgets.Get(ctx, "b1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/location-budgets/b1")
}

func TestSalesLedgerReport(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"rows":[{"id":"m1","summary_id":"s1","order_id":"o1","order_number":"SO-9","movement_type":"sale","delta_base_price":"6500.00","delta_upgrades_amount":"0.00","delta_removed_upgrades_amount":"0.00","delta_amount":"6500.00","new_base_price":"6500.00","new_upgrades_amount":"0.00","new_removed_upgrades_amount":"0.00","new_order_amount":"6500.00","new_ledger_amount":"6500.00","cancelled":false,"sale_date":"2026-08-03T00:00:00Z","ledger_date":"2026-08-03T10:00:00Z","created_at":"2026-08-03T10:00:00Z","product_size":"10x12"}],"pagination":{"page":1,"page_size":50,"total_items":1},"sum_delta_amount":"6500.00"}`)
	rep, err := c.SalesLedger.Report(context.Background(), SalesLedgerRequest{SaleFrom: "2026-08-01", SaleTo: "2026-08-31", StoreID: "s1", MovementMode: "non_zero", IncludeSums: true})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/sales-ledger")
	b := rec.json(t)
	if b["sale_from"] != "2026-08-01" || b["store_id"] != "s1" || b["movement_mode"] != "non_zero" || b["include_sums"] != true || b["ledger_from"] != nil {
		t.Fatalf("body = %v", b)
	}
	if len(rep.Rows) != 1 || rep.Rows[0].DeltaAmount != "6500.00" || rep.Rows[0].ProductSize != "10x12" || rep.Pagination.TotalItems != 1 || rep.SumDeltaAmount != "6500.00" {
		t.Fatalf("report = %+v", rep)
	}
}

func TestRevenueForecast(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"forecast":[{"label":"This week","amount":12500.5},{"label":"Next week","amount":9000}]}`)
	f, err := c.RevenueForecast.Get(context.Background(), RevenueForecastParams{Year: 2026, LocationIDs: []string{"a", "b"}, WeekStart: "2026-08-03", Lang: "en"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/revenue-forecast")
	if q(rec, "year") != "2026" || q(rec, "locationIds") != "a,b" || q(rec, "weekStart") != "2026-08-03" || q(rec, "lang") != "en" {
		t.Fatalf("query = %v (the ERP reads camelCase here)", rec.Query)
	}
	if len(f.Forecast) != 2 || f.Forecast[0].Amount != 12500.5 || f.Forecast[1].Label != "Next week" {
		t.Fatalf("forecast = %+v", f)
	}
}

func TestSiteEvents(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 202, `{"success":true,"accepted":2}`)
	res, err := c.SiteEvents.Track(ctx, SiteEventsTrackRequest{
		SessionID: "sess-1", VisitorID: "vis-1", SiteHost: "example-sheds.com", ClientIP: "203.0.113.9",
		UTM: map[string]string{"source": "google"},
		Events: []SiteEventInput{
			{EventType: "page.view", Page: "/barns", ProductID: "p1"},
			{EventType: "customer.profile", Customer: &SiteEventCustomer{FirstName: "Sam", Email: "sam@example.com"}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/site-events")
	body := rec.json(t)
	for _, k := range []string{"session_id", "visitor_id", "site_host", "client_ip", "events", "utm"} {
		if body[k] == nil {
			t.Fatalf("body lacks snake_case key %q: %s", k, rec.Body)
		}
	}
	ev := body["events"].([]any)[0].(map[string]any)
	if ev["event_type"] != "page.view" || ev["product_id"] != "p1" {
		t.Fatalf("event = %v", ev)
	}
	if strings.Contains(string(rec.Body), "sessionId") {
		t.Fatal("Track must stay snake_case")
	}
	if !res.Success || res.Accepted != 2 {
		t.Fatalf("res = %+v", res)
	}

	c, rec = fakeERP(t, 200, `{"data":[{"eventId":"e1","eventType":"page.view","occurredAt":"2026-08-01T00:00:00Z","sessionId":"sess-1","source":"marketing","geoCity":"Austin"}],"hasMore":false}`)
	list, err := c.SiteEvents.List(ctx, SiteEventListParams{Limit: 50, SessionID: "sess-1", Types: []string{"page.view", "quote.start"}, From: "2026-08-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/site-events")
	if q(rec, "sessionId") != "sess-1" || q(rec, "types") != "page.view,quote.start" || q(rec, "limit") != "50" || q(rec, "from") == "" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(list.Data) != 1 || list.Data[0].EventID != "e1" || list.Data[0].GeoCity != "Austin" || list.HasMore {
		t.Fatalf("list = %+v", list)
	}
}

func TestCustomersAndPayments(t *testing.T) {
	ctx := context.Background()
	cu := `{"id":"c1","name":"Sam River","email":"sam@example.com","preferred_language":"en","status":"active","created_at":"2026-08-01T00:00:00Z","updated_at":"2026-08-01T00:00:00Z"}`
	c, rec := fakeERP(t, 200, `{"customers":[`+cu+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	page, err := c.Customers.List(ctx, CustomerListParams{ListParams: ListParams{Q: "sam"}})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/customers")
	if q(rec, "q") != "sam" || page.Items[0].Email != "sam@example.com" {
		t.Fatalf("page = %+v", page)
	}
	if _, err := c.Customers.Get(ctx, "c1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/customers/c1")
	ph := "+15550199"
	if _, err := c.Customers.Update(ctx, "c1", CustomerUpdateRequest{Phone: &ph}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/customers/c1")

	pay := `{"id":"pm1","customer_id":"c1","payment_number":"P-1","payment_date":"2026-08-19","amount":"500.00","allocated_amount":"0.00","unapplied_amount":"500.00","refunded_amount":"0.00","method":"check","status":"posted","reference_type":"order","reference_id":"o1","created_at":"2026-08-19T00:00:00Z","updated_at":"2026-08-19T00:00:00Z"}`
	c, rec = fakeERP(t, 201, pay)
	p, err := c.Customers.CreatePayment(ctx, "c1", PaymentCreateRequest{Amount: "500.00", Method: "check", PaymentDate: "2026-08-19", Reference: "CHK-1042", ReferenceType: "order", ReferenceID: "o1"}, WithIdempotencyKey("pay-1"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/customers/c1/payments")
	if b := rec.json(t); b["amount"] != "500.00" || b["reference_type"] != "order" || rec.Header.Get("Idempotency-Key") != "pay-1" {
		t.Fatalf("body = %v", b)
	}
	if p.Amount != "500.00" || p.UnappliedAmount != "500.00" {
		t.Fatalf("payment = %+v", p)
	}

	c, rec = fakeERP(t, 200, `{"payments":[`+pay+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	pp, err := c.Payments.List(ctx, PaymentListParams{ReferenceType: "order", ReferenceID: "o1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/payments")
	if q(rec, "reference_type") != "order" || q(rec, "reference_id") != "o1" || len(pp.Items) != 1 {
		t.Fatalf("query = %v", rec.Query)
	}
	if _, err := c.Customers.Payments(ctx, "c1", ListParams{}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/customers/c1/payments")
	if _, err := c.Payments.Get(ctx, "pm1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/payments/pm1")

	c, rec = fakeERP(t, 200, `{"customer_id":"c1","total_charged":"7036.25","total_paid":"500.00","total_allocated":"0.00","credit_balance":"0.00","pending_refunds":"0.00","open_balance":"6536.25","net_balance":"6536.25"}`)
	l, err := c.Customers.Ledger(ctx, "c1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/customers/c1/ledger")
	if l.OpenBalance != "6536.25" || l.NetBalance != "6536.25" {
		t.Fatalf("ledger = %+v", l)
	}
}

func TestLocationsProductsAgreements(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"locations":[{"id":"loc1","code":"AUS","name":"Austin","location_type":"store","phone":"555","is_active":true,"consider_ledger":true,"is_stock_location":true,"is_visible":true,"related_count":0,"lot_stock_limit":null,"rental_return_limit":"3","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}],"pagination":{"page":1,"page_size":50,"total":1}}`)
	lp, err := c.Locations.List(ctx, LocationListParams{ListParams: ListParams{Q: "aus"}})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/locations")
	if lp.Items[0].Code != "AUS" || lp.Items[0].LotStockLimit != nil || *lp.Items[0].RentalReturnLimit != "3" {
		t.Fatalf("locations = %+v", lp.Items)
	}
	active := false
	if _, err := c.Locations.Create(ctx, LocationCreateRequest{Code: "DAL", Name: "Dallas", LocationType: "store", IsActive: &active}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/locations")
	if rec.json(t)["is_active"] != false {
		t.Fatal("is_active=false must be sent")
	}
	nm := "Dallas 2"
	if _, err := c.Locations.Update(ctx, "loc1", LocationUpdateRequest{Name: &nm}); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/locations/loc1")

	c, rec = fakeERP(t, 200, `{"storefronts":[{"id":"sf1","apex_domain":"example-sheds.com","subdomain":"www","draft":false,"sort_order":0,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}]}`)
	d, err := c.Locations.Domains(ctx, "loc1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/locations/loc1/domains")
	if len(d) != 1 || d[0].ApexDomain != "example-sheds.com" {
		t.Fatalf("domains = %+v", d)
	}

	c, rec = fakeERP(t, 200, `{"products":[{"id":"p1","sku":"BRN-1012","name":"Barn 10x12","product_type":"finished_good","uom_code":"EA","tracking":"serial","is_stocked":true,"hide_sku":false,"hide_on_configurator":false,"standard_cost":"4000.00","list_price":"6500.00","is_active":true,"is_system":false,"bom_mode":"snapshot","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}],"pagination":{"page":1,"page_size":50,"total":1}}`)
	pp, err := c.Products.List(ctx, ProductListParams{ListParams: ListParams{Q: "barn"}})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/products")
	if pp.Items[0].ListPrice != "6500.00" || pp.Items[0].SKU != "BRN-1012" {
		t.Fatalf("products = %+v", pp.Items)
	}
	if _, err := c.Products.Get(ctx, "p1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/products/p1")

	c, rec = fakeERP(t, 200, `{"price":"6500.00","source":"item"}`)
	pr, err := c.Products.ResolvePrice(ctx, PricingResolveParams{ProductID: "p1", LocationID: "loc1", WidthFt: "10", LengthFt: "12"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/pricing/resolve")
	if q(rec, "product_id") != "p1" || q(rec, "width_ft") != "10" || pr["source"] != "item" {
		t.Fatalf("query=%v pr=%v", rec.Query, pr)
	}

	c, rec = fakeERP(t, 200, `{"agreements":[{"id":"ag1","counterparty_tenant_id":"t2","type":"rto","title":"Program","status":"active","rate_config_version":2,"pending_rate_config_version":0,"is_inhouse":false,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}],"pagination":{"total":1,"page":1,"page_size":25}}`)
	ag, err := c.Agreements.Active(ctx, AgreementListParams{Type: "rto"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/agreements/active")
	if q(rec, "type") != "rto" || ag.Items[0].Status != "active" {
		t.Fatalf("agreements = %+v", ag.Items)
	}
	c, rec = fakeERP(t, 200, `{"state_legal":[{"id":"sl1","agreement_id":"ag1","state_code":"TX","include_customer_data_sheet":true,"include_state_rpa_pages":false,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}]}`)
	sl, err := c.Agreements.ListStateLegal(ctx, "ag1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/agreements/ag1/state-legal")
	if len(sl) != 1 || sl[0].StateCode != "TX" {
		t.Fatalf("state legal = %+v", sl)
	}
	if _, err := c.Agreements.GetStateLegal(ctx, "ag1", "TX"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/agreements/ag1/state-legal/TX")
}

func TestUsersDocumentsEventsConfigurator(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"members":[{"user":{"id":"u1","display_name":"Pat","email":"pat@example.com","avatar_url":"","active":true},"membership":{"id":"m1","tenant_id":"t","user_id":"u1","status":"active","is_default":true,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"},"roles":[{"id":"r1","slug":"sales","name":"Sales","is_system":false,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}],"status":"active","access_blocked":false,"has_permission_extras":false}],"pagination":{"page":1,"page_size":50,"total":1,"total_pages":1}}`)
	us, err := c.Users.List(ctx, UserListParams{Role: "sales", Status: "active"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/users")
	if q(rec, "role") != "sales" || us.Items[0].User.Email != "pat@example.com" || us.Items[0].Roles[0].Slug != "sales" {
		t.Fatalf("users = %+v", us.Items)
	}
	c, rec = fakeERP(t, 201, `{"id":"inv1","tenant_id":"t","email":"new@example.com","role_id":"r1","status":"pending","created_at":"2026-08-01T00:00:00Z"}`)
	send := false
	inv, err := c.Users.Invite(ctx, UserInviteRequest{Email: "new@example.com", Role: "sales", SendEmail: &send})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/users")
	if rec.json(t)["send_email"] != false || inv.Status != "pending" {
		t.Fatalf("body = %s inv = %+v", rec.Body, inv)
	}
	c, rec = fakeERP(t, 200, `{"roles":[{"id":"r1","slug":"sales","name":"Sales","is_system":false,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z","permissions":[]}]}`)
	roles, err := c.Users.Roles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/roles")
	if len(roles) != 1 || roles[0].Slug != "sales" {
		t.Fatalf("roles = %+v", roles)
	}

	c, rec = fakeERP(t, 200, `{"storefronts":[{"id":"sf1","apex_domain":"x.com","subdomain":"","draft":true,"sort_order":1,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}]}`)
	if sf, err := c.Domains.List(ctx); err != nil || len(sf) != 1 || !sf[0].Draft {
		t.Fatalf("domains = %+v err=%v", sf, err)
	}
	expect(t, rec, "GET", "/partner/v1/domains")

	c, rec = fakeERP(t, 200, `{"documents":[{"id":"f1","link_id":"l1","entity_type":"order","entity_id":"o1","role":"contract","position":0,"is_primary":true,"file_name":"contract.pdf","created_at":"2026-08-01T00:00:00Z"}]}`)
	docs, err := c.Documents.List(ctx, DocumentFilesParams{EntityType: "order", EntityID: "o1", Role: "contract"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/documents")
	if q(rec, "entity_type") != "order" || q(rec, "entity_id") != "o1" || docs[0].FileName != "contract.pdf" {
		t.Fatalf("docs = %+v query = %v", docs, rec.Query)
	}
	c, rec = fakeERP(t, 200, `{"url":"https://files.example/x?sig=1","method":"GET","expires_at":"2026-08-01T01:00:00Z"}`)
	dl, err := c.Documents.Download(ctx, "f1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/documents/f1/download")
	if dl.URL == "" || dl.Method != "GET" {
		t.Fatalf("download = %+v", dl)
	}

	c, rec = fakeERP(t, 200, `{"events":[{"id":"ev1","entity_type":"quote","entity_id":"q1","action":"sold","created_at":"2026-08-01T00:00:00Z"}]}`)
	evs, err := c.Events.List(ctx, EventListParams{Limit: 10, After: "2026-08-01T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/events")
	if q(rec, "limit") != "10" || evs[0].Action != "sold" {
		t.Fatalf("events = %+v", evs)
	}
	c, rec = fakeERP(t, 200, `{"deliveries":[{"id":"d1","subscription_id":"s1","event_id":"ev1","status":"dead","attempt_count":10,"response_code":410,"created_at":"2026-08-01T00:00:00Z"}],"catalog":["quote.sold"]}`)
	dv, err := c.Events.Deliveries(ctx, WebhookDeliveryListParams{Status: "dead"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/webhook-deliveries")
	if q(rec, "status") != "dead" || dv.Deliveries[0].AttemptCount != 10 || *dv.Deliveries[0].ResponseCode != 410 {
		t.Fatalf("deliveries = %+v", dv)
	}
	c, rec = fakeERP(t, 200, `{"event_id":"ev1","enqueued":2,"event_key":"quote.sold"}`)
	rd, err := c.Events.Redeliver(ctx, "ev1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/events/ev1/redeliver")
	if rd.Enqueued != 2 || rd.EventKey != "quote.sold" {
		t.Fatalf("redeliver = %+v", rd)
	}

	c, rec = fakeERP(t, 201, `{"url":"https://threed.example/?company_id=t&product_id=p1","product_id":"p1","company_id":"t"}`)
	cs, err := c.ConfiguratorSessions.Create(ctx, ConfiguratorSessionCreateRequest{ProductID: "p1", ReturnURL: "https://dealer.example/return"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/configurator-sessions")
	if rec.json(t)["product_id"] != "p1" || cs.CompanyID != "t" || !strings.Contains(cs.URL, "product_id=p1") {
		t.Fatalf("session = %+v body=%s", cs, rec.Body)
	}
}
