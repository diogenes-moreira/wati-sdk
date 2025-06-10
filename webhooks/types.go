package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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
	ContactCreated        WebhookEventType = "contact_created"
	ContactUpdated        WebhookEventType = "contact_updated"
	ChatbotStarted        WebhookEventType = "chatbot_started"
	ChatbotStopped        WebhookEventType = "chatbot_stopped"
	ChatStatusChanged     WebhookEventType = "chat_status_changed"
)

// WebhookEvent representa un evento de webhook
type WebhookEvent struct {
	ID        string           `json:"id"`
	Type      WebhookEventType `json:"type"`
	Timestamp string           `json:"timestamp"`
	Data      interface{}      `json:"data"`
	Source    string           `json:"source,omitempty"`
	Version   string           `json:"version,omitempty"`
}

// WebhookHandler es una función que maneja eventos de webhook
type WebhookHandler func(event *WebhookEvent) error

// MessageReceivedData representa los datos de un mensaje recibido
type MessageReceivedData struct {
	MessageID      string                 `json:"messageId"`
	From           string                 `json:"from"`
	To             string                 `json:"to"`
	MessageType    string                 `json:"messageType"`
	Text           string                 `json:"text,omitempty"`
	Media          *WebhookMediaInfo      `json:"media,omitempty"`
	Location       *WebhookLocationInfo   `json:"location,omitempty"`
	Contact        *WebhookContactInfo    `json:"contact,omitempty"`
	Interactive    *WebhookInteractiveInfo `json:"interactive,omitempty"`
	Timestamp      string                 `json:"timestamp"`
	ContactProfile *WebhookContactProfile `json:"contactProfile,omitempty"`
}

// MessageSentData representa los datos de un mensaje enviado
type MessageSentData struct {
	MessageID     string                 `json:"messageId"`
	From          string                 `json:"from"`
	To            string                 `json:"to"`
	MessageType   string                 `json:"messageType"`
	TemplateName  string                 `json:"templateName,omitempty"`
	Status        string                 `json:"status"`
	Timestamp     string                 `json:"timestamp"`
	ErrorCode     string                 `json:"errorCode,omitempty"`
	ErrorMessage  string                 `json:"errorMessage,omitempty"`
	ContactProfile *WebhookContactProfile `json:"contactProfile,omitempty"`
}

// MessageStatusData representa los datos de cambio de estado de mensaje
type MessageStatusData struct {
	MessageID   string `json:"messageId"`
	From        string `json:"from"`
	To          string `json:"to"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp"`
	ErrorCode   string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// ContactEventData representa los datos de eventos de contacto
type ContactEventData struct {
	ContactID      string                 `json:"contactId"`
	WhatsappNumber string                 `json:"whatsappNumber"`
	FirstName      string                 `json:"firstName"`
	LastName       string                 `json:"lastName,omitempty"`
	FullName       string                 `json:"fullName"`
	Email          string                 `json:"email,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	CustomParams   []WebhookCustomParam   `json:"customParams,omitempty"`
	Source         string                 `json:"source,omitempty"`
	Timestamp      string                 `json:"timestamp"`
	Changes        map[string]interface{} `json:"changes,omitempty"`
}

// ChatbotEventData representa los datos de eventos de chatbot
type ChatbotEventData struct {
	ChatbotID      string `json:"chatbotId"`
	ChatbotName    string `json:"chatbotName"`
	WhatsappNumber string `json:"whatsappNumber"`
	SessionID      string `json:"sessionId,omitempty"`
	Status         string `json:"status"`
	Timestamp      string `json:"timestamp"`
	Reason         string `json:"reason,omitempty"`
}

// ChatStatusEventData representa los datos de cambio de estado de chat
type ChatStatusEventData struct {
	WhatsappNumber string `json:"whatsappNumber"`
	OldStatus      string `json:"oldStatus"`
	NewStatus      string `json:"newStatus"`
	AssignedTo     string `json:"assignedTo,omitempty"`
	AssignedBy     string `json:"assignedBy,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Notes          string `json:"notes,omitempty"`
	Timestamp      string `json:"timestamp"`
}

// WebhookMediaInfo representa información de media en webhook
type WebhookMediaInfo struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
	Caption  string `json:"caption,omitempty"`
}

// WebhookLocationInfo representa información de ubicación
type WebhookLocationInfo struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// WebhookContactInfo representa información de contacto compartido
type WebhookContactInfo struct {
	Name         string                `json:"name"`
	PhoneNumbers []WebhookPhoneNumber  `json:"phoneNumbers,omitempty"`
	Emails       []WebhookEmail        `json:"emails,omitempty"`
	URLs         []WebhookURL          `json:"urls,omitempty"`
	Addresses    []WebhookAddress      `json:"addresses,omitempty"`
	Organization string                `json:"organization,omitempty"`
	Birthday     string                `json:"birthday,omitempty"`
}

// WebhookInteractiveInfo representa información de mensaje interactivo
type WebhookInteractiveInfo struct {
	Type        string                      `json:"type"`
	ButtonReply *WebhookButtonReply         `json:"buttonReply,omitempty"`
	ListReply   *WebhookListReply           `json:"listReply,omitempty"`
}

// WebhookContactProfile representa el perfil del contacto
type WebhookContactProfile struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

// WebhookCustomParam representa un parámetro personalizado
type WebhookCustomParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// WebhookPhoneNumber representa un número de teléfono
type WebhookPhoneNumber struct {
	Phone string `json:"phone"`
	Type  string `json:"type,omitempty"`
}

// WebhookEmail representa un email
type WebhookEmail struct {
	Email string `json:"email"`
	Type  string `json:"type,omitempty"`
}

// WebhookURL representa una URL
type WebhookURL struct {
	URL  string `json:"url"`
	Type string `json:"type,omitempty"`
}

// WebhookAddress representa una dirección
type WebhookAddress struct {
	Street      string `json:"street,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Zip         string `json:"zip,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	Type        string `json:"type,omitempty"`
}

// WebhookButtonReply representa la respuesta de un botón
type WebhookButtonReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// WebhookListReply representa la respuesta de una lista
type WebhookListReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// WebhookConfig representa la configuración de un webhook
type WebhookConfig struct {
	URL         string             `json:"url"`
	Events      []WebhookEventType `json:"events"`
	Secret      string             `json:"secret,omitempty"`
	IsActive    bool               `json:"isActive"`
	Description string             `json:"description,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

// WebhookRegistration representa una petición de registro de webhook
type WebhookRegistration struct {
	URL         string             `json:"url"`
	Events      []WebhookEventType `json:"events"`
	Secret      string             `json:"secret,omitempty"`
	Description string             `json:"description,omitempty"`
}

// WebhooksResponse representa la respuesta de lista de webhooks
type WebhooksResponse struct {
	BaseResponse
	Webhooks []WebhookConfig `json:"webhooks"`
}

// WebhookServer representa un servidor de webhooks
type WebhookServer struct {
	Port     int                                    `json:"port"`
	Handlers map[WebhookEventType]WebhookHandler   `json:"-"`
	Secret   string                                 `json:"secret,omitempty"`
	server   *http.Server                          `json:"-"`
	IsRunning bool                                  `json:"isRunning"`
}

// BaseResponse representa la respuesta base de la API
type BaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Validate valida la configuración de registro de webhook
func (r *WebhookRegistration) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("webhook URL is required")
	}
	
	if len(r.Events) == 0 {
		return fmt.Errorf("at least one event type is required")
	}
	
	// Validar que los tipos de evento sean válidos
	validEvents := map[WebhookEventType]bool{
		MessageReceived:       true,
		NewContactMessage:     true,
		SessionMessageSent:    true,
		TemplateMessageSent:   true,
		MessageDelivered:      true,
		MessageRead:           true,
		MessageReplied:        true,
		TemplateMessageFailed: true,
		ContactCreated:        true,
		ContactUpdated:        true,
		ChatbotStarted:        true,
		ChatbotStopped:        true,
		ChatStatusChanged:     true,
	}
	
	for _, event := range r.Events {
		if !validEvents[event] {
			return fmt.Errorf("invalid event type: %s", event)
		}
	}
	
	return nil
}

// ParseWebhookEvent parsea un evento de webhook desde JSON
func ParseWebhookEvent(payload []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("error parsing webhook event: %w", err)
	}
	
	// Parsear los datos específicos según el tipo de evento
	if err := parseEventData(&event); err != nil {
		return nil, fmt.Errorf("error parsing event data: %w", err)
	}
	
	return &event, nil
}

// parseEventData parsea los datos específicos del evento
func parseEventData(event *WebhookEvent) error {
	if event.Data == nil {
		return nil
	}
	
	// Convertir a JSON y luego al tipo específico
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	
	switch event.Type {
	case MessageReceived, NewContactMessage:
		var data MessageReceivedData
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return err
		}
		event.Data = data
		
	case SessionMessageSent, TemplateMessageSent, TemplateMessageFailed:
		var data MessageSentData
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return err
		}
		event.Data = data
		
	case MessageDelivered, MessageRead, MessageReplied:
		var data MessageStatusData
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return err
		}
		event.Data = data
		
	case ContactCreated, ContactUpdated:
		var data ContactEventData
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return err
		}
		event.Data = data
		
	case ChatbotStarted, ChatbotStopped:
		var data ChatbotEventData
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return err
		}
		event.Data = data
		
	case ChatStatusChanged:
		var data ChatStatusEventData
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return err
		}
		event.Data = data
	}
	
	return nil
}

// ValidateSignature valida la firma de un webhook
func ValidateSignature(payload []byte, signature string, secret string) bool {
	if secret == "" {
		return true // Si no hay secreto configurado, no validamos
	}
	
	// Calcular HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	
	// Comparar firmas
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// GetMessageText extrae el texto de un mensaje recibido
func (d *MessageReceivedData) GetMessageText() string {
	switch d.MessageType {
	case "text":
		return d.Text
	case "interactive":
		if d.Interactive != nil {
			if d.Interactive.ButtonReply != nil {
				return d.Interactive.ButtonReply.Title
			}
			if d.Interactive.ListReply != nil {
				return d.Interactive.ListReply.Title
			}
		}
	}
	return ""
}

// IsTextMessage verifica si es un mensaje de texto
func (d *MessageReceivedData) IsTextMessage() bool {
	return d.MessageType == "text"
}

// IsMediaMessage verifica si es un mensaje de media
func (d *MessageReceivedData) IsMediaMessage() bool {
	return d.Media != nil
}

// IsLocationMessage verifica si es un mensaje de ubicación
func (d *MessageReceivedData) IsLocationMessage() bool {
	return d.Location != nil
}

// IsContactMessage verifica si es un mensaje de contacto
func (d *MessageReceivedData) IsContactMessage() bool {
	return d.Contact != nil
}

// IsInteractiveMessage verifica si es un mensaje interactivo
func (d *MessageReceivedData) IsInteractiveMessage() bool {
	return d.Interactive != nil
}

// IsButtonReply verifica si es una respuesta de botón
func (d *MessageReceivedData) IsButtonReply() bool {
	return d.Interactive != nil && d.Interactive.ButtonReply != nil
}

// IsListReply verifica si es una respuesta de lista
func (d *MessageReceivedData) IsListReply() bool {
	return d.Interactive != nil && d.Interactive.ListReply != nil
}

// GetContactName obtiene el nombre del contacto
func (d *MessageReceivedData) GetContactName() string {
	if d.ContactProfile != nil {
		return d.ContactProfile.Name
	}
	return ""
}

// IsDelivered verifica si el mensaje fue entregado
func (d *MessageSentData) IsDelivered() bool {
	return d.Status == "delivered"
}

// IsRead verifica si el mensaje fue leído
func (d *MessageSentData) IsRead() bool {
	return d.Status == "read"
}

// IsFailed verifica si el mensaje falló
func (d *MessageSentData) IsFailed() bool {
	return d.Status == "failed"
}

// HasError verifica si hay un error
func (d *MessageSentData) HasError() bool {
	return d.ErrorCode != "" || d.ErrorMessage != ""
}

// GetErrorInfo obtiene información del error
func (d *MessageSentData) GetErrorInfo() string {
	if d.ErrorMessage != "" {
		return d.ErrorMessage
	}
	if d.ErrorCode != "" {
		return fmt.Sprintf("Error code: %s", d.ErrorCode)
	}
	return ""
}

