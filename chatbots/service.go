package chatbots

import (
	"context"
	"fmt"
)

// HTTPClient define la interfaz para realizar peticiones HTTP
type HTTPClient interface {
	DoRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error
}

// Service implementa ChatbotsService
type Service struct {
	client HTTPClient
}

// NewService crea una nueva instancia del servicio de chatbots
func NewService(client HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// GetChatbots obtiene la lista de todos los chatbots
func (s *Service) GetChatbots(ctx context.Context) (*ChatbotsResponse, error) {
	var response ChatbotsResponse
	err := s.client.DoRequest(ctx, "GET", "/api/v1/chatbots", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting chatbots: %w", err)
	}
	
	return &response, nil
}

// GetChatbot obtiene un chatbot específico por ID
func (s *Service) GetChatbot(ctx context.Context, id string) (*Chatbot, error) {
	if id == "" {
		return nil, fmt.Errorf("chatbot ID is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/chatbots/%s", id)
	
	var response struct {
		BaseResponse
		Chatbot Chatbot `json:"chatbot"`
	}
	
	err := s.client.DoRequest(ctx, "GET", endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting chatbot %s: %w", id, err)
	}
	
	return &response.Chatbot, nil
}

// StartChatbot inicia un chatbot para un contacto específico
func (s *Service) StartChatbot(ctx context.Context, req *StartChatbotRequest) (*ChatbotResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	var response ChatbotResponse
	err := s.client.DoRequest(ctx, "POST", "/api/v1/startChatbot", req, &response)
	if err != nil {
		return nil, fmt.Errorf("error starting chatbot: %w", err)
	}
	
	return &response, nil
}

// StopChatbot detiene un chatbot para un contacto específico
func (s *Service) StopChatbot(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("chatbot ID is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/stopChatbot/%s", id)
	
	var response BaseResponse
	err := s.client.DoRequest(ctx, "POST", endpoint, nil, &response)
	if err != nil {
		return fmt.Errorf("error stopping chatbot %s: %w", id, err)
	}
	
	return nil
}

// UpdateChatStatus actualiza el estado de un chat
func (s *Service) UpdateChatStatus(ctx context.Context, req *UpdateChatStatusRequest) (*ChatStatusResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	var response ChatStatusResponse
	err := s.client.DoRequest(ctx, "POST", "/api/v1/updateChatStatus", req, &response)
	if err != nil {
		return nil, fmt.Errorf("error updating chat status: %w", err)
	}
	
	return &response, nil
}

// CreateChatbot crea un nuevo chatbot
func (s *Service) CreateChatbot(ctx context.Context, req *CreateChatbotRequest) (*Chatbot, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	var response struct {
		BaseResponse
		Chatbot Chatbot `json:"chatbot"`
	}
	
	err := s.client.DoRequest(ctx, "POST", "/api/v1/chatbots", req, &response)
	if err != nil {
		return nil, fmt.Errorf("error creating chatbot: %w", err)
	}
	
	return &response.Chatbot, nil
}

// UpdateChatbot actualiza un chatbot existente
func (s *Service) UpdateChatbot(ctx context.Context, id string, req *UpdateChatbotRequest) (*Chatbot, error) {
	if id == "" {
		return nil, fmt.Errorf("chatbot ID is required")
	}
	
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/chatbots/%s", id)
	
	var response struct {
		BaseResponse
		Chatbot Chatbot `json:"chatbot"`
	}
	
	err := s.client.DoRequest(ctx, "PUT", endpoint, req, &response)
	if err != nil {
		return nil, fmt.Errorf("error updating chatbot %s: %w", id, err)
	}
	
	return &response.Chatbot, nil
}

// DeleteChatbot elimina un chatbot
func (s *Service) DeleteChatbot(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("chatbot ID is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/chatbots/%s", id)
	
	var response BaseResponse
	err := s.client.DoRequest(ctx, "DELETE", endpoint, nil, &response)
	if err != nil {
		return fmt.Errorf("error deleting chatbot %s: %w", id, err)
	}
	
	return nil
}

// GetActiveChatbots obtiene solo los chatbots activos
func (s *Service) GetActiveChatbots(ctx context.Context) ([]Chatbot, error) {
	response, err := s.GetChatbots(ctx)
	if err != nil {
		return nil, err
	}
	
	var activeChatbots []Chatbot
	for _, chatbot := range response.Chatbots {
		if chatbot.IsActive() {
			activeChatbots = append(activeChatbots, chatbot)
		}
	}
	
	return activeChatbots, nil
}

// ActivateChatbot activa un chatbot
func (s *Service) ActivateChatbot(ctx context.Context, id string) (*Chatbot, error) {
	isActive := true
	req := &UpdateChatbotRequest{
		IsActive: &isActive,
	}
	
	return s.UpdateChatbot(ctx, id, req)
}

// DeactivateChatbot desactiva un chatbot
func (s *Service) DeactivateChatbot(ctx context.Context, id string) (*Chatbot, error) {
	isActive := false
	req := &UpdateChatbotRequest{
		IsActive: &isActive,
	}
	
	return s.UpdateChatbot(ctx, id, req)
}

// StartChatbotForContact inicia un chatbot para un contacto específico
func (s *Service) StartChatbotForContact(ctx context.Context, chatbotID, whatsappNumber string) (*ChatbotResponse, error) {
	req := &StartChatbotRequest{
		ChatbotID:      chatbotID,
		WhatsappNumber: whatsappNumber,
	}
	
	return s.StartChatbot(ctx, req)
}

// StartChatbotWithMessage inicia un chatbot con un mensaje inicial
func (s *Service) StartChatbotWithMessage(ctx context.Context, chatbotID, whatsappNumber, initialMessage string) (*ChatbotResponse, error) {
	req := &StartChatbotRequest{
		ChatbotID:      chatbotID,
		WhatsappNumber: whatsappNumber,
		InitialMessage: initialMessage,
	}
	
	return s.StartChatbot(ctx, req)
}

// AssignChatToUser asigna un chat a un usuario específico
func (s *Service) AssignChatToUser(ctx context.Context, whatsappNumber, userID string) (*ChatStatusResponse, error) {
	req := &UpdateChatStatusRequest{
		WhatsappNumber: whatsappNumber,
		Status:         string(ChatStatusAssigned),
		AssignedTo:     userID,
	}
	
	return s.UpdateChatStatus(ctx, req)
}

// TransferChatToHuman transfiere un chat de bot a humano
func (s *Service) TransferChatToHuman(ctx context.Context, whatsappNumber, userID string, notes string) (*ChatStatusResponse, error) {
	req := &UpdateChatStatusRequest{
		WhatsappNumber: whatsappNumber,
		Status:         string(ChatStatusAssigned),
		AssignedTo:     userID,
		Notes:          notes,
	}
	
	return s.UpdateChatStatus(ctx, req)
}

// CloseChatSession cierra una sesión de chat
func (s *Service) CloseChatSession(ctx context.Context, whatsappNumber string, notes string) (*ChatStatusResponse, error) {
	req := &UpdateChatStatusRequest{
		WhatsappNumber: whatsappNumber,
		Status:         string(ChatStatusClosed),
		Notes:          notes,
	}
	
	return s.UpdateChatStatus(ctx, req)
}

// ResolveChatSession marca un chat como resuelto
func (s *Service) ResolveChatSession(ctx context.Context, whatsappNumber string, notes string) (*ChatStatusResponse, error) {
	req := &UpdateChatStatusRequest{
		WhatsappNumber: whatsappNumber,
		Status:         string(ChatStatusResolved),
		Notes:          notes,
	}
	
	return s.UpdateChatStatus(ctx, req)
}

// AddTagsToChat agrega etiquetas a un chat
func (s *Service) AddTagsToChat(ctx context.Context, whatsappNumber string, tags []string) (*ChatStatusResponse, error) {
	// Primero obtener el estado actual para no sobrescribir otros campos
	req := &UpdateChatStatusRequest{
		WhatsappNumber: whatsappNumber,
		Status:         string(ChatStatusOpen), // Mantener estado actual o usar uno por defecto
		Tags:           tags,
	}
	
	return s.UpdateChatStatus(ctx, req)
}

// GetChatbotByName busca un chatbot por nombre
func (s *Service) GetChatbotByName(ctx context.Context, name string) (*Chatbot, error) {
	if name == "" {
		return nil, fmt.Errorf("chatbot name is required")
	}
	
	response, err := s.GetChatbots(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting chatbots: %w", err)
	}
	
	for _, chatbot := range response.Chatbots {
		if chatbot.Name == name {
			return &chatbot, nil
		}
	}
	
	return nil, fmt.Errorf("chatbot with name '%s' not found", name)
}

// GetChatbotsByKeyword busca chatbots que contengan una palabra clave específica
func (s *Service) GetChatbotsByKeyword(ctx context.Context, keyword string) ([]Chatbot, error) {
	if keyword == "" {
		return nil, fmt.Errorf("keyword is required")
	}
	
	response, err := s.GetChatbots(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting chatbots: %w", err)
	}
	
	var matchingChatbots []Chatbot
	for _, chatbot := range response.Chatbots {
		if chatbot.HasKeyword(keyword) {
			matchingChatbots = append(matchingChatbots, chatbot)
		}
	}
	
	return matchingChatbots, nil
}

// UpdateChatbotKeywords actualiza las palabras clave de un chatbot
func (s *Service) UpdateChatbotKeywords(ctx context.Context, id string, keywords []string) (*Chatbot, error) {
	req := &UpdateChatbotRequest{
		Keywords: keywords,
	}
	
	return s.UpdateChatbot(ctx, id, req)
}

// UpdateChatbotResponses actualiza las respuestas de un chatbot
func (s *Service) UpdateChatbotResponses(ctx context.Context, id string, responses []Response) (*Chatbot, error) {
	req := &UpdateChatbotRequest{
		Responses: responses,
	}
	
	return s.UpdateChatbot(ctx, id, req)
}

