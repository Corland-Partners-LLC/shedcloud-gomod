package partnerapi

import (
	"context"
	"encoding/json"
	"testing"
)

func TestMRPListRuns(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"runs":[{"id":"r1","tenant_id":"t1","status":"completed","note":"nightly","created_at":"2026-01-02T03:04:05Z"}],"pagination":{"page":2,"page_size":10,"total":11}}`)
	m := &MRPService{c: c}
	pg, err := m.ListRuns(context.Background(), MRPRunListParams{ListParams: ListParams{Page: 2, PageSize: 10}})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/mrp/runs")
	if q(rec, "page") != "2" || q(rec, "page_size") != "10" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(pg.Items) != 1 || pg.Items[0].Status != "completed" || pg.Items[0].Note != "nightly" || pg.Pagination.Total != 11 {
		t.Fatalf("page = %+v", pg)
	}
}

func TestMRPGetRun(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"id":"r1","tenant_id":"t1","status":"completed","created_at":"2026-01-02T03:04:05Z","planned_orders":[{"id":"p1","tenant_id":"t1","run_id":"r1","product_id":"prod1","action":"buy","qty":"12.5000","converted_requisition_id":"rq1","created_at":"2026-01-02T03:04:05Z"}]}`)
	m := &MRPService{c: c}
	run, err := m.GetRun(context.Background(), "r1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/mrp/runs/r1")
	if len(run.PlannedOrders) != 1 {
		t.Fatalf("run = %+v", run)
	}
	po := run.PlannedOrders[0]
	if po.Action != "buy" || po.Qty != "12.5000" || po.ProductID != "prod1" || po.ConvertedRequisitionID == nil || *po.ConvertedRequisitionID != "rq1" || po.ConvertedWorkOrderID != nil {
		t.Fatalf("planned order = %+v", po)
	}
}

func TestMRPRun(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"id":"r2","tenant_id":"t1","status":"running","note":"manual","created_at":"2026-01-02T03:04:05Z"}`)
	m := &MRPService{c: c}
	run, err := m.Run(context.Background(), MRPRunRequest{Note: "manual"}, WithIdempotencyKey("k-mrp"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/mrp/runs")
	if rec.Header.Get("Idempotency-Key") != "k-mrp" {
		t.Fatalf("idempotency header = %q", rec.Header.Get("Idempotency-Key"))
	}
	if b := rec.json(t); b["note"] != "manual" || len(b) != 1 {
		t.Fatalf("body = %v", b)
	}
	if run.ID != "r2" || run.Status != "running" {
		t.Fatalf("run = %+v", run)
	}
}

const ncrJSON = `{"id":"n1","tenant_id":"t1","number":"NCR-1","source_type":"work_order","source_id":"wo1","severity":"major","disposition":"rework","capa":"retrain","status":"open","created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-03T03:04:05Z"}`

func TestQualityListNCRs(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"ncrs":[`+ncrJSON+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	s := &QualityService{c: c}
	pg, err := s.ListNCRs(context.Background(), QualityNCRListParams{ListParams: ListParams{Q: "NCR"}, Status: "open"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/quality/ncrs")
	if q(rec, "status") != "open" || q(rec, "q") != "NCR" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(pg.Items) != 1 || pg.Items[0].Number != "NCR-1" || pg.Items[0].CAPA != "retrain" || pg.Items[0].Severity != "major" || pg.Pagination.Total != 1 {
		t.Fatalf("page = %+v", pg)
	}
}

func TestQualityGetNCR(t *testing.T) {
	c, rec := fakeERP(t, 200, ncrJSON)
	s := &QualityService{c: c}
	n, err := s.GetNCR(context.Background(), "n 1")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Method != "GET" || rec.Path != "/partner/v1/quality/ncrs/n 1" {
		t.Fatalf("request = %s %s", rec.Method, rec.Path)
	}
	if n.SourceType != "work_order" || n.SourceID != "wo1" || n.Disposition != "rework" || n.Status != "open" {
		t.Fatalf("ncr = %+v", n)
	}
}

func TestQualityCreateNCR(t *testing.T) {
	c, rec := fakeERP(t, 201, ncrJSON)
	s := &QualityService{c: c}
	n, err := s.CreateNCR(context.Background(), QualityNCRCreateRequest{SourceType: "work_order", SourceID: "wo1", Severity: "major", Disposition: "rework", CAPA: "retrain"}, WithIdempotencyKey("k-ncr"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/quality/ncrs")
	if rec.Header.Get("Idempotency-Key") != "k-ncr" {
		t.Fatalf("idempotency header = %q", rec.Header.Get("Idempotency-Key"))
	}
	b := rec.json(t)
	if b["source_type"] != "work_order" || b["source_id"] != "wo1" || b["severity"] != "major" || b["disposition"] != "rework" || b["capa"] != "retrain" || len(b) != 5 {
		t.Fatalf("body = %v", b)
	}
	if n.ID != "n1" {
		t.Fatalf("ncr = %+v", n)
	}
}

const exportJSON = `{"id":"j1","tenant_id":"t1","dataset":"sales_ledger","status":"done","storage_key":"exports/j1.csv","created_at":"2026-01-02T03:04:05Z","finished_at":"2026-01-02T03:05:05Z"}`

func TestAnalyticsListExports(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"jobs":[`+exportJSON+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	s := &AnalyticsService{c: c}
	pg, err := s.ListExports(context.Background(), AnalyticsExportListParams{Status: "done"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/analytics/exports")
	if q(rec, "status") != "done" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(pg.Items) != 1 || pg.Items[0].Dataset != "sales_ledger" || pg.Items[0].StorageKey == nil || *pg.Items[0].StorageKey != "exports/j1.csv" || pg.Items[0].FinishedAt == nil {
		t.Fatalf("page = %+v", pg)
	}
}

func TestAnalyticsGetExport(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"id":"j2","tenant_id":"t1","dataset":"orders","status":"failed","error":"boom","created_at":"2026-01-02T03:04:05Z"}`)
	s := &AnalyticsService{c: c}
	j, err := s.GetExport(context.Background(), "j2")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/analytics/exports/j2")
	if j.Status != "failed" || j.Error != "boom" || j.StorageKey != nil || j.FinishedAt != nil {
		t.Fatalf("job = %+v", j)
	}
}

const ediJSON = `{"id":"e1","tenant_id":"t1","transaction_set":"850","name":"PO in","mapping":{"po":"BEG03"},"is_active":true,"created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-02T03:04:05Z"}`

func TestEDIMapsList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"edi_maps":[`+ediJSON+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	s := &EDIMapsService{c: c}
	pg, err := s.List(context.Background(), EDIMapListParams{ListParams: ListParams{Q: "po"}})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/edi-maps")
	if q(rec, "q") != "po" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(pg.Items) != 1 || pg.Items[0].TransactionSet != "850" || !pg.Items[0].IsActive || string(pg.Items[0].Mapping) != `{"po":"BEG03"}` {
		t.Fatalf("page = %+v", pg)
	}
}

func TestEDIMapsGet(t *testing.T) {
	c, rec := fakeERP(t, 200, ediJSON)
	s := &EDIMapsService{c: c}
	m, err := s.Get(context.Background(), "e1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/edi-maps/e1")
	if m.Name != "PO in" || m.ID != "e1" {
		t.Fatalf("map = %+v", m)
	}
}

func TestEDIMapsCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, ediJSON)
	s := &EDIMapsService{c: c}
	m, err := s.Create(context.Background(), EDIMapCreateRequest{TransactionSet: "850", Name: "PO in", Mapping: json.RawMessage(`{"po":"BEG03"}`)}, WithIdempotencyKey("k-edi"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/edi-maps")
	if rec.Header.Get("Idempotency-Key") != "k-edi" {
		t.Fatalf("idempotency header = %q", rec.Header.Get("Idempotency-Key"))
	}
	b := rec.json(t)
	mp, _ := b["mapping"].(map[string]any)
	if b["transaction_set"] != "850" || b["name"] != "PO in" || mp["po"] != "BEG03" || len(b) != 3 {
		t.Fatalf("body = %v", b)
	}
	if m.ID != "e1" {
		t.Fatalf("map = %+v", m)
	}
}

const ecoJSON = `{"id":"c1","tenant_id":"t1","number":"ECO-1","product_id":"prod1","status":"pending","bom_snapshot":[{"sku":"A","qty":"2"}],"options_snapshot":{"a":1},"approved_by":"u1","approved_at":"2026-01-03T03:04:05Z","created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-03T03:04:05Z"}`

func TestEngineeringChangesList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"engineering_changes":[`+ecoJSON+`],"pagination":{"page":1,"page_size":50,"total":1}}`)
	s := &EngineeringChangesService{c: c}
	pg, err := s.List(context.Background(), EngineeringChangeListParams{Status: "pending"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/engineering-changes")
	if q(rec, "status") != "pending" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(pg.Items) != 1 || pg.Items[0].Number != "ECO-1" || pg.Items[0].ProductID != "prod1" || pg.Items[0].ApprovedAt == nil {
		t.Fatalf("page = %+v", pg)
	}
}

func TestEngineeringChangesGet(t *testing.T) {
	c, rec := fakeERP(t, 200, ecoJSON)
	s := &EngineeringChangesService{c: c}
	e, err := s.Get(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/engineering-changes/c1")
	if e.Status != "pending" || string(e.BOMSnapshot) != `[{"sku":"A","qty":"2"}]` || string(e.OptionsSnapshot) != `{"a":1}` || e.ApprovedBy == nil || *e.ApprovedBy != "u1" {
		t.Fatalf("eco = %+v", e)
	}
}

func TestEngineeringChangesCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, ecoJSON)
	s := &EngineeringChangesService{c: c}
	e, err := s.Create(context.Background(), EngineeringChangeCreateRequest{
		ProductID:       "prod1",
		BOMSnapshot:     json.RawMessage(`[{"sku":"A","qty":"2"}]`),
		OptionsSnapshot: json.RawMessage(`{"a":1}`),
	}, WithIdempotencyKey("k-eco"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/engineering-changes")
	if rec.Header.Get("Idempotency-Key") != "k-eco" {
		t.Fatalf("idempotency header = %q", rec.Header.Get("Idempotency-Key"))
	}
	b := rec.json(t)
	bom, _ := b["bom_snapshot"].([]any)
	opts, _ := b["options_snapshot"].(map[string]any)
	if b["product_id"] != "prod1" || len(bom) != 1 || opts["a"] != float64(1) || len(b) != 3 {
		t.Fatalf("body = %v", b)
	}
	if e.ID != "c1" {
		t.Fatalf("eco = %+v", e)
	}
}
