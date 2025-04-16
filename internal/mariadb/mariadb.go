package mariadb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func InitDB() (*sql.DB, error) {
	godotenv.Load()
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	database, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("no s'ha pogut connectar a la base de dades: %w", err)
	}

	err = database.Ping()
	if err != nil {
		return nil, fmt.Errorf("no s'ha pogut fer ping a la base de dades: %w", err)
	}
	log.Println("Connexió a la base de dades establerta correctament")
	return database, nil
}

func UserExists(db *sql.DB, userID int64) (bool, error) {
	var id int64
	err := db.QueryRow("SELECT id FROM usuaris WHERE id = ?", userID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error en comprovar l'usuari a la base de dades: %w", err)
	}
	return true, nil
}
