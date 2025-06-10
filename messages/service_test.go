package messages

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockHTTPClient implementa HTTPClient para testing
type MockHTTPClient struct {
	DoRequestFunc func(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error
}

func (m *MockHTTPClient) DoRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	if m.DoRequestFunc != nil {
		return m.DoRequestFunc(ctx, method, endpoint, body, result)
	}
	return nil
}

func TestNewService(t *testing.T) {
	mockClient := &MockHTTPClient{}
	service := NewService(mockClient)
	
	if service == nil {
		t.Error("NewService() returned nil")
	}
	
	if service.client != mockClient {
		t.Error("Service client not set correctly")
	}
}

func TestSendTemplateMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		request *SendTemplateMessageRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &SendTemplateMessageRequest{
				WhatsappNumber: "1234567890",
				TemplateName:   "hello_world",
				BroadcastName:  "test_broadcast",
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			request: nil,
			wantErr: true,
		},
		{
			name: "missing phone number",
			request: &SendTemplateMessageRequest{
				TemplateName:  "hello_world",
				BroadcastName: "test_broadcast",
			},
			wantErr: true,
		},
		{
			name: "missing template name",
			request: &SendTemplateMessageRequest{
				WhatsappNumber: "1234567890",
				BroadcastName:  "test_broadcast",
			},
			wantErr: true,
		},
		{
			name: "missing broadcast name",
			request: &SendTemplateMessageRequest{
				WhatsappNumber: "1234567890",
				TemplateName:   "hello_world",
			},
			wantErr: true,
		},
		{
			name: "invalid phone number",
			request: &SendTemplateMessageRequest{
				WhatsappNumber: "123",
				TemplateName:   "hello_world",
				BroadcastName:  "test_broadcast",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.request != nil {
				err := tt.request.Validate()
				if (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestSendTemplateMessage(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
			// Verificar método y endpoint
			if method != "POST" {
				t.Errorf("Expected POST method, got %s", method)
			}
			
			if endpoint != "/api/v1/sendTemplateMessage" {
				t.Errorf("Expected endpoint '/api/v1/sendTemplateMessage', got %s", endpoint)
			}
			
			// Verificar que el body es del tipo correcto
			if _, ok := body.(*SendTemplateMessageRequest); !ok {
				t.Errorf("Expected SendTemplateMessageRequest body, got %T", body)
			}
			
			// Simular respuesta exitosa
			if response, ok := result.(*MessageResponse); ok {
				response.BaseResponse.Result = true
				response.PhoneNumber = "1234567890"
				response.TemplateName = "hello_world"
				response.ValidWhatsAppNumber = true
			}
			
			return nil
		},
	}

	service := NewService(mockClient)
	ctx := context.Background()
	
	request := &SendTemplateMessageRequest{
		WhatsappNumber: "1234567890",
		TemplateName:   "hello_world",
		BroadcastName:  "test_broadcast",
	}
	
	response, err := service.SendTemplateMessage(ctx, request)
	if err != nil {
		t.Errorf("SendTemplateMessage() error = %v", err)
		return
	}
	
	if response == nil {
		t.Error("SendTemplateMessage() returned nil response")
		return
	}
	
	if response.PhoneNumber != "1234567890" {
		t.Errorf("Expected phone number '1234567890', got %s", response.PhoneNumber)
	}
	
	if !response.ValidWhatsAppNumber {
		t.Error("Expected ValidWhatsAppNumber to be true")
	}
}

func TestSendTemplateMessagesValidation(t *testing.T) {
	tests := []struct {
		name    string
		request *SendTemplateMessagesRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &SendTemplateMessagesRequest{
				TemplateName:  "hello_world",
				BroadcastName: "test_broadcast",
				Recipients: []TemplateMessageRecipient{
					{WhatsappNumber: "1234567890"},
					{WhatsappNumber: "0987654321"},
				},
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			request: nil,
			wantErr: true,
		},
		{
			name: "no recipients",
			request: &SendTemplateMessagesRequest{
				TemplateName:  "hello_world",
				BroadcastName: "test_broadcast",
				Recipients:    []TemplateMessageRecipient{},
			},
			wantErr: true,
		},
		{
			name: "too many recipients",
			request: &SendTemplateMessagesRequest{
				TemplateName:  "hello_world",
				BroadcastName: "test_broadcast",
				Recipients:    make([]TemplateMessageRecipient, 101), // Más de 100
			},
			wantErr: true,
		},
		{
			name: "invalid recipient phone",
			request: &SendTemplateMessagesRequest{
				TemplateName:  "hello_world",
				BroadcastName: "test_broadcast",
				Recipients: []TemplateMessageRecipient{
					{WhatsappNumber: "123"}, // Muy corto
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.request != nil {
				err := tt.request.Validate()
				if (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestInteractiveListMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		request *InteractiveListMessageRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &InteractiveListMessageRequest{
				WhatsappNumber: "1234567890",
				Body:           InteractiveBody{Text: "Choose an option"},
				Action: InteractiveListAction{
					Button: "Options",
					Sections: []InteractiveSection{
						{
							Title: "Products",
							Rows: []InteractiveListRow{
								{ID: "1", Title: "Product 1"},
								{ID: "2", Title: "Product 2"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing phone number",
			request: &InteractiveListMessageRequest{
				Body: InteractiveBody{Text: "Choose an option"},
				Action: InteractiveListAction{
					Button: "Options",
					Sections: []InteractiveSection{
						{
							Title: "Products",
							Rows: []InteractiveListRow{
								{ID: "1", Title: "Product 1"},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing body text",
			request: &InteractiveListMessageRequest{
				WhatsappNumber: "1234567890",
				Action: InteractiveListAction{
					Button: "Options",
					Sections: []InteractiveSection{
						{
							Title: "Products",
							Rows: []InteractiveListRow{
								{ID: "1", Title: "Product 1"},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no sections",
			request: &InteractiveListMessageRequest{
				WhatsappNumber: "1234567890",
				Body:           InteractiveBody{Text: "Choose an option"},
				Action: InteractiveListAction{
					Button:   "Options",
					Sections: []InteractiveSection{},
				},
			},
			wantErr: true,
		},
		{
			name: "section without rows",
			request: &InteractiveListMessageRequest{
				WhatsappNumber: "1234567890",
				Body:           InteractiveBody{Text: "Choose an option"},
				Action: InteractiveListAction{
					Button: "Options",
					Sections: []InteractiveSection{
						{
							Title: "Products",
							Rows:  []InteractiveListRow{},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInteractiveButtonMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		request *InteractiveButtonMessageRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &InteractiveButtonMessageRequest{
				WhatsappNumber: "1234567890",
				Body:           InteractiveBody{Text: "Choose an option"},
				Action: InteractiveButtonAction{
					Buttons: []InteractiveButton{
						{
							Type:  "reply",
							Reply: InteractiveButtonReply{ID: "1", Title: "Yes"},
						},
						{
							Type:  "reply",
							Reply: InteractiveButtonReply{ID: "2", Title: "No"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no buttons",
			request: &InteractiveButtonMessageRequest{
				WhatsappNumber: "1234567890",
				Body:           InteractiveBody{Text: "Choose an option"},
				Action: InteractiveButtonAction{
					Buttons: []InteractiveButton{},
				},
			},
			wantErr: true,
		},
		{
			name: "too many buttons",
			request: &InteractiveButtonMessageRequest{
				WhatsappNumber: "1234567890",
				Body:           InteractiveBody{Text: "Choose an option"},
				Action: InteractiveButtonAction{
					Buttons: []InteractiveButton{
						{Type: "reply", Reply: InteractiveButtonReply{ID: "1", Title: "Button 1"}},
						{Type: "reply", Reply: InteractiveButtonReply{ID: "2", Title: "Button 2"}},
						{Type: "reply", Reply: InteractiveButtonReply{ID: "3", Title: "Button 3"}},
						{Type: "reply", Reply: InteractiveButtonReply{ID: "4", Title: "Button 4"}}, // Más de 3
					},
				},
			},
			wantErr: true,
		},
		{
			name: "button without ID",
			request: &InteractiveButtonMessageRequest{
				WhatsappNumber: "1234567890",
				Body:           InteractiveBody{Text: "Choose an option"},
				Action: InteractiveButtonAction{
					Buttons: []InteractiveButton{
						{
							Type:  "reply",
							Reply: InteractiveButtonReply{Title: "Yes"}, // Sin ID
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetMessagesParams(t *testing.T) {
	params := &GetMessagesParams{
		PageSize:   10,
		PageNumber: 2,
		Phone:      "1234567890",
		FromDate:   "2024-01-01",
		ToDate:     "2024-01-31",
	}
	
	queryMap := params.ToMap()
	
	expectedParams := map[string]string{
		"pageSize":   "10",
		"pageNumber": "2",
		"phone":      "1234567890",
		"fromDate":   "2024-01-01",
		"toDate":     "2024-01-31",
	}
	
	for key, expectedValue := range expectedParams {
		if value, exists := queryMap[key]; !exists {
			t.Errorf("Expected parameter %s not found", key)
		} else if value != expectedValue {
			t.Errorf("Expected %s = %s, got %s", key, expectedValue, value)
		}
	}
}

func TestGetMessagesParamsDefaults(t *testing.T) {
	params := &GetMessagesParams{}
	params.SetDefaults()
	
	if params.PageSize != 20 {
		t.Errorf("Expected default PageSize 20, got %d", params.PageSize)
	}
	
	if params.PageNumber != 1 {
		t.Errorf("Expected default PageNumber 1, got %d", params.PageNumber)
	}
}

func TestSendQuickReplyButtons(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
			if endpoint != "/api/v1/sendInteractiveButtonMessage" {
				t.Errorf("Expected endpoint '/api/v1/sendInteractiveButtonMessage', got %s", endpoint)
			}
			
			// Simular respuesta exitosa
			if response, ok := result.(*MessageResponse); ok {
				response.BaseResponse.Result = true
			}
			
			return nil
		},
	}

	service := NewService(mockClient)
	ctx := context.Background()
	
	buttonTitles := []string{"Yes", "No", "Maybe"}
	
	response, err := service.SendQuickReplyButtons(ctx, "1234567890", "Do you agree?", buttonTitles)
	if err != nil {
		t.Errorf("SendQuickReplyButtons() error = %v", err)
		return
	}
	
	if response == nil {
		t.Error("SendQuickReplyButtons() returned nil response")
	}
}

func TestSendQuickReplyButtonsValidation(t *testing.T) {
	service := NewService(&MockHTTPClient{})
	ctx := context.Background()
	
	// Test con demasiados botones
	tooManyButtons := []string{"1", "2", "3", "4"} // Más de 3
	
	_, err := service.SendQuickReplyButtons(ctx, "1234567890", "Choose", tooManyButtons)
	if err == nil {
		t.Error("Expected error for too many buttons, got nil")
	}
	
	// Test sin botones
	_, err = service.SendQuickReplyButtons(ctx, "1234567890", "Choose", []string{})
	if err == nil {
		t.Error("Expected error for no buttons, got nil")
	}
}

func TestSendListMenu(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
			if endpoint != "/api/v1/sendInteractiveListMessage" {
				t.Errorf("Expected endpoint '/api/v1/sendInteractiveListMessage', got %s", endpoint)
			}
			
			// Verificar que el body contiene las secciones correctas
			if req, ok := body.(*InteractiveListMessageRequest); ok {
				if len(req.Action.Sections) == 0 {
					t.Error("Expected sections in list message")
				}
			}
			
			// Simular respuesta exitosa
			if response, ok := result.(*MessageResponse); ok {
				response.BaseResponse.Result = true
			}
			
			return nil
		},
	}

	service := NewService(mockClient)
	ctx := context.Background()
	
	menuItems := map[string][]string{
		"Products": {"Phone", "Tablet"},
		"Services": {"Support", "Warranty"},
	}
	
	response, err := service.SendListMenu(ctx, "1234567890", "What do you need?", "Options", menuItems)
	if err != nil {
		t.Errorf("SendListMenu() error = %v", err)
		return
	}
	
	if response == nil {
		t.Error("SendListMenu() returned nil response")
	}
}

// Benchmark para medir performance del servicio
func BenchmarkSendTemplateMessage(b *testing.B) {
	mockClient := &MockHTTPClient{
		DoRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
			if response, ok := result.(*MessageResponse); ok {
				response.BaseResponse.Result = true
				response.PhoneNumber = "1234567890"
			}
			return nil
		},
	}

	service := NewService(mockClient)
	ctx := context.Background()
	
	request := &SendTemplateMessageRequest{
		WhatsappNumber: "1234567890",
		TemplateName:   "hello_world",
		BroadcastName:  "test_broadcast",
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		service.SendTemplateMessage(ctx, request)
	}
}

