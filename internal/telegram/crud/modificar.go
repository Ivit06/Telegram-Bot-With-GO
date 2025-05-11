package crud

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func ModificarUsuari(db *sql.DB, id int64, nombre, apellido, segundoApellido, rol string) error {
	query := "UPDATE usuaris SET"
	var args []interface{}
	var setClauses []string

	if nombre != "" {
		setClauses = append(setClauses, "nom = ?")
		args = append(args, nombre)
	}
	if apellido != "" {
		setClauses = append(setClauses, "pcognom = ?")
		args = append(args, apellido)
	}
	if segundoApellido != "" {
		setClauses = append(setClauses, "scognom = ?")
		args = append(args, segundoApellido)
	}
	if rol != "" {
		setClauses = append(setClauses, "role = ?")
		args = append(args, rol)
	}

	if len(setClauses) == 0 {
		return nil
	}

	query += " " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, id)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error al preparar la consulta de modificación: %v", err)
		return fmt.Errorf("error al preparar la consulta de modificació: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Printf("Error al ejecutar la consulta de modificación: %v", err)
		return fmt.Errorf("error al executar la consulta de modificació: %w", err)
	}

	return nil
}

func CheckUserExists(db *sql.DB, userID int64) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuaris WHERE id = ?)", userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error al verificar la existencia del usuario: %w", err)
	}
	return exists, nil
}
