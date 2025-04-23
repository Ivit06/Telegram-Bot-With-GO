package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Telegram-Bot-With-GO/internal/mariadb"
	"Telegram-Bot-With-GO/internal/telegram"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	godotenv.Load()
	
	bot, err := telegram.InitBot()
	if err != nil {
		log.Fatalf("Error en inicialitzar el bot de Telegram: %v", err)
	}

	err = telegram.SetWebhook(bot)
	if err != nil {
		log.Fatalf("Error en configurar el webhook: %v", err)
	}

	database, err := mariadb.InitDB()
	if err != nil {
		log.Fatalf("Error en inicialitzar la base de dades: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error en tancar la base de dades: %v", err)
		}
	}()

	crudDatabase, err := mariadb.InitDBCRUD()
	if err != nil {
		log.Fatalf("Error en inicialitzar la base de dades del CRUD: %v", err)
	}
	defer func() {
		if err := crudDatabase.Close(); err != nil {
			log.Printf("Error en tancar la base de dades del CRUD: %v", err)
		}
	}()

	http.HandleFunc("/", telegram.HandleWebhook(bot, database, crudDatabase))

	port := os.Getenv("PORT")
	serverAddress := fmt.Sprintf(":%s", port)
	webhookURL := os.Getenv("NGROK_URL")
	log.Printf("Bot escoltant a %s amb webhook: %s", serverAddress, webhookURL)
	err = http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Fatalf("Error en iniciar el servidor: %v", err)
	}
}