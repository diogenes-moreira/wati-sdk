package wati

import (
	"context"
	"io"
	
	"github.com/tu-usuario/go-wati/contacts"
	"github.com/tu-usuario/go-wati/messages"
	"github.com/tu-usuario/go-wati/chatbots"
)

// ContactsService define la interfaz para el servicio de contactos
type ContactsService interface {
	// CRUD operations
	GetContacts(ctx context.Context, params *contacts.GetContactsParams) (*contacts.ContactsResponse, error)
	GetContact(ctx context.Context, id string) (*contacts.Contact, error)
	AddContact(ctx context.Context, contact *contacts.CreateContactRequest) (*contacts.Contact, error)
	UpdateContact(ctx context.Context, id string, contact *contacts.UpdateContactRequest) (*contacts.Contact, error)
	DeleteContact(ctx context.Context, id string) error
	
	// Búsqueda y filtrado
	SearchContacts(ctx context.Context, query string) (*contacts.ContactsResponse, error)
	FilterContacts(ctx context.Context, filter *contacts.ContactFilter) (*contacts.ContactsResponse, error)
	
	// Operaciones en lote
	AddContacts(ctx context.Context, contacts []*contacts.CreateContactRequest) (*contacts.BulkContactResponse, error)
}

// MessagesService define la interfaz para el servicio de mensajes
type MessagesService interface {
	// Mensajes de plantilla
	SendTemplateMessage(ctx context.Context, req *messages.SendTemplateMessageRequest) (*messages.MessageResponse, error)
	SendTemplateMessages(ctx context.Context, req *messages.SendTemplateMessagesRequest) (*messages.BulkMessageResponse, error)
	
	// Mensajes interactivos
	SendInteractiveListMessage(ctx context.Context, req *messages.InteractiveListMessageRequest) (*messages.MessageResponse, error)
	SendInteractiveButtonMessage(ctx context.Context, req *messages.InteractiveButtonMessageRequest) (*messages.MessageResponse, error)
	
	// Gestión de plantillas
	GetMessageTemplates(ctx context.Context) (*messages.TemplatesResponse, error)
	GetMessageTemplate(ctx context.Context, name string) (*messages.Template, error)
	
	// Historial de mensajes
	GetMessages(ctx context.Context, params *messages.GetMessagesParams) (*messages.MessagesResponse, error)
	GetMessage(ctx context.Context, id string) (*messages.Message, error)
	
	// Estado de mensajes
	GetMessageStatus(ctx context.Context, id string) (*MessageStatus, error)
}

// ChatbotsService define la interfaz para el servicio de chatbots
type ChatbotsService interface {
	GetChatbots(ctx context.Context) (*chatbots.ChatbotsResponse, error)
	GetChatbot(ctx context.Context, id string) (*chatbots.Chatbot, error)
	StartChatbot(ctx context.Context, req *chatbots.StartChatbotRequest) (*chatbots.ChatbotResponse, error)
	StopChatbot(ctx context.Context, id string) error
	UpdateChatStatus(ctx context.Context, req *chatbots.UpdateChatStatusRequest) (*chatbots.ChatStatusResponse, error)
}

// MediaService define la interfaz para el servicio de media
type MediaService interface {
	GetMediaByFileName(ctx context.Context, fileName string) (*MediaResponse, error)
	UploadMedia(ctx context.Context, file io.Reader, fileName string, mediaType string) (*UploadResponse, error)
	DeleteMedia(ctx context.Context, fileName string) error
	GetMediaURL(ctx context.Context, fileName string) (string, error)
}

// WebhooksService define la interfaz para el servicio de webhooks
type WebhooksService interface {
	// Configuración de webhooks
	RegisterWebhook(ctx context.Context, url string, events []WebhookEventType) error
	UnregisterWebhook(ctx context.Context, url string) error
	ListWebhooks(ctx context.Context) (*WebhooksResponse, error)
	
	// Manejo de eventos
	HandleWebhook(payload []byte, signature string) (*WebhookEvent, error)
	ValidateWebhookSignature(payload []byte, signature string) bool
	
	// Servidor de webhooks
	StartWebhookServer(port int, handlers map[WebhookEventType]WebhookHandler) error
	StopWebhookServer() error
}

