package handlers

import (
	"github_wb/application"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PullRequestEvent(ctx *gin.Context) {
	eventType := ctx.GetHeader("X-GitHub-Event")

	log.Printf("Webhook recibido: \nEvento=%s", eventType)

	payload, err := ctx.GetRawData()
	if err != nil {
		log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el cuerpo"})
		return
	}

	var statusCode int

	switch eventType {
	case "pull_request":
		statusCode = application.ProcessPullRequest(payload)
	case "push":
		statusCode = application.ProcessPush(payload)  // Manejo de eventos push
	default:
		log.Printf("Evento no manejado: %s", eventType)
		ctx.JSON(http.StatusNotImplemented, gin.H{"status": "Evento no manejado"})
		return
	}

	ctx.JSON(statusCode, gin.H{"status": "Evento procesado"})
}

