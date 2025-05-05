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

func GetStorageUsage(bot *tgbotapi.BotAPI, chatID int64, instance string) {
	prometheusURL := os.Getenv("PROMETHEUS_URL")

	sizeQuery := url.QueryEscape(fmt.Sprintf("node_filesystem_size_bytes{instance=\"%s\", job=\"dinf-node-exporter\", mountpoint=\"/\"}", instance))
	availQuery := url.QueryEscape(fmt.Sprintf("node_filesystem_avail_bytes{instance=\"%s\", job=\"dinf-node-exporter\", mountpoint=\"/\"}", instance))

	sizeAPIURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, sizeQuery)
	availAPIURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, availQuery)

	var totalBytes float64
	var availableBytes float64

	respSize, err := http.Get(sizeAPIURL)
	if err != nil {
		log.Printf("Error al obtenir la mida total del disc / de %s: %v", instance, err)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error en obtenir la mida total del disc / de %s: %v", instance, err)))
		return
	}
	defer respSize.Body.Close()

	var prometheusResponseSize PrometheusResponse
	json.NewDecoder(respSize.Body).Decode(&prometheusResponseSize)
	valueStrSize, _ := prometheusResponseSize.Data.Result[0].Value[1].(string)
	totalBytes, _ = strconv.ParseFloat(valueStrSize, 64)

	respAvail, err := http.Get(availAPIURL)
	if err != nil {
		log.Printf("Error al obtenir l'espai lliure del disc de %s: %v", instance, err)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error en obtenir l'espai lliure del disc de %s: %v", instance, err)))
		return
	}
	defer respAvail.Body.Close()

	var prometheusResponseAvail PrometheusResponse
	json.NewDecoder(respAvail.Body).Decode(&prometheusResponseAvail)
	valueStrAvail, _ := prometheusResponseAvail.Data.Result[0].Value[1].(string)
	availableBytes, _ = strconv.ParseFloat(valueStrAvail, 64)

	usedBytes := totalBytes - availableBytes
	totalGB := totalBytes / (1024 * 1024 * 1024)
	usedGB := usedBytes / (1024 * 1024 * 1024)
	availableGB := availableBytes / (1024 * 1024 * 1024)
	usagePct := (usedBytes / totalBytes) * 100

	message := fmt.Sprintf("La instància %s té %.2f GB\nTe usats %.2f GB (%.2f%%)\nI te lliures %.2f GB", instance, totalGB, usedGB, usagePct, availableGB)
	bot.Send(tgbotapi.NewMessage(chatID, message))
}
