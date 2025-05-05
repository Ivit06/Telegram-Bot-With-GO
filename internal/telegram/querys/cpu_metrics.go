package querys

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetCPUUsagePercentage(bot *tgbotapi.BotAPI, chatID int64, instance string) {
	prometheusURL := os.Getenv("PROMETHEUS_URL")

	query := url.QueryEscape(fmt.Sprintf("100 - (avg by(instance) (rate(node_cpu_seconds_total{job=\"dinf-node-exporter\", instance=\"%s\", mode=\"idle\"}[2m])) * 100)", instance))
	apiURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, query)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error al consultar Prometheus per a la CPU de %s: %v", instance, err)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error en obtenir l'ús de la CPU de %s: %v", instance, err)))
		return
	}
	defer resp.Body.Close()

	var prometheusResponse PrometheusResponse

	json.NewDecoder(resp.Body).Decode(&prometheusResponse)

	if len(prometheusResponse.Data.Result) > 0 {
		result := prometheusResponse.Data.Result[0]
		if len(result.Value) > 1 {
			if value, ok := result.Value[1].(string); ok {
				cpuUsage, _ := strconv.ParseFloat(value, 64)
				message := fmt.Sprintf("Ús de la CPU per a %s: %.2f%% (mitjana dels últims 2 minuts)", instance, cpuUsage)
				bot.Send(tgbotapi.NewMessage(chatID, message))
			}
		}
	}
}