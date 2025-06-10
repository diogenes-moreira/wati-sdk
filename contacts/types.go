package contacts

import (
	"fmt"
	"strconv"
	"time"
)

// Contact representa un contacto en WATI
type Contact struct {
	ID                string        `json:"id"`
	WAId              string        `json:"wAid"`
	FirstName         string        `json:"firstName"`
	LastName          string        `json:"lastName,omitempty"`
	FullName          string        `json:"fullName"`
	Phone             string        `json:"phone"`
	Email             string        `json:"email,omitempty"`
	Source            string        `json:"source,omitempty"`
	ContactStatus     string        `json:"contactStatus"`
	Photo             string        `json:"photo,omitempty"`
	Created           string        `json:"created"`
	Tags              []string      `json:"tags"`
	CustomParams      []CustomParam `json:"customParams"`
	OptedIn           bool          `json:"optedIn"`
	IsDeleted         bool          `json:"isDeleted"`
	LastUpdated       string        `json:"lastUpdated"`
	AllowBroadcast    bool          `json:"allowBroadcast"`
	AllowSMS          bool          `json:"allowSMS"`
	TeamIds           []string      `json:"teamIds"`
	IsInFlow          bool          `json:"isInFlow"`
	LastFlowId        string        `json:"lastFlowId,omitempty"`
	CurrentFlowNodeId string        `json:"currentFlowNodeId,omitempty"`
}

// CustomParam representa un parámetro personalizado del contacto
type CustomParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetContactsParams representa los parámetros para obtener contactos
type GetContactsParams struct {
	PageSize    int    `json:"pageSize,omitempty"`
	PageNumber  int    `json:"pageNumber,omitempty"`
	Name        string `json:"name,omitempty"`
	Attribute   string `json:"attribute,omitempty"`
	CreatedDate string `json:"createdDate,omitempty"`
}

// ContactsResponse representa la respuesta de la lista de contactos
type ContactsResponse struct {
	BaseResponse
	PaginatedResponse
	Contacts []Contact `json:"contacts"`
}

// CreateContactRequest representa la petición para crear un contacto
type CreateContactRequest struct {
	FirstName      string        `json:"firstName"`
	LastName       string        `json:"lastName,omitempty"`
	Phone          string        `json:"phone"`
	Email          string        `json:"email,omitempty"`
	CustomParams   []CustomParam `json:"customParams,omitempty"`
	Tags           []string      `json:"tags,omitempty"`
	AllowBroadcast bool          `json:"allowBroadcast"`
	AllowSMS       bool          `json:"allowSMS"`
}

// UpdateContactRequest representa la petición para actualizar un contacto
type UpdateContactRequest struct {
	FirstName      *string       `json:"firstName,omitempty"`
	LastName       *string       `json:"lastName,omitempty"`
	Email          *string       `json:"email,omitempty"`
	CustomParams   []CustomParam `json:"customParams,omitempty"`
	Tags           []string      `json:"tags,omitempty"`
	AllowBroadcast *bool         `json:"allowBroadcast,omitempty"`
	AllowSMS       *bool         `json:"allowSMS,omitempty"`
}

// ContactFilter representa filtros para búsqueda de contactos
type ContactFilter struct {
	Name          string    `json:"name,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	Email         string    `json:"email,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	ContactStatus string    `json:"contactStatus,omitempty"`
	CreatedAfter  time.Time `json:"createdAfter,omitempty"`
	CreatedBefore time.Time `json:"createdBefore,omitempty"`
	OptedIn       *bool     `json:"optedIn,omitempty"`
	AllowBroadcast *bool    `json:"allowBroadcast,omitempty"`
}

// BulkContactResponse representa la respuesta de operaciones en lote
type BulkContactResponse struct {
	BaseResponse
	SuccessCount int       `json:"successCount"`
	FailureCount int       `json:"failureCount"`
	Contacts     []Contact `json:"contacts"`
	Errors       []struct {
		Index   int    `json:"index"`
		Error   string `json:"error"`
		Contact CreateContactRequest `json:"contact"`
	} `json:"errors,omitempty"`
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

// Validate valida los datos del contacto
func (c *CreateContactRequest) Validate() error {
	if c.FirstName == "" {
		return fmt.Errorf("firstName is required")
	}
	
	if c.Phone == "" {
		return fmt.Errorf("phone is required")
	}
	
	// Validación básica del número de teléfono
	if len(c.Phone) < 10 {
		return fmt.Errorf("phone number must be at least 10 digits")
	}
	
	return nil
}

// ToMap convierte GetContactsParams a un mapa para query parameters
func (p *GetContactsParams) ToMap() map[string]string {
	params := make(map[string]string)
	
	if p.PageSize > 0 {
		params["pageSize"] = strconv.Itoa(p.PageSize)
	}
	
	if p.PageNumber > 0 {
		params["pageNumber"] = strconv.Itoa(p.PageNumber)
	}
	
	if p.Name != "" {
		params["name"] = p.Name
	}
	
	if p.Attribute != "" {
		params["attribute"] = p.Attribute
	}
	
	if p.CreatedDate != "" {
		params["createdDate"] = p.CreatedDate
	}
	
	return params
}

// SetDefaults establece valores por defecto para GetContactsParams
func (p *GetContactsParams) SetDefaults() {
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	
	if p.PageNumber <= 0 {
		p.PageNumber = 1
	}
}

