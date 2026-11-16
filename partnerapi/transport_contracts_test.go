package partnerapi

import (
	"context"
	"testing"
)

func TestDeliveryOffersList(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"offers":[{"id":"o1","manufacturer_tenant_id":"m1","carrier_tenant_id":"c1","work_order_id":"w1","kind":"delivery","status":"offered","offered_at":"2026-09-01T10:00:00Z","haul_fee":"250.00","display":{"wo_number":"WO-1","lat":30.5},"charges":{"delivery_charge":"200.00","fuel_surcharge":"10.00","additional_charges":[{"name":"Tarp","amount":"15.00"}],"total":"225.00"},"created_at":"2026-09-01T10:00:00Z","updated_at":"2026-09-01T10:00:00Z"}]}`)
	s := &DeliveryOffersService{c}
	out, err := s.List(ctx, DeliveryOfferListParams{WorkOrderID: "w1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/delivery-offers")
	if q(rec, "work_order_id") != "w1" {
		t.Fatalf("query = %v", rec.Query)
	}
	o := out[0]
	if len(out) != 1 || o.ID != "o1" || o.HaulFee != "250.00" || o.Display.WONumber != "WO-1" || *o.Display.Lat != 30.5 ||
		o.Charges.Total != "225.00" || o.Charges.AdditionalCharges[0].Name != "Tarp" {
		t.Fatalf("offers = %+v", out)
	}
}

func TestDeliveryOffersCreate(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 201, `{"id":"o2","work_order_id":"w9","status":"offered","carrier_tenant_id":"c1"}`)
	s := &DeliveryOffersService{c}
	o, err := s.Create(ctx, "w9", DeliveryOfferCreateRequest{AgreementCode: "AG-1", Message: "hi", HaulFee: "300.00"}, WithIdempotencyKey("k1"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/work-orders/w9/delivery-offers")
	if rec.Header.Get("Idempotency-Key") != "k1" {
		t.Fatal("idempotency key missing")
	}
	if b := rec.json(t); b["agreement_code"] != "AG-1" || b["message"] != "hi" || b["haul_fee"] != "300.00" || len(b) != 3 {
		t.Fatalf("body = %v", b)
	}
	if o.ID != "o2" || o.Status != "offered" {
		t.Fatalf("offer = %+v", o)
	}
}

func TestDeliveryOffersAccept(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"offer":{"id":"o3","status":"accepted"},"public_track_token":"tok123"}`)
	s := &DeliveryOffersService{c}
	r, err := s.Accept(ctx, "o3")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/delivery-offers/o3/accept")
	if len(rec.Body) != 0 {
		t.Fatalf("body = %q", rec.Body)
	}
	if r.Offer.Status != "accepted" || r.PublicTrackToken != "tok123" {
		t.Fatalf("result = %+v", r)
	}
}

func TestDeliveryOffersTransition(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"id":"o4","status":"failed","failed_reason":"gate locked"}`)
	s := &DeliveryOffersService{c}
	o, err := s.Transition(ctx, "o4", DeliveryOfferTransitionRequest{Status: "failed", Reason: "gate locked"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/delivery-offers/o4/transition")
	if b := rec.json(t); b["status"] != "failed" || b["reason"] != "gate locked" || len(b) != 2 {
		t.Fatalf("body = %v", b)
	}
	if o.Status != "failed" || o.FailedReason != "gate locked" {
		t.Fatalf("offer = %+v", o)
	}
}

func TestDeliveryOffersSubmitPOD(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 201, `{"id":"p1","signature":"sig","notes":"ok","geo_lat":"30.1","geo_lng":"-97.2","created_at":"2026-09-02T12:00:00Z"}`)
	s := &DeliveryOffersService{c}
	pod, err := s.SubmitPOD(ctx, "o5", DeliveryOfferPODRequest{Signature: "sig", Notes: "ok", GeoLat: "30.1", GeoLng: "-97.2", Photos: []any{"data:image/png;base64,AAA"}})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/delivery-offers/o5/pod")
	b := rec.json(t)
	if b["signature"] != "sig" || b["geo_lat"] != "30.1" || b["checklist"] != nil || len(b["photos"].([]any)) != 1 {
		t.Fatalf("body = %v", b)
	}
	if pod.ID != "p1" || pod.GeoLng != "-97.2" || pod.CreatedAt.IsZero() {
		t.Fatalf("pod = %+v", pod)
	}
}

const rtoAgreementJSON = `{"id":"a1","tenant_id":"t1","customer_id":"c1","agreement_number":"RTO-1","status":"draft","term_months":24,"monthly_payment":"150.00","security_deposit":"100.00","additional_deposit":"0.00","customer_reserve":"0.00","damage_waiver":true,"damage_waiver_fee":"10.00","rent":"140.00","rate_factor":"1.85","customer_reserve_months":0,"use_damage_waiver_pct":false,"contract_kind":"rto","schedule_kind":"installments","created_at":"2026-09-01T00:00:00Z","updated_at":"2026-09-01T00:00:00Z","installments":[{"id":"i1","tenant_id":"t1","rto_agreement_id":"a1","sequence_no":1,"line_kind":"monthly","due_date":"2026-10-01","amount":"150.00","status":"pending","created_at":"2026-09-01T00:00:00Z","updated_at":"2026-09-01T00:00:00Z"}],"total_scheduled":"3600.00","balance_due":"3600.00"}`

func TestRTOAgreementsList(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"agreements":[`+rtoAgreementJSON+`],"pagination":{"page":2,"page_size":10,"total":11}}`)
	s := &RTOAgreementsService{c}
	page, err := s.List(ctx, RTOAgreementListParams{ListParams: ListParams{Page: 2, PageSize: 10, Q: "smith"}, CustomerID: "c1", QuoteID: "q1", Status: "draft"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/rto-agreements")
	if q(rec, "page") != "2" || q(rec, "page_size") != "10" || q(rec, "q") != "smith" || q(rec, "customer_id") != "c1" || q(rec, "quote_id") != "q1" || q(rec, "status") != "draft" {
		t.Fatalf("query = %v", rec.Query)
	}
	a := page.Items[0]
	if len(page.Items) != 1 || page.Pagination.Total != 11 || a.AgreementNumber != "RTO-1" || a.MonthlyPayment != "150.00" || !a.DamageWaiver ||
		a.Installments[0].Amount != "150.00" || a.TotalScheduled != "3600.00" {
		t.Fatalf("page = %+v", page)
	}
}

func TestRTOAgreementsCreate(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 201, rtoAgreementJSON)
	s := &RTOAgreementsService{c}
	a, err := s.Create(ctx, RTOAgreementCreateRequest{CustomerID: "c1", QuoteID: "q1", TermMonths: 24, MonthlyPayment: "150.00", DamageWaiver: true, StartDate: "2026-10-01"}, WithIdempotencyKey("k2"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/rto-agreements")
	if rec.Header.Get("Idempotency-Key") != "k2" {
		t.Fatal("idempotency key missing")
	}
	if b := rec.json(t); b["customer_id"] != "c1" || b["term_months"] != float64(24) || b["monthly_payment"] != "150.00" || b["damage_waiver"] != true || b["rent"] != nil || len(b) != 6 {
		t.Fatalf("body = %v", b)
	}
	if a.ID != "a1" || a.Status != "draft" || a.TermMonths != 24 {
		t.Fatalf("agreement = %+v", a)
	}
}

func TestRTOAgreementsGetActivate(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, rtoAgreementJSON)
	s := &RTOAgreementsService{c}
	a, err := s.Get(ctx, "a1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/rto-agreements/a1")
	if a.ID != "a1" || a.ContractKind != "rto" || a.RateFactor != "1.85" {
		t.Fatalf("agreement = %+v", a)
	}

	if _, err := s.Activate(ctx, "a1"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/rto-agreements/a1/activate")
	if len(rec.Body) != 0 {
		t.Fatalf("body = %q", rec.Body)
	}
}

const signatureFollowJSON = `{"signature_process_status":"out_for_signature","contract_pdf_status":"pending","started_at":"2026-09-01T09:00:00Z","is_rto":true,"corenter_required":false,"sheet":{"order_id":"ord1","full_name":"Jane Doe","ssn_masked":"***-**-1234","has_ssn":true,"references":[{"name":"Bob","phone":"555","relationship":"friend"}]},"signatures":[{"id":"s1","order_id":"ord1","role":"customer","storage_key":"k","signed_at":"2026-09-01T10:00:00Z"}],"missing":[{"party":"customer","section_id":"sec1","field":"email","label":"Email"}],"parties_complete":false,"next_party":"customer","has_customer_share":true,"has_corenter_share":false,"customer_signed":false,"corenter_signed":false,"salesperson_signed":false,"signature_provider":"native","boldsign_connected":false,"docusign_connected":false,"available_providers":["native"]}`

func checkFollow(t *testing.T, f *SignatureProcessFollow) {
	t.Helper()
	if f.ProcessStatus != "out_for_signature" || f.ContractPDFStatus != "pending" || !f.IsRTO || f.StartedAt == nil ||
		f.Sheet.FullName != "Jane Doe" || f.Sheet.SSNMasked != "***-**-1234" || f.Sheet.References[0].Name != "Bob" ||
		f.Signatures[0].Role != "customer" || f.Missing[0].Field != "email" || f.NextParty != "customer" ||
		f.SignatureProvider != "native" || len(f.AvailableProviders) != 1 {
		t.Fatalf("follow = %+v", f)
	}
}

func TestOrdersSignatureProcess(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, signatureFollowJSON)
	s := newOrdersService(c)
	f, err := s.SignatureProcess(ctx, "ord1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/orders/ord1/signature-process")
	checkFollow(t, f)
}

func TestOrdersStartSignature(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, signatureFollowJSON)
	s := newOrdersService(c)
	f, err := s.StartSignature(ctx, "ord1", SignatureStartRequest{Send: true, Provider: "native"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/orders/ord1/signature-process/start")
	if b := rec.json(t); b["send"] != true || b["provider"] != "native" || len(b) != 2 {
		t.Fatalf("body = %v", b)
	}
	checkFollow(t, f)
}

func TestOrdersSendSignature(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"id":"ord1","number":"ORD-1","status":"open"}`)
	s := newOrdersService(c)
	o, err := s.SendSignature(ctx, "ord1", SignatureSendRequest{To: "corenter"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/orders/ord1/signature-process/send")
	if b := rec.json(t); b["to"] != "corenter" || len(b) != 1 {
		t.Fatalf("body = %v", b)
	}
	if o.ID != "ord1" {
		t.Fatalf("order = %+v", o)
	}
}

func TestOrdersCancelSignature(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"signature_process_status":"cancelled","contract_pdf_status":"none","signatures":[],"missing":[]}`)
	s := newOrdersService(c)
	f, err := s.CancelSignature(ctx, "ord1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/orders/ord1/signature-process/cancel")
	if len(rec.Body) != 0 {
		t.Fatalf("body = %q", rec.Body)
	}
	if f.ProcessStatus != "cancelled" {
		t.Fatalf("follow = %+v", f)
	}
}

func TestOrdersRestartSignature(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, signatureFollowJSON)
	s := newOrdersService(c)
	f, err := s.RestartSignature(ctx, "ord1", SignatureRestartRequest{ClearForms: true})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/orders/ord1/signature-process/restart")
	if b := rec.json(t); b["clear_forms"] != true || len(b) != 1 {
		t.Fatalf("body = %v", b)
	}
	checkFollow(t, f)
}

func TestOrdersRetryContractPDF(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"signature_process_status":"generating_pdf","contract_pdf_status":"generating","signatures":[],"missing":[]}`)
	s := newOrdersService(c)
	f, err := s.RetryContractPDF(ctx, "ord1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/orders/ord1/signature-process/retry-pdf")
	if len(rec.Body) != 0 {
		t.Fatalf("body = %q", rec.Body)
	}
	if f.ContractPDFStatus != "generating" {
		t.Fatalf("follow = %+v", f)
	}
}

func TestSoldUpgradesReport(t *testing.T) {
	ctx := context.Background()
	c, rec := fakeERP(t, 200, `{"rows":[{"sold_at":"2026-09-01T10:00:00Z","sale_date":"2026-09-01","order_id":"o1","order_number":"ORD-1","customer_name":"Jane","store_location_id":"l1","store_location_name":"Austin","salesperson_name":"Sam","product_name":"Barn 10x12","upgrade_label":"Loft","amount":"350.00","action":"added"}],"pagination":{"page":1,"page_size":50,"total_items":1},"sum_amount":"350.00","by_label":[{"label":"Loft","count":1,"sum_amount":"350.00"}]}`)
	s := &SoldUpgradesService{c}
	r, err := s.Report(ctx, SoldUpgradesReportRequest{Action: "added", SaleFrom: "2026-09-01", SaleTo: "2026-09-30", StoreLocationID: "l1", Q: "loft", Page: 1, PageSize: 50, IncludeSums: true})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/sold-upgrades")
	b := rec.json(t)
	if b["action"] != "added" || b["sale_from"] != "2026-09-01" || b["sale_to"] != "2026-09-30" || b["store_location_id"] != "l1" ||
		b["q"] != "loft" || b["page"] != float64(1) || b["page_size"] != float64(50) || b["include_sums"] != true || len(b) != 8 {
		t.Fatalf("body = %v", b)
	}
	row := r.Rows[0]
	if len(r.Rows) != 1 || row.OrderNumber != "ORD-1" || row.UpgradeLabel != "Loft" || row.Amount != "350.00" || row.SoldAt == nil ||
		r.Pagination.TotalItems != 1 || r.SumAmount != "350.00" || r.ByLabel[0].Count != 1 || r.ByLabel[0].SumAmount != "350.00" {
		t.Fatalf("report = %+v", r)
	}
}
