package telegram

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"Telegram-Bot-With-GO/internal/mariadb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func InitBot() (*tgbotapi.BotAPI, error) {
	godotenv.Load()
	botToken := os.Getenv("TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("no s'ha pogut inicialitzar el bot: %w", err)
	}
	log.Printf("Bot iniciat com a %s", bot.Self.UserName)
	return bot, nil
}

func SetWebhook(bot *tgbotapi.BotAPI) error {
	webhookURL := os.Getenv("NGROK_URL")
	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return fmt.Errorf("no s'ha pogut crear el webhook: %w", err)
	}
	_, err = bot.Request(wh)
	if err != nil {
		return fmt.Errorf("no s'ha pogut configurar el webhook: %w", err)
	}
	return nil
}

func HandleWebhook(bot *tgbotapi.BotAPI, database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		update, err := bot.HandleUpdate(r)
		if err != nil {
			log.Printf("Error en gestionar l'actualització: %v", err)
			return
		}

		if update.Message != nil {
			chatID := update.Message.Chat.ID
			userID := update.Message.From.ID
			receivedText := update.Message.Text
			firstName := update.Message.From.FirstName
			userLanguageCode := update.Message.From.LanguageCode

			exists, err := mariadb.UserExists(database, userID)
			if err != nil {
				log.Printf("Error en comprovar l'usuari: %v", err)
				return
			}

			if !exists {
				notAuthorizedMessage := getUnauthorizedMessage(userLanguageCode)
				msg := tgbotapi.NewMessage(chatID, notAuthorizedMessage)
				_, err = bot.Send(msg)
				if err != nil {
					log.Printf("Error en enviar el missatge de no autoritzat: %v", err)
				}
				return
			}

			var response string

			switch userLanguageCode {
			case "es":
				response = fmt.Sprintf("¡Hola, %s! Has dicho: %s", firstName, receivedText)
			case "ca":
				response = fmt.Sprintf("¡Hola, %s! Has dit: %s", firstName, receivedText)
			default:
				response = fmt.Sprintf("¡Hello, %s! You say: %s", firstName, receivedText)
			}

			msg := tgbotapi.NewMessage(chatID, response)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func getUnauthorizedMessage(languageCode string) string {
	switch languageCode {
	case "es":
		return "Lo siento, no estás autorizado para utilizar este bot."
	case "ca":
		return "Ho sento, no estàs autoritzat per utilitzar aquest bot."
	default:
		return "Sorry, you are not authorized to use this bot."
	}
}