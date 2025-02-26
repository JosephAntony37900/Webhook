package infrastructure

import (
	"github/pull_request_webhook/application"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlePullRequestEvent(ctx *gin.Context) {

	eventType := ctx.GetHeader("X-GitHub-Event")
	deliveryD := ctx.GetHeader("X-GitHub-Delivery")

	log.Printf("Nuevo evento: %s con ID: %s", eventType, deliveryD)

	rawData, err := ctx.GetRawData()

	if err != nil {
		log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "leer datos"})
	}

	var statusCode int

	switch eventType {
	case "ping":
		ctx.JSON(http.StatusOK, gin.H{"status": "success"})
	case "pull_request":
		message := application.ProcessPullRequestEvent(rawData)

		if message == "ERROR" {
			statusCode = 500
		} else {
			statusCode = post_discord(message)
		}
	}

	switch statusCode {
	case 200:
		ctx.JSON(http.StatusOK, gin.H{"success": "Pull Request procesado con éxito"})
	case 403:
		log.Printf("Error al deserializar el payload del pull request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el payload del pull request"})
	case 500:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
	default:
		ctx.JSON(http.StatusOK, gin.H{"success": "Normal"})
	}

}
