package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// HTTPClient define la interfaz para realizar peticiones HTTP
type HTTPClient interface {
	DoRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error
}

// Service implementa WebhooksService
type Service struct {
	client HTTPClient
	server *WebhookServer
	mutex  sync.RWMutex
}

// NewService crea una nueva instancia del servicio de webhooks
func NewService(client HTTPClient) *Service {
	return &Service{
		client: client,
		server: &WebhookServer{
			Handlers:  make(map[WebhookEventType]WebhookHandler),
			IsRunning: false,
		},
	}
}

// RegisterWebhook registra un webhook en WATI
func (s *Service) RegisterWebhook(ctx context.Context, url string, events []WebhookEventType) error {
	registration := &WebhookRegistration{
		URL:    url,
		Events: events,
	}
	
	return s.RegisterWebhookWithConfig(ctx, registration)
}

// RegisterWebhookWithConfig registra un webhook con configuración completa
func (s *Service) RegisterWebhookWithConfig(ctx context.Context, config *WebhookRegistration) error {
	if config == nil {
		return fmt.Errorf("webhook configuration is required")
	}
	
	if err := config.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	
	var response BaseResponse
	err := s.client.DoRequest(ctx, "POST", "/api/v1/webhooks", config, &response)
	if err != nil {
		return fmt.Errorf("error registering webhook: %w", err)
	}
	
	return nil
}

// UnregisterWebhook desregistra un webhook
func (s *Service) UnregisterWebhook(ctx context.Context, url string) error {
	if url == "" {
		return fmt.Errorf("webhook URL is required")
	}
	
	requestBody := struct {
		URL string `json:"url"`
	}{
		URL: url,
	}
	
	var response BaseResponse
	err := s.client.DoRequest(ctx, "DELETE", "/api/v1/webhooks", requestBody, &response)
	if err != nil {
		return fmt.Errorf("error unregistering webhook: %w", err)
	}
	
	return nil
}

// ListWebhooks obtiene la lista de webhooks registrados
func (s *Service) ListWebhooks(ctx context.Context) (*WebhooksResponse, error) {
	var response WebhooksResponse
	err := s.client.DoRequest(ctx, "GET", "/api/v1/webhooks", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error listing webhooks: %w", err)
	}
	
	return &response, nil
}

// HandleWebhook procesa un evento de webhook
func (s *Service) HandleWebhook(payload []byte, signature string) (*WebhookEvent, error) {
	// Parsear el evento
	event, err := ParseWebhookEvent(payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing webhook event: %w", err)
	}
	
	// Validar firma si hay un secreto configurado
	s.mutex.RLock()
	secret := s.server.Secret
	s.mutex.RUnlock()
	
	if !ValidateSignature(payload, signature, secret) {
		return nil, fmt.Errorf("invalid webhook signature")
	}
	
	// Ejecutar handler si existe
	s.mutex.RLock()
	handler, exists := s.server.Handlers[event.Type]
	s.mutex.RUnlock()
	
	if exists && handler != nil {
		if err := handler(event); err != nil {
			return event, fmt.Errorf("error executing webhook handler: %w", err)
		}
	}
	
	return event, nil
}

// ValidateWebhookSignature valida la firma de un webhook
func (s *Service) ValidateWebhookSignature(payload []byte, signature string) bool {
	s.mutex.RLock()
	secret := s.server.Secret
	s.mutex.RUnlock()
	
	return ValidateSignature(payload, signature, secret)
}

// StartWebhookServer inicia el servidor de webhooks
func (s *Service) StartWebhookServer(port int, handlers map[WebhookEventType]WebhookHandler) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if s.server.IsRunning {
		return fmt.Errorf("webhook server is already running")
	}
	
	// Configurar handlers
	if handlers != nil {
		s.server.Handlers = handlers
	}
	
	s.server.Port = port
	
	// Crear servidor HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", s.handleWebhookRequest)
	mux.HandleFunc("/health", s.handleHealthCheck)
	
	s.server.server = &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Iniciar servidor en goroutine
	go func() {
		log.Printf("Starting webhook server on port %d", port)
		if err := s.server.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Webhook server error: %v", err)
		}
	}()
	
	s.server.IsRunning = true
	return nil
}

// StopWebhookServer detiene el servidor de webhooks
func (s *Service) StopWebhookServer() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if !s.server.IsRunning {
		return fmt.Errorf("webhook server is not running")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := s.server.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping webhook server: %w", err)
	}
	
	s.server.IsRunning = false
	log.Println("Webhook server stopped")
	return nil
}

// RegisterHandler registra un handler para un tipo de evento específico
func (s *Service) RegisterHandler(eventType WebhookEventType, handler WebhookHandler) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.server.Handlers[eventType] = handler
}

// UnregisterHandler desregistra un handler
func (s *Service) UnregisterHandler(eventType WebhookEventType) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	delete(s.server.Handlers, eventType)
}

// SetSecret establece el secreto para validación de firmas
func (s *Service) SetSecret(secret string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.server.Secret = secret
}

// GetServerStatus obtiene el estado del servidor
func (s *Service) GetServerStatus() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	return s.server.IsRunning
}

// GetServerPort obtiene el puerto del servidor
func (s *Service) GetServerPort() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	return s.server.Port
}

// handleWebhookRequest maneja las peticiones de webhook
func (s *Service) handleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Leer el cuerpo de la petición
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading webhook body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	// Obtener firma del header
	signature := r.Header.Get("X-Webhook-Signature")
	if signature == "" {
		signature = r.Header.Get("X-Hub-Signature-256")
	}
	
	// Procesar webhook
	event, err := s.HandleWebhook(body, signature)
	if err != nil {
		log.Printf("Error handling webhook: %v", err)
		http.Error(w, "Error processing webhook", http.StatusBadRequest)
		return
	}
	
	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"status":    "success",
		"eventId":   event.ID,
		"eventType": event.Type,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleHealthCheck maneja las peticiones de health check
func (s *Service) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"server": map[string]interface{}{
			"port":      s.GetServerPort(),
			"running":   s.GetServerStatus(),
			"handlers":  len(s.server.Handlers),
		},
	}
	
	json.NewEncoder(w).Encode(response)
}

// CreateMessageHandler crea un handler para mensajes recibidos
func CreateMessageHandler(handler func(data MessageReceivedData) error) WebhookHandler {
	return func(event *WebhookEvent) error {
		if data, ok := event.Data.(MessageReceivedData); ok {
			return handler(data)
		}
		return fmt.Errorf("invalid data type for message event")
	}
}

// CreateMessageStatusHandler crea un handler para cambios de estado de mensaje
func CreateMessageStatusHandler(handler func(data MessageStatusData) error) WebhookHandler {
	return func(event *WebhookEvent) error {
		if data, ok := event.Data.(MessageStatusData); ok {
			return handler(data)
		}
		return fmt.Errorf("invalid data type for message status event")
	}
}

// CreateContactHandler crea un handler para eventos de contacto
func CreateContactHandler(handler func(data ContactEventData) error) WebhookHandler {
	return func(event *WebhookEvent) error {
		if data, ok := event.Data.(ContactEventData); ok {
			return handler(data)
		}
		return fmt.Errorf("invalid data type for contact event")
	}
}

// CreateChatbotHandler crea un handler para eventos de chatbot
func CreateChatbotHandler(handler func(data ChatbotEventData) error) WebhookHandler {
	return func(event *WebhookEvent) error {
		if data, ok := event.Data.(ChatbotEventData); ok {
			return handler(data)
		}
		return fmt.Errorf("invalid data type for chatbot event")
	}
}

// CreateChatStatusHandler crea un handler para cambios de estado de chat
func CreateChatStatusHandler(handler func(data ChatStatusEventData) error) WebhookHandler {
	return func(event *WebhookEvent) error {
		if data, ok := event.Data.(ChatStatusEventData); ok {
			return handler(data)
		}
		return fmt.Errorf("invalid data type for chat status event")
	}
}

// RegisterMessageHandlers registra handlers comunes para mensajes
func (s *Service) RegisterMessageHandlers(
	onMessageReceived func(MessageReceivedData) error,
	onMessageDelivered func(MessageStatusData) error,
	onMessageRead func(MessageStatusData) error,
) {
	if onMessageReceived != nil {
		s.RegisterHandler(MessageReceived, CreateMessageHandler(onMessageReceived))
		s.RegisterHandler(NewContactMessage, CreateMessageHandler(onMessageReceived))
	}
	
	if onMessageDelivered != nil {
		s.RegisterHandler(MessageDelivered, CreateMessageStatusHandler(onMessageDelivered))
	}
	
	if onMessageRead != nil {
		s.RegisterHandler(MessageRead, CreateMessageStatusHandler(onMessageRead))
	}
}

// RegisterAllEventHandlers registra un handler genérico para todos los eventos
func (s *Service) RegisterAllEventHandlers(handler WebhookHandler) {
	events := []WebhookEventType{
		MessageReceived,
		NewContactMessage,
		SessionMessageSent,
		TemplateMessageSent,
		MessageDelivered,
		MessageRead,
		MessageReplied,
		TemplateMessageFailed,
		ContactCreated,
		ContactUpdated,
		ChatbotStarted,
		ChatbotStopped,
		ChatStatusChanged,
	}
	
	for _, eventType := range events {
		s.RegisterHandler(eventType, handler)
	}
}

// TestWebhook envía un evento de prueba al webhook
func (s *Service) TestWebhook(ctx context.Context, webhookURL string) error {
	testEvent := &WebhookEvent{
		ID:        "test-" + strconv.FormatInt(time.Now().Unix(), 10),
		Type:      MessageReceived,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data: MessageReceivedData{
			MessageID:   "test-message-id",
			From:        "1234567890",
			To:          "0987654321",
			MessageType: "text",
			Text:        "This is a test message from WATI webhook",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		},
		Source:  "wati-webhook-test",
		Version: "1.0",
	}
	
	// Convertir a JSON
	payload, err := json.Marshal(testEvent)
	if err != nil {
		return fmt.Errorf("error marshaling test event: %w", err)
	}
	
	// Enviar petición HTTP al webhook
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("error sending test webhook: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook test failed with status: %d", resp.StatusCode)
	}
	
	return nil
}

