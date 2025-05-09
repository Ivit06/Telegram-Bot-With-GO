package querys

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetActivePorts(bot *tgbotapi.BotAPI, chatID int64, instance string) {
	prometheusURL := os.Getenv("PROMETHEUS_URL")

	parts := strings.Split(instance, ":")
	instanceIP := parts[0]
	instanceToCheck := fmt.Sprintf("%s:8000", instanceIP)

	query := url.QueryEscape(fmt.Sprintf("instance_open_ports{job=\"dinf-port-exporter\", instance=\"%s\"}", instanceToCheck))
	apiURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, query)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error en consultar Prometheus per als ports actius de %s: %v", instance, err)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error en obtenir els ports actius de %s.", instance)))
		return
	}
	defer resp.Body.Close()

	var prometheusResponse PrometheusResponse

	json.NewDecoder(resp.Body).Decode(&prometheusResponse)

	var openPorts []string
	if prometheusResponse.Status == "success" && prometheusResponse.Data.ResultType == "vector" {
		for _, result := range prometheusResponse.Data.Result {
			if portStr, ok := result.Metric["port"]; ok {
				openPorts = append(openPorts, portStr)
			}
		}
	} else {
		log.Printf("Resposta inesperada de Prometheus per als ports actius de %s: %+v", instance, prometheusResponse)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("No s'ha pogut obtenir la llista de ports actius de %s.", instance)))
		return
	}

	if len(openPorts) > 0 {
		message := fmt.Sprintf("Ports actius a %s:\n", instance)
		for _, port := range openPorts {
			message += fmt.Sprintf("- %s\n", port)
		}
		bot.Send(tgbotapi.NewMessage(chatID, message))
	} else {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("No s'han trobat ports actius a %s.", instance)))
	}
}
