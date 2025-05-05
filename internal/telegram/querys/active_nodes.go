package querys

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func QueryActiveNodes(bot *tgbotapi.BotAPI, chatID int64) {
	prometheusURL := os.Getenv("PROMETHEUS_URL")

	if prometheusURL == "" {
		msg := tgbotapi.NewMessage(chatID, "La URL de Prometheus no està configurada.")
		bot.Send(msg)
		return
	}

	query := url.QueryEscape("up{job=\"dinf-node-exporter\"} == 1")
	apiURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, query)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error al consultar els nodes actius: %v", err)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error en obtenir la llista de nodes actius: %v", err)))
		return
	}
	defer resp.Body.Close()

	var prometheusResponse PrometheusResponse

	json.NewDecoder(resp.Body).Decode(&prometheusResponse)

	var keyboardRows [][]tgbotapi.InlineKeyboardButton
	for _, result := range prometheusResponse.Data.Result {
		instanceWithPort, ok := result.Metric["instance"]
		if ok {
			parts := strings.Split(instanceWithPort, ":")
			instanceIP := parts[0]
			buttonText := fmt.Sprintf("Instància %s", instanceIP)
			button := tgbotapi.NewInlineKeyboardButtonData(buttonText, fmt.Sprintf("node_%s", instanceWithPort))
			keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(button))
		}
	}

	if len(keyboardRows) > 0 {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
		msg := tgbotapi.NewMessage(chatID, "Instàncies Actives:")
		msg.ReplyMarkup = &keyboard
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "No s'han trobat instàncies de nodes actius.")
		bot.Send(msg)
	}
}