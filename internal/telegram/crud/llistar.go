package crud

import (
	"database/sql"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LlistarElements(bot *tgbotapi.BotAPI, chatID int64, db *sql.DB) {
	rows, err := db.Query("SELECT id, role, nom, pcognom, scognom FROM usuaris")
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error al llistar els usuaris: %v", err))
		bot.Send(msg)
		return
	}
	defer rows.Close()

	var message string = "Llista d'usuaris:\n\n"
	var count int
	for rows.Next() {
		count++
		var id int
		var role, nom string
		var pcognom, scognom sql.NullString
		err = rows.Scan(&id, &role, &nom, &pcognom, &scognom)
		if err != nil {
			fmt.Printf("Error al llegir la fila: %v\n", err)
			continue
		}

		primerCognom := "null"
		if pcognom.Valid {
			primerCognom = pcognom.String
		}

		segonCognom := "null"
		if scognom.Valid {
			segonCognom = scognom.String
		}

		message += fmt.Sprintf("ID: <code>%d</code>\n", id)
		message += fmt.Sprintf("Rol: %s\n", role)
		message += fmt.Sprintf("Nom: %s\n", nom)
		message += fmt.Sprintf("Cognom 1: %s\n", primerCognom)
		message += fmt.Sprintf("Cognom 2: %s\n\n", segonCognom)
	}

	if count > 0 {
		message = fmt.Sprintf("%s\nTotal d'usuaris: %d", message, count)
	} else {
		message = "No hi ha cap usuari registrat."
	}

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}