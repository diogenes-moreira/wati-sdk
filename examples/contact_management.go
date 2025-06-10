package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tu-usuario/go-wati"
	"github.com/tu-usuario/go-wati/contacts"
)

func main() {
	// Configurar cliente WATI
	client := wati.NewClient(
		"https://live-server-12345.wati.io",
		"tu-token-aqui",
		wati.WithTimeout(30),
	)

	ctx := context.Background()

	// Ejemplo 1: Crear un nuevo contacto
	fmt.Println("=== Creando nuevo contacto ===")
	
	newContact := &contacts.CreateContactRequest{
		WhatsappNumber: "1234567890",
		FirstName:      "Juan",
		LastName:       "Pérez",
		Email:          "juan.perez@email.com",
		Tags:           []string{"cliente", "premium"},
		CustomParams: []contacts.CustomParam{
			{Name: "empresa", Value: "Tech Corp"},
			{Name: "cargo", Value: "Gerente"},
		},
	}
	
	createdContact, err := client.Contacts().AddContact(ctx, newContact)
	if err != nil {
		log.Printf("Error creando contacto: %v", err)
	} else {
		fmt.Printf("Contacto creado: %s (%s)\n", createdContact.FullName, createdContact.Phone)
	}

	// Ejemplo 2: Buscar contacto por teléfono
	fmt.Println("\n=== Buscando contacto por teléfono ===")
	
	foundContact, err := client.Contacts().GetContactByPhone(ctx, "1234567890")
	if err != nil {
		log.Printf("Error buscando contacto: %v", err)
	} else {
		fmt.Printf("Contacto encontrado: %s\n", foundContact.FullName)
		fmt.Printf("Email: %s\n", foundContact.Email)
		fmt.Printf("Tags: %v\n", foundContact.Tags)
	}

	// Ejemplo 3: Obtener lista de contactos con paginación
	fmt.Println("\n=== Obteniendo lista de contactos ===")
	
	contactsList, err := client.Contacts().GetContacts(ctx, &contacts.GetContactsParams{
		PageSize:   20,
		PageNumber: 1,
	})
	
	if err != nil {
		log.Printf("Error obteniendo contactos: %v", err)
	} else {
		fmt.Printf("Contactos encontrados: %d de %d total\n", 
			len(contactsList.Contacts), contactsList.TotalCount)
		
		for i, contact := range contactsList.Contacts {
			if i >= 5 { // Mostrar solo los primeros 5
				break
			}
			fmt.Printf("- %s (%s)\n", contact.FullName, contact.Phone)
		}
	}

	// Ejemplo 4: Buscar contactos por nombre
	fmt.Println("\n=== Buscando contactos por nombre ===")
	
	searchResults, err := client.Contacts().SearchContacts(ctx, "Juan")
	if err != nil {
		log.Printf("Error en búsqueda: %v", err)
	} else {
		fmt.Printf("Resultados de búsqueda: %d\n", len(searchResults.Contacts))
		for _, contact := range searchResults.Contacts {
			fmt.Printf("- %s (%s)\n", contact.FullName, contact.Phone)
		}
	}

	// Ejemplo 5: Filtrar contactos por fecha
	fmt.Println("\n=== Filtrando contactos por fecha ===")
	
	filter := &contacts.ContactFilter{
		CreatedAfter: time.Now().AddDate(0, -1, 0), // Últimos 30 días
	}
	
	filteredContacts, err := client.Contacts().FilterContacts(ctx, filter)
	if err != nil {
		log.Printf("Error filtrando contactos: %v", err)
	} else {
		fmt.Printf("Contactos recientes: %d\n", len(filteredContacts.Contacts))
	}

	// Ejemplo 6: Actualizar contacto
	fmt.Println("\n=== Actualizando contacto ===")
	
	if foundContact != nil {
		updateData := &contacts.UpdateContactRequest{
			Email: "juan.perez.nuevo@email.com",
			Tags:  []string{"cliente", "premium", "actualizado"},
			CustomParams: []contacts.CustomParam{
				{Name: "empresa", Value: "New Tech Corp"},
				{Name: "cargo", Value: "Director"},
				{Name: "fecha_actualizacion", Value: time.Now().Format("2006-01-02")},
			},
		}
		
		updatedContact, err := client.Contacts().UpdateContact(ctx, foundContact.ID, updateData)
		if err != nil {
			log.Printf("Error actualizando contacto: %v", err)
		} else {
			fmt.Printf("Contacto actualizado: %s\n", updatedContact.Email)
		}
	}

	// Ejemplo 7: Agregar múltiples contactos
	fmt.Println("\n=== Agregando múltiples contactos ===")
	
	bulkContacts := []*contacts.CreateContactRequest{
		{
			WhatsappNumber: "1111111111",
			FirstName:      "María",
			LastName:       "García",
			Email:          "maria.garcia@email.com",
			Tags:           []string{"cliente", "nuevo"},
		},
		{
			WhatsappNumber: "2222222222",
			FirstName:      "Carlos",
			LastName:       "López",
			Email:          "carlos.lopez@email.com",
			Tags:           []string{"prospecto"},
		},
		{
			WhatsappNumber: "3333333333",
			FirstName:      "Ana",
			LastName:       "Martínez",
			Email:          "ana.martinez@email.com",
			Tags:           []string{"cliente", "premium"},
		},
	}
	
	bulkResult, err := client.Contacts().AddContacts(ctx, bulkContacts)
	if err != nil {
		log.Printf("Error agregando contactos en lote: %v", err)
	} else {
		fmt.Printf("Contactos agregados exitosamente: %d\n", bulkResult.SuccessCount)
		if bulkResult.FailureCount > 0 {
			fmt.Printf("Fallos: %d\n", bulkResult.FailureCount)
		}
	}

	// Ejemplo 8: Obtener todos los contactos (con paginación automática)
	fmt.Println("\n=== Obteniendo todos los contactos ===")
	
	allContacts, err := client.Contacts().GetAllContacts(ctx)
	if err != nil {
		log.Printf("Error obteniendo todos los contactos: %v", err)
	} else {
		fmt.Printf("Total de contactos obtenidos: %d\n", len(allContacts))
		
		// Estadísticas por tags
		tagCount := make(map[string]int)
		for _, contact := range allContacts {
			for _, tag := range contact.Tags {
				if tagStr, ok := tag.(string); ok {
					tagCount[tagStr]++
				}
			}
		}
		
		fmt.Println("Distribución por tags:")
		for tag, count := range tagCount {
			fmt.Printf("- %s: %d contactos\n", tag, count)
		}
	}

	// Ejemplo 9: Actualizar solo las etiquetas de un contacto
	fmt.Println("\n=== Actualizando etiquetas ===")
	
	if foundContact != nil {
		newTags := []string{"cliente", "vip", "actualizado_" + time.Now().Format("2006-01-02")}
		
		updatedContact, err := client.Contacts().UpdateContactTags(ctx, foundContact.ID, newTags)
		if err != nil {
			log.Printf("Error actualizando etiquetas: %v", err)
		} else {
			fmt.Printf("Etiquetas actualizadas para: %s\n", updatedContact.FullName)
		}
	}

	fmt.Println("\n=== Ejemplo de gestión de contactos completado ===")
}

