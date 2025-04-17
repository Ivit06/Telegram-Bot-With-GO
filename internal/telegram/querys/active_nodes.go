package querys

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func QueryActiveNodes(bot *tgbotapi.BotAPI, chatID int64) {
	prometheusURL := os.Getenv("PROMETHEUS_URL")

	if prometheusURL == "" {
		msg := tgbotapi.NewMessage(chatID, "La URL de Prometheus no està configurada.")
		bot.Send(msg)
		return
	}

	query := url.QueryEscape("up{job=\"fuji\"} == 1")
	apiURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, query)

	resp, err := http.Get(apiURL)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error al consultar Prometheus: %v", err))
		bot.Send(msg)
		return
	}
	defer resp.Body.Close()

	var prometheusResponse struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Value  []interface{}       `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&prometheusResponse); err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error al decodificar la resposta de Prometheus: %v", err))
		bot.Send(msg)
		return
	}

	if prometheusResponse.Status != "success" {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error en la resposta de Prometheus: %s", prometheusResponse.Status))
		bot.Send(msg)
		return
	}

	if len(prometheusResponse.Data.Result) == 0 {
		msg := tgbotapi.NewMessage(chatID, "No s'han trobat nodes actius.")
		bot.Send(msg)
		return
	}

	var keyboardRows [][]tgbotapi.InlineKeyboardButton
	for _, result := range prometheusResponse.Data.Result {
		instance, ok := result.Metric["instance"]
		if ok {
			button := tgbotapi.NewInlineKeyboardButtonData(instance, fmt.Sprintf("node_%s", instance))
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