package partnerapi

import (
	"context"
	"testing"
)

func TestServiceTicketsList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"service_tickets":[{"id":"t1","work_order_id":"w1","priority":"high","status":"open","balance_due":"12.50","sla_breached":false,"contact":{"name":"Ann"}}]}`)
	s := &ServiceTicketsService{c}
	page, err := s.List(context.Background(), ServiceTicketListParams{
		ListParams:  ListParams{Page: 2, PageSize: 10, Q: "roof"},
		Status:      "open",
		WorkOrderID: "w1",
		Sort:        "created_at",
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/service-tickets")
	for k, v := range map[string]string{"page": "2", "page_size": "10", "q": "roof", "status": "open", "work_order_id": "w1", "sort": "created_at"} {
		if q(rec, k) != v {
			t.Fatalf("query %s = %q, want %q", k, q(rec, k), v)
		}
	}
	if len(page.Items) != 1 || page.Items[0].BalanceDue != "12.50" || page.Items[0].Priority != "high" || string(page.Items[0].Contact) != `{"name":"Ann"}` {
		t.Fatalf("items = %+v", page.Items)
	}
}

func TestServiceTicketsGet(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"id":"t1","ticket_number":"ST-1","status":"open","priority":"low","balance_due":"0.00","warranty_months":12,"display":{"customer_name":"Ann"}}`)
	got, err := (&ServiceTicketsService{c}).Get(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/service-tickets/t1")
	if got.TicketNumber != "ST-1" || got.WarrantyMonths == nil || *got.WarrantyMonths != 12 || got.Display == nil || got.Display.CustomerName != "Ann" {
		t.Fatalf("got = %+v", got)
	}
}

func TestServiceTicketsCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"id":"t2","status":"new","priority":"medium","balance_due":"0"}`)
	got, err := (&ServiceTicketsService{c}).Create(context.Background(), ServiceTicketCreateRequest{
		WorkOrderID:   "w1",
		Description:   "leak",
		ScheduledDate: "2026-12-01",
	}, WithIdempotencyKey("k1"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/service-tickets")
	if rec.Header.Get("Idempotency-Key") != "k1" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	if len(body) != 3 || body["work_order_id"] != "w1" || body["description"] != "leak" || body["scheduled_date"] != "2026-12-01" {
		t.Fatalf("body = %v", body)
	}
	if got.ID != "t2" || got.Status != "new" {
		t.Fatalf("got = %+v", got)
	}
}

func TestServiceTicketsUpdate(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"id":"t1","status":"open","priority":"high","balance_due":"5.00"}`)
	empty, bal, sla := "", "5.00", true
	got, err := (&ServiceTicketsService{c}).Update(context.Background(), "t1", ServiceTicketUpdateRequest{
		TechName:    &empty,
		BalanceDue:  &bal,
		SLABreached: &sla,
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/service-tickets/t1")
	body := rec.json(t)
	if len(body) != 3 || body["tech_name"] != "" || body["balance_due"] != "5.00" || body["sla_breached"] != true {
		t.Fatalf("body = %v", body)
	}
	if got.BalanceDue != "5.00" {
		t.Fatalf("got = %+v", got)
	}
}

func TestServiceTicketsUpdateStatus(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"id":"t1","status":"solved","prior_status":"open","priority":"low","balance_due":"0"}`)
	got, err := (&ServiceTicketsService{c}).UpdateStatus(context.Background(), "t1", "solved")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/service-tickets/t1/status")
	body := rec.json(t)
	if len(body) != 1 || body["status"] != "solved" {
		t.Fatalf("body = %v", body)
	}
	if got.Status != "solved" || got.PriorStatus != "open" {
		t.Fatalf("got = %+v", got)
	}
}

func TestNotesList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"notes":[{"id":"n1","entity_type":"order","entity_id":"o1","body":"hi","note_type":"general","show_in_bill_of_lading":true,"assignments":[{"id":"a1","subject_type":"user","subject_id":"u1"}],"created_at":"2026-11-01T10:00:00Z"}]}`)
	got, err := (&NotesService{c}).List(context.Background(), NoteListParams{EntityType: "order", EntityID: "o1"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/notes")
	if q(rec, "entity_type") != "order" || q(rec, "entity_id") != "o1" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(got) != 1 || got[0].Body != "hi" || !got[0].ShowInBillOfLading || len(got[0].Assignments) != 1 || got[0].Assignments[0].SubjectType != "user" {
		t.Fatalf("got = %+v", got)
	}
}

func TestNotesCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"id":"n2","entity_type":"order","entity_id":"o1","body":"call back","note_type":"follow_up","due_date":"2026-12-02"}`)
	got, err := (&NotesService{c}).Create(context.Background(), NoteCreateRequest{
		EntityType:  "order",
		EntityID:    "o1",
		Body:        "call back",
		NoteType:    "follow_up",
		DueDate:     "2026-12-02",
		AssigneeIDs: []string{"u1", "u2"},
	}, WithIdempotencyKey("nk"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/notes")
	if rec.Header.Get("Idempotency-Key") != "nk" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	ids, _ := body["assignee_ids"].([]any)
	if len(body) != 6 || body["note_type"] != "follow_up" || len(ids) != 2 || ids[1] != "u2" {
		t.Fatalf("body = %v", body)
	}
	if got.ID != "n2" || got.DueDate != "2026-12-02" {
		t.Fatalf("got = %+v", got)
	}
}

func TestNotesUpdate(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"id":"n2","body":"edited","note_type":"general"}`)
	got, err := (&NotesService{c}).Update(context.Background(), "n2", "edited")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "PATCH", "/partner/v1/notes/n2")
	body := rec.json(t)
	if len(body) != 1 || body["body"] != "edited" {
		t.Fatalf("body = %v", body)
	}
	if got.Body != "edited" {
		t.Fatalf("got = %+v", got)
	}
}

func TestNotesDelete(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"status":"deleted"}`)
	if err := (&NotesService{c}).Delete(context.Background(), "n2"); err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/notes/n2/delete")
	if len(rec.Body) != 0 {
		t.Fatalf("body = %q, want empty", rec.Body)
	}
}

func TestAuditList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"events":[{"id":"e1","entity_type":"order","entity_id":"o1","action":"update","source":"user","actor_label":"Ann","changes":[{"field":"status","old":"draft","new":"confirmed"}],"created_at":"2026-11-01T10:00:00Z"}]}`)
	got, err := (&AuditService{c}).List(context.Background(), AuditListParams{EntityType: "order", EntityID: "o1", Action: "update", Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/audit")
	for k, v := range map[string]string{"entity_type": "order", "entity_id": "o1", "action": "update", "limit": "5"} {
		if q(rec, k) != v {
			t.Fatalf("query %s = %q, want %q", k, q(rec, k), v)
		}
	}
	if len(got) != 1 || got[0].ActorLabel != "Ann" || len(got[0].Changes) != 1 || got[0].Changes[0].New != "confirmed" {
		t.Fatalf("got = %+v", got)
	}
}

func TestCommunicationsSettings(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"connected":true,"has_auth_token":true,"has_voice_key":false,"recording_enabled":true,"a2p_status":"approved","campaign_status":"verified","welcome_text":"Hi","handoff_keywords":["human"],"after_hours_enabled":true,"after_hours_message":"Closed","last_error":"","updated_at":"2026-11-01T10:00:00Z"}`)
	got, err := (&CommunicationsService{c}).Settings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/communications/settings")
	if !got.Connected || got.A2PStatus != "approved" || len(got.HandoffKeywords) != 1 || !got.AfterHoursEnabled || got.UpdatedAt.IsZero() {
		t.Fatalf("got = %+v", got)
	}
}

func TestCommunicationsNumbers(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"numbers":[{"id":"p1","e164":"+15125550100","sms_enabled":true,"purpose":"sales","is_default_sms":true,"inbound_route":"user"}]}`)
	got, err := (&CommunicationsService{c}).Numbers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/communications/numbers")
	if len(got) != 1 || got[0].E164 != "+15125550100" || !got[0].IsDefaultSMS || got[0].InboundRoute != "user" {
		t.Fatalf("got = %+v", got)
	}
}

func TestCommunicationsInbox(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"conversations":[{"id":"c1","channel":"sms","peer_e164":"+15125550101","status":"open","last_message_at":"2026-11-01T10:00:00Z"}]}`)
	got, err := (&CommunicationsService{c}).Inbox(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/communications/inbox")
	if len(got) != 1 || got[0].PeerE164 != "+15125550101" || got[0].LastMessageAt.IsZero() {
		t.Fatalf("got = %+v", got)
	}
}

func TestCommunicationsListConversations(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"conversations":[{"id":"c1","customer_id":"cu1","channel":"sms","peer_e164":"+1","status":"open"}]}`)
	got, err := (&CommunicationsService{c}).ListConversations(context.Background(), "cu1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/customers/cu1/conversations")
	if len(got) != 1 || got[0].CustomerID != "cu1" {
		t.Fatalf("got = %+v", got)
	}
}

func TestCommunicationsStartConversation(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"conversation":{"id":"c9","channel":"sms","peer_e164":"+1","status":"open"},"message":{"id":"m1","conversation_id":"c9","direction":"outbound","actor":"user","body":"Hello"}}`)
	got, err := (&CommunicationsService{c}).StartConversation(context.Background(), "cu1", ConversationStartRequest{Body: "Hello", Channel: "sms"}, WithIdempotencyKey("ck"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/customers/cu1/conversations")
	if rec.Header.Get("Idempotency-Key") != "ck" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	if len(body) != 2 || body["body"] != "Hello" || body["channel"] != "sms" {
		t.Fatalf("body = %v", body)
	}
	if got.Conversation.ID != "c9" || got.Message.ConversationID != "c9" || got.Message.Body != "Hello" {
		t.Fatalf("got = %+v", got)
	}
}

func TestCommunicationsListMessages(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"messages":[{"id":"m1","conversation_id":"c1","direction":"inbound","actor":"customer","body":"yo","created_at":"2026-11-01T10:00:00Z"}]}`)
	got, err := (&CommunicationsService{c}).ListMessages(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/conversations/c1/messages")
	if len(got) != 1 || got[0].Direction != "inbound" || got[0].Body != "yo" {
		t.Fatalf("got = %+v", got)
	}
}

func TestCommunicationsSendMessage(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"id":"m2","conversation_id":"c1","direction":"outbound","actor":"user","body":"on it"}`)
	got, err := (&CommunicationsService{c}).SendMessage(context.Background(), "c1", MessageSendRequest{Body: "on it", Channel: "sms"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/conversations/c1/messages")
	body := rec.json(t)
	if len(body) != 2 || body["body"] != "on it" || body["channel"] != "sms" {
		t.Fatalf("body = %v", body)
	}
	if got.ID != "m2" || got.Direction != "outbound" {
		t.Fatalf("got = %+v", got)
	}
}

func TestOpportunitiesList(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"opportunities":[{"id":"op1","name":"Deal","stage":"proposal","amount":"1500.00","close_date":"2026-12-31","owner_id":"u1"}],"pagination":{"page":1,"page_size":50,"total":1}}`)
	page, err := (&OpportunitiesService{c}).List(context.Background(), OpportunityListParams{ListParams: ListParams{Page: 1, Q: "deal"}, Status: "proposal"})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/opportunities")
	if q(rec, "page") != "1" || q(rec, "q") != "deal" || q(rec, "status") != "proposal" {
		t.Fatalf("query = %v", rec.Query)
	}
	if len(page.Items) != 1 || page.Items[0].Amount != "1500.00" || page.Items[0].CloseDate == nil || *page.Items[0].CloseDate != "2026-12-31" || page.Pagination.Total != 1 {
		t.Fatalf("page = %+v", page)
	}
}

func TestOpportunitiesGet(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"id":"op1","name":"Deal","stage":"won","amount":"10.00","lead_id":"l1"}`)
	got, err := (&OpportunitiesService{c}).Get(context.Background(), "op1")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "GET", "/partner/v1/opportunities/op1")
	if got.Stage != "won" || got.LeadID == nil || *got.LeadID != "l1" {
		t.Fatalf("got = %+v", got)
	}
}

func TestOpportunitiesCreate(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"id":"op2","name":"New","stage":"qualify","amount":"99.99"}`)
	got, err := (&OpportunitiesService{c}).Create(context.Background(), OpportunityCreateRequest{
		Name:      "New",
		Amount:    "99.99",
		CloseDate: "2026-12-31",
		LeadID:    "l1",
	}, WithIdempotencyKey("ok"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/opportunities")
	if rec.Header.Get("Idempotency-Key") != "ok" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	if len(body) != 4 || body["amount"] != "99.99" || body["close_date"] != "2026-12-31" || body["lead_id"] != "l1" {
		t.Fatalf("body = %v", body)
	}
	if got.ID != "op2" || got.Amount != "99.99" {
		t.Fatalf("got = %+v", got)
	}
}

func TestDocumentsUploadURL(t *testing.T) {
	c, rec := fakeERP(t, 200, `{"storage_key":"tenants/t/file.pdf","upload":{"url":"https://r2.example/put","method":"PUT","expires_at":"2026-11-01T10:15:00Z","headers":{"Content-Type":"application/pdf"}}}`)
	got, err := (&DocumentsService{c}).UploadURL(context.Background(), DocumentUploadURLRequest{
		EntityType:  "order",
		EntityID:    "o1",
		Role:        "attachment",
		FileName:    "file.pdf",
		ContentType: "application/pdf",
	})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/documents/upload-url")
	body := rec.json(t)
	if len(body) != 5 || body["file_name"] != "file.pdf" || body["content_type"] != "application/pdf" || body["role"] != "attachment" {
		t.Fatalf("body = %v", body)
	}
	if got.StorageKey != "tenants/t/file.pdf" || got.Upload.Method != "PUT" || got.Upload.URL != "https://r2.example/put" || got.Upload.Headers["Content-Type"] != "application/pdf" || got.Upload.ExpiresAt.IsZero() {
		t.Fatalf("got = %+v", got)
	}
}

func TestDocumentsAttach(t *testing.T) {
	c, rec := fakeERP(t, 201, `{"id":"d1","link_id":"l1","entity_type":"order","entity_id":"o1","role":"attachment","position":1,"is_primary":true,"file_name":"file.pdf","mime_type":"application/pdf","size_bytes":2048,"storage_key":"tenants/t/file.pdf","url":"https://r2.example/get","access":"edit","created_at":"2026-11-01T10:00:00Z"}`)
	got, err := (&DocumentsService{c}).Attach(context.Background(), DocumentAttachRequest{
		EntityType: "order",
		EntityID:   "o1",
		Role:       "attachment",
		FileName:   "file.pdf",
		MimeType:   "application/pdf",
		SizeBytes:  2048,
		StorageKey: "tenants/t/file.pdf",
		IsPrimary:  true,
		Metadata:   map[string]any{"src": "sdk"},
	}, WithIdempotencyKey("dk"))
	if err != nil {
		t.Fatal(err)
	}
	expect(t, rec, "POST", "/partner/v1/documents")
	if rec.Header.Get("Idempotency-Key") != "dk" {
		t.Fatalf("Idempotency-Key = %q", rec.Header.Get("Idempotency-Key"))
	}
	body := rec.json(t)
	if len(body) != 9 || body["storage_key"] != "tenants/t/file.pdf" || body["size_bytes"] != float64(2048) || body["is_primary"] != true {
		t.Fatalf("body = %v", body)
	}
	if got.ID != "d1" || got.FileName != "file.pdf" || got.SizeBytes != 2048 || got.StorageKey != "tenants/t/file.pdf" || got.URL != "https://r2.example/get" || got.Access != "edit" || !got.IsPrimary {
		t.Fatalf("got = %+v", got)
	}
}
