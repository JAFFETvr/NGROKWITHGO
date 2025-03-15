package handlers

import (
	"github_wb/application"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GitHubEventHandler(ctx *gin.Context) {
	eventType := ctx.GetHeader("X-GitHub-Event")
	log.Printf("Evento recibido: %s", eventType)

	payload, err := ctx.GetRawData()
	if err != nil {
		log.Printf("Error leyendo el cuerpo de la solicitud: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el cuerpo"})
		return
	}

	var status int

	switch eventType {
	case "pull_request":
		status = application.HandlePullRequestEvent(payload)
	case "push":
		status = application.HandlePushEvent(payload)
	default:
		log.Printf("Evento no manejado: %s", eventType)
		ctx.JSON(http.StatusNotImplemented, gin.H{"status": "Evento no soportado"})
		return
	}

	ctx.JSON(status, gin.H{"status": "Evento procesado correctamente"})
}
