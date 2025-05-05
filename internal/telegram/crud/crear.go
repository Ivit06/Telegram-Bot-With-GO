package crud

import (
	"database/sql"
	"log"
)

func CrearUsuari(db *sql.DB, id int64, rol, nombre, apellido, segundoApellido string) error {
	stmt, err := db.Prepare("INSERT INTO usuaris (id, role, nom, pcognom, scognom) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error al preparar la consulta de creación de usuario: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, rol, nombre, apellido, segundoApellido)
	if err != nil {
		log.Printf("Error al ejecutar la consulta de creación de usuario: %v", err)
		return err
	}
	return nil
}