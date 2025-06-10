package messages

import (
	"context"
	"fmt"
	"strings"
)

// HTTPClient define la interfaz para realizar peticiones HTTP
type HTTPClient interface {
	DoRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error
}

// Service implementa MessagesService
type Service struct {
	client HTTPClient
}

// NewService crea una nueva instancia del servicio de mensajes
func NewService(client HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// SendTemplateMessage envía un mensaje de plantilla a un contacto
func (s *Service) SendTemplateMessage(ctx context.Context, req *SendTemplateMessageRequest) (*MessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	var response MessageResponse
	err := s.client.DoRequest(ctx, "POST", "/api/v1/sendTemplateMessage", req, &response)
	if err != nil {
		return nil, fmt.Errorf("error sending template message: %w", err)
	}
	
	return &response, nil
}

// SendTemplateMessages envía mensajes de plantilla a múltiples contactos
func (s *Service) SendTemplateMessages(ctx context.Context, req *SendTemplateMessagesRequest) (*BulkMessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	var response BulkMessageResponse
	err := s.client.DoRequest(ctx, "POST", "/api/v1/sendTemplateMessages", req, &response)
	if err != nil {
		return nil, fmt.Errorf("error sending template messages: %w", err)
	}
	
	return &response, nil
}

// SendInteractiveListMessage envía un mensaje de lista interactiva
func (s *Service) SendInteractiveListMessage(ctx context.Context, req *InteractiveListMessageRequest) (*MessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	var response MessageResponse
	err := s.client.DoRequest(ctx, "POST", "/api/v1/sendInteractiveListMessage", req, &response)
	if err != nil {
		return nil, fmt.Errorf("error sending interactive list message: %w", err)
	}
	
	return &response, nil
}

// SendInteractiveButtonMessage envía un mensaje de botones interactivos
func (s *Service) SendInteractiveButtonMessage(ctx context.Context, req *InteractiveButtonMessageRequest) (*MessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	var response MessageResponse
	err := s.client.DoRequest(ctx, "POST", "/api/v1/sendInteractiveButtonMessage", req, &response)
	if err != nil {
		return nil, fmt.Errorf("error sending interactive button message: %w", err)
	}
	
	return &response, nil
}

// GetMessageTemplates obtiene todas las plantillas de mensajes disponibles
func (s *Service) GetMessageTemplates(ctx context.Context) (*TemplatesResponse, error) {
	var response TemplatesResponse
	err := s.client.DoRequest(ctx, "GET", "/api/v1/getMessageTemplates", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting message templates: %w", err)
	}
	
	return &response, nil
}

// GetMessageTemplate obtiene una plantilla específica por nombre
func (s *Service) GetMessageTemplate(ctx context.Context, name string) (*Template, error) {
	if name == "" {
		return nil, fmt.Errorf("template name is required")
	}
	
	// Obtener todas las plantillas y filtrar por nombre
	templates, err := s.GetMessageTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting templates: %w", err)
	}
	
	for _, template := range templates.Templates {
		if template.Name == name {
			return &template, nil
		}
	}
	
	return nil, fmt.Errorf("template '%s' not found", name)
}

// GetMessages obtiene el historial de mensajes con parámetros opcionales
func (s *Service) GetMessages(ctx context.Context, params *GetMessagesParams) (*MessagesResponse, error) {
	if params == nil {
		params = &GetMessagesParams{}
	}
	
	params.SetDefaults()
	
	// Construir endpoint con query parameters
	endpoint := "/api/v1/getMessages"
	queryParams := params.ToMap()
	
	if len(queryParams) > 0 {
		var parts []string
		for key, value := range queryParams {
			parts = append(parts, fmt.Sprintf("%s=%s", key, value))
		}
		endpoint += "?" + strings.Join(parts, "&")
	}
	
	var response MessagesResponse
	err := s.client.DoRequest(ctx, "GET", endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting messages: %w", err)
	}
	
	return &response, nil
}

// GetMessage obtiene un mensaje específico por ID
func (s *Service) GetMessage(ctx context.Context, id string) (*Message, error) {
	if id == "" {
		return nil, fmt.Errorf("message ID is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/getMessage/%s", id)
	
	var response struct {
		BaseResponse
		Message Message `json:"message"`
	}
	
	err := s.client.DoRequest(ctx, "GET", endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting message %s: %w", id, err)
	}
	
	return &response.Message, nil
}

// GetMessageStatus obtiene el estado de un mensaje específico
func (s *Service) GetMessageStatus(ctx context.Context, id string) (*MessageStatus, error) {
	if id == "" {
		return nil, fmt.Errorf("message ID is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/getMessageStatus/%s", id)
	
	var response struct {
		BaseResponse
		Status MessageStatus `json:"status"`
	}
	
	err := s.client.DoRequest(ctx, "GET", endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting message status %s: %w", id, err)
	}
	
	return &response.Status, nil
}

// GetMessagesByPhone obtiene mensajes de un número de teléfono específico
func (s *Service) GetMessagesByPhone(ctx context.Context, phone string, params *GetMessagesParams) (*MessagesResponse, error) {
	if phone == "" {
		return nil, fmt.Errorf("phone number is required")
	}
	
	if params == nil {
		params = &GetMessagesParams{}
	}
	
	params.Phone = phone
	return s.GetMessages(ctx, params)
}

// GetMessagesByDateRange obtiene mensajes en un rango de fechas
func (s *Service) GetMessagesByDateRange(ctx context.Context, fromDate, toDate string, params *GetMessagesParams) (*MessagesResponse, error) {
	if fromDate == "" || toDate == "" {
		return nil, fmt.Errorf("both fromDate and toDate are required")
	}
	
	if params == nil {
		params = &GetMessagesParams{}
	}
	
	params.FromDate = fromDate
	params.ToDate = toDate
	return s.GetMessages(ctx, params)
}

// SendSimpleTemplateMessage envía un mensaje de plantilla simple sin parámetros
func (s *Service) SendSimpleTemplateMessage(ctx context.Context, phone, templateName, broadcastName string) (*MessageResponse, error) {
	req := &SendTemplateMessageRequest{
		WhatsappNumber: phone,
		TemplateName:   templateName,
		BroadcastName:  broadcastName,
	}
	
	return s.SendTemplateMessage(ctx, req)
}

// SendTemplateMessageWithParams envía un mensaje de plantilla con parámetros
func (s *Service) SendTemplateMessageWithParams(ctx context.Context, phone, templateName, broadcastName string, params map[string]string) (*MessageResponse, error) {
	var parameters []Parameter
	for name, value := range params {
		parameters = append(parameters, Parameter{
			Name:  name,
			Value: value,
		})
	}
	
	req := &SendTemplateMessageRequest{
		WhatsappNumber: phone,
		TemplateName:   templateName,
		BroadcastName:  broadcastName,
		Parameters:     parameters,
	}
	
	return s.SendTemplateMessage(ctx, req)
}

// CreateSimpleListMessage crea un mensaje de lista interactiva simple
func (s *Service) CreateSimpleListMessage(phone, bodyText, buttonText string, sections []InteractiveSection) *InteractiveListMessageRequest {
	return &InteractiveListMessageRequest{
		WhatsappNumber: phone,
		Body: InteractiveBody{
			Text: bodyText,
		},
		Action: InteractiveListAction{
			Button:   buttonText,
			Sections: sections,
		},
	}
}

// CreateSimpleButtonMessage crea un mensaje de botones interactivos simple
func (s *Service) CreateSimpleButtonMessage(phone, bodyText string, buttons []InteractiveButton) *InteractiveButtonMessageRequest {
	return &InteractiveButtonMessageRequest{
		WhatsappNumber: phone,
		Body: InteractiveBody{
			Text: bodyText,
		},
		Action: InteractiveButtonAction{
			Buttons: buttons,
		},
	}
}

// SendQuickReplyButtons envía botones de respuesta rápida
func (s *Service) SendQuickReplyButtons(ctx context.Context, phone, bodyText string, buttonTitles []string) (*MessageResponse, error) {
	if len(buttonTitles) == 0 || len(buttonTitles) > 3 {
		return nil, fmt.Errorf("must provide 1-3 button titles, got %d", len(buttonTitles))
	}
	
	var buttons []InteractiveButton
	for i, title := range buttonTitles {
		buttons = append(buttons, InteractiveButton{
			Type: "reply",
			Reply: InteractiveButtonReply{
				ID:    fmt.Sprintf("btn_%d", i+1),
				Title: title,
			},
		})
	}
	
	req := s.CreateSimpleButtonMessage(phone, bodyText, buttons)
	return s.SendInteractiveButtonMessage(ctx, req)
}

// SendListMenu envía un menú de lista con opciones
func (s *Service) SendListMenu(ctx context.Context, phone, bodyText, buttonText string, menuItems map[string][]string) (*MessageResponse, error) {
	var sections []InteractiveSection
	
	for sectionTitle, items := range menuItems {
		var rows []InteractiveListRow
		for i, item := range items {
			rows = append(rows, InteractiveListRow{
				ID:    fmt.Sprintf("%s_%d", strings.ToLower(strings.ReplaceAll(sectionTitle, " ", "_")), i+1),
				Title: item,
			})
		}
		
		sections = append(sections, InteractiveSection{
			Title: sectionTitle,
			Rows:  rows,
		})
	}
	
	req := s.CreateSimpleListMessage(phone, bodyText, buttonText, sections)
	return s.SendInteractiveListMessage(ctx, req)
}

// GetTemplatesByCategory obtiene plantillas filtradas por categoría
func (s *Service) GetTemplatesByCategory(ctx context.Context, category string) ([]Template, error) {
	templates, err := s.GetMessageTemplates(ctx)
	if err != nil {
		return nil, err
	}
	
	var filtered []Template
	for _, template := range templates.Templates {
		if template.Category == category {
			filtered = append(filtered, template)
		}
	}
	
	return filtered, nil
}

// GetActiveTemplates obtiene solo las plantillas activas
func (s *Service) GetActiveTemplates(ctx context.Context) ([]Template, error) {
	templates, err := s.GetMessageTemplates(ctx)
	if err != nil {
		return nil, err
	}
	
	var active []Template
	for _, template := range templates.Templates {
		if template.Status == "APPROVED" || template.Status == "ACTIVE" {
			active = append(active, template)
		}
	}
	
	return active, nil
}

