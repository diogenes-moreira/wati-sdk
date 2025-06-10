# Go WATI - Librería Go para WATI API

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-available-brightgreen.svg)](DOCUMENTATION.md)

Una librería completa y robusta en Go para interactuar con la API de WATI (WhatsApp Business API). Proporciona todas las funcionalidades necesarias para gestionar mensajes, contactos, chatbots, media y webhooks de manera sencilla y eficiente.

## 🚀 Características Principales

- **📨 Mensajería Completa**: Mensajes de plantilla, interactivos (botones/listas), y gestión de historial
- **👥 Gestión de Contactos**: CRUD completo, búsqueda avanzada y operaciones en lote
- **🤖 Chatbots Inteligentes**: Creación, gestión y control automatizado de chatbots
- **📁 Gestión de Media**: Subida, descarga y gestión de archivos multimedia
- **🔗 Webhooks Avanzados**: Servidor integrado con manejo de eventos en tiempo real
- **🛡️ Seguridad**: Autenticación OAuth2, validación de firmas HMAC-SHA256
- **⚡ Performance**: Rate limiting automático, reintentos inteligentes, timeouts configurables
- **🔧 Developer-Friendly**: Tipado fuerte, validaciones completas, logging detallado

## 📦 Instalación

```bash
go get github.com/tu-usuario/go-wati
```

## 🏃‍♂️ Inicio Rápido

```go
package main

import (
    "context"
    "log"
    "github.com/tu-usuario/go-wati"
    "github.com/tu-usuario/go-wati/messages"
)

func main() {
    // Crear cliente
    client := wati.NewClient(
        "https://live-server-12345.wati.io", // Tu endpoint de WATI
        "tu-token-aqui",                     // Tu token de API
        wati.WithTimeout(30),                // Configuraciones opcionales
        wati.WithRetries(3),
    )

    ctx := context.Background()

    // Enviar mensaje de plantilla
    response, err := client.Messages().SendTemplateMessage(ctx, &messages.SendTemplateMessageRequest{
        WhatsappNumber: "1234567890",
        TemplateName:   "hello_world",
        BroadcastName:  "mi_broadcast",
    })

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Mensaje enviado exitosamente: %s", response.PhoneNumber)
}
```

## 📚 Ejemplos de Uso

### Envío de Mensajes Interactivos

```go
// Botones de respuesta rápida
client.Messages().SendQuickReplyButtons(
    ctx,
    "1234567890",
    "¿Te interesa nuestro producto?",
    []string{"Sí, me interesa", "No, gracias", "Más información"},
)

// Menú de lista
menuItems := map[string][]string{
    "Productos": {"Smartphone", "Tablet", "Laptop"},
    "Servicios": {"Soporte", "Garantía", "Reparación"},
}

client.Messages().SendListMenu(
    ctx,
    "1234567890",
    "¿En qué podemos ayudarte?",
    "Ver opciones",
    menuItems,
)
```

### Gestión de Contactos

```go
// Crear contacto
newContact := &contacts.CreateContactRequest{
    WhatsappNumber: "1234567890",
    FirstName:      "Juan",
    LastName:       "Pérez",
    Email:          "juan@email.com",
    Tags:           []string{"cliente", "premium"},
}

contact, err := client.Contacts().AddContact(ctx, newContact)

// Buscar contactos
results, err := client.Contacts().SearchContacts(ctx, "Juan")
```

### Servidor de Webhooks

```go
// Configurar handlers
webhookService := client.Webhooks()

onMessage := func(data webhooks.MessageReceivedData) error {
    fmt.Printf("Mensaje de %s: %s\n", data.From, data.GetMessageText())
    return nil
}

webhookService.RegisterHandler(
    webhooks.MessageReceived,
    webhooks.CreateMessageHandler(onMessage),
)

// Iniciar servidor
webhookService.StartWebhookServer(8080, nil)
```

### Gestión de Chatbots

```go
// Crear chatbot
newBot := &chatbots.CreateChatbotRequest{
    Name:        "Asistente de Ventas",
    Keywords:    []string{"ventas", "productos", "precio"},
    IsActive:    true,
}

bot, err := client.Chatbots().CreateChatbot(ctx, newBot)

// Iniciar chatbot para un contacto
client.Chatbots().StartChatbotForContact(ctx, bot.ID, "1234567890")
```

## 📖 Documentación Completa

Para documentación detallada, ejemplos completos y guías avanzadas, consulta:

- **[📚 Documentación Completa](DOCUMENTATION.md)** - Guía completa con todos los detalles
- **[📁 Ejemplos](examples/)** - Ejemplos prácticos y casos de uso
- **[🔧 API Reference](https://pkg.go.dev/github.com/tu-usuario/go-wati)** - Referencia completa de la API

## 🏗️ Arquitectura

```
go-wati/
├── client.go          # Cliente principal y configuración
├── config.go          # Configuraciones y opciones
├── types.go           # Tipos comunes
├── errors.go          # Manejo de errores
├── interfaces.go      # Interfaces de servicios
├── contacts/          # Módulo de gestión de contactos
├── messages/          # Módulo de mensajería
├── chatbots/          # Módulo de chatbots
├── media/             # Módulo de gestión de media
├── webhooks/          # Módulo de webhooks
├── examples/          # Ejemplos de uso
└── tests/             # Tests unitarios e integración
```

## 🛠️ Configuración Avanzada

```go
client := wati.NewClient(
    endpoint,
    token,
    wati.WithTimeout(45),              // Timeout personalizado
    wati.WithRetries(5),               // Reintentos automáticos
    wati.WithRateLimit(200, 60),       // Rate limiting personalizado
    wati.WithUserAgent("MiApp/2.0"),   // User agent personalizado
    wati.WithDebug(true),              // Logging detallado
)
```

## 🚨 Manejo de Errores

```go
response, err := client.Messages().SendTemplateMessage(ctx, request)
if err != nil {
    switch e := err.(type) {
    case *wati.APIError:
        log.Printf("Error de API: %s (código: %d)", e.Message, e.Code)
    case *wati.RateLimitError:
        log.Printf("Rate limit excedido. Reintentar en: %v", e.RetryAfter)
    case *wati.ValidationError:
        log.Printf("Error de validación: %s", e.Message)
    default:
        log.Printf("Error: %v", err)
    }
}
```

## 📊 Características Técnicas

- **Go Version**: 1.21+
- **Dependencias Mínimas**: Solo `golang.org/x/time` para rate limiting
- **Thread-Safe**: Seguro para uso concurrente
- **Context Support**: Soporte completo para context.Context
- **Rate Limiting**: Automático basado en planes de WATI
- **Reintentos**: Inteligentes con backoff exponencial
- **Validaciones**: Completas antes del envío a la API
- **Testing**: Cobertura completa con mocks incluidos

## 🤝 Contribución

¡Las contribuciones son bienvenidas! Por favor lee nuestras [guías de contribución](CONTRIBUTING.md).

1. Fork el proyecto
2. Crea tu feature branch (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la branch (`git push origin feature/nueva-funcionalidad`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está licenciado bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 🔗 Enlaces

- [Documentación de WATI API](https://docs.wati.io)
- [WhatsApp Business API](https://developers.facebook.com/docs/whatsapp)
- [Issues y Soporte](https://github.com/tu-usuario/go-wati/issues)

---

**Desarrollado con ❤️ para la comunidad Go**

