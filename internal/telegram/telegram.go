package telegram

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"Telegram-Bot-With-GO/internal/mariadb"
	"Telegram-Bot-With-GO/internal/telegram/querys"
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
			userLanguageCode := update.Message.From.LanguageCode

			role, err := mariadb.GetUserRole(database, userID)
			if err != nil {
				log.Printf("Error en obtenir el rol de l'usuari: %v", err)
				return
			}

			if role == "" {
				notAuthorizedMessage := getUnauthorizedMessage(userLanguageCode)
				msg := tgbotapi.NewMessage(chatID, notAuthorizedMessage)
				_, err = bot.Send(msg)
				if err != nil {
					log.Printf("Error en enviar el missatge de no autoritzat: %v", err)
				}
				return
			}

			var keyboard tgbotapi.InlineKeyboardMarkup
			var messageText string

			switch role {
			case "admin":
				messageText = "Selecciona una opció:"
				keyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Instàncies Actives", "show_active_instances"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Accedir al CRUD", "access_crud"),
					),
				)
			case "worker":
				messageText = "Selecciona una opció:"
				keyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Instàncies Actives", "show_active_instances"),
					),
				)
			}

			msg := tgbotapi.NewMessage(chatID, messageText)
			msg.ReplyMarkup = &keyboard
			_, err = bot.Send(msg)
			if err != nil {
				log.Printf("Error en enviar el missatge amb el teclat: %v", err)
			}
		} else if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			chatID := callback.Message.Chat.ID
			data := callback.Data

			callbackResponse := tgbotapi.NewCallback(callback.ID, "")
			_, err := bot.Request(callbackResponse)
			if err != nil {
				fmt.Printf("Error en respondre al callback: %v\n", err)
			}

			switch {
			case data == "show_active_instances":
				querys.QueryActiveNodes(bot, chatID)
			case data == "access_crud":
				msg := tgbotapi.NewMessage(chatID, "L'accés al CRUD encara no està implementat.")
				bot.Send(msg)
			case strings.HasPrefix(data, "node_"):
				instance := strings.TrimPrefix(data, "node_")
				responseMessage := fmt.Sprintf("Les funcionalitats de mètriques per a la instància '%s' estan en procés de desenvolupament.", instance)
				msg := tgbotapi.NewMessage(chatID, responseMessage)
				bot.Send(msg)
			}
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