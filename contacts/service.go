package contacts

import (
	"context"
	"fmt"
	"strings"
)

// HTTPClient define la interfaz para realizar peticiones HTTP
type HTTPClient interface {
	DoRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error
}

// Service implementa ContactsService
type Service struct {
	client HTTPClient
}

// NewService crea una nueva instancia del servicio de contactos
func NewService(client HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// GetContacts obtiene la lista de contactos con parámetros opcionales
func (s *Service) GetContacts(ctx context.Context, params *GetContactsParams) (*ContactsResponse, error) {
	if params == nil {
		params = &GetContactsParams{}
	}
	
	params.SetDefaults()
	
	// Construir endpoint con query parameters
	endpoint := "/api/v1/getContacts"
	queryParams := params.ToMap()
	
	if len(queryParams) > 0 {
		var parts []string
		for key, value := range queryParams {
			parts = append(parts, fmt.Sprintf("%s=%s", key, value))
		}
		endpoint += "?" + strings.Join(parts, "&")
	}
	
	var response ContactsResponse
	err := s.client.DoRequest(ctx, "GET", endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting contacts: %w", err)
	}
	
	return &response, nil
}

// GetContact obtiene un contacto específico por ID
func (s *Service) GetContact(ctx context.Context, id string) (*Contact, error) {
	if id == "" {
		return nil, fmt.Errorf("contact ID is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/getContact/%s", id)
	
	var response struct {
		BaseResponse
		Contact Contact `json:"contact"`
	}
	
	err := s.client.DoRequest(ctx, "GET", endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting contact %s: %w", id, err)
	}
	
	return &response.Contact, nil
}

// AddContact crea un nuevo contacto
func (s *Service) AddContact(ctx context.Context, contact *CreateContactRequest) (*Contact, error) {
	if contact == nil {
		return nil, fmt.Errorf("contact data is required")
	}
	
	if err := contact.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	var response struct {
		BaseResponse
		Contact Contact `json:"contact"`
	}
	
	err := s.client.DoRequest(ctx, "POST", "/api/v1/addContact", contact, &response)
	if err != nil {
		return nil, fmt.Errorf("error adding contact: %w", err)
	}
	
	return &response.Contact, nil
}

// UpdateContact actualiza un contacto existente
func (s *Service) UpdateContact(ctx context.Context, id string, contact *UpdateContactRequest) (*Contact, error) {
	if id == "" {
		return nil, fmt.Errorf("contact ID is required")
	}
	
	if contact == nil {
		return nil, fmt.Errorf("contact update data is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/updateContact/%s", id)
	
	var response struct {
		BaseResponse
		Contact Contact `json:"contact"`
	}
	
	err := s.client.DoRequest(ctx, "PUT", endpoint, contact, &response)
	if err != nil {
		return nil, fmt.Errorf("error updating contact %s: %w", id, err)
	}
	
	return &response.Contact, nil
}

// DeleteContact elimina un contacto
func (s *Service) DeleteContact(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("contact ID is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/deleteContact/%s", id)
	
	var response BaseResponse
	err := s.client.DoRequest(ctx, "DELETE", endpoint, nil, &response)
	if err != nil {
		return fmt.Errorf("error deleting contact %s: %w", id, err)
	}
	
	return nil
}

// SearchContacts busca contactos por query
func (s *Service) SearchContacts(ctx context.Context, query string) (*ContactsResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}
	
	params := &GetContactsParams{
		Name: query,
	}
	
	return s.GetContacts(ctx, params)
}

// FilterContacts filtra contactos según criterios específicos
func (s *Service) FilterContacts(ctx context.Context, filter *ContactFilter) (*ContactsResponse, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is required")
	}
	
	// Convertir filtro a parámetros de consulta
	params := &GetContactsParams{
		Name: filter.Name,
	}
	
	// Si hay filtros por fecha, convertir a string
	if !filter.CreatedAfter.IsZero() {
		params.CreatedDate = filter.CreatedAfter.Format("2006-01-02")
	}
	
	return s.GetContacts(ctx, params)
}

// AddContacts agrega múltiples contactos en una operación
func (s *Service) AddContacts(ctx context.Context, contacts []*CreateContactRequest) (*BulkContactResponse, error) {
	if len(contacts) == 0 {
		return nil, fmt.Errorf("at least one contact is required")
	}
	
	// Validar todos los contactos antes de enviar
	for i, contact := range contacts {
		if err := contact.Validate(); err != nil {
			return nil, fmt.Errorf("validation error for contact %d: %w", i, err)
		}
	}
	
	// WATI permite hasta 100 contactos por llamada
	if len(contacts) > 100 {
		return nil, fmt.Errorf("maximum 100 contacts allowed per request, got %d", len(contacts))
	}
	
	requestBody := struct {
		Contacts []*CreateContactRequest `json:"contacts"`
	}{
		Contacts: contacts,
	}
	
	var response BulkContactResponse
	err := s.client.DoRequest(ctx, "POST", "/api/v1/addContacts", requestBody, &response)
	if err != nil {
		return nil, fmt.Errorf("error adding contacts: %w", err)
	}
	
	return &response, nil
}

// GetContactsByPage obtiene contactos de una página específica
func (s *Service) GetContactsByPage(ctx context.Context, page, pageSize int) (*ContactsResponse, error) {
	params := &GetContactsParams{
		PageNumber: page,
		PageSize:   pageSize,
	}
	
	return s.GetContacts(ctx, params)
}

// GetAllContacts obtiene todos los contactos paginando automáticamente
func (s *Service) GetAllContacts(ctx context.Context) ([]Contact, error) {
	var allContacts []Contact
	page := 1
	pageSize := 50
	
	for {
		response, err := s.GetContactsByPage(ctx, page, pageSize)
		if err != nil {
			return nil, fmt.Errorf("error getting contacts page %d: %w", page, err)
		}
		
		allContacts = append(allContacts, response.Contacts...)
		
		// Si no hay más páginas, terminar
		if page >= response.TotalPages || len(response.Contacts) == 0 {
			break
		}
		
		page++
	}
	
	return allContacts, nil
}

// GetContactByPhone busca un contacto por número de teléfono
func (s *Service) GetContactByPhone(ctx context.Context, phone string) (*Contact, error) {
	if phone == "" {
		return nil, fmt.Errorf("phone number is required")
	}
	
	// Buscar usando el endpoint de búsqueda
	params := &GetContactsParams{
		PageSize: 1,
	}
	
	// Construir endpoint con el teléfono como filtro
	endpoint := fmt.Sprintf("/api/v1/getContacts?phone=%s&pageSize=%d", phone, params.PageSize)
	
	var response ContactsResponse
	err := s.client.DoRequest(ctx, "GET", endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error searching contact by phone %s: %w", phone, err)
	}
	
	if len(response.Contacts) == 0 {
		return nil, fmt.Errorf("contact with phone %s not found", phone)
	}
	
	return &response.Contacts[0], nil
}

// UpdateContactTags actualiza solo las etiquetas de un contacto
func (s *Service) UpdateContactTags(ctx context.Context, id string, tags []string) (*Contact, error) {
	updateReq := &UpdateContactRequest{
		Tags: tags,
	}
	
	return s.UpdateContact(ctx, id, updateReq)
}

// UpdateContactCustomParams actualiza solo los parámetros personalizados de un contacto
func (s *Service) UpdateContactCustomParams(ctx context.Context, id string, customParams []CustomParam) (*Contact, error) {
	updateReq := &UpdateContactRequest{
		CustomParams: customParams,
	}
	
	return s.UpdateContact(ctx, id, updateReq)
}

