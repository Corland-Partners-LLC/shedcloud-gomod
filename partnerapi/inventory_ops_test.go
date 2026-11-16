package partnerapi

import (
	"context"
	"testing"
)

func TestInventorySettings(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"tenant_id":"t1","partner_sold_visible_days":30,"partner_sold_visible_mode":"hide_after_days","inventory_limit_mode":"warn","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-02T00:00:00Z"}`)
	svc := &InventoryService{c: c}

	got, err := svc.GetSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/inventory-settings")
	if got.PartnerSoldVisibleDays != 30 || got.PartnerSoldVisibleMode != "hide_after_days" || got.InventoryLimitMode != "warn" {
		t.Fatalf("settings = %+v", got)
	}

	days, mode := 45, "always_show_sold"
	got, err = svc.UpdateSettings(context.Background(), InventorySettingsUpdateRequest{PartnerSoldVisibleDays: &days, PartnerSoldVisibleMode: &mode})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PUT", "/partner/v1/inventory-settings")
	body := rec.json(t)
	if body["partner_sold_visible_days"] != float64(45) || body["partner_sold_visible_mode"] != "always_show_sold" || len(body) != 2 {
		t.Fatalf("body = %v", body)
	}
	if got.TenantID != "t1" {
		t.Fatalf("settings = %+v", got)
	}

	// Only non-nil fields are sent.
	if _, err := svc.UpdateSettings(context.Background(), InventorySettingsUpdateRequest{PartnerSoldVisibleDays: &days}); err != nil {
		t.Fatal(err)
	}
	if body := rec.json(t); len(body) != 1 || body["partner_sold_visible_days"] != float64(45) {
		t.Fatalf("partial body = %v", body)
	}
}

func TestInventoryAdjustments(t *testing.T) {
	const detail = `{"adjustment":{"id":"a1","adjustment_no":"ADJ-1","location_id":"l1","notes":"n","line_count":1,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"},"lines":[{"id":"ln1","product_id":"p1","quantity":"-2.0000","qty_before":"10.0000","qty_after":"8.0000","unit_cost":"4.5000","movement_type":"adjustment","occurred_at":"2026-01-01T00:00:00Z"}]}`

	c, rec := fakeERP(t, 200, `{"adjustments":[{"id":"a1","adjustment_no":"ADJ-1","location_id":"l1","location_code":"MAIN","line_count":2}],"pagination":{"page":2,"page_size":10,"total":11}}`)
	svc := &InventoryService{c: c}
	page, err := svc.ListAdjustments(context.Background(), InventoryAdjustmentListParams{
		ListParams: ListParams{Page: 2, PageSize: 10, Q: "adj"}, LocationID: "l1", From: "2026-01-01", To: "2026-01-31",
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/inventory/adjustments")
	if q(rec, "page") != "2" || q(rec, "page_size") != "10" || q(rec, "q") != "adj" || q(rec, "location_id") != "l1" || q(rec, "from") != "2026-01-01" || q(rec, "to") != "2026-01-31" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(page.Items) != 1 || page.Items[0].AdjustmentNo != "ADJ-1" || page.Items[0].LocationCode != "MAIN" || page.Items[0].LineCount != 2 || page.Pagination.Total != 11 {
		t.Fatalf("page = %+v", page)
	}

	c, rec = fakeERP(t, 201, detail)
	svc = &InventoryService{c: c}
	created, err := svc.CreateAdjustment(context.Background(), InventoryAdjustmentCreateRequest{
		LocationID: "l1", Notes: "recount",
		Lines: []InventoryAdjustmentLineRequest{{ProductID: "p1", Quantity: "-2", UnitCost: "4.5", Notes: "damaged"}},
	}, WithIdempotencyKey("idem-1"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/inventory/adjustments")
	if rec.Header.Get("Idempotency-Key") != "idem-1" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	lines, _ := body["lines"].([]any)
	if body["location_id"] != "l1" || body["notes"] != "recount" || len(lines) != 1 {
		t.Fatalf("body = %v", body)
	}
	ln, _ := lines[0].(map[string]any)
	if ln["product_id"] != "p1" || ln["quantity"] != "-2" || ln["unit_cost"] != "4.5" || ln["notes"] != "damaged" {
		t.Fatalf("line = %v", ln)
	}
	if created.Adjustment.ID != "a1" || len(created.Lines) != 1 || created.Lines[0].QtyAfter != "8.0000" || created.Lines[0].MovementType != "adjustment" {
		t.Fatalf("created = %+v", created)
	}

	c, rec = fakeERP(t, 200, detail)
	svc = &InventoryService{c: c}
	got, err := svc.GetAdjustment(context.Background(), "a1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/inventory/adjustments/a1")
	if got.Adjustment.AdjustmentNo != "ADJ-1" || got.Lines[0].Quantity != "-2.0000" {
		t.Fatalf("got = %+v", got)
	}

	c, rec = fakeERP(t, 200, `{"rows":[{"adjustment_id":"a1","adjustment_no":"ADJ-1","movement_id":"m1","product_sku":"SKU1","quantity":"-2.0000","extended_cost":"-9.0000","document_notes":"d"}],"pagination":{"page":1,"page_size":25,"total":1}}`)
	svc = &InventoryService{c: c}
	rep, err := svc.AdjustmentsReport(context.Background(), InventoryAdjustmentReportParams{
		ListParams: ListParams{Q: "sku"}, LocationIDs: []string{"l1", "l2"}, From: "2026-01-01",
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/inventory/adjustments/report")
	if q(rec, "location_id") != "l1,l2" || q(rec, "from") != "2026-01-01" || q(rec, "q") != "sku" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(rep.Items) != 1 || rep.Items[0].ProductSKU != "SKU1" || rep.Items[0].ExtendedCost != "-9.0000" || rep.Items[0].DocumentNotes != "d" || rep.Pagination.Total != 1 {
		t.Fatalf("report = %+v", rep)
	}
}

func TestInventoryTransfers(t *testing.T) {
	const detail = `{"transfer":{"id":"t1","transfer_no":"TR-1","from_location_id":"l1","to_location_id":"l2","from_location_code":"A","to_location_name":"Lot B","line_count":1,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"},"lines":[{"id":"ln1","product_id":"p1","quantity":"3.0000","unit_cost":"2.0000","extended_cost":"6.0000","from_qty_before":"10.0000","from_qty_after":"7.0000","to_qty_before":"0.0000","to_qty_after":"3.0000","occurred_at":"2026-01-01T00:00:00Z"}]}`

	c, rec := fakeERP(t, 200, `{"transfers":[{"id":"t1","transfer_no":"TR-1","from_location_id":"l1","to_location_id":"l2","line_count":1}],"pagination":{"page":1,"page_size":25,"total":1}}`)
	svc := &InventoryService{c: c}
	page, err := svc.ListTransfers(context.Background(), InventoryTransferListParams{
		ListParams: ListParams{Page: 1, Q: "tr"}, FromLocationID: "l1", ToLocationID: "l2", From: "2026-01-01", To: "2026-01-31",
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/inventory/transfers")
	if q(rec, "from_location_id") != "l1" || q(rec, "to_location_id") != "l2" || q(rec, "from") != "2026-01-01" || q(rec, "to") != "2026-01-31" || q(rec, "q") != "tr" || q(rec, "page") != "1" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(page.Items) != 1 || page.Items[0].TransferNo != "TR-1" || page.Items[0].ToLocationID != "l2" || page.Pagination.Total != 1 {
		t.Fatalf("page = %+v", page)
	}

	c, rec = fakeERP(t, 201, detail)
	svc = &InventoryService{c: c}
	created, err := svc.CreateTransfer(context.Background(), InventoryTransferCreateRequest{
		FromLocationID: "l1", ToLocationID: "l2", Notes: "rebalance",
		Lines: []InventoryTransferLineRequest{{ProductID: "p1", Quantity: "3", Notes: "x"}},
	}, WithIdempotencyKey("idem-2"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/inventory/transfers")
	if rec.Header.Get("Idempotency-Key") != "idem-2" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	lines, _ := body["lines"].([]any)
	if body["from_location_id"] != "l1" || body["to_location_id"] != "l2" || body["notes"] != "rebalance" || len(lines) != 1 {
		t.Fatalf("body = %v", body)
	}
	ln, _ := lines[0].(map[string]any)
	if ln["product_id"] != "p1" || ln["quantity"] != "3" || ln["notes"] != "x" {
		t.Fatalf("line = %v", ln)
	}
	if created.Transfer.ID != "t1" || created.Transfer.FromLocationCode != "A" || len(created.Lines) != 1 || created.Lines[0].ToQtyAfter != "3.0000" {
		t.Fatalf("created = %+v", created)
	}

	c, rec = fakeERP(t, 200, detail)
	svc = &InventoryService{c: c}
	got, err := svc.GetTransfer(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/inventory/transfers/t1")
	if got.Transfer.TransferNo != "TR-1" || got.Lines[0].ExtendedCost != "6.0000" {
		t.Fatalf("got = %+v", got)
	}

	c, rec = fakeERP(t, 200, `{"rows":[{"transfer_id":"t1","transfer_no":"TR-1","movement_id":"m1","from_location_code":"A","to_location_code":"B","quantity":"3.0000","extended_cost":"6.0000"}],"pagination":{"page":1,"page_size":25,"total":1}}`)
	svc = &InventoryService{c: c}
	rep, err := svc.TransfersReport(context.Background(), InventoryTransferReportParams{LocationIDs: []string{"l1"}, To: "2026-02-01"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/inventory/transfers/report")
	if q(rec, "location_id") != "l1" || q(rec, "to") != "2026-02-01" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(rep.Items) != 1 || rep.Items[0].ToLocationCode != "B" || rep.Items[0].ExtendedCost != "6.0000" || rep.Pagination.Total != 1 {
		t.Fatalf("report = %+v", rep)
	}
}

func TestStockSaleRequests(t *testing.T) {
	const one = `{"id":"r1","stock_unit_id":"u1","requester_id":"usr1","assignee_id":"usr2","message":"sell it","status":"pending","serial_number":"SN1","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`

	c, rec := fakeERP(t, 201, one)
	svc := &StockSaleRequestsService{c: c}
	created, err := svc.Create(context.Background(), "u1", StockSaleRequestCreateRequest{AssigneeUserID: "usr2", Message: "sell it"}, WithIdempotencyKey("idem-3"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/stock-units/u1/sale-requests")
	if rec.Header.Get("Idempotency-Key") != "idem-3" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	if body := rec.json(t); body["assignee_user_id"] != "usr2" || body["message"] != "sell it" || len(body) != 2 {
		t.Fatalf("body = %v", body)
	}
	if created.ID != "r1" || created.StockUnitID != "u1" || created.Status != "pending" || created.SerialNumber != "SN1" {
		t.Fatalf("created = %+v", created)
	}

	c, rec = fakeERP(t, 200, `{"sale_requests":[`+one+`]}`)
	svc = &StockSaleRequestsService{c: c}
	list, err := svc.List(context.Background(), StockSaleRequestListParams{Status: "pending", StockUnitID: "u1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/stock-sale-requests")
	if q(rec, "status") != "pending" || q(rec, "stock_unit_id") != "u1" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(list) != 1 || list[0].AssigneeID != "usr2" {
		t.Fatalf("list = %+v", list)
	}

	c, rec = fakeERP(t, 200, `{"id":"r1","stock_unit_id":"u1","requester_id":"usr1","assignee_id":"usr2","status":"approved","reason":"ok","resolved_at":"2026-01-02T00:00:00Z","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-02T00:00:00Z"}`)
	svc = &StockSaleRequestsService{c: c}
	res, err := svc.Resolve(context.Background(), "r1", StockSaleRequestResolveRequest{Decision: "approve", Reason: "ok"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/stock-sale-requests/r1/resolve")
	if body := rec.json(t); body["decision"] != "approve" || body["reason"] != "ok" || len(body) != 2 {
		t.Fatalf("body = %v", body)
	}
	if res.Status != "approved" || res.Reason != "ok" || res.ResolvedAt == nil {
		t.Fatalf("resolved = %+v", res)
	}
}

func TestSalesReturns(t *testing.T) {
	const header = `{"id":"sr1","number":"RMA-1","status":"draft","order_id":"o1","customer_id":"c1","restock_total":"100.00","credit_total":"80.00","order_number":"ORD-9","line_count":1,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`
	const line = `{"id":"srl1","return_id":"sr1","order_line_id":"ol1","stock_unit_id":"u1","disposition":"restock","qty":"1.0000","amount":"80.00","created_at":"2026-01-01T00:00:00Z"}`

	c, rec := fakeERP(t, 200, `{"sales_returns":[`+header+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	svc := &SalesReturnsService{c: c}
	page, err := svc.List(context.Background(), SalesReturnListParams{ListParams: ListParams{Q: "rma"}, Status: "draft", OrderID: "o1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/sales-returns")
	if q(rec, "q") != "rma" || q(rec, "status") != "draft" || q(rec, "order_id") != "o1" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(page.Items) != 1 || page.Items[0].Number != "RMA-1" || page.Items[0].CreditTotal != "80.00" || page.Items[0].OrderNumber != "ORD-9" || page.Pagination.Total != 1 {
		t.Fatalf("page = %+v", page)
	}

	c, rec = fakeERP(t, 201, `{"sales_return":`+header+`,"lines":[`+line+`]}`)
	svc = &SalesReturnsService{c: c}
	created, err := svc.Create(context.Background(), SalesReturnCreateRequest{
		OrderID: "o1", CustomerID: "c1", StoreLocationID: "l1", Reason: "defect", Notes: "n",
		Lines: []SalesReturnLineRequest{{OrderLineID: "ol1", StockUnitID: "u1", ProductID: "p1", Disposition: "restock", Qty: "1", Amount: "80.00", Reason: "r"}},
	}, WithIdempotencyKey("idem-4"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/sales-returns")
	if rec.Header.Get("Idempotency-Key") != "idem-4" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	lines, _ := body["lines"].([]any)
	if body["order_id"] != "o1" || body["customer_id"] != "c1" || body["store_location_id"] != "l1" || body["reason"] != "defect" || body["notes"] != "n" || len(lines) != 1 {
		t.Fatalf("body = %v", body)
	}
	ln, _ := lines[0].(map[string]any)
	if ln["order_line_id"] != "ol1" || ln["stock_unit_id"] != "u1" || ln["product_id"] != "p1" || ln["disposition"] != "restock" || ln["qty"] != "1" || ln["amount"] != "80.00" || ln["reason"] != "r" {
		t.Fatalf("line = %v", ln)
	}
	if created.SalesReturn.ID != "sr1" || len(created.Lines) != 1 || created.Lines[0].Disposition != "restock" {
		t.Fatalf("created = %+v", created)
	}

	c, rec = fakeERP(t, 200, `{"id":"sr1","number":"RMA-1","status":"draft","order_id":"o1","lines":[`+line+`],"sales_return":`+header+`}`)
	svc = &SalesReturnsService{c: c}
	got, err := svc.Get(context.Background(), "sr1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/sales-returns/sr1")
	if got.ID != "sr1" || got.Number != "RMA-1" || got.OrderID != "o1" || got.SalesReturn.RestockTotal != "100.00" || len(got.Lines) != 1 || got.Lines[0].Amount != "80.00" {
		t.Fatalf("got = %+v", got)
	}
}

func TestStockReconciliations(t *testing.T) {
	const sess = `{"id":"s1","location_id":"l1","all_locations":false,"status":"active","started_by":"usr1","reviewer_ids":["usr2"],"item_count":12,"decided_count":3,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`

	c, rec := fakeERP(t, 200, `{"sessions":[`+sess+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	svc := &StockReconciliationsService{c: c}
	page, err := svc.List(context.Background(), StockReconciliationListParams{ListParams: ListParams{PageSize: 50}, Status: "active", LocationID: "l1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/stock-reconciliations")
	if q(rec, "status") != "active" || q(rec, "location_id") != "l1" || q(rec, "page_size") != "50" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(page.Items) != 1 || page.Items[0].ItemCount != 12 || page.Items[0].ReviewerIDs[0] != "usr2" || page.Pagination.Total != 1 {
		t.Fatalf("page = %+v", page)
	}

	c, rec = fakeERP(t, 201, sess)
	svc = &StockReconciliationsService{c: c}
	started, err := svc.Start(context.Background(), StockReconciliationStartRequest{LocationID: "l1", AllLocations: true, ReviewerIDs: []string{"usr2"}}, WithIdempotencyKey("idem-5"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/stock-reconciliations")
	if rec.Header.Get("Idempotency-Key") != "idem-5" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	ids, _ := body["reviewer_ids"].([]any)
	if body["location_id"] != "l1" || body["all_locations"] != true || len(ids) != 1 || ids[0] != "usr2" {
		t.Fatalf("body = %v", body)
	}
	if started.ID != "s1" || started.Status != "active" || started.DecidedCount != 3 {
		t.Fatalf("started = %+v", started)
	}

	c, rec = fakeERP(t, 200, sess)
	svc = &StockReconciliationsService{c: c}
	got, err := svc.Get(context.Background(), "s1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/stock-reconciliations/s1")
	if got.ID != "s1" || got.LocationID != "l1" || got.StartedBy != "usr1" {
		t.Fatalf("got = %+v", got)
	}
}

func TestCycleCounts(t *testing.T) {
	const sess = `{"id":"cc1","location_id":"l1","status":"open","started_by":"usr1","line_count":1,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`
	const detail = `{"session":` + sess + `,"lines":[{"id":"ccl1","session_id":"cc1","product_id":"p1","system_qty":"10.0000","counted_qty":"9.0000","variance_qty":"-1.0000","created_at":"2026-01-01T00:00:00Z"}]}`

	c, rec := fakeERP(t, 200, `{"cycle_counts":[`+sess+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	svc := &CycleCountsService{c: c}
	page, err := svc.List(context.Background(), CycleCountListParams{LocationID: "l1", Status: "open"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/cycle-counts")
	if q(rec, "location_id") != "l1" || q(rec, "status") != "open" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(page.Items) != 1 || page.Items[0].ID != "cc1" || page.Items[0].LineCount != 1 || page.Pagination.Total != 1 {
		t.Fatalf("page = %+v", page)
	}

	c, rec = fakeERP(t, 201, detail)
	svc = &CycleCountsService{c: c}
	started, err := svc.Start(context.Background(), CycleCountStartRequest{LocationID: "l1"}, WithIdempotencyKey("idem-6"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/cycle-counts")
	if rec.Header.Get("Idempotency-Key") != "idem-6" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	if body := rec.json(t); body["location_id"] != "l1" || len(body) != 1 {
		t.Fatalf("body = %v", body)
	}
	if started.Session.ID != "cc1" || len(started.Lines) != 1 || started.Lines[0].VarianceQty != "-1.0000" {
		t.Fatalf("started = %+v", started)
	}

	c, rec = fakeERP(t, 200, detail)
	svc = &CycleCountsService{c: c}
	got, err := svc.Get(context.Background(), "cc1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/cycle-counts/cc1")
	if got.Session.Status != "open" || got.Lines[0].SystemQty != "10.0000" || got.Lines[0].CountedQty != "9.0000" {
		t.Fatalf("got = %+v", got)
	}
}

func TestWarehouseTasks(t *testing.T) {
	const task = `{"id":"w1","kind":"putaway","status":"open","location_id":"l1","product_id":"p1","qty":"5.0000","source_type":"receipt","source_id":"rc1","assigned_to":"usr1","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`

	c, rec := fakeERP(t, 200, `{"warehouse_tasks":[`+task+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	svc := &WarehouseTasksService{c: c}
	page, err := svc.List(context.Background(), WarehouseTaskListParams{ListParams: ListParams{Page: 1}, Status: "open"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/warehouse-tasks")
	if q(rec, "status") != "open" || q(rec, "page") != "1" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(page.Items) != 1 || page.Items[0].Kind != "putaway" || page.Items[0].Qty != "5.0000" || page.Items[0].AssignedTo != "usr1" || page.Pagination.Total != 1 {
		t.Fatalf("page = %+v", page)
	}

	c, rec = fakeERP(t, 201, task)
	svc = &WarehouseTasksService{c: c}
	created, err := svc.Create(context.Background(), WarehouseTaskCreateRequest{
		Kind: "putaway", LocationID: "l1", ProductID: "p1", Qty: "5", SourceType: "receipt", SourceID: "rc1", AssignedTo: "usr1",
	}, WithIdempotencyKey("idem-7"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/warehouse-tasks")
	if rec.Header.Get("Idempotency-Key") != "idem-7" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	if body["kind"] != "putaway" || body["location_id"] != "l1" || body["product_id"] != "p1" || body["qty"] != "5" || body["source_type"] != "receipt" || body["source_id"] != "rc1" || body["assigned_to"] != "usr1" || len(body) != 7 {
		t.Fatalf("body = %v", body)
	}
	if created.ID != "w1" || created.SourceID != "rc1" || created.Status != "open" {
		t.Fatalf("created = %+v", created)
	}

	c, rec = fakeERP(t, 200, task)
	svc = &WarehouseTasksService{c: c}
	got, err := svc.Get(context.Background(), "w1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/warehouse-tasks/w1")
	if got.ID != "w1" || got.LocationID != "l1" || got.ProductID != "p1" {
		t.Fatalf("got = %+v", got)
	}
}
