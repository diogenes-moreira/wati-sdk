package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tu-usuario/go-wati"
	"github.com/tu-usuario/go-wati/webhooks"
)

func main() {
	// Configurar cliente WATI
	client := wati.NewClient(
		"https://live-server-12345.wati.io",
		"tu-token-aqui",
		wati.WithTimeout(30),
	)

	ctx := context.Background()

	// Ejemplo 1: Configurar handlers para diferentes tipos de eventos
	fmt.Println("=== Configurando handlers de webhooks ===")

	// Handler para mensajes recibidos
	onMessageReceived := func(data webhooks.MessageReceivedData) error {
		fmt.Printf("\n📨 Mensaje recibido de %s:\n", data.From)
		fmt.Printf("   Tipo: %s\n", data.MessageType)
		
		if data.IsTextMessage() {
			fmt.Printf("   Texto: %s\n", data.Text)
		} else if data.IsMediaMessage() {
			fmt.Printf("   Media: %s (%s)\n", data.Media.FileName, data.Media.MimeType)
		} else if data.IsInteractiveMessage() {
			if data.IsButtonReply() {
				fmt.Printf("   Botón seleccionado: %s\n", data.Interactive.ButtonReply.Title)
			} else if data.IsListReply() {
				fmt.Printf("   Opción de lista: %s\n", data.Interactive.ListReply.Title)
			}
		}
		
		if contactName := data.GetContactName(); contactName != "" {
			fmt.Printf("   Contacto: %s\n", contactName)
		}
		
		// Aquí puedes agregar lógica de respuesta automática
		// Por ejemplo, responder a ciertos mensajes o palabras clave
		
		return nil
	}

	// Handler para estado de mensajes
	onMessageStatus := func(data webhooks.MessageStatusData) error {
		fmt.Printf("\n📊 Estado de mensaje %s: %s\n", data.MessageID, data.Status)
		
		if data.Status == "delivered" {
			fmt.Printf("   ✅ Mensaje entregado a %s\n", data.To)
		} else if data.Status == "read" {
			fmt.Printf("   👁️ Mensaje leído por %s\n", data.To)
		} else if data.Status == "failed" {
			fmt.Printf("   ❌ Mensaje falló: %s\n", data.ErrorMessage)
		}
		
		return nil
	}

	// Handler para eventos de contacto
	onContactEvent := func(data webhooks.ContactEventData) error {
		fmt.Printf("\n👤 Evento de contacto: %s\n", data.ContactID)
		fmt.Printf("   Nombre: %s\n", data.FullName)
		fmt.Printf("   Teléfono: %s\n", data.WhatsappNumber)
		
		if len(data.Tags) > 0 {
			fmt.Printf("   Tags: %v\n", data.Tags)
		}
		
		return nil
	}

	// Handler para eventos de chatbot
	onChatbotEvent := func(data webhooks.ChatbotEventData) error {
		fmt.Printf("\n🤖 Evento de chatbot: %s\n", data.ChatbotName)
		fmt.Printf("   Estado: %s\n", data.Status)
		fmt.Printf("   Usuario: %s\n", data.WhatsappNumber)
		
		if data.Reason != "" {
			fmt.Printf("   Razón: %s\n", data.Reason)
		}
		
		return nil
	}

	// Handler para cambios de estado de chat
	onChatStatusChange := func(data webhooks.ChatStatusEventData) error {
		fmt.Printf("\n💬 Cambio de estado de chat: %s → %s\n", data.OldStatus, data.NewStatus)
		fmt.Printf("   Usuario: %s\n", data.WhatsappNumber)
		
		if data.AssignedTo != "" {
			fmt.Printf("   Asignado a: %s\n", data.AssignedTo)
		}
		
		return nil
	}

	// Ejemplo 2: Registrar handlers en el servicio de webhooks
	webhookService := client.Webhooks()
	
	// Configurar secreto para validación de firmas
	webhookService.SetSecret("mi-secreto-super-seguro")
	
	// Registrar handlers individuales
	webhookService.RegisterHandler(webhooks.MessageReceived, webhooks.CreateMessageHandler(onMessageReceived))
	webhookService.RegisterHandler(webhooks.NewContactMessage, webhooks.CreateMessageHandler(onMessageReceived))
	webhookService.RegisterHandler(webhooks.MessageDelivered, webhooks.CreateMessageStatusHandler(onMessageStatus))
	webhookService.RegisterHandler(webhooks.MessageRead, webhooks.CreateMessageStatusHandler(onMessageStatus))
	webhookService.RegisterHandler(webhooks.ContactCreated, webhooks.CreateContactHandler(onContactEvent))
	webhookService.RegisterHandler(webhooks.ContactUpdated, webhooks.CreateContactHandler(onContactEvent))
	webhookService.RegisterHandler(webhooks.ChatbotStarted, webhooks.CreateChatbotHandler(onChatbotEvent))
	webhookService.RegisterHandler(webhooks.ChatbotStopped, webhooks.CreateChatbotHandler(onChatbotEvent))
	webhookService.RegisterHandler(webhooks.ChatStatusChanged, webhooks.CreateChatStatusHandler(onChatStatusChange))

	// Ejemplo 3: Iniciar servidor de webhooks
	fmt.Println("\n=== Iniciando servidor de webhooks ===")
	
	port := 8080
	err := webhookService.StartWebhookServer(port, nil) // Los handlers ya están registrados
	if err != nil {
		log.Fatalf("Error iniciando servidor de webhooks: %v", err)
	}
	
	fmt.Printf("🚀 Servidor de webhooks iniciado en puerto %d\n", port)
	fmt.Printf("📡 Endpoint: http://localhost:%d/webhook\n", port)
	fmt.Printf("🏥 Health check: http://localhost:%d/health\n", port)

	// Ejemplo 4: Registrar webhook en WATI (opcional)
	fmt.Println("\n=== Registrando webhook en WATI ===")
	
	webhookURL := fmt.Sprintf("https://tu-dominio.com/webhook") // Cambia por tu URL pública
	events := []webhooks.WebhookEventType{
		webhooks.MessageReceived,
		webhooks.NewContactMessage,
		webhooks.MessageDelivered,
		webhooks.MessageRead,
		webhooks.ContactCreated,
		webhooks.ContactUpdated,
		webhooks.ChatbotStarted,
		webhooks.ChatbotStopped,
		webhooks.ChatStatusChanged,
	}
	
	err = webhookService.RegisterWebhook(ctx, webhookURL, events)
	if err != nil {
		log.Printf("⚠️ Error registrando webhook en WATI: %v", err)
		log.Printf("   Esto es normal si no tienes una URL pública configurada")
	} else {
		fmt.Printf("✅ Webhook registrado exitosamente en WATI\n")
	}

	// Ejemplo 5: Listar webhooks registrados
	fmt.Println("\n=== Listando webhooks registrados ===")
	
	webhooksList, err := webhookService.ListWebhooks(ctx)
	if err != nil {
		log.Printf("Error obteniendo lista de webhooks: %v", err)
	} else {
		fmt.Printf("Webhooks registrados: %d\n", len(webhooksList.Webhooks))
		for _, webhook := range webhooksList.Webhooks {
			fmt.Printf("- %s (activo: %v)\n", webhook.URL, webhook.IsActive)
			fmt.Printf("  Eventos: %v\n", webhook.Events)
		}
	}

	// Ejemplo 6: Mostrar información del servidor
	fmt.Println("\n=== Estado del servidor ===")
	fmt.Printf("Puerto: %d\n", webhookService.GetServerPort())
	fmt.Printf("Estado: %v\n", webhookService.GetServerStatus())

	// Ejemplo 7: Configurar manejo de señales para cierre limpio
	fmt.Println("\n=== Servidor listo para recibir webhooks ===")
	fmt.Println("Presiona Ctrl+C para detener el servidor")
	
	// Canal para señales del sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine para mostrar estadísticas periódicas
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				if webhookService.GetServerStatus() {
					fmt.Printf("\n📊 Servidor activo - Puerto: %d\n", webhookService.GetServerPort())
				}
			case <-sigChan:
				return
			}
		}
	}()

	// Esperar señal de cierre
	<-sigChan
	
	fmt.Println("\n🛑 Cerrando servidor de webhooks...")
	
	// Detener servidor de webhooks
	err = webhookService.StopWebhookServer()
	if err != nil {
		log.Printf("Error deteniendo servidor: %v", err)
	} else {
		fmt.Println("✅ Servidor detenido exitosamente")
	}

	// Opcional: Desregistrar webhook de WATI
	if webhookURL != "" {
		fmt.Println("🗑️ Desregistrando webhook de WATI...")
		err = webhookService.UnregisterWebhook(ctx, webhookURL)
		if err != nil {
			log.Printf("Error desregistrando webhook: %v", err)
		} else {
			fmt.Println("✅ Webhook desregistrado exitosamente")
		}
	}

	fmt.Println("\n=== Ejemplo de webhooks completado ===")
}

// Ejemplo de función auxiliar para procesar mensajes automáticamente
func processIncomingMessage(client wati.WATIClient, data webhooks.MessageReceivedData) error {
	ctx := context.Background()
	
	// Ejemplo de respuesta automática basada en el contenido del mensaje
	messageText := data.GetMessageText()
	
	switch messageText {
	case "hola", "Hola", "HOLA":
		// Responder con botones de opciones
		return client.Messages().SendQuickReplyButtons(
			ctx,
			data.From,
			"¡Hola! ¿En qué puedo ayudarte?",
			[]string{"Información", "Soporte", "Ventas"},
		)
		
	case "menu", "Menu", "MENU":
		// Responder with lista de opciones
		menuItems := map[string][]string{
			"Productos": {"Smartphones", "Tablets", "Laptops"},
			"Servicios": {"Soporte", "Garantía", "Reparación"},
			"Empresa": {"Sobre nosotros", "Contacto", "Ubicación"},
		}
		
		return client.Messages().SendListMenu(
			ctx,
			data.From,
			"¿Qué te interesa conocer?",
			"Ver opciones",
			menuItems,
		)
		
	case "gracias", "Gracias", "GRACIAS":
		// Respuesta simple
		return client.Messages().SendSimpleTemplateMessage(
			ctx,
			data.From,
			"thank_you",
			"agradecimientos",
		)
	}
	
	return nil
}

