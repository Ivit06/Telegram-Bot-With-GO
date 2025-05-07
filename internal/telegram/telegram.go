package telegram

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
	"regexp"

	"Telegram-Bot-With-GO/internal/mariadb"
	"Telegram-Bot-With-GO/internal/telegram/crud"
	"Telegram-Bot-With-GO/internal/telegram/querys"
	"Telegram-Bot-With-GO/internal/telegram/discover"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var waitingForUserID = make(map[int64]bool)
var deleteAttempts = make(map[int64]int)

var createUserStep = make(map[int64]string)
var createUserState = make(map[int64]map[string]string)

var modifyUserState = make(map[int64]map[string]string)
var modifyUserStep = make(map[int64]string)
var waitingForUserToModifyID = make(map[int64]bool)

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

func HandleWebhook(bot *tgbotapi.BotAPI, database *sql.DB, crudDB *sql.DB) http.HandlerFunc {
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

			if update.Message.Command() == "start" {
				role, err := mariadb.GetUserRole(database, userID)
				if err != nil {
					log.Printf("Error en obtenir el rol de l'usuari per a /start: %v", err)
					return
				}
				if role != "" {
					messageText := "Selecciona una opció:"
					var keyboard tgbotapi.InlineKeyboardMarkup
					switch role {
					case "admin":
						keyboard = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Instàncies Actives", "show_active_instances"),
							),
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Accedir al CRUD", "access_crud"),
							),
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Autodescobriment", "access_discover"),
							),
						)
					case "worker":
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
						log.Printf("Error en enviar el missatge amb el teclat per a /start: %v", err)
					}
				} else {
					notAuthorizedMessage := getUnauthorizedMessage(userLanguageCode)
					msg := tgbotapi.NewMessage(chatID, notAuthorizedMessage)
					_, err = bot.Send(msg)
					if err != nil {
						log.Printf("Error en enviar el missatge de no autoritzat per a /start: %v", err)
					}
				}
				return
			}

			if update.Message.Command() == "help" {
				userID := update.Message.From.ID
				chatID := update.Message.Chat.ID

				role, err := mariadb.GetUserRole(database, userID)
				if err != nil {
					log.Printf("Error en obtenir el rol de l'usuari per a /help: %v", err)
					msg := tgbotapi.NewMessage(chatID, "Error en obtenir la informació d'ajuda.")
					_, err = bot.Send(msg)
					if err != nil {
						log.Printf("Error en enviar el missatge d'error d'ajuda: %v", err)
					}
					return
				}

				var helpText string
				messagesDir := "messeges/"

				switch role {
				case "admin":
					adminHelpBytes, err := os.ReadFile(messagesDir + "admin_help.json")
					if err != nil {
						log.Printf("Error en llegir el fitxer d'ajuda d'admin: %v", err)
						helpText = "Les ajudes per a administradors no estan disponibles en aquest moment."
					} else {
						var data map[string]string
						err = json.Unmarshal(adminHelpBytes, &data)
						if err != nil {
							log.Printf("Error en parsejar el fitxer JSON d'ajuda d'admin: %v", err)
							helpText = "Error en carregar l'ajuda per a administradors."
						} else if msg, ok := data["ca"]; ok {
							helpText = msg
						} else {
							helpText = "Ajuda per a administradors no disponible en català."
						}
					}
				case "worker":
					workerHelpBytes, err := os.ReadFile(messagesDir + "worker_help.json")
					if err != nil {
						log.Printf("Error en llegir el fitxer d'ajuda de worker: %v", err)
						helpText = "Les ajudes per a treballadors no estan disponibles en aquest moment."
					} else {
						var data map[string]string
						err = json.Unmarshal(workerHelpBytes, &data)
						if err != nil {
							log.Printf("Error en parsejar el fitxer JSON d'ajuda de worker: %v", err)
							helpText = "Error en carregar l'ajuda per a treballadors."
						} else if msg, ok := data["ca"]; ok {
							helpText = msg
						} else {
							helpText = "Ajuda per a treballadors no disponible en català."
						}
					}
				}

				msg := tgbotapi.NewMessage(chatID, helpText)
				_, err = bot.Send(msg)
				if err != nil {
					log.Printf("Error en enviar el missatge d'ajuda: %v", err)
				}
				return
			}

			if waitingForUserID[chatID] {
				idToDelete, err := strconv.Atoi(update.Message.Text)
				if err != nil {
						msg := tgbotapi.NewMessage(chatID, "El ID introduït no és vàlid. Si us plau, introduïu un ID vàlid (1234)")
						bot.Send(msg)
						delete(waitingForUserID, chatID)
						delete(deleteAttempts, chatID)
						return
				}

				rowsAffected, err := crud.EliminarUsuari(crudDB, int64(idToDelete))
				attempts := deleteAttempts[chatID]
				deleteAttempts[chatID] = attempts + 1
				if err != nil {
						log.Printf("Error al eliminar l'usuari: %v", err)
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error al eliminar l'usuari amb ID %d.", idToDelete))
						bot.Send(msg)

				} else if rowsAffected > 0 {
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Usuari amb ID %d eliminat correctament.", idToDelete))
						bot.Send(msg)
				} else {
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("No s'ha trobat cap usuari amb ID %d per eliminar. Et queden %d intents.", idToDelete, 2-attempts))
						bot.Send(msg)
						if attempts >= 2 {
								msg := tgbotapi.NewMessage(chatID, "Torna a intentar-ho.")
								bot.Send(msg)
								delete(waitingForUserID, chatID)
								delete(deleteAttempts, chatID)
								return
						}
				}
				return
			}

			if step, ok := createUserStep[chatID]; ok {
				switch step {
				case "ask_id":
						idStr := update.Message.Text
						id, err := strconv.ParseInt(idStr, 10, 64)
						if err != nil {
								msg := tgbotapi.NewMessage(chatID, "La ID introduïda no és vàlida. Si us plau, introdueix una ID vàlida:")
								bot.Send(msg)
								return
						}
						createUserState[chatID]["id"] = strconv.FormatInt(id, 10)
						createUserStep[chatID] = "ask_nombre"
						msg := tgbotapi.NewMessage(chatID, "Si us plau, introdueix el nom de l'usuari:")
						bot.Send(msg)

				case "ask_nombre":
						nombre := update.Message.Text
						matched, err := regexp.MatchString(`^[a-zA-Z\s]+$`, nombre)
						if err != nil {
								log.Printf("Error al validar el nombre: %v", err)
								msg := tgbotapi.NewMessage(chatID, "Error interno al validar el nom.")
								bot.Send(msg)
								return
						}
						if !matched {
								msg := tgbotapi.NewMessage(chatID, "El nom introduït no és vàlid (només lletres i espais). Si us plau, introdueix el nom de nou:")
								bot.Send(msg)
								return
						}
						createUserState[chatID]["nombre"] = update.Message.Text
						createUserStep[chatID] = "ask_apellido"
						msg := tgbotapi.NewMessage(chatID, "Si us plau, introdueix el primer cognom de l'usuari (opcional, /skip per ometre):")
						bot.Send(msg)

				case "ask_apellido":
						apellido := update.Message.Text

						if strings.ToLower(apellido) == "/skip" {
								createUserState[chatID]["apellido"] = "NULL"
								createUserStep[chatID] = "ask_segundo_apellido"
								msg := tgbotapi.NewMessage(chatID, "Si us plau, introdueix el segon cognom de l'usuari (opcional, /skip per ometre):")
								bot.Send(msg)
								return
						}

						matched, err := regexp.MatchString(`^[a-zA-Z]+$`, apellido)
						if err != nil {
								log.Printf("Error al validar el cognom: %v", err)
								msg := tgbotapi.NewMessage(chatID, "Error interno al validar el cognom.")
								bot.Send(msg)
								return
						}
						if !matched {
								msg := tgbotapi.NewMessage(chatID, "El cognom introduït no és vàlid (només lletres i sense espais). Si us plau, introdueix el primer cognom de nou (/skip per ometre):")
								bot.Send(msg)
								return
						}
						createUserState[chatID]["apellido"] = update.Message.Text
						createUserStep[chatID] = "ask_segundo_apellido"
						msg := tgbotapi.NewMessage(chatID, "Si us plau, introdueix el segon cognom de l'usuari (opcional, /skip per ometre):")
						bot.Send(msg)

				case "ask_segundo_apellido":
						segundoApellido := update.Message.Text

						if strings.ToLower(segundoApellido) == "/skip" {
								createUserState[chatID]["segundo_apellido"] = "NULL"
								markup := tgbotapi.NewRemoveKeyboard(true)
								msg := tgbotapi.NewMessage(chatID, "Si us plau, introdueix el rol de l'usuari (admin o worker):")
								msg.ReplyMarkup = &markup
								createUserStep[chatID] = "ask_rol"
								bot.Send(msg)
								return
						}

						if segundoApellido != "" {
								matched, err := regexp.MatchString(`^[a-zA-Z]+$`, segundoApellido)
								if err != nil {
										log.Printf("Error al validar el segon cognom: %v", err)
										msg := tgbotapi.NewMessage(chatID, "Error interno al validar el segon cognom.")
										bot.Send(msg)
										return
								}
								if !matched {
										msg := tgbotapi.NewMessage(chatID, "El segon cognom introduït no és vàlid (només lletres i sense espais). Si us plau, introdueix el segon cognom de nou (/skip per ometre):")
										bot.Send(msg)
										return
								}
						}
						createUserState[chatID]["segundo_apellido"] = update.Message.Text
						markup := tgbotapi.NewRemoveKeyboard(true)
						msg := tgbotapi.NewMessage(chatID, "Si us plau, introdueix el rol de l'usuari (admin o worker):")
						msg.ReplyMarkup = &markup
						createUserStep[chatID] = "ask_rol"
						bot.Send(msg)

				case "ask_rol":
						rol := strings.ToLower(update.Message.Text)
						if rol == "admin" || rol == "worker" {
								createUserState[chatID]["rol"] = rol
								id, _ := strconv.ParseInt(createUserState[chatID]["id"], 10, 64)
								nombre := createUserState[chatID]["nombre"]
								apellido := createUserState[chatID]["apellido"]
								segundoApellido := createUserState[chatID]["segundo_apellido"]

								err := crud.CrearUsuari(crudDB, id, rol, nombre, apellido, segundoApellido)
								if err != nil {
										msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error al crear l'usuari amb ID %d, nom %s %s %s amb rol %s: %v", id, nombre, apellido, segundoApellido, rol, err))
										bot.Send(msg)
								} else {
										msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Usuari amb ID %d, nom %s %s %s amb rol %s creat correctament.", id, nombre, apellido, segundoApellido, rol))
										bot.Send(msg)
								}

								delete(createUserState, chatID)
								delete(createUserStep, chatID)

						} else {
								msg := tgbotapi.NewMessage(chatID, "El rol ha de ser 'admin' o 'worker'. Si us plau, introdueix el rol de nou:")
								bot.Send(msg)
						}
				}
				return
			}
		}

		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			chatID := callback.Message.Chat.ID
			userIDCallback := callback.From.ID

			exists, err := crud.CheckUserExists(crudDB, userIDCallback)
			if err != nil {
				log.Printf("Error al verificar si existe el usuario: %v", err)
				msg := tgbotapi.NewMessage(chatID, "Error al verificar l'usuari. Si us plau, intenta-ho de nou.")
				bot.Send(msg)
				return
			}
			if !exists {
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("No s'ha trobat cap usuari amb la ID %d. Torna a intentar-ho", userIDCallback))
				bot.Send(msg)
				return
			}

			data := callback.Data

			callbackResponse := tgbotapi.NewCallback(callback.ID, "")
			_, err = bot.Request(callbackResponse)
			if err != nil {
				fmt.Printf("Error en respondre al callback: %v\n", err)
			}

			switch {
			case data == "show_active_instances":
				querys.QueryActiveNodes(bot, chatID)
			case data == "access_crud":
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("CREAR", "crud_crear"),
						tgbotapi.NewInlineKeyboardButtonData("LLISTAR", "crud_llistar"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("ACTUALITZAR", "crud_actualitzar"),
						tgbotapi.NewInlineKeyboardButtonData("ELIMINAR", "crud_eliminar"),
					),
				)
				msg := tgbotapi.NewMessage(chatID, "Selecciona una acció del CRUD:")
				msg.ReplyMarkup = &keyboard
				bot.Send(msg)
			case data == "crud_llistar":
				crud.LlistarElements(bot, chatID, crudDB)
			case data == "crud_eliminar":
				waitingForUserID[chatID] = true
				msg := tgbotapi.NewMessage(chatID, "Si us plau, introduïu l'ID de l'usuari que voleu eliminar.")
				bot.Send(msg)
			case data == "crud_crear":
				createUserState[chatID] = make(map[string]string)
				createUserStep[chatID] = "ask_id"
				msg := tgbotapi.NewMessage(chatID, "Si us plau, introdueix la ID de l'usuari:")
				bot.Send(msg)
			case data == "crud_actualitzar":
				msg := tgbotapi.NewMessage(chatID, "Aquesta funcionalitat del CRUD encara no està implementada.")
				bot.Send(msg)
			case strings.HasPrefix(data, "node_"):
				instance := strings.TrimPrefix(data, "node_")
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("CPU", fmt.Sprintf("get_cpu_info_%s", instance)),
						tgbotapi.NewInlineKeyboardButtonData("RAM", fmt.Sprintf("get_ram_info_%s", instance)),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("STORAGE", fmt.Sprintf("get_storage_info_%s", instance)),
						tgbotapi.NewInlineKeyboardButtonData("NETWORK", fmt.Sprintf("get_network_info_%s", instance)),
					),
				)
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Selecciona la mètrica per a %s:", instance))
				msg.ReplyMarkup = &keyboard
				bot.Send(msg)
			case strings.HasPrefix(data, "get_cpu_info_"):
				instance := strings.TrimPrefix(data, "get_cpu_info_")
				querys.GetCPUUsagePercentage(bot, chatID, instance)
			case strings.HasPrefix(data, "get_ram_info_"):
				instance := strings.TrimPrefix(data, "get_ram_info_")
				querys.GetRAMUsagePercentage(bot, chatID, instance)
			case strings.HasPrefix(data, "get_storage_info_"):
				instance := strings.TrimPrefix(data, "get_storage_info_")
				querys.GetStorageUsage(bot, chatID, instance)
			case strings.HasPrefix(data, "get_network_info_"):
				instance := strings.TrimPrefix(data, "get_network_info_")
				querys.GetActivePorts(bot, chatID, instance)
			case data == "access_discover":
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Node Exporter", "discover_node_exporter"),
						tgbotapi.NewInlineKeyboardButtonData("Port Exporter", "discover_port_exporter"),
					),
				)
				msg := tgbotapi.NewMessage(chatID, "De què vols fer l'autodescobriment:")
				msg.ReplyMarkup = &keyboard
				bot.Send(msg)
			case data == "discover_node_exporter":
				msg := tgbotapi.NewMessage(chatID, "Descobriment per a Node Exporter inicialitzat.")
				bot.Send(msg)

				output, err := discover.ExecuteDiscoverNodeExporter()
				if err != nil {
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error en executar l'script: %v", err))
					bot.Send(msg)
					return
				}

				msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("%s", output))
				bot.Send(msg)
			case data == "discover_port_exporter":
				msg := tgbotapi.NewMessage(chatID, "Descobriment per a Port Exporter inicialitzat.")
				bot.Send(msg)

				output, err := discover.ExecuteDiscoverPortExporter()
				if err != nil {
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error en executar l'script: %v", err))
					bot.Send(msg)
					return
				}

				msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("%s", output))
				bot.Send(msg)
			}
		}
	}
}
