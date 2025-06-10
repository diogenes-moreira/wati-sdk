package wati

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
	"github.com/tu-usuario/go-wati/contacts"
	"github.com/tu-usuario/go-wati/messages"
	"github.com/tu-usuario/go-wati/chatbots"
	"github.com/tu-usuario/go-wati/media"
	"github.com/tu-usuario/go-wati/webhooks"
)

// WATIClient es la interfaz principal del cliente WATI
type WATIClient interface {
	// Servicios
	Contacts() ContactsService
	Messages() MessagesService
	Chatbots() ChatbotsService
	Media() MediaService
	Webhooks() WebhooksService
	
	// Configuración
	SetAPIEndpoint(endpoint string)
	SetToken(token string)
	GetConfig() *Config
	
	// Utilidades
	ValidateToken() error
	RotateToken() (*TokenResponse, error)
	
	// HTTP client interno
	DoRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error
}

// Client implementa WATIClient
type Client struct {
	config   *Config
	client   *http.Client
	limiter  *rate.Limiter
	
	// Servicios
	contacts  ContactsService
	messages  MessagesService
	chatbots  ChatbotsService
	media     MediaService
	webhooks  WebhooksService
}

// NewClient crea una nueva instancia del cliente WATI
func NewClient(apiEndpoint, token string, options ...ClientOption) WATIClient {
	config := DefaultConfig()
	config.APIEndpoint = strings.TrimSuffix(apiEndpoint, "/")
	config.Token = token
	
	// Aplicar opciones
	for _, option := range options {
		option(config)
	}
	
	// Crear rate limiter
	rateLimiter := rate.NewLimiter(
		rate.Limit(config.RateLimit.RequestsPerSecond),
		config.RateLimit.BurstSize,
	)
	
	// Crear cliente HTTP
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}
	
	client := &Client{
		config:      config,
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
	}
	
	// Inicializar servicios
	client.initServices()
	
	return client
}

// initServices inicializa todos los servicios
func (c *Client) initServices() {
	c.contacts = contacts.NewService(c)
	c.messages = messages.NewService(c)
	c.chatbots = chatbots.NewService(c)
	c.media = media.NewService(c)
	c.webhooks = webhooks.NewService(c)
}

// Contacts retorna el servicio de contactos
func (c *Client) Contacts() ContactsService {
	return c.contacts
}

// Messages retorna el servicio de mensajes
func (c *Client) Messages() MessagesService {
	return c.messages
}

// Chatbots retorna el servicio de chatbots
func (c *Client) Chatbots() ChatbotsService {
	return c.chatbots
}

// Media retorna el servicio de media
func (c *Client) Media() MediaService {
	return c.media
}

// Webhooks retorna el servicio de webhooks
func (c *Client) Webhooks() WebhooksService {
	return c.webhooks
}

// SetAPIEndpoint establece el endpoint de la API
func (c *Client) SetAPIEndpoint(endpoint string) {
	c.config.APIEndpoint = strings.TrimSuffix(endpoint, "/")
}

// SetToken establece el token de autenticación
func (c *Client) SetToken(token string) {
	c.config.Token = token
}

// GetConfig retorna la configuración actual
func (c *Client) GetConfig() *Config {
	return c.config
}

// ValidateToken valida el token actual
func (c *Client) ValidateToken() error {
	ctx := context.Background()
	
	// Intentar hacer una petición simple para validar el token
	var result BaseResponse
	err := c.DoRequest(ctx, "GET", "/api/v1/chatbots", nil, &result)
	if err != nil {
		if watiErr, ok := err.(*WATIError); ok && watiErr.IsAuthenticationError() {
			return ErrInvalidToken
		}
		return err
	}
	
	return nil
}

// RotateToken rota el token de autenticación
func (c *Client) RotateToken() (*TokenResponse, error) {
	ctx := context.Background()
	
	var result TokenResponse
	err := c.DoRequest(ctx, "POST", "/api/v1/rotateToken", nil, &result)
	if err != nil {
		return nil, err
	}
	
	// Actualizar el token en la configuración
	c.config.Token = result.Token
	
	return &result, nil
}

// DoRequest realiza una petición HTTP a la API de WATI
func (c *Client) DoRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	// Aplicar rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limiter error: %w", err)
	}
	
	// Construir URL completa
	fullURL := c.config.APIEndpoint + endpoint
	
	// Preparar el cuerpo de la petición
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}
	
	// Crear la petición
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	// Establecer headers
	req.Header.Set("Authorization", "Bearer "+c.config.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "go-wati/1.0.0")
	
	// Realizar la petición con reintentos
	var resp *http.Response
	var lastErr error
	
	for attempt := 0; attempt <= c.config.RetryCount; attempt++ {
		if attempt > 0 {
			// Esperar antes del reintento
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		
		resp, lastErr = c.httpClient.Do(req)
		if lastErr != nil {
			if attempt == c.config.RetryCount {
				return &NetworkError{
					Operation: fmt.Sprintf("%s %s", method, endpoint),
					Err:       lastErr,
				}
			}
			continue
		}
		
		// Si la respuesta es exitosa o no es reintentable, salir del bucle
		if resp.StatusCode < 500 && resp.StatusCode != 429 {
			break
		}
		
		resp.Body.Close()
		
		// Si es el último intento, no cerrar la respuesta aquí
		if attempt == c.config.RetryCount {
			break
		}
	}
	
	if resp == nil {
		return &NetworkError{
			Operation: fmt.Sprintf("%s %s", method, endpoint),
			Err:       lastErr,
		}
	}
	
	defer resp.Body.Close()
	
	// Leer el cuerpo de la respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	
	// Verificar el código de estado
	if resp.StatusCode >= 400 {
		// Intentar parsear el error de la API
		var apiError struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		
		if json.Unmarshal(respBody, &apiError) == nil && apiError.Error != "" {
			return NewWATIError(resp.StatusCode, apiError.Error)
		}
		
		if json.Unmarshal(respBody, &apiError) == nil && apiError.Message != "" {
			return NewWATIError(resp.StatusCode, apiError.Message)
		}
		
		return NewWATIError(resp.StatusCode, string(respBody))
	}
	
	// Parsear la respuesta exitosa
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("error unmarshaling response: %w", err)
		}
	}
	
	return nil
}

// buildURL construye una URL con parámetros de consulta
func (c *Client) buildURL(endpoint string, params map[string]string) string {
	u, _ := url.Parse(c.config.APIEndpoint + endpoint)
	
	if len(params) > 0 {
		q := u.Query()
		for key, value := range params {
			if value != "" {
				q.Set(key, value)
			}
		}
		u.RawQuery = q.Encode()
	}
	
	return u.String()
}

