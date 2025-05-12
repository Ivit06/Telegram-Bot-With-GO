package mariadb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InitDBCRUD() (*sql.DB, error) {
	godotenv.Load()
	dbCRUDUser := os.Getenv("DB_USER_CRUD")
	dbCRUDPass := os.Getenv("DB_PASS_CRUD")
	dbCRUDName := os.Getenv("DB_NAME_CRUD")
	dbCRUDHost := os.Getenv("DB_HOST_CRUD")
	dbCRUDPort := os.Getenv("DB_PORT_CRUD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbCRUDUser, dbCRUDPass, dbCRUDHost, dbCRUDPort, dbCRUDName)

	var err error
	masterDatabase, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("no s'ha pogut connectar a la base de dades del CRUD: %w", err)
	}

	err = masterDatabase.Ping()
	if err != nil {
		return nil, fmt.Errorf("no s'ha pogut fer ping a la base de dades del CRUD: %w", err)
	}
	log.Println("Connexió a la base de dades del CRUD establerta correctament")
	return masterDatabase, nil
}
