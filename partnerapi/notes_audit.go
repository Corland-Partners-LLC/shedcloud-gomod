package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// Note is a polymorphic note attached to an entity (/partner/v1/notes).
type Note struct {
	ID                 string           `json:"id"`
	TenantID           string           `json:"tenant_id"`
	EntityType         string           `json:"entity_type"`
	EntityID           string           `json:"entity_id"`
	ParentID           string           `json:"parent_id,omitempty"`
	Body               string           `json:"body"`
	NoteType           string           `json:"note_type"`
	ShowInBillOfLading bool             `json:"show_in_bill_of_lading"`
	DueDate            string           `json:"due_date,omitempty"`
	SourceEntityType   string           `json:"source_entity_type,omitempty"`
	SourceEntityID     string           `json:"source_entity_id,omitempty"`
	SourceNoteID       string           `json:"source_note_id,omitempty"`
	CreatedBy          string           `json:"created_by,omitempty"`
	AuthorLabel        string           `json:"author_label,omitempty"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	Attachments        []NoteAttachment `json:"attachments,omitempty"`
	Assignments        []NoteAssignment `json:"assignments,omitempty"`
}

// NoteAttachment is a file attached to a note.
type NoteAttachment struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	NoteID      string    `json:"note_id"`
	StorageKey  string    `json:"storage_key"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type,omitempty"`
	SizeBytes   int64     `json:"size_bytes"`
	UploadedBy  string    `json:"uploaded_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// NoteAssignment assigns a note to a user or group (SubjectType user | group).
type NoteAssignment struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	NoteID       string    `json:"note_id"`
	SubjectType  string    `json:"subject_type"`
	SubjectID    string    `json:"subject_id"`
	SubjectLabel string    `json:"subject_label,omitempty"`
	CreatedBy    string    `json:"created_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// NoteListParams selects the entity whose notes are listed.
type NoteListParams struct {
	EntityType string `json:"entity_type,omitempty"`
	EntityID   string `json:"entity_id,omitempty"`
}

// NoteCreateRequest is POST /partner/v1/notes. NoteType is one of general,
// internal, customer, follow_up, issue, approved, denied. DueDate is
// YYYY-MM-DD.
type NoteCreateRequest struct {
	EntityType       string   `json:"entity_type,omitempty"`
	EntityID         string   `json:"entity_id,omitempty"`
	ParentID         string   `json:"parent_id,omitempty"`
	Body             string   `json:"body,omitempty"`
	NoteType         string   `json:"note_type,omitempty"`
	DueDate          string   `json:"due_date,omitempty"`
	AssigneeIDs      []string `json:"assignee_ids,omitempty"`
	AssigneeGroupIDs []string `json:"assignee_group_ids,omitempty"`
}

// NotesService calls /partner/v1/notes.
type NotesService struct{ c *Client }

// List returns the notes of one entity (not paginated). Scope
// partner-api.notes.read.
func (s *NotesService) List(ctx context.Context, params NoteListParams) ([]Note, error) {
	return listKey[Note](ctx, s.c, http.MethodGet, "/partner/v1/notes", "notes", params)
}

// Create adds a note. Scope partner-api.notes.write. Pass WithIdempotencyKey
// to make retries safe.
func (s *NotesService) Create(ctx context.Context, body NoteCreateRequest, opts ...RequestOption) (*Note, error) {
	return call[Note](ctx, s.c, http.MethodPost, "/partner/v1/notes", nil, body, opts...)
}

// Update replaces a note's body. Scope partner-api.notes.write.
func (s *NotesService) Update(ctx context.Context, noteID, body string) (*Note, error) {
	req := struct {
		Body string `json:"body"`
	}{body}
	return call[Note](ctx, s.c, http.MethodPatch, "/partner/v1/notes/"+seg(noteID), nil, req)
}

// Delete soft-deletes a note (POST /notes/{noteId}/delete). Scope
// partner-api.notes.write.
func (s *NotesService) Delete(ctx context.Context, noteID string) error {
	return s.c.http.request(ctx, http.MethodPost, "/partner/v1/notes/"+seg(noteID)+"/delete", nil, nil, nil)
}

// AuditChange is one field change of an audit event. Old and New are
// whatever JSON the field held.
type AuditChange struct {
	Field string `json:"field"`
	Old   any    `json:"old,omitempty"`
	New   any    `json:"new,omitempty"`
}

// AuditEvent is one entry of an entity audit trail (/partner/v1/audit).
// Source is "user" or "system".
type AuditEvent struct {
	ID         string         `json:"id"`
	TenantID   string         `json:"tenant_id"`
	EntityType string         `json:"entity_type"`
	EntityID   string         `json:"entity_id"`
	Action     string         `json:"action"`
	ActorID    string         `json:"actor_id,omitempty"`
	ActorLabel string         `json:"actor_label,omitempty"`
	Source     string         `json:"source"`
	Summary    string         `json:"summary,omitempty"`
	Changes    []AuditChange  `json:"changes,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
}

// AuditListParams filters GET /partner/v1/audit.
type AuditListParams struct {
	EntityType string `json:"entity_type,omitempty"`
	EntityID   string `json:"entity_id,omitempty"`
	Action     string `json:"action,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

// AuditService calls /partner/v1/audit.
type AuditService struct{ c *Client }

// List returns audit events (not paginated). Scope partner-api.audit.read.
func (s *AuditService) List(ctx context.Context, params AuditListParams) ([]AuditEvent, error) {
	return listKey[AuditEvent](ctx, s.c, http.MethodGet, "/partner/v1/audit", "events", params)
}
