package partnerapi

import (
	"context"
	"net/http"
	"time"
)

// CommunicationSettings is the tenant messaging/voice configuration
// (GET /partner/v1/communications/settings). Secrets are never returned, only
// the has_* flags.
type CommunicationSettings struct {
	Connected         bool      `json:"connected"`
	HasAuthToken      bool      `json:"has_auth_token"`
	HasVoiceKey       bool      `json:"has_voice_key"`
	RecordingEnabled  bool      `json:"recording_enabled"`
	A2PStatus         string    `json:"a2p_status"`
	CampaignStatus    string    `json:"campaign_status,omitempty"`
	WelcomeText       string    `json:"welcome_text,omitempty"`
	HandoffKeywords   []string  `json:"handoff_keywords"`
	AfterHoursEnabled bool      `json:"after_hours_enabled"`
	AfterHoursMessage string    `json:"after_hours_message,omitempty"`
	LastError         string    `json:"last_error,omitempty"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CommunicationNumber is a phone number of the tenant.
type CommunicationNumber struct {
	ID                  string    `json:"id"`
	TenantID            string    `json:"tenant_id"`
	TwilioSID           string    `json:"twilio_sid,omitempty"`
	E164                string    `json:"e164"`
	FriendlyName        string    `json:"friendly_name,omitempty"`
	SMSEnabled          bool      `json:"sms_enabled"`
	VoiceEnabled        bool      `json:"voice_enabled"`
	MMSEnabled          bool      `json:"mms_enabled"`
	Purpose             string    `json:"purpose"`
	LocationID          string    `json:"location_id,omitempty"`
	AIEnabled           bool      `json:"ai_enabled"`
	AssignedUserID      string    `json:"assigned_user_id,omitempty"`
	IsDefaultSMS        bool      `json:"is_default_sms"`
	IsDefaultVoice      bool      `json:"is_default_voice"`
	InboundRoute        string    `json:"inbound_route"`
	MessagingServiceSID string    `json:"messaging_service_sid,omitempty"`
	CampaignSID         string    `json:"campaign_sid,omitempty"`
	A2PAttached         bool      `json:"a2p_attached"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Conversation is a message thread with one phone number
// (/partner/v1/communications/inbox, /customers/{id}/conversations).
type Conversation struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	CustomerID     string     `json:"customer_id,omitempty"`
	NumberID       string     `json:"number_id,omitempty"`
	Channel        string     `json:"channel"`
	PeerE164       string     `json:"peer_e164"`
	Status         string     `json:"status"`
	AssignedUserID string     `json:"assigned_user_id,omitempty"`
	AssistantSID   string     `json:"assistant_sid,omitempty"`
	IgnoredAt      *time.Time `json:"ignored_at,omitempty"`
	LastMessageAt  time.Time  `json:"last_message_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ConversationMessage is one message of a conversation.
type ConversationMessage struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	ConversationID string    `json:"conversation_id"`
	Direction      string    `json:"direction"`
	Actor          string    `json:"actor"`
	Body           string    `json:"body"`
	MediaURL       string    `json:"media_url,omitempty"`
	TwilioSID      string    `json:"twilio_sid,omitempty"`
	UserID         string    `json:"user_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ConversationStartRequest is POST /customers/{customerId}/conversations.
type ConversationStartRequest struct {
	Body    string `json:"body,omitempty"`
	Channel string `json:"channel,omitempty"`
}

// ConversationStartResult is the reply of StartConversation.
type ConversationStartResult struct {
	Conversation Conversation        `json:"conversation"`
	Message      ConversationMessage `json:"message"`
}

// MessageSendRequest is POST /conversations/{conversationId}/messages.
type MessageSendRequest struct {
	Body    string `json:"body,omitempty"`
	Channel string `json:"channel,omitempty"`
}

// CommunicationsService calls the /partner/v1/communications, conversations
// and customer conversation routes.
type CommunicationsService struct{ c *Client }

// Settings returns the messaging configuration. Scope
// partner-api.communications.read.
func (s *CommunicationsService) Settings(ctx context.Context) (*CommunicationSettings, error) {
	return call[CommunicationSettings](ctx, s.c, http.MethodGet, "/partner/v1/communications/settings", nil, nil)
}

// Numbers lists the tenant's phone numbers (not paginated). Scope
// partner-api.communications.read.
func (s *CommunicationsService) Numbers(ctx context.Context) ([]CommunicationNumber, error) {
	return listKey[CommunicationNumber](ctx, s.c, http.MethodGet, "/partner/v1/communications/numbers", "numbers", nil)
}

// Inbox lists the inbox conversations (not paginated). Scope
// partner-api.communications.read.
func (s *CommunicationsService) Inbox(ctx context.Context) ([]Conversation, error) {
	return listKey[Conversation](ctx, s.c, http.MethodGet, "/partner/v1/communications/inbox", "conversations", nil)
}

// ListConversations lists a customer's conversations (not paginated). Scope
// partner-api.communications.read.
func (s *CommunicationsService) ListConversations(ctx context.Context, customerID string) ([]Conversation, error) {
	return listKey[Conversation](ctx, s.c, http.MethodGet, "/partner/v1/customers/"+seg(customerID)+"/conversations", "conversations", nil)
}

// StartConversation opens a conversation with a customer and sends the first
// message. Scope partner-api.communications.write.
func (s *CommunicationsService) StartConversation(ctx context.Context, customerID string, body ConversationStartRequest, opts ...RequestOption) (*ConversationStartResult, error) {
	return call[ConversationStartResult](ctx, s.c, http.MethodPost, "/partner/v1/customers/"+seg(customerID)+"/conversations", nil, body, opts...)
}

// ListMessages lists the messages of a conversation (not paginated). Scope
// partner-api.communications.read.
func (s *CommunicationsService) ListMessages(ctx context.Context, conversationID string) ([]ConversationMessage, error) {
	return listKey[ConversationMessage](ctx, s.c, http.MethodGet, "/partner/v1/conversations/"+seg(conversationID)+"/messages", "messages", nil)
}

// SendMessage sends a message in a conversation. Scope
// partner-api.communications.write.
func (s *CommunicationsService) SendMessage(ctx context.Context, conversationID string, body MessageSendRequest, opts ...RequestOption) (*ConversationMessage, error) {
	return call[ConversationMessage](ctx, s.c, http.MethodPost, "/partner/v1/conversations/"+seg(conversationID)+"/messages", nil, body, opts...)
}
