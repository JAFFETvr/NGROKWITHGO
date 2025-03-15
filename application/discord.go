package application

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type DiscordPayload struct {
	Content string `json:"content"`
}

// Cargar variables de entorno
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env")
	}
}

// Función genérica para enviar mensajes a Discord
func postDiscordMessage(webhookURL, message string) int {
	if webhookURL == "" {
		log.Println("Error: Webhook de Discord no configurado.")
		return 500
	}

	body, _ := json.Marshal(DiscordPayload{Content: message})
	response, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println("Error al enviar mensaje a Discord:", err)
		return 500
	}

	defer response.Body.Close()
	log.Println("Mensaje enviado a Discord con estado:", response.StatusCode)
	return response.StatusCode
}

// Manejo de eventos PUSH
func HandlePushEvent(payload []byte) int {
	loadEnv()
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Println("Error al parsear el payload:", err)
		return 500
	}

	commits, ok := data["commits"].([]interface{})
	if !ok {
		log.Println("Error: No se encontraron commits en el payload")
		return 500
	}

	message := "**Cambios recientes en el repositorio**\n"
	for _, commit := range commits {
		commitData := commit.(map[string]interface{})
		author := commitData["author"].(map[string]interface{})["name"].(string)
		commitMessage := commitData["message"].(string)
		message += "- " + author + ": " + commitMessage + "\n"
	}

	return postDiscordMessage(webhookURL, message)
}

// Manejo de eventos WORKFLOW_RUN
func HandleWorkflowRunEvent(payload []byte) int {
	loadEnv()
	webhookURL := os.Getenv("DISCORD_WEBHOOK_TESTS")

	var data struct {
		Action      string `json:"action"`
		WorkflowRun struct {
			Name       string `json:"name"`
			Status     string `json:"status"`
			Conclusion string `json:"conclusion"`
			URL        string `json:"html_url"`
		} `json:"workflow_run"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		log.Println("Error al parsear el payload de workflow_run:", err)
		return 500
	}

	// Determinar estado del workflow
	var estado string
	switch data.WorkflowRun.Status {
	case "completed":
		if data.WorkflowRun.Conclusion == "success" {
			estado = "✅ Exitoso"
		} else {
			estado = "❌ Falló"
		}
	case "in_progress":
		estado = "⏳ En progreso"
	default:
		estado = "❓ Desconocido"
	}

	// Construir mensaje para Discord
	message := "**GitHub Actions Workflow Ejecutado**\n" +
		"📂 Repositorio: " + data.Repository.FullName + "\n" +
		"🔄 Workflow: " + data.WorkflowRun.Name + "\n" +
		"📌 Estado: " + estado + "\n" +
		"🔗 [Ver detalles](" + data.WorkflowRun.URL + ")"

	return postDiscordMessage(webhookURL, message)
}
