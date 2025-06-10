package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tu-usuario/go-wati"
	"github.com/tu-usuario/go-wati/messages"
)

func main() {
	// Configurar cliente WATI
	client := wati.NewClient(
		"https://live-server-12345.wati.io", // Tu endpoint de WATI
		"tu-token-aqui",                     // Tu token de API
		wati.WithTimeout(30),                // Timeout de 30 segundos
		wati.WithRetries(3),                 // 3 reintentos
	)

	ctx := context.Background()

	// Ejemplo 1: Enviar mensaje de plantilla simple
	fmt.Println("=== Enviando mensaje de plantilla simple ===")
	
	response, err := client.Messages().SendTemplateMessage(ctx, &messages.SendTemplateMessageRequest{
		WhatsappNumber: "1234567890",
		TemplateName:   "hello_world",
		BroadcastName:  "mi_broadcast",
	})
	
	if err != nil {
		log.Printf("Error enviando mensaje: %v", err)
	} else {
		fmt.Printf("Mensaje enviado exitosamente. ID: %s\n", response.PhoneNumber)
	}

	// Ejemplo 2: Enviar mensaje de plantilla con parámetros
	fmt.Println("\n=== Enviando mensaje con parámetros ===")
	
	templateParams := map[string]string{
		"name":    "Juan Pérez",
		"product": "Smartphone XYZ",
		"price":   "$299.99",
	}
	
	response2, err := client.Messages().SendTemplateMessageWithParams(
		ctx,
		"1234567890",
		"order_confirmation",
		"confirmaciones",
		templateParams,
	)
	
	if err != nil {
		log.Printf("Error enviando mensaje con parámetros: %v", err)
	} else {
		fmt.Printf("Mensaje con parámetros enviado. Válido: %v\n", response2.ValidWhatsAppNumber)
	}

	// Ejemplo 3: Enviar botones de respuesta rápida
	fmt.Println("\n=== Enviando botones de respuesta rápida ===")
	
	buttonTitles := []string{"Sí, me interesa", "No, gracias", "Más información"}
	
	response3, err := client.Messages().SendQuickReplyButtons(
		ctx,
		"1234567890",
		"¿Te interesa nuestro nuevo producto?",
		buttonTitles,
	)
	
	if err != nil {
		log.Printf("Error enviando botones: %v", err)
	} else {
		fmt.Printf("Botones enviados exitosamente\n")
	}

	// Ejemplo 4: Enviar menú de lista
	fmt.Println("\n=== Enviando menú de lista ===")
	
	menuItems := map[string][]string{
		"Productos": {"Smartphone", "Tablet", "Laptop"},
		"Servicios": {"Soporte técnico", "Garantía", "Reparación"},
		"Información": {"Horarios", "Ubicación", "Contacto"},
	}
	
	response4, err := client.Messages().SendListMenu(
		ctx,
		"1234567890",
		"¿En qué podemos ayudarte hoy?",
		"Ver opciones",
		menuItems,
	)
	
	if err != nil {
		log.Printf("Error enviando menú: %v", err)
	} else {
		fmt.Printf("Menú enviado exitosamente\n")
	}

	// Ejemplo 5: Obtener plantillas disponibles
	fmt.Println("\n=== Obteniendo plantillas disponibles ===")
	
	templates, err := client.Messages().GetActiveTemplates(ctx)
	if err != nil {
		log.Printf("Error obteniendo plantillas: %v", err)
	} else {
		fmt.Printf("Plantillas activas encontradas: %d\n", len(templates))
		for _, template := range templates {
			fmt.Printf("- %s (%s)\n", template.Name, template.Language)
		}
	}

	// Ejemplo 6: Obtener historial de mensajes
	fmt.Println("\n=== Obteniendo historial de mensajes ===")
	
	messagesHistory, err := client.Messages().GetMessagesByPhone(ctx, "1234567890", &messages.GetMessagesParams{
		PageSize: 10,
	})
	
	if err != nil {
		log.Printf("Error obteniendo historial: %v", err)
	} else {
		fmt.Printf("Mensajes encontrados: %d\n", len(messagesHistory.Messages))
		for _, msg := range messagesHistory.Messages {
			fmt.Printf("- %s: %s (%s)\n", msg.From, msg.Content, msg.Status)
		}
	}

	fmt.Println("\n=== Ejemplo completado ===")
}

