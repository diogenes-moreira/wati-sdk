package messages

import (
	"fmt"
	"strconv"
)

// Message representa un mensaje en WATI
type Message struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Content     string          `json:"content"`
	From        string          `json:"from"`
	To          string          `json:"to"`
	Timestamp   string          `json:"timestamp"`
	Status      string          `json:"status"`
	Direction   string          `json:"direction"`
	MessageType string          `json:"messageType"`
	Media       *MediaInfo      `json:"media,omitempty"`
	Template    *TemplateInfo   `json:"template,omitempty"`
	Interactive *InteractiveInfo `json:"interactive,omitempty"`
}

// MediaInfo representa información de media en un mensaje
type MediaInfo struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
	Caption  string `json:"caption,omitempty"`
}

// TemplateInfo representa información de una plantilla
type TemplateInfo struct {
	Name       string      `json:"name"`
	Language   string      `json:"language"`
	Parameters []Parameter `json:"parameters,omitempty"`
}

// Parameter representa un parámetro de plantilla
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// InteractiveInfo representa información de mensaje interactivo
type InteractiveInfo struct {
	Type    string      `json:"type"`
	Header  interface{} `json:"header,omitempty"`
	Body    interface{} `json:"body,omitempty"`
	Footer  interface{} `json:"footer,omitempty"`
	Action  interface{} `json:"action,omitempty"`
}

// SendTemplateMessageRequest representa la petición para enviar un mensaje de plantilla
type SendTemplateMessageRequest struct {
	WhatsappNumber string      `json:"whatsappNumber"`
	TemplateName   string      `json:"template_name"`
	BroadcastName  string      `json:"broadcast_name"`
	Parameters     []Parameter `json:"parameters,omitempty"`
}

// SendTemplateMessagesRequest representa la petición para enviar múltiples mensajes de plantilla
type SendTemplateMessagesRequest struct {
	TemplateName   string                        `json:"template_name"`
	BroadcastName  string                        `json:"broadcast_name"`
	Recipients     []TemplateMessageRecipient    `json:"recipients"`
}

// TemplateMessageRecipient representa un destinatario de mensaje de plantilla
type TemplateMessageRecipient struct {
	WhatsappNumber string      `json:"whatsappNumber"`
	Parameters     []Parameter `json:"parameters,omitempty"`
}

// MessageResponse representa la respuesta de envío de mensaje
type MessageResponse struct {
	BaseResponse
	PhoneNumber         string    `json:"phone_number"`
	TemplateName        string    `json:"template_name"`
	Parameters          []Parameter `json:"parameteres"` // Nota: typo en la API original
	Contact             Contact   `json:"contact"`
	Model               Model     `json:"model"`
	ValidWhatsAppNumber bool      `json:"validWhatsAppNumber"`
}

// BulkMessageResponse representa la respuesta de envío múltiple
type BulkMessageResponse struct {
	BaseResponse
	SuccessCount int             `json:"successCount"`
	FailureCount int             `json:"failureCount"`
	Messages     []MessageResponse `json:"messages"`
	Errors       []struct {
		Index     int    `json:"index"`
		Error     string `json:"error"`
		Recipient TemplateMessageRecipient `json:"recipient"`
	} `json:"errors,omitempty"`
}

// Contact representa un contacto en la respuesta de mensaje
type Contact struct {
	ID                string        `json:"id"`
	WAId              string        `json:"wAid"`
	FirstName         string        `json:"firstName"`
	FullName          string        `json:"fullName"`
	Phone             string        `json:"phone"`
	Source            interface{}   `json:"source"`
	ContactStatus     string        `json:"contactStatus"`
	Photo             interface{}   `json:"photo"`
	Created           string        `json:"created"`
	Tags              []interface{} `json:"tags"`
	CustomParams      []CustomParam `json:"customParams"`
	OptedIn           bool          `json:"optedIn"`
	IsDeleted         bool          `json:"isDeleted"`
	LastUpdated       string        `json:"lastUpdated"`
	AllowBroadcast    bool          `json:"allowBroadcast"`
	AllowSMS          bool          `json:"allowSMS"`
	TeamIds           []string      `json:"teamIds"`
	IsInFlow          bool          `json:"isInFlow"`
	LastFlowId        string        `json:"lastFlowId"`
	CurrentFlowNodeId string        `json:"currentFlowNodeId"`
}

// CustomParam representa un parámetro personalizado
type CustomParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Model representa el modelo en la respuesta
type Model struct {
	IDs []string `json:"ids"`
}

// InteractiveListMessageRequest representa la petición para mensaje de lista interactiva
type InteractiveListMessageRequest struct {
	WhatsappNumber string                `json:"whatsappNumber"`
	Header         *InteractiveHeader    `json:"header,omitempty"`
	Body           InteractiveBody       `json:"body"`
	Footer         *InteractiveFooter    `json:"footer,omitempty"`
	Action         InteractiveListAction `json:"action"`
}

// InteractiveButtonMessageRequest representa la petición para mensaje de botones interactivos
type InteractiveButtonMessageRequest struct {
	WhatsappNumber string                  `json:"whatsappNumber"`
	Header         *InteractiveHeader      `json:"header,omitempty"`
	Body           InteractiveBody         `json:"body"`
	Footer         *InteractiveFooter      `json:"footer,omitempty"`
	Action         InteractiveButtonAction `json:"action"`
}

// InteractiveHeader representa el header de un mensaje interactivo
type InteractiveHeader struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// InteractiveBody representa el cuerpo de un mensaje interactivo
type InteractiveBody struct {
	Text string `json:"text"`
}

// InteractiveFooter representa el footer de un mensaje interactivo
type InteractiveFooter struct {
	Text string `json:"text"`
}

// InteractiveListAction representa la acción de lista interactiva
type InteractiveListAction struct {
	Button   string                 `json:"button"`
	Sections []InteractiveSection   `json:"sections"`
}

// InteractiveButtonAction representa la acción de botones interactivos
type InteractiveButtonAction struct {
	Buttons []InteractiveButton `json:"buttons"`
}

// InteractiveSection representa una sección de lista interactiva
type InteractiveSection struct {
	Title string                `json:"title"`
	Rows  []InteractiveListRow  `json:"rows"`
}

// InteractiveListRow representa una fila de lista interactiva
type InteractiveListRow struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// InteractiveButton representa un botón interactivo
type InteractiveButton struct {
	Type  string                `json:"type"`
	Reply InteractiveButtonReply `json:"reply"`
}

// InteractiveButtonReply representa la respuesta de un botón interactivo
type InteractiveButtonReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Template representa una plantilla de mensaje
type Template struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Language    string              `json:"language"`
	Status      string              `json:"status"`
	Category    string              `json:"category"`
	Components  []TemplateComponent `json:"components"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
}

// TemplateComponent representa un componente de plantilla
type TemplateComponent struct {
	Type       string                 `json:"type"`
	Format     string                 `json:"format,omitempty"`
	Text       string                 `json:"text,omitempty"`
	Parameters []TemplateParameter    `json:"parameters,omitempty"`
	Buttons    []TemplateButton       `json:"buttons,omitempty"`
}

// TemplateParameter representa un parámetro de plantilla
type TemplateParameter struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// TemplateButton representa un botón de plantilla
type TemplateButton struct {
	Type string `json:"type"`
	Text string `json:"text"`
	URL  string `json:"url,omitempty"`
}

// TemplatesResponse representa la respuesta de plantillas
type TemplatesResponse struct {
	BaseResponse
	Templates []Template `json:"templates"`
}

// GetMessagesParams representa los parámetros para obtener mensajes
type GetMessagesParams struct {
	PageSize   int    `json:"pageSize,omitempty"`
	PageNumber int    `json:"pageNumber,omitempty"`
	Phone      string `json:"phone,omitempty"`
	FromDate   string `json:"fromDate,omitempty"`
	ToDate     string `json:"toDate,omitempty"`
}

// MessagesResponse representa la respuesta de mensajes
type MessagesResponse struct {
	BaseResponse
	PaginatedResponse
	Messages []Message `json:"messages"`
}

// MessageStatus representa el estado de un mensaje
type MessageStatus struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Error     string `json:"error,omitempty"`
}

// BaseResponse representa la respuesta base de la API
type BaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
	TotalCount int `json:"totalCount"`
}

// Validate valida la petición de mensaje de plantilla
func (r *SendTemplateMessageRequest) Validate() error {
	if r.WhatsappNumber == "" {
		return fmt.Errorf("whatsappNumber is required")
	}
	
	if r.TemplateName == "" {
		return fmt.Errorf("template_name is required")
	}
	
	if r.BroadcastName == "" {
		return fmt.Errorf("broadcast_name is required")
	}
	
	// Validación básica del número de teléfono
	if len(r.WhatsappNumber) < 10 {
		return fmt.Errorf("whatsappNumber must be at least 10 digits")
	}
	
	return nil
}

// Validate valida la petición de múltiples mensajes de plantilla
func (r *SendTemplateMessagesRequest) Validate() error {
	if r.TemplateName == "" {
		return fmt.Errorf("template_name is required")
	}
	
	if r.BroadcastName == "" {
		return fmt.Errorf("broadcast_name is required")
	}
	
	if len(r.Recipients) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	
	// WATI permite hasta 100 destinatarios por llamada
	if len(r.Recipients) > 100 {
		return fmt.Errorf("maximum 100 recipients allowed per request, got %d", len(r.Recipients))
	}
	
	// Validar cada destinatario
	for i, recipient := range r.Recipients {
		if recipient.WhatsappNumber == "" {
			return fmt.Errorf("whatsappNumber is required for recipient %d", i)
		}
		
		if len(recipient.WhatsappNumber) < 10 {
			return fmt.Errorf("whatsappNumber must be at least 10 digits for recipient %d", i)
		}
	}
	
	return nil
}

// Validate valida la petición de mensaje de lista interactiva
func (r *InteractiveListMessageRequest) Validate() error {
	if r.WhatsappNumber == "" {
		return fmt.Errorf("whatsappNumber is required")
	}
	
	if len(r.WhatsappNumber) < 10 {
		return fmt.Errorf("whatsappNumber must be at least 10 digits")
	}
	
	if r.Body.Text == "" {
		return fmt.Errorf("body text is required")
	}
	
	if r.Action.Button == "" {
		return fmt.Errorf("action button text is required")
	}
	
	if len(r.Action.Sections) == 0 {
		return fmt.Errorf("at least one section is required")
	}
	
	// Validar secciones
	for i, section := range r.Action.Sections {
		if section.Title == "" {
			return fmt.Errorf("section title is required for section %d", i)
		}
		
		if len(section.Rows) == 0 {
			return fmt.Errorf("at least one row is required for section %d", i)
		}
		
		// Validar filas
		for j, row := range section.Rows {
			if row.ID == "" {
				return fmt.Errorf("row ID is required for section %d, row %d", i, j)
			}
			
			if row.Title == "" {
				return fmt.Errorf("row title is required for section %d, row %d", i, j)
			}
		}
	}
	
	return nil
}

// Validate valida la petición de mensaje de botones interactivos
func (r *InteractiveButtonMessageRequest) Validate() error {
	if r.WhatsappNumber == "" {
		return fmt.Errorf("whatsappNumber is required")
	}
	
	if len(r.WhatsappNumber) < 10 {
		return fmt.Errorf("whatsappNumber must be at least 10 digits")
	}
	
	if r.Body.Text == "" {
		return fmt.Errorf("body text is required")
	}
	
	if len(r.Action.Buttons) == 0 {
		return fmt.Errorf("at least one button is required")
	}
	
	if len(r.Action.Buttons) > 3 {
		return fmt.Errorf("maximum 3 buttons allowed, got %d", len(r.Action.Buttons))
	}
	
	// Validar botones
	for i, button := range r.Action.Buttons {
		if button.Reply.ID == "" {
			return fmt.Errorf("button ID is required for button %d", i)
		}
		
		if button.Reply.Title == "" {
			return fmt.Errorf("button title is required for button %d", i)
		}
	}
	
	return nil
}

// ToMap convierte GetMessagesParams a un mapa para query parameters
func (p *GetMessagesParams) ToMap() map[string]string {
	params := make(map[string]string)
	
	if p.PageSize > 0 {
		params["pageSize"] = strconv.Itoa(p.PageSize)
	}
	
	if p.PageNumber > 0 {
		params["pageNumber"] = strconv.Itoa(p.PageNumber)
	}
	
	if p.Phone != "" {
		params["phone"] = p.Phone
	}
	
	if p.FromDate != "" {
		params["fromDate"] = p.FromDate
	}
	
	if p.ToDate != "" {
		params["toDate"] = p.ToDate
	}
	
	return params
}

// SetDefaults establece valores por defecto para GetMessagesParams
func (p *GetMessagesParams) SetDefaults() {
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	
	if p.PageNumber <= 0 {
		p.PageNumber = 1
	}
}

