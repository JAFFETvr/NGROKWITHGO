package application

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type DiscordMessage struct {
	Content string `json:"content"`
}

func SendDiscordMessage(message string) {
	godotenv.Load()
	webhookURL := os.Getenv("DISCORD_wH_URL")

	if webhookURL == "" {
		log.Println(" No se encontró la URL del webhook de Discord en .env")
		return
	}

	payload := DiscordMessage{Content: message}
	jsonValue, _ := json.Marshal(payload)

	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("Error al enviar mensaje a Discord:", err)
	}
}

func ProcessPush(payload []byte) int {
	webhookURL := os.Getenv("DISCORD_wH_URL")

	if webhookURL == "" {
		log.Println("Error: Webhook de Discord no configurado en .env")
		return 500
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Println("Error al parsear el payload de push:", err)
		return 500
	}

	commits := data["commits"].([]interface{})
	message := "Nuevo push detectado:\n"

	for _, c := range commits {
		commit := c.(map[string]interface{})
		author := commit["author"].(map[string]interface{})["name"]
		msg := commit["message"]

		message += "- " + author.(string) + ": " + msg.(string) + "\n"
	}

	return sendDiscordMessage(webhookURL, message)
}

func sendDiscordMessage(url, message string) int {
	body, _ := json.Marshal(DiscordMessage{Content: message})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println("Error al enviar mensaje a Discord:", err)
		return 500
	}

	defer resp.Body.Close()
	log.Println("Mensaje enviado a Discord con estado:", resp.StatusCode)
	return resp.StatusCode
}

