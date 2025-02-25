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

func SendDiscordNotification(content string) {
	godotenv.Load()
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	if webhookURL == "" {
		log.Println("No se encontró la URL del webhook de Discord en .env")
		return
	}

	payload := DiscordPayload{Content: content}
	jsonData, _ := json.Marshal(payload)

	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error al enviar mensaje a Discord:", err)
	}
}

func HandlePushEvent(payload []byte) int {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	if webhookURL == "" {
		log.Println("Error: Webhook de Discord no configurado en .env")
		return 500
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Println("Error al parsear el payload:", err)
		return 500
	}

	commits := data["commits"].([]interface{})
	message := "Cambios recientes:\n"

	for _, commit := range commits {
		commitData := commit.(map[string]interface{})
		author := commitData["author"].(map[string]interface{})["name"]
		commitMessage := commitData["message"]
		message += "- " + author.(string) + ": " + commitMessage.(string) + "\n"
	}

	return postDiscordMessage(webhookURL, message)
}

func postDiscordMessage(url, message string) int {
	body, _ := json.Marshal(DiscordPayload{Content: message})
	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println("Error al enviar mensaje a Discord:", err)
		return 500
	}

	defer response.Body.Close()
	log.Println("Mensaje enviado a Discord con estado:", response.StatusCode)
	return response.StatusCode
}