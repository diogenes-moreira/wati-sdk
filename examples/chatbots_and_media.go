package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tu-usuario/go-wati"
	"github.com/tu-usuario/go-wati/chatbots"
	"github.com/tu-usuario/go-wati/media"
)

func main() {
	// Configurar cliente WATI
	client := wati.NewClient(
		"https://live-server-12345.wati.io",
		"tu-token-aqui",
		wati.WithTimeout(30),
	)

	ctx := context.Background()

	// PARTE 1: GESTIÓN DE MEDIA
	fmt.Println("=== GESTIÓN DE MEDIA ===")

	// Ejemplo 1: Subir una imagen
	fmt.Println("\n📸 Subiendo imagen...")
	
	// Simular archivo de imagen (en la práctica, abrirías un archivo real)
	imageFile := strings.NewReader("contenido-simulado-de-imagen")
	
	uploadResponse, err := client.Media().UploadImage(
		ctx,
		imageFile,
		"producto-destacado.jpg",
		"Nuestro producto más vendido",
	)
	
	if err != nil {
		log.Printf("Error subiendo imagen: %v", err)
	} else {
		fmt.Printf("✅ Imagen subida: %s\n", uploadResponse.Media.FileName)
		fmt.Printf("   URL: %s\n", uploadResponse.Media.URL)
		fmt.Printf("   Tamaño: %s\n", uploadResponse.Media.FormatFileSize())
	}

	// Ejemplo 2: Subir un documento
	fmt.Println("\n📄 Subiendo documento...")
	
	docFile := strings.NewReader("contenido-simulado-de-documento-pdf")
	
	docResponse, err := client.Media().UploadDocument(
		ctx,
		docFile,
		"catalogo-productos.pdf",
		"Catálogo completo de productos 2024",
	)
	
	if err != nil {
		log.Printf("Error subiendo documento: %v", err)
	} else {
		fmt.Printf("✅ Documento subido: %s\n", docResponse.Media.FileName)
	}

	// Ejemplo 3: Listar archivos de media
	fmt.Println("\n📁 Listando archivos de media...")
	
	mediaList, err := client.Media().ListMedia(ctx, &media.GetMediaParams{
		PageSize: 10,
	})
	
	if err != nil {
		log.Printf("Error listando media: %v", err)
	} else {
		fmt.Printf("Archivos encontrados: %d\n", len(mediaList.Media))
		for _, file := range mediaList.Media {
			fmt.Printf("- %s (%s) - %s\n", 
				file.FileName, 
				file.FormatFileSize(),
				file.Status,
			)
		}
	}

	// Ejemplo 4: Obtener estadísticas de media
	fmt.Println("\n📊 Estadísticas de media...")
	
	stats, err := client.Media().GetMediaStats(ctx)
	if err != nil {
		log.Printf("Error obteniendo estadísticas: %v", err)
	} else {
		fmt.Printf("Total de archivos: %d\n", stats.Stats.TotalFiles)
		fmt.Printf("Imágenes: %d\n", stats.Stats.ImageCount)
		fmt.Printf("Videos: %d\n", stats.Stats.VideoCount)
		fmt.Printf("Documentos: %d\n", stats.Stats.DocumentCount)
		fmt.Printf("Audio: %d\n", stats.Stats.AudioCount)
	}

	// Ejemplo 5: Buscar archivos por tipo
	fmt.Println("\n🔍 Buscando imágenes...")
	
	images, err := client.Media().GetImages(ctx, &media.GetMediaParams{
		PageSize: 5,
	})
	
	if err != nil {
		log.Printf("Error obteniendo imágenes: %v", err)
	} else {
		fmt.Printf("Imágenes encontradas: %d\n", len(images.Media))
		for _, img := range images.Media {
			fmt.Printf("- %s (%dx%d)\n", img.FileName, img.Width, img.Height)
		}
	}

	// PARTE 2: GESTIÓN DE CHATBOTS
	fmt.Println("\n\n=== GESTIÓN DE CHATBOTS ===")

	// Ejemplo 6: Listar chatbots disponibles
	fmt.Println("\n🤖 Listando chatbots...")
	
	chatbotsList, err := client.Chatbots().GetChatbots(ctx)
	if err != nil {
		log.Printf("Error obteniendo chatbots: %v", err)
	} else {
		fmt.Printf("Chatbots encontrados: %d\n", len(chatbotsList.Chatbots))
		for _, bot := range chatbotsList.Chatbots {
			status := "❌ Inactivo"
			if bot.IsActive() {
				status = "✅ Activo"
			}
			fmt.Printf("- %s (%s) %s\n", bot.Name, bot.ID, status)
			if len(bot.Keywords) > 0 {
				fmt.Printf("  Palabras clave: %v\n", bot.Keywords)
			}
		}
	}

	// Ejemplo 7: Crear un nuevo chatbot
	fmt.Println("\n➕ Creando nuevo chatbot...")
	
	newChatbot := &chatbots.CreateChatbotRequest{
		Name:        "Asistente de Ventas",
		Description: "Chatbot para atender consultas de ventas y productos",
		Keywords:    []string{"ventas", "productos", "precio", "comprar"},
		Responses: []chatbots.Response{
			{
				Trigger:  "ventas",
				Message:  "¡Hola! Soy tu asistente de ventas. ¿En qué producto estás interesado?",
				IsActive: true,
			},
			{
				Trigger:  "precio",
				Message:  "Te ayudo con información de precios. ¿Qué producto te interesa?",
				IsActive: true,
			},
		},
		IsActive: true,
	}
	
	createdBot, err := client.Chatbots().CreateChatbot(ctx, newChatbot)
	if err != nil {
		log.Printf("Error creando chatbot: %v", err)
	} else {
		fmt.Printf("✅ Chatbot creado: %s (ID: %s)\n", createdBot.Name, createdBot.ID)
	}

	// Ejemplo 8: Obtener chatbots activos
	fmt.Println("\n🟢 Obteniendo chatbots activos...")
	
	activeBots, err := client.Chatbots().GetActiveChatbots(ctx)
	if err != nil {
		log.Printf("Error obteniendo chatbots activos: %v", err)
	} else {
		fmt.Printf("Chatbots activos: %d\n", len(activeBots))
		for _, bot := range activeBots {
			fmt.Printf("- %s\n", bot.Name)
			activeRules := bot.GetActiveRules()
			fmt.Printf("  Reglas activas: %d\n", len(activeRules))
		}
	}

	// Ejemplo 9: Iniciar chatbot para un contacto
	fmt.Println("\n🚀 Iniciando chatbot para contacto...")
	
	if len(activeBots) > 0 {
		startRequest := &chatbots.StartChatbotRequest{
			ChatbotID:      activeBots[0].ID,
			WhatsappNumber: "1234567890",
			InitialMessage: "¡Hola! He activado el asistente automático para ayudarte.",
		}
		
		startResponse, err := client.Chatbots().StartChatbot(ctx, startRequest)
		if err != nil {
			log.Printf("Error iniciando chatbot: %v", err)
		} else {
			fmt.Printf("✅ Chatbot iniciado para contacto\n")
			fmt.Printf("   Estado: %s\n", startResponse.Status)
			if startResponse.SessionID != "" {
				fmt.Printf("   Sesión: %s\n", startResponse.SessionID)
			}
		}
	}

	// Ejemplo 10: Gestionar estado de chat
	fmt.Println("\n💬 Gestionando estado de chat...")
	
	// Asignar chat a un agente humano
	assignRequest := &chatbots.UpdateChatStatusRequest{
		WhatsappNumber: "1234567890",
		Status:         string(chatbots.ChatStatusAssigned),
		AssignedTo:     "agente@empresa.com",
		Notes:          "Cliente requiere atención personalizada",
		Tags:           []string{"vip", "consulta_compleja"},
	}
	
	statusResponse, err := client.Chatbots().UpdateChatStatus(ctx, assignRequest)
	if err != nil {
		log.Printf("Error actualizando estado de chat: %v", err)
	} else {
		fmt.Printf("✅ Chat asignado a agente\n")
		fmt.Printf("   Estado: %s\n", statusResponse.Status)
		fmt.Printf("   Asignado a: %s\n", statusResponse.AssignedTo)
	}

	// Ejemplo 11: Buscar chatbot por nombre
	fmt.Println("\n🔍 Buscando chatbot por nombre...")
	
	foundBot, err := client.Chatbots().GetChatbotByName(ctx, "Asistente de Ventas")
	if err != nil {
		log.Printf("Error buscando chatbot: %v", err)
	} else {
		fmt.Printf("✅ Chatbot encontrado: %s\n", foundBot.Name)
		fmt.Printf("   Descripción: %s\n", foundBot.Description)
		fmt.Printf("   Respuestas activas: %d\n", len(foundBot.GetActiveResponses()))
	}

	// Ejemplo 12: Actualizar palabras clave de chatbot
	fmt.Println("\n🔧 Actualizando chatbot...")
	
	if foundBot != nil {
		newKeywords := []string{"ventas", "productos", "precio", "comprar", "ofertas", "descuentos"}
		
		updatedBot, err := client.Chatbots().UpdateChatbotKeywords(ctx, foundBot.ID, newKeywords)
		if err != nil {
			log.Printf("Error actualizando palabras clave: %v", err)
		} else {
			fmt.Printf("✅ Palabras clave actualizadas\n")
			fmt.Printf("   Nuevas palabras: %v\n", updatedBot.Keywords)
		}
	}

	// Ejemplo 13: Transferir chat a humano
	fmt.Println("\n👤 Transfiriendo chat a humano...")
	
	transferResponse, err := client.Chatbots().TransferChatToHuman(
		ctx,
		"1234567890",
		"supervisor@empresa.com",
		"Cliente solicita hablar con supervisor",
	)
	
	if err != nil {
		log.Printf("Error transfiriendo chat: %v", err)
	} else {
		fmt.Printf("✅ Chat transferido exitosamente\n")
		fmt.Printf("   Nuevo estado: %s\n", transferResponse.Status)
	}

	// Ejemplo 14: Cerrar sesión de chat
	fmt.Println("\n🔚 Cerrando sesión de chat...")
	
	closeResponse, err := client.Chatbots().CloseChatSession(
		ctx,
		"1234567890",
		"Consulta resuelta satisfactoriamente",
	)
	
	if err != nil {
		log.Printf("Error cerrando chat: %v", err)
	} else {
		fmt.Printf("✅ Chat cerrado exitosamente\n")
		fmt.Printf("   Estado final: %s\n", closeResponse.Status)
	}

	// PARTE 3: INTEGRACIÓN MEDIA + CHATBOTS
	fmt.Println("\n\n=== INTEGRACIÓN AVANZADA ===")

	// Ejemplo 15: Validar archivo antes de subir
	fmt.Println("\n🔍 Validando archivo...")
	
	err = client.Media().ValidateUpload("catalogo.pdf", 2*1024*1024, "application/pdf")
	if err != nil {
		log.Printf("❌ Archivo no válido: %v", err)
	} else {
		fmt.Printf("✅ Archivo válido para subir\n")
	}

	// Ejemplo 16: Buscar chatbots por palabra clave
	fmt.Println("\n🔍 Buscando chatbots por palabra clave...")
	
	keywordBots, err := client.Chatbots().GetChatbotsByKeyword(ctx, "ventas")
	if err != nil {
		log.Printf("Error buscando por palabra clave: %v", err)
	} else {
		fmt.Printf("Chatbots con palabra 'ventas': %d\n", len(keywordBots))
		for _, bot := range keywordBots {
			fmt.Printf("- %s\n", bot.Name)
		}
	}

	fmt.Println("\n=== Ejemplo de chatbots y media completado ===")
}

// Función auxiliar para demostrar el flujo completo de un chatbot
func demonstrateChatbotFlow(client wati.WATIClient, ctx context.Context) {
	fmt.Println("\n=== FLUJO COMPLETO DE CHATBOT ===")
	
	// 1. Crear chatbot especializado
	salesBot := &chatbots.CreateChatbotRequest{
		Name:        "Bot de Soporte",
		Description: "Chatbot para soporte técnico y FAQ",
		Keywords:    []string{"ayuda", "soporte", "problema", "error"},
		Responses: []chatbots.Response{
			{
				Trigger:  "ayuda",
				Message:  "¡Hola! Soy el asistente de soporte. ¿En qué puedo ayudarte?",
				IsActive: true,
			},
			{
				Trigger:  "problema",
				Message:  "Entiendo que tienes un problema. ¿Podrías describirlo brevemente?",
				IsActive: true,
			},
		},
		IsActive: true,
	}
	
	bot, err := client.Chatbots().CreateChatbot(ctx, salesBot)
	if err != nil {
		log.Printf("Error creando bot de soporte: %v", err)
		return
	}
	
	fmt.Printf("✅ Bot de soporte creado: %s\n", bot.ID)
	
	// 2. Activar para un usuario
	_, err = client.Chatbots().StartChatbotForContact(ctx, bot.ID, "1234567890")
	if err != nil {
		log.Printf("Error activando bot: %v", err)
		return
	}
	
	fmt.Printf("✅ Bot activado para usuario\n")
	
	// 3. Simular escalamiento a humano
	_, err = client.Chatbots().TransferChatToHuman(
		ctx,
		"1234567890",
		"soporte@empresa.com",
		"Problema complejo que requiere atención humana",
	)
	
	if err != nil {
		log.Printf("Error transfiriendo: %v", err)
		return
	}
	
	fmt.Printf("✅ Chat transferido a soporte humano\n")
}

