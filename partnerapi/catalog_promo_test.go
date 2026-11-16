package partnerapi

import (
	"context"
	"testing"
)

func TestProductsCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"id":"p1","sku":"S-1","name":"Shed","product_type":"finished","list_price":"1999.00","is_active":true}`)
	stocked := false
	p, err := (&ProductsService{c}).Create(context.Background(), ProductCreateRequest{
		SKU: "S-1", Name: "Shed", ListPrice: "1999.00", IsStocked: &stocked,
		ManufacturerUpcs: []ProductManufacturerUpc{{Manufacturer: "Acme", Upc: "123"}},
	}, WithIdempotencyKey("k1"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/products")
	if rec.Header.Get("Idempotency-Key") != "k1" {
		t.Fatalf("idempotency key = %q", rec.Header.Get("Idempotency-Key"))
	}
	b := rec.json(t)
	if b["sku"] != "S-1" || b["list_price"] != "1999.00" || b["is_stocked"] != false {
		t.Fatalf("body = %v", b)
	}
	if _, ok := b["description"]; ok {
		t.Fatalf("empty field sent: %v", b)
	}
	if u := b["manufacturer_upcs"].([]any)[0].(map[string]any); u["upc"] != "123" {
		t.Fatalf("upcs = %v", b["manufacturer_upcs"])
	}
	if p.ID != "p1" || p.ListPrice != "1999.00" || !p.IsActive {
		t.Fatalf("product = %+v", p)
	}
}

func TestProductsUpdate(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"id":"p1","sku":"S-1","name":"New","list_price":"2000.00"}`)
	name, hide := "New", true
	p, err := (&ProductsService{c}).Update(context.Background(), "p1", ProductUpdateRequest{Name: &name, HideSKU: &hide})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/products/p1")
	b := rec.json(t)
	if len(b) != 2 || b["name"] != "New" || b["hide_sku"] != true {
		t.Fatalf("body = %v", b)
	}
	if p.Name != "New" || p.ListPrice != "2000.00" {
		t.Fatalf("product = %+v", p)
	}
}

func TestProductsPricingMatrix(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"product_id":"p1","price_list_id":"pl","price_list_code":"RETAIL","list_price":"100.00","cells":[{"width_ft":"10","length_ft":"12","dimension_key":"10x12","sell_price":"1500.00","has_tier":true}]}`)
	m, err := (&ProductsService{c}).PricingMatrix(context.Background(), "p1", ProductPricingMatrixParams{LocationID: "loc1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/products/p1/pricing-matrix")
	if q(rec, "location_id") != "loc1" {
		t.Fatalf("query = %v", rec.Query)
	}
	if m.PriceListCode != "RETAIL" || len(m.Cells) != 1 || m.Cells[0].SellPrice != "1500.00" || !m.Cells[0].HasTier || m.Cells[0].DimensionKey != "10x12" {
		t.Fatalf("matrix = %+v", m)
	}
}

func TestCatalogTemplates(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"templates":[{"id":"t1","slug":"barn","name":"Barn","family":"barns","status":"available","available":true,"acquisition_price":"0.00","highlights":["a"],"added":true,"adopted_product_id":"p9","size_range":{"unit":"ft","min_width":8,"max_width":12},"rating_avg":4.5,"option_count":3}]}`)
	list, err := (&CatalogService{c}).Templates(context.Background(), CatalogTemplateListParams{Family: "barns", Status: "available", Q: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/catalog/templates")
	if q(rec, "family") != "barns" || q(rec, "status") != "available" || q(rec, "q") != "bar" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(list) != 1 || list[0].Slug != "barn" || !list[0].Added || list[0].AdoptedProductID != "p9" ||
		list[0].AcquisitionPrice != "0.00" || list[0].SizeRange == nil || *list[0].SizeRange.MaxWidth != 12 || list[0].RatingAvg != 4.5 {
		t.Fatalf("templates = %+v", list)
	}
}

func TestCatalogTemplate(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"id":"t1","slug":"barn","name":"Barn","family":"barns","status":"available","acquisition_price":"10.00","highlights":[],"images":[{"id":"i1","url":"https://x/i.png","is_primary":true,"sort_order":0}],"options":[{"id":"o1","group_key":"size","key":"width","label":"Width","value_type":"number","min_value":"8","required":true,"choices":["a"]}],"bom":[{"component":{"id":"m1","sku":"LUM","name":"Lumber","kind":"raw","category":"wood","uom_code":"EA","default_cost":"3.50"},"quantity":"4","uom_code":"EA","qty_basis":"fixed","scrap_pct":"0","sequence":1}]}`)
	d, err := (&CatalogService{c}).Template(context.Background(), "barn")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/catalog/templates/barn")
	if d.Slug != "barn" || len(d.Images) != 1 || !d.Images[0].IsPrimary {
		t.Fatalf("detail = %+v", d)
	}
	if len(d.Options) != 1 || d.Options[0].Key != "width" || *d.Options[0].MinValue != "8" || string(d.Options[0].Choices) != `["a"]` {
		t.Fatalf("options = %+v", d.Options)
	}
	if len(d.BOM) != 1 || d.BOM[0].Component.SKU != "LUM" || d.BOM[0].Component.DefaultCost != "3.50" || d.BOM[0].Quantity != "4" {
		t.Fatalf("bom = %+v", d.BOM)
	}
}

func TestCatalogKits(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"kits":[{"id":"k1","sku":"KIT-D","name":"Door kit","category":"doors","line":"door","uom_code":"EA","default_cost":"120.00","component_count":5,"used_by_count":2,"added":false}]}`)
	kits, err := (&CatalogService{c}).Kits(context.Background(), CatalogKitListParams{Category: "doors", Line: "door", Q: "kit"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/catalog/kits")
	if q(rec, "category") != "doors" || q(rec, "line") != "door" || q(rec, "q") != "kit" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(kits) != 1 || kits[0].SKU != "KIT-D" || kits[0].DefaultCost != "120.00" || kits[0].UsedByCount != 2 {
		t.Fatalf("kits = %+v", kits)
	}
}

const promoJSON = `{"id":"pc1","tenant_id":"t","code":"SAVE","status":"active","starts_at":"2026-01-01T00:00:00Z","rules":{"version":1},"max_redemptions":10,"redemption_count":3,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-02T00:00:00Z"}`

func TestPromoCodesList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"promo_codes":[`+promoJSON+`]}`)
	list, err := (&PromoCodesService{c}).List(context.Background(), PromoCodeListParams{Status: "active", Q: "SAVE"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/promo-codes")
	if q(rec, "status") != "active" || q(rec, "q") != "SAVE" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(list) != 1 || list[0].Code != "SAVE" || *list[0].MaxRedemptions != 10 || list[0].RedemptionCount != 3 ||
		list[0].StartsAt == nil || string(list[0].Rules) != `{"version":1}` {
		t.Fatalf("promos = %+v", list)
	}
}

func TestPromoCodesCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, promoJSON)
	p, err := (&PromoCodesService{c}).Create(context.Background(), PromoCodeCreateRequest{
		Code: "SAVE", Status: "active", Rules: []byte(`{"version":1}`), MaxRedemptions: 10, StartsAt: "2026-01-01T00:00:00Z",
	}, WithIdempotencyKey("k2"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/promo-codes")
	if rec.Header.Get("Idempotency-Key") != "k2" {
		t.Fatalf("idempotency key = %q", rec.Header.Get("Idempotency-Key"))
	}
	b := rec.json(t)
	if b["code"] != "SAVE" || b["max_redemptions"] != float64(10) || b["starts_at"] != "2026-01-01T00:00:00Z" {
		t.Fatalf("body = %v", b)
	}
	if rules, _ := b["rules"].(map[string]any); rules["version"] != float64(1) {
		t.Fatalf("rules = %v", b["rules"])
	}
	if _, ok := b["ends_at"]; ok {
		t.Fatalf("empty field sent: %v", b)
	}
	if p.ID != "pc1" {
		t.Fatalf("promo = %+v", p)
	}
}

func TestPromoCodesValidate(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"valid":true,"promo":`+promoJSON+`,"effects":{"matched":true,"adjustments":[{"field":"security_deposit","operator":"set","value_unit":"usd","value":"0.00"}]}}`)
	v, err := (&PromoCodesService{c}).Validate(context.Background(), PromoCodeValidateRequest{
		Code: "SAVE", QuoteID: "q1", Subtotal: "5000.00", TermMonths: 36, LineSources: []string{"configurator"},
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/promo-codes/validate")
	b := rec.json(t)
	if b["code"] != "SAVE" || b["subtotal"] != "5000.00" || b["term_months"] != float64(36) || b["line_sources"].([]any)[0] != "configurator" {
		t.Fatalf("body = %v", b)
	}
	if _, ok := b["customer_id"]; ok {
		t.Fatalf("empty field sent: %v", b)
	}
	if !v.Valid || v.Promo == nil || v.Promo.Code != "SAVE" || v.Effects == nil || !v.Effects.Matched ||
		len(v.Effects.Adjustments) != 1 || v.Effects.Adjustments[0].Value != "0.00" || v.Effects.Adjustments[0].Operator != "set" {
		t.Fatalf("validation = %+v", v)
	}
}

func TestPromoCodesValidateInvalid(t *testing.T) {
	c, _ := fakeERP(t, 200, `{"valid":false,"reason":"expired"}`)
	v, err := (&PromoCodesService{c}).Validate(context.Background(), PromoCodeValidateRequest{Code: "OLD"})
	if err != nil {
		t.Fatal(err)
	}
	if v.Valid || v.Reason != "expired" || v.Promo != nil || v.Effects != nil {
		t.Fatalf("validation = %+v", v)
	}
}

func TestPromoCodesGet(t *testing.T) {
	c, rec := fakeERP(t, 200, promoJSON)
	p, err := (&PromoCodesService{c}).Get(context.Background(), "pc1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/promo-codes/pc1")
	if p.Code != "SAVE" || p.Status != "active" {
		t.Fatalf("promo = %+v", p)
	}
}

func TestPromoCodesUpdate(t *testing.T) {
	c, rec := fakeERP(t, 200, promoJSON)
	status, zero, ends := "paused", 0, ""
	p, err := (&PromoCodesService{c}).Update(context.Background(), "pc1", PromoCodeUpdateRequest{Status: &status, MaxRedemptions: &zero, EndsAt: &ends})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/promo-codes/pc1")
	b := rec.json(t)
	if len(b) != 3 || b["status"] != "paused" || b["max_redemptions"] != float64(0) || b["ends_at"] != "" {
		t.Fatalf("body = %v", b)
	}
	if p.ID != "pc1" {
		t.Fatalf("promo = %+v", p)
	}
}
