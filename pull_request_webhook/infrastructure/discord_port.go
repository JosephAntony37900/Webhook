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

// post_discord envía un mensaje a Discord y maneja rate limits (429)
func post_discord(msg string) int {
	discord_webhook_url := os.Getenv("DISCORD_WEBHOOK_URL")
	log.Println("📌 URL de Discord Webhook recibida:", discord_webhook_url)

	if discord_webhook_url == "" {
		log.Println("❌ Error: El link hacia Discord no existe o no está en .env")
		return 500
	}

	payload := map[string]string{"content": msg}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("❌ Error al serializar JSON:", err)
		return 500
	}

	maxRetries := 3 // Intentar máximo 3 veces
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Post(discord_webhook_url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Println("❌ Error al mandar el mensaje a Discord:", err)
			return 500
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case 200, 204:
			log.Println("✅ Mensaje enviado correctamente a Discord")
			return 200
		case 429: // Rate Limited
			retryAfter := resp.Header.Get("Retry-After")
			ms, err := strconv.Atoi(retryAfter)
			if err != nil {
				log.Println("🚨 Error 429: Rate Limited. Reintentar después de 45s por defecto.")
				ms = 45000 // Espera 45s por defecto si `Retry-After` no se puede leer
			}
			log.Printf("🚨 Error 429: Rate Limited. Esperando %d ms antes de reintentar... (%d/%d)", ms, i+1, maxRetries)
			time.Sleep(time.Duration(ms) * time.Millisecond)
		default:
			log.Printf("❌ Error al mandar mensaje, código: %d", resp.StatusCode)
			return 400
		}
	}

	log.Println("❌ Error: Se alcanzó el límite de reintentos.")
	return 429
}
