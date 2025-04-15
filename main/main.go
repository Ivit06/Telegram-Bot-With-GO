package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	godotenv.Load()
	botToken := os.Getenv("TELEGRAM_APITOKEN")
	bot, _ := tgbotapi.NewBotAPI(botToken)

	log.Printf("Bot started as %s", bot.Self.UserName)

	webhookURL := os.Getenv("NGROK_URL")
	wh, _ := tgbotapi.NewWebhook(webhookURL)
	bot.Request(wh)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		update, _ := bot.HandleUpdate(r)
		if update.Message == nil {
			return
		}

		chatID := update.Message.Chat.ID
		receivedText := update.Message.Text
		firstName := update.Message.From.FirstName
		userLanguageCode := update.Message.From.LanguageCode

		var response string

		switch userLanguageCode {
		case "es":
			response = fmt.Sprintf("¡Hola, %s! Has dicho: %s", firstName, receivedText)
		case "ca":
			response = fmt.Sprintf("Hola, %s! Has dit: %s", firstName, receivedText)
		default:
			response = fmt.Sprintf("Hello, %s! You said: %s", firstName, receivedText)
		}

		msg := tgbotapi.NewMessage(chatID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	})

	port := os.Getenv("PORT")
	serverAddress := fmt.Sprintf(":%s", port)

	log.Printf("Bot listening on %s with webhook: %s", serverAddress, webhookURL)
	http.ListenAndServe(serverAddress, nil)
}