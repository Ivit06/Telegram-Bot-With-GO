package mariadb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

func GetUserRole(db *sql.DB, userID int64) (string, error) {
	var role string
	err := db.QueryRow("SELECT role FROM usuaris WHERE id = ?", userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("error verificant el rol de l'usuari a la base de dades: %w", err)
	}
	return role, nil
}
