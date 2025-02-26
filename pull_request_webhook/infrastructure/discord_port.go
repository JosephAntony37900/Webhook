package infrastructure

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Enviar mensaje a Discord con manejo de errores y rate limit
func post_discord(msg string) int {
	// Leer URL desde .env y verificar que existe
	discord_webhook_url := os.Getenv("DISCORD_WEBHOOK_URL")
	log.Println("📌 URL de Discord Webhook recibida:", discord_webhook_url) // DEBUG

	if discord_webhook_url == "" {
		log.Println("❌ Error: El link hacia Discord no existe o no está en .env")
		return 500
	}

	// Crear payload
	payload := map[string]string{"content": msg}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("❌ Error al serializar JSON:", err)
		return 500
	}

	// Enviar solicitud a Discord
	resp, err := http.Post(discord_webhook_url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("❌ Error al mandar el mensaje a Discord:", err)
		return 500
	}
	defer resp.Body.Close()

	// Manejo de respuestas de Discord
	switch resp.StatusCode {
	case 200, 204:
		log.Println("✅ Mensaje enviado correctamente a Discord")
		return 200
	case 429: // Rate Limit Exceeded
		retryAfter := resp.Header.Get("Retry-After")
		log.Println("🚨 Error 429: Rate Limited. Reintentar después de:", retryAfter, "ms")

		// Convertir `Retry-After` a número y esperar antes de reintentar
		if ms, err := strconv.Atoi(retryAfter); err == nil {
			time.Sleep(time.Duration(ms) * time.Millisecond)
			return post_discord(msg) // Reintentar
		}
		return 429
	default:
		log.Printf("❌ Error al mandar mensaje, código: %d", resp.StatusCode)
		return 400
	}
}
