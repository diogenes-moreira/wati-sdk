# Go WATI - Librería Go para WATI API

Una librería completa en Go para interactuar con la API de WATI (WhatsApp Business API), que proporciona todas las funcionalidades necesarias para gestionar mensajes, contactos, chatbots, media y webhooks.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Configuración Inicial](#configuración-inicial)
- [Guía de Uso](#guía-de-uso)
  - [Mensajes](#mensajes)
  - [Contactos](#contactos)
  - [Chatbots](#chatbots)
  - [Media](#media)
  - [Webhooks](#webhooks)
- [Ejemplos Completos](#ejemplos-completos)
- [Configuración Avanzada](#configuración-avanzada)
- [Manejo de Errores](#manejo-de-errores)
- [Mejores Prácticas](#mejores-prácticas)
- [Contribución](#contribución)
- [Licencia](#licencia)

## ✨ Características

### 🚀 Funcionalidades Principales

- **Mensajería Completa**: Envío de mensajes de plantilla, mensajes interactivos (botones y listas)
- **Gestión de Contactos**: CRUD completo, búsqueda, filtrado y operaciones en lote
- **Chatbots Inteligentes**: Creación, gestión y control de chatbots automatizados
- **Gestión de Media**: Subida, descarga y gestión de archivos multimedia
- **Webhooks Avanzados**: Servidor integrado con manejo de eventos en tiempo real
- **Rate Limiting**: Control automático de límites de velocidad según el plan de WATI
- **Reintentos Automáticos**: Manejo robusto de errores con reintentos configurables
- **Validaciones**: Validación completa de datos antes del envío
- **Tipado Fuerte**: Tipos de datos completos y seguros

### 🛡️ Características de Seguridad

- Autenticación OAuth2 con Bearer tokens
- Validación de firmas de webhooks con HMAC-SHA256
- Manejo seguro de secretos y configuraciones
- Timeouts configurables para prevenir bloqueos

### 📊 Características de Monitoreo

- Logging detallado de operaciones
- Métricas de uso y estadísticas
- Health checks para webhooks
- Manejo de estados de mensajes y chats

## 📦 Instalación

```bash
go get github.com/tu-usuario/go-wati
```

### Dependencias

La librería utiliza las siguientes dependencias:

```go
require (
    golang.org/x/time v0.5.0
)
```

## ⚙️ Configuración Inicial

### 1. Obtener Credenciales de WATI

Antes de usar la librería, necesitas:

1. **Endpoint de API**: Tu URL de servidor WATI (ej: `https://live-server-12345.wati.io`)
2. **Token de API**: Tu token de autenticación de WATI

### 2. Configuración Básica

```go
package main

import (
    "github.com/tu-usuario/go-wati"
)

func main() {
    // Configuración básica
    client := wati.NewClient(
        "https://live-server-12345.wati.io", // Tu endpoint
        "tu-token-aqui",                     // Tu token
    )
    
    // Usar el cliente...
}
```

### 3. Configuración Avanzada

```go
client := wati.NewClient(
    "https://live-server-12345.wati.io",
    "tu-token-aqui",
    wati.WithTimeout(30),           // Timeout de 30 segundos
    wati.WithRetries(3),            // 3 reintentos automáticos
    wati.WithRateLimit(100, 60),    // 100 requests por minuto
    wati.WithUserAgent("MiApp/1.0"), // User agent personalizado
)
```

## 📖 Guía de Uso

### 📨 Mensajes

#### Envío de Mensajes de Plantilla

```go
import (
    "context"
    "github.com/tu-usuario/go-wati/messages"
)

ctx := context.Background()

// Mensaje simple
response, err := client.Messages().SendTemplateMessage(ctx, &messages.SendTemplateMessageRequest{
    WhatsappNumber: "1234567890",
    TemplateName:   "hello_world",
    BroadcastName:  "mi_broadcast",
})

// Mensaje con parámetros
params := map[string]string{
    "name": "Juan Pérez",
    "product": "Smartphone XYZ",
}

response, err := client.Messages().SendTemplateMessageWithParams(
    ctx,
    "1234567890",
    "order_confirmation",
    "confirmaciones",
    params,
)
```

#### Mensajes Interactivos

```go
// Botones de respuesta rápida
response, err := client.Messages().SendQuickReplyButtons(
    ctx,
    "1234567890",
    "¿Te interesa nuestro producto?",
    []string{"Sí", "No", "Más info"},
)

// Menú de lista
menuItems := map[string][]string{
    "Productos": {"Smartphone", "Tablet", "Laptop"},
    "Servicios": {"Soporte", "Garantía", "Reparación"},
}

response, err := client.Messages().SendListMenu(
    ctx,
    "1234567890",
    "¿En qué podemos ayudarte?",
    "Ver opciones",
    menuItems,
)
```

#### Gestión de Plantillas

```go
// Obtener todas las plantillas
templates, err := client.Messages().GetMessageTemplates(ctx)

// Obtener plantillas activas
activeTemplates, err := client.Messages().GetActiveTemplates(ctx)

// Buscar plantilla específica
template, err := client.Messages().GetMessageTemplate(ctx, "hello_world")
```

#### Historial de Mensajes

```go
// Obtener mensajes con paginación
messages, err := client.Messages().GetMessages(ctx, &messages.GetMessagesParams{
    PageSize:   20,
    PageNumber: 1,
})

// Mensajes de un contacto específico
messages, err := client.Messages().GetMessagesByPhone(ctx, "1234567890", nil)

// Mensajes en un rango de fechas
messages, err := client.Messages().GetMessagesByDateRange(
    ctx,
    "2024-01-01",
    "2024-01-31",
    nil,
)
```

### 👥 Contactos

#### Operaciones CRUD

```go
import "github.com/tu-usuario/go-wati/contacts"

// Crear contacto
newContact := &contacts.CreateContactRequest{
    WhatsappNumber: "1234567890",
    FirstName:      "Juan",
    LastName:       "Pérez",
    Email:          "juan@email.com",
    Tags:           []string{"cliente", "premium"},
    CustomParams: []contacts.CustomParam{
        {Name: "empresa", Value: "Tech Corp"},
        {Name: "cargo", Value: "Gerente"},
    },
}

contact, err := client.Contacts().AddContact(ctx, newContact)

// Obtener contacto
contact, err := client.Contacts().GetContact(ctx, "contact-id")

// Buscar por teléfono
contact, err := client.Contacts().GetContactByPhone(ctx, "1234567890")

// Actualizar contacto
updateData := &contacts.UpdateContactRequest{
    Email: "nuevo@email.com",
    Tags:  []string{"cliente", "vip"},
}

contact, err := client.Contacts().UpdateContact(ctx, "contact-id", updateData)

// Eliminar contacto
err := client.Contacts().DeleteContact(ctx, "contact-id")
```

#### Búsqueda y Filtrado

```go
// Búsqueda por nombre
results, err := client.Contacts().SearchContacts(ctx, "Juan")

// Filtrado avanzado
filter := &contacts.ContactFilter{
    CreatedAfter: time.Now().AddDate(0, -1, 0), // Últimos 30 días
}

results, err := client.Contacts().FilterContacts(ctx, filter)

// Obtener todos los contactos (con paginación automática)
allContacts, err := client.Contacts().GetAllContacts(ctx)
```

#### Operaciones en Lote

```go
// Agregar múltiples contactos
bulkContacts := []*contacts.CreateContactRequest{
    {WhatsappNumber: "1111111111", FirstName: "María", LastName: "García"},
    {WhatsappNumber: "2222222222", FirstName: "Carlos", LastName: "López"},
}

result, err := client.Contacts().AddContacts(ctx, bulkContacts)
fmt.Printf("Éxitos: %d, Fallos: %d\n", result.SuccessCount, result.FailureCount)
```

### 🤖 Chatbots

#### Gestión de Chatbots

```go
import "github.com/tu-usuario/go-wati/chatbots"

// Listar chatbots
chatbots, err := client.Chatbots().GetChatbots(ctx)

// Obtener chatbots activos
activeBots, err := client.Chatbots().GetActiveChatbots(ctx)

// Crear nuevo chatbot
newBot := &chatbots.CreateChatbotRequest{
    Name:        "Asistente de Ventas",
    Description: "Bot para consultas de ventas",
    Keywords:    []string{"ventas", "productos", "precio"},
    Responses: []chatbots.Response{
        {
            Trigger:  "ventas",
            Message:  "¡Hola! ¿En qué producto estás interesado?",
            IsActive: true,
        },
    },
    IsActive: true,
}

bot, err := client.Chatbots().CreateChatbot(ctx, newBot)
```

#### Control de Chatbots

```go
// Iniciar chatbot para un contacto
response, err := client.Chatbots().StartChatbotForContact(ctx, "bot-id", "1234567890")

// Iniciar con mensaje personalizado
response, err := client.Chatbots().StartChatbotWithMessage(
    ctx,
    "bot-id",
    "1234567890",
    "¡Hola! He activado el asistente automático.",
)

// Detener chatbot
err := client.Chatbots().StopChatbot(ctx, "bot-id")

// Activar/Desactivar chatbot
bot, err := client.Chatbots().ActivateChatbot(ctx, "bot-id")
bot, err := client.Chatbots().DeactivateChatbot(ctx, "bot-id")
```

#### Gestión de Estados de Chat

```go
// Asignar chat a agente
response, err := client.Chatbots().AssignChatToUser(ctx, "1234567890", "agente@empresa.com")

// Transferir a humano
response, err := client.Chatbots().TransferChatToHuman(
    ctx,
    "1234567890",
    "supervisor@empresa.com",
    "Cliente requiere atención especializada",
)

// Cerrar chat
response, err := client.Chatbots().CloseChatSession(
    ctx,
    "1234567890",
    "Consulta resuelta satisfactoriamente",
)

// Marcar como resuelto
response, err := client.Chatbots().ResolveChatSession(
    ctx,
    "1234567890",
    "Problema solucionado",
)
```

### 📁 Media

#### Subida de Archivos

```go
import (
    "os"
    "github.com/tu-usuario/go-wati/media"
)

// Subir imagen
file, err := os.Open("imagen.jpg")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

response, err := client.Media().UploadImage(
    ctx,
    file,
    "producto-destacado.jpg",
    "Nuestro producto más vendido",
)

// Subir documento
docFile, err := os.Open("catalogo.pdf")
if err != nil {
    log.Fatal(err)
}
defer docFile.Close()

response, err := client.Media().UploadDocument(
    ctx,
    docFile,
    "catalogo-productos.pdf",
    "Catálogo completo 2024",
)

// Subir video
videoFile, err := os.Open("demo.mp4")
if err != nil {
    log.Fatal(err)
}
defer videoFile.Close()

response, err := client.Media().UploadVideo(
    ctx,
    videoFile,
    "demo-producto.mp4",
    "Demostración del producto",
)
```

#### Gestión de Archivos

```go
// Listar archivos
mediaList, err := client.Media().ListMedia(ctx, &media.GetMediaParams{
    PageSize:  20,
    MediaType: string(media.MediaTypeImage),
})

// Obtener archivo específico
mediaFile, err := client.Media().GetMediaByFileName(ctx, "imagen.jpg")

// Obtener URL de archivo
url, err := client.Media().GetMediaURL(ctx, "imagen.jpg")

// Eliminar archivo
err := client.Media().DeleteMedia(ctx, "imagen.jpg")
```

#### Filtrado y Búsqueda

```go
// Obtener solo imágenes
images, err := client.Media().GetImages(ctx, nil)

// Obtener solo documentos
docs, err := client.Media().GetDocuments(ctx, nil)

// Buscar archivos por nombre
results, err := client.Media().SearchMedia(ctx, "producto", nil)

// Obtener estadísticas
stats, err := client.Media().GetMediaStats(ctx)
fmt.Printf("Total: %d, Imágenes: %d, Videos: %d\n", 
    stats.Stats.TotalFiles,
    stats.Stats.ImageCount,
    stats.Stats.VideoCount,
)
```

#### Validación de Archivos

```go
// Validar antes de subir
err := client.Media().ValidateUpload("archivo.pdf", 5*1024*1024, "application/pdf")
if err != nil {
    log.Printf("Archivo no válido: %v", err)
}

// Esperar a que el archivo esté listo
mediaFile, err := client.Media().WaitForMediaReady(ctx, "archivo.pdf", 60) // 60 segundos máximo
```

### 🔗 Webhooks

#### Configuración de Servidor

```go
import "github.com/tu-usuario/go-wati/webhooks"

// Configurar handlers
webhookService := client.Webhooks()

// Handler para mensajes recibidos
onMessage := func(data webhooks.MessageReceivedData) error {
    fmt.Printf("Mensaje de %s: %s\n", data.From, data.GetMessageText())
    
    // Lógica de respuesta automática
    if data.GetMessageText() == "hola" {
        return client.Messages().SendQuickReplyButtons(
            ctx,
            data.From,
            "¡Hola! ¿En qué puedo ayudarte?",
            []string{"Información", "Soporte", "Ventas"},
        )
    }
    
    return nil
}

// Registrar handlers
webhookService.RegisterHandler(
    webhooks.MessageReceived,
    webhooks.CreateMessageHandler(onMessage),
)

// Configurar secreto para validación
webhookService.SetSecret("mi-secreto-super-seguro")

// Iniciar servidor
err := webhookService.StartWebhookServer(8080, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Servidor de webhooks iniciado en puerto 8080")
```

#### Manejo de Eventos

```go
// Handler para estado de mensajes
onMessageStatus := func(data webhooks.MessageStatusData) error {
    switch data.Status {
    case "delivered":
        fmt.Printf("✅ Mensaje %s entregado\n", data.MessageID)
    case "read":
        fmt.Printf("👁️ Mensaje %s leído\n", data.MessageID)
    case "failed":
        fmt.Printf("❌ Mensaje %s falló: %s\n", data.MessageID, data.ErrorMessage)
    }
    return nil
}

// Handler para eventos de contacto
onContactEvent := func(data webhooks.ContactEventData) error {
    fmt.Printf("👤 Contacto %s: %s\n", data.ContactID, data.FullName)
    return nil
}

// Registrar múltiples handlers
webhookService.RegisterMessageHandlers(onMessage, onMessageStatus, onMessageStatus)
webhookService.RegisterHandler(webhooks.ContactCreated, webhooks.CreateContactHandler(onContactEvent))
```

#### Gestión de Webhooks en WATI

```go
// Registrar webhook en WATI
events := []webhooks.WebhookEventType{
    webhooks.MessageReceived,
    webhooks.MessageDelivered,
    webhooks.MessageRead,
    webhooks.ContactCreated,
}

err := webhookService.RegisterWebhook(ctx, "https://tu-dominio.com/webhook", events)

// Listar webhooks registrados
webhooksList, err := webhookService.ListWebhooks(ctx)

// Desregistrar webhook
err := webhookService.UnregisterWebhook(ctx, "https://tu-dominio.com/webhook")
```

## 📚 Ejemplos Completos

La librería incluye ejemplos completos en el directorio `examples/`:

- `basic_messaging.go` - Envío básico de mensajes
- `contact_management.go` - Gestión completa de contactos
- `webhook_server.go` - Servidor de webhooks completo
- `chatbots_and_media.go` - Chatbots y gestión de media

Para ejecutar un ejemplo:

```bash
cd examples
go run basic_messaging.go
```

## ⚙️ Configuración Avanzada

### Opciones de Cliente

```go
client := wati.NewClient(
    endpoint,
    token,
    // Timeout personalizado
    wati.WithTimeout(45), // 45 segundos
    
    // Reintentos automáticos
    wati.WithRetries(5), // 5 reintentos
    
    // Rate limiting personalizado
    wati.WithRateLimit(200, 60), // 200 requests por minuto
    
    // User agent personalizado
    wati.WithUserAgent("MiAplicacion/2.0"),
    
    // Cliente HTTP personalizado
    wati.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
    }),
)
```

### Configuración de Rate Limiting

```go
// Rate limiting automático basado en el plan de WATI
client := wati.NewClient(endpoint, token, wati.WithPlan(wati.PlanProfessional))

// Rate limiting personalizado
client := wati.NewClient(endpoint, token, wati.WithRateLimit(100, 60)) // 100 req/min
```

### Configuración de Logging

```go
// Habilitar logging detallado
client := wati.NewClient(endpoint, token, wati.WithDebug(true))

// Logger personalizado
client := wati.NewClient(endpoint, token, wati.WithLogger(myLogger))
```

## 🚨 Manejo de Errores

### Tipos de Error

La librería define varios tipos de error específicos:

```go
import "github.com/tu-usuario/go-wati"

// Verificar tipo de error
response, err := client.Messages().SendTemplateMessage(ctx, request)
if err != nil {
    switch e := err.(type) {
    case *wati.APIError:
        fmt.Printf("Error de API: %s (código: %d)\n", e.Message, e.Code)
    case *wati.RateLimitError:
        fmt.Printf("Rate limit excedido. Reintentar en: %v\n", e.RetryAfter)
    case *wati.ValidationError:
        fmt.Printf("Error de validación: %s\n", e.Message)
    case *wati.NetworkError:
        fmt.Printf("Error de red: %s\n", e.Message)
    default:
        fmt.Printf("Error desconocido: %v\n", err)
    }
}
```

### Reintentos Automáticos

```go
// Los reintentos se manejan automáticamente para errores temporales
client := wati.NewClient(endpoint, token, wati.WithRetries(3))

// Configuración avanzada de reintentos
retryConfig := &wati.RetryConfig{
    MaxRetries:    5,
    InitialDelay:  time.Second,
    MaxDelay:      30 * time.Second,
    Multiplier:    2.0,
    RetryableErrors: []int{429, 500, 502, 503, 504},
}

client := wati.NewClient(endpoint, token, wati.WithRetryConfig(retryConfig))
```

### Timeouts y Context

```go
// Timeout por operación
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := client.Messages().SendTemplateMessage(ctx, request)

// Cancelación manual
ctx, cancel := context.WithCancel(context.Background())

// En otra goroutine
go func() {
    time.Sleep(10 * time.Second)
    cancel() // Cancelar operación
}()

response, err := client.Messages().SendTemplateMessage(ctx, request)
```

## 💡 Mejores Prácticas

### 1. Gestión de Configuración

```go
// Usar variables de entorno
endpoint := os.Getenv("WATI_ENDPOINT")
token := os.Getenv("WATI_TOKEN")

if endpoint == "" || token == "" {
    log.Fatal("WATI_ENDPOINT y WATI_TOKEN son requeridos")
}

client := wati.NewClient(endpoint, token)
```

### 2. Manejo de Context

```go
// Siempre usar context con timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Propagar context en funciones
func sendWelcomeMessage(ctx context.Context, client wati.WATIClient, phone string) error {
    return client.Messages().SendSimpleTemplateMessage(ctx, phone, "welcome", "onboarding")
}
```

### 3. Validación de Datos

```go
// Validar antes de enviar
func sendMessage(client wati.WATIClient, phone, template string) error {
    if phone == "" {
        return fmt.Errorf("número de teléfono requerido")
    }
    
    if len(phone) < 10 {
        return fmt.Errorf("número de teléfono inválido")
    }
    
    ctx := context.Background()
    return client.Messages().SendSimpleTemplateMessage(ctx, phone, template, "default")
}
```

### 4. Logging y Monitoreo

```go
import "log/slog"

// Configurar logging estructurado
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// Log de operaciones importantes
func sendMessageWithLogging(client wati.WATIClient, phone, template string) error {
    logger.Info("Enviando mensaje",
        "phone", phone,
        "template", template,
    )
    
    ctx := context.Background()
    err := client.Messages().SendSimpleTemplateMessage(ctx, phone, template, "default")
    
    if err != nil {
        logger.Error("Error enviando mensaje",
            "phone", phone,
            "template", template,
            "error", err,
        )
        return err
    }
    
    logger.Info("Mensaje enviado exitosamente",
        "phone", phone,
        "template", template,
    )
    
    return nil
}
```

### 5. Gestión de Webhooks

```go
// Configuración robusta de webhooks
func setupWebhooks(client wati.WATIClient) error {
    webhookService := client.Webhooks()
    
    // Configurar secreto desde variable de entorno
    secret := os.Getenv("WEBHOOK_SECRET")
    if secret == "" {
        return fmt.Errorf("WEBHOOK_SECRET requerido")
    }
    webhookService.SetSecret(secret)
    
    // Handler con logging
    onMessage := func(data webhooks.MessageReceivedData) error {
        logger.Info("Mensaje recibido",
            "from", data.From,
            "type", data.MessageType,
            "text", data.GetMessageText(),
        )
        
        // Procesar mensaje...
        return nil
    }
    
    webhookService.RegisterHandler(webhooks.MessageReceived, webhooks.CreateMessageHandler(onMessage))
    
    // Iniciar servidor con manejo de errores
    port := 8080
    if err := webhookService.StartWebhookServer(port, nil); err != nil {
        return fmt.Errorf("error iniciando servidor de webhooks: %w", err)
    }
    
    logger.Info("Servidor de webhooks iniciado", "port", port)
    return nil
}
```

### 6. Operaciones en Lote

```go
// Procesar contactos en lotes
func addContactsInBatches(client wati.WATIClient, contacts []*contacts.CreateContactRequest) error {
    const batchSize = 50 // WATI permite hasta 100, pero usamos 50 para ser conservadores
    
    ctx := context.Background()
    
    for i := 0; i < len(contacts); i += batchSize {
        end := i + batchSize
        if end > len(contacts) {
            end = len(contacts)
        }
        
        batch := contacts[i:end]
        
        logger.Info("Procesando lote de contactos",
            "batch", i/batchSize+1,
            "size", len(batch),
        )
        
        result, err := client.Contacts().AddContacts(ctx, batch)
        if err != nil {
            return fmt.Errorf("error en lote %d: %w", i/batchSize+1, err)
        }
        
        logger.Info("Lote procesado",
            "success", result.SuccessCount,
            "failures", result.FailureCount,
        )
        
        // Pausa entre lotes para evitar rate limiting
        time.Sleep(1 * time.Second)
    }
    
    return nil
}
```

## 🔧 Desarrollo y Testing

### Configuración de Desarrollo

```bash
# Clonar repositorio
git clone https://github.com/tu-usuario/go-wati.git
cd go-wati

# Instalar dependencias
go mod download

# Ejecutar tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...

# Ejecutar linter
golangci-lint run
```

### Testing

```go
// Ejemplo de test unitario
func TestSendTemplateMessage(t *testing.T) {
    client := wati.NewClient("https://test.wati.io", "test-token")
    
    ctx := context.Background()
    request := &messages.SendTemplateMessageRequest{
        WhatsappNumber: "1234567890",
        TemplateName:   "test_template",
        BroadcastName:  "test_broadcast",
    }
    
    response, err := client.Messages().SendTemplateMessage(ctx, request)
    
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, "1234567890", response.PhoneNumber)
}
```

### Mocking para Tests

```go
// Interface para mocking
type MockWATIClient struct {
    messages  MockMessagesService
    contacts  MockContactsService
    // ...
}

func (m *MockWATIClient) Messages() MessagesService {
    return &m.messages
}

// Usar en tests
func TestBusinessLogic(t *testing.T) {
    mockClient := &MockWATIClient{}
    
    // Configurar expectativas...
    
    err := businessFunction(mockClient)
    assert.NoError(t, err)
}
```

## 📄 Licencia

Este proyecto está licenciado bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 🤝 Contribución

¡Las contribuciones son bienvenidas! Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Abre un Pull Request

### Guías de Contribución

- Seguir las convenciones de código de Go
- Agregar tests para nuevas funcionalidades
- Actualizar documentación cuando sea necesario
- Usar commits descriptivos

## 📞 Soporte

- **Documentación**: [Documentación completa](https://github.com/tu-usuario/go-wati/docs)
- **Issues**: [GitHub Issues](https://github.com/tu-usuario/go-wati/issues)
- **Discusiones**: [GitHub Discussions](https://github.com/tu-usuario/go-wati/discussions)

## 🔗 Enlaces Útiles

- [Documentación oficial de WATI API](https://docs.wati.io)
- [WhatsApp Business API](https://developers.facebook.com/docs/whatsapp)
- [Go Documentation](https://golang.org/doc/)

---

**Desarrollado con ❤️ para la comunidad Go**

