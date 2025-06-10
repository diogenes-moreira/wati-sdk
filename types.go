package wati

import (
	"time"
)

// BaseResponse representa la respuesta base de la API de WATI
type BaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	BaseResponse
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
	TotalCount int `json:"totalCount"`
}

// CustomParam representa un parámetro personalizado
type CustomParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Parameter representa un parámetro de plantilla
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MediaInfo representa información de media
type MediaInfo struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
}

// TemplateInfo representa información de una plantilla
type TemplateInfo struct {
	Name       string      `json:"name"`
	Language   string      `json:"language"`
	Parameters []Parameter `json:"parameters,omitempty"`
}

// InteractiveInfo representa información de mensaje interactivo
type InteractiveInfo struct {
	Type    string      `json:"type"`
	Header  interface{} `json:"header,omitempty"`
	Body    interface{} `json:"body,omitempty"`
	Footer  interface{} `json:"footer,omitempty"`
	Action  interface{} `json:"action,omitempty"`
}

// TokenResponse representa la respuesta de rotación de token
type TokenResponse struct {
	BaseResponse
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// WebhookEventType representa el tipo de evento de webhook
type WebhookEventType string

const (
	MessageReceived       WebhookEventType = "message_received"
	NewContactMessage     WebhookEventType = "new_contact_message"
	SessionMessageSent    WebhookEventType = "session_message_sent"
	TemplateMessageSent   WebhookEventType = "template_message_sent"
	MessageDelivered      WebhookEventType = "message_delivered"
	MessageRead           WebhookEventType = "message_read"
	MessageReplied        WebhookEventType = "message_replied"
	TemplateMessageFailed WebhookEventType = "template_message_failed"
)

// WebhookEvent representa un evento de webhook
type WebhookEvent struct {
	Type      WebhookEventType `json:"type"`
	Timestamp string           `json:"timestamp"`
	Data      interface{}      `json:"data"`
}

// WebhookHandler es una función que maneja eventos de webhook
type WebhookHandler func(event *WebhookEvent) error

// MessageStatus representa el estado de un mensaje
type MessageStatus struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Error     string `json:"error,omitempty"`
}

