package querys

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetRAMUsagePercentage(bot *tgbotapi.BotAPI, chatID int64, instance string) {
	prometheusURL := os.Getenv("PROMETHEUS_URL")

	query := url.QueryEscape(fmt.Sprintf("(1 - (node_memory_MemAvailable_bytes{instance=\"%s\", job=\"dinf-node-exporter\"} / node_memory_MemTotal_bytes{instance=\"%s\", job=\"dinf-node-exporter\"})) * 100", instance, instance))
	apiURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, query)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error al consultar Prometheus per a la RAM de %s: %v", instance, err)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error en obtenir l'ús de la RAM de %s: %v", instance, err)))
		return
	}
	defer resp.Body.Close()

	var prometheusResponse PrometheusResponse

	json.NewDecoder(resp.Body).Decode(&prometheusResponse)

	if len(prometheusResponse.Data.Result) > 0 {
		result := prometheusResponse.Data.Result[0]
		if len(result.Value) > 1 {
			if value, ok := result.Value[1].(string); ok {
				ramUsage, _ := strconv.ParseFloat(value, 64)
				message := fmt.Sprintf("Ús de la RAM per a %s: %.2f%%", instance, ramUsage)
				bot.Send(tgbotapi.NewMessage(chatID, message))
			}
		}
	}
}