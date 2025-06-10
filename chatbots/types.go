package chatbots

import (
	"fmt"
	"time"
)

// Chatbot representa un chatbot en WATI
type Chatbot struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Created     string    `json:"created"`
	Updated     string    `json:"updated,omitempty"`
	Rules       []Rule    `json:"rules,omitempty"`
	Keywords    []string  `json:"keywords,omitempty"`
	Responses   []Response `json:"responses,omitempty"`
}

// Rule representa una regla de chatbot
type Rule struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Trigger     Trigger     `json:"trigger"`
	Actions     []Action    `json:"actions"`
	Conditions  []Condition `json:"conditions,omitempty"`
	IsActive    bool        `json:"isActive"`
	Priority    int         `json:"priority"`
}

// Trigger representa un disparador de regla
type Trigger struct {
	Type     string   `json:"type"`
	Keywords []string `json:"keywords,omitempty"`
	Pattern  string   `json:"pattern,omitempty"`
	Event    string   `json:"event,omitempty"`
}

// Action representa una acción de chatbot
type Action struct {
	Type         string                 `json:"type"`
	Message      string                 `json:"message,omitempty"`
	Template     string                 `json:"template,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	Delay        int                    `json:"delay,omitempty"`
	AssignTo     string                 `json:"assignTo,omitempty"`
	TagsToAdd    []string               `json:"tagsToAdd,omitempty"`
	TagsToRemove []string               `json:"tagsToRemove,omitempty"`
}

// Condition representa una condición para ejecutar una regla
type Condition struct {
	Type     string      `json:"type"`
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// Response representa una respuesta de chatbot
type Response struct {
	ID       string `json:"id"`
	Trigger  string `json:"trigger"`
	Message  string `json:"message"`
	IsActive bool   `json:"isActive"`
}

// StartChatbotRequest representa la petición para iniciar un chatbot
type StartChatbotRequest struct {
	ChatbotID      string `json:"chatbotId"`
	WhatsappNumber string `json:"whatsappNumber"`
	InitialMessage string `json:"initialMessage,omitempty"`
}

// ChatbotResponse representa la respuesta de operaciones de chatbot
type ChatbotResponse struct {
	BaseResponse
	Chatbot   Chatbot `json:"chatbot"`
	SessionID string  `json:"sessionId,omitempty"`
	Status    string  `json:"status"`
}

// UpdateChatStatusRequest representa la petición para actualizar estado de chat
type UpdateChatStatusRequest struct {
	WhatsappNumber string `json:"whatsappNumber"`
	Status         string `json:"status"`
	AssignedTo     string `json:"assignedTo,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

// ChatStatusResponse representa la respuesta de actualización de estado
type ChatStatusResponse struct {
	BaseResponse
	WhatsappNumber string    `json:"whatsappNumber"`
	Status         string    `json:"status"`
	AssignedTo     string    `json:"assignedTo,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// ChatbotsResponse representa la respuesta de lista de chatbots
type ChatbotsResponse struct {
	BaseResponse
	Chatbots []Chatbot `json:"chatbots"`
}

// CreateChatbotRequest representa la petición para crear un chatbot
type CreateChatbotRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Keywords    []string   `json:"keywords,omitempty"`
	Responses   []Response `json:"responses,omitempty"`
	IsActive    bool       `json:"isActive"`
}

// UpdateChatbotRequest representa la petición para actualizar un chatbot
type UpdateChatbotRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Keywords    []string   `json:"keywords,omitempty"`
	Responses   []Response `json:"responses,omitempty"`
	IsActive    *bool      `json:"isActive,omitempty"`
}

// ChatSession representa una sesión de chat
type ChatSession struct {
	ID             string    `json:"id"`
	WhatsappNumber string    `json:"whatsappNumber"`
	ChatbotID      string    `json:"chatbotId"`
	Status         string    `json:"status"`
	StartedAt      time.Time `json:"startedAt"`
	EndedAt        *time.Time `json:"endedAt,omitempty"`
	LastActivity   time.Time `json:"lastActivity"`
	MessageCount   int       `json:"messageCount"`
	CurrentStep    string    `json:"currentStep,omitempty"`
	Variables      map[string]interface{} `json:"variables,omitempty"`
}

// ChatFlow representa un flujo de conversación
type ChatFlow struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Steps       []FlowStep `json:"steps"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// FlowStep representa un paso en un flujo de conversación
type FlowStep struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Message     string                 `json:"message,omitempty"`
	Options     []FlowOption           `json:"options,omitempty"`
	Conditions  []Condition            `json:"conditions,omitempty"`
	Actions     []Action               `json:"actions,omitempty"`
	NextStep    string                 `json:"nextStep,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
}

// FlowOption representa una opción en un paso de flujo
type FlowOption struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Value    string `json:"value"`
	NextStep string `json:"nextStep"`
}

// BaseResponse representa la respuesta base de la API
type BaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ChatStatus representa los posibles estados de un chat
type ChatStatus string

const (
	ChatStatusOpen     ChatStatus = "OPEN"
	ChatStatusAssigned ChatStatus = "ASSIGNED"
	ChatStatusResolved ChatStatus = "RESOLVED"
	ChatStatusClosed   ChatStatus = "CLOSED"
	ChatStatusBot      ChatStatus = "BOT"
)

// ChatbotStatus representa los posibles estados de un chatbot
type ChatbotStatus string

const (
	ChatbotStatusActive   ChatbotStatus = "ACTIVE"
	ChatbotStatusInactive ChatbotStatus = "INACTIVE"
	ChatbotStatusPaused   ChatbotStatus = "PAUSED"
)

// TriggerType representa los tipos de disparadores
type TriggerType string

const (
	TriggerTypeKeyword     TriggerType = "KEYWORD"
	TriggerTypePattern     TriggerType = "PATTERN"
	TriggerTypeEvent       TriggerType = "EVENT"
	TriggerTypeSchedule    TriggerType = "SCHEDULE"
	TriggerTypeInactivity  TriggerType = "INACTIVITY"
)

// ActionType representa los tipos de acciones
type ActionType string

const (
	ActionTypeSendMessage    ActionType = "SEND_MESSAGE"
	ActionTypeSendTemplate   ActionType = "SEND_TEMPLATE"
	ActionTypeAssignUser     ActionType = "ASSIGN_USER"
	ActionTypeAddTag         ActionType = "ADD_TAG"
	ActionTypeRemoveTag      ActionType = "REMOVE_TAG"
	ActionTypeSetVariable    ActionType = "SET_VARIABLE"
	ActionTypeWait           ActionType = "WAIT"
	ActionTypeTransferToHuman ActionType = "TRANSFER_TO_HUMAN"
)

// Validate valida la petición de inicio de chatbot
func (r *StartChatbotRequest) Validate() error {
	if r.ChatbotID == "" {
		return fmt.Errorf("chatbotId is required")
	}
	
	if r.WhatsappNumber == "" {
		return fmt.Errorf("whatsappNumber is required")
	}
	
	// Validación básica del número de teléfono
	if len(r.WhatsappNumber) < 10 {
		return fmt.Errorf("whatsappNumber must be at least 10 digits")
	}
	
	return nil
}

// Validate valida la petición de actualización de estado de chat
func (r *UpdateChatStatusRequest) Validate() error {
	if r.WhatsappNumber == "" {
		return fmt.Errorf("whatsappNumber is required")
	}
	
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	
	// Validar que el estado sea válido
	validStatuses := []string{
		string(ChatStatusOpen),
		string(ChatStatusAssigned),
		string(ChatStatusResolved),
		string(ChatStatusClosed),
		string(ChatStatusBot),
	}
	
	isValid := false
	for _, status := range validStatuses {
		if r.Status == status {
			isValid = true
			break
		}
	}
	
	if !isValid {
		return fmt.Errorf("invalid status: %s. Valid statuses are: %v", r.Status, validStatuses)
	}
	
	// Validación básica del número de teléfono
	if len(r.WhatsappNumber) < 10 {
		return fmt.Errorf("whatsappNumber must be at least 10 digits")
	}
	
	return nil
}

// Validate valida la petición de creación de chatbot
func (r *CreateChatbotRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if len(r.Keywords) == 0 && len(r.Responses) == 0 {
		return fmt.Errorf("at least one keyword or response is required")
	}
	
	// Validar respuestas
	for i, response := range r.Responses {
		if response.Trigger == "" {
			return fmt.Errorf("trigger is required for response %d", i)
		}
		
		if response.Message == "" {
			return fmt.Errorf("message is required for response %d", i)
		}
	}
	
	return nil
}

// IsActive verifica si el chatbot está activo
func (c *Chatbot) IsActive() bool {
	return c.Status == "active"
}

// GetActiveRules retorna solo las reglas activas del chatbot
func (c *Chatbot) GetActiveRules() []Rule {
	var activeRules []Rule
	for _, rule := range c.Rules {
		if rule.IsActive {
			activeRules = append(activeRules, rule)
		}
	}
	return activeRules
}

// GetActiveResponses retorna solo las respuestas activas del chatbot
func (c *Chatbot) GetActiveResponses() []Response {
	var activeResponses []Response
	for _, response := range c.Responses {
		if response.IsActive {
			activeResponses = append(activeResponses, response)
		}
	}
	return activeResponses
}

// HasKeyword verifica si el chatbot tiene una palabra clave específica
func (c *Chatbot) HasKeyword(keyword string) bool {
	for _, k := range c.Keywords {
		if k == keyword {
			return true
		}
	}
	return false
}

// AddKeyword agrega una palabra clave al chatbot
func (c *Chatbot) AddKeyword(keyword string) {
	if !c.HasKeyword(keyword) {
		c.Keywords = append(c.Keywords, keyword)
	}
}

// RemoveKeyword elimina una palabra clave del chatbot
func (c *Chatbot) RemoveKeyword(keyword string) {
	for i, k := range c.Keywords {
		if k == keyword {
			c.Keywords = append(c.Keywords[:i], c.Keywords[i+1:]...)
			break
		}
	}
}

