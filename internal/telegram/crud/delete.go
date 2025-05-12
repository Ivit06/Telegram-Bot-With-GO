package crud

import (
	"database/sql"
	"fmt"
	_ "log"
	_ "strconv"
)

func EliminarUsuari(db *sql.DB, usuariID int64) (int64, error) {
	query := "DELETE FROM usuaris WHERE id = ?"
	result, err := db.Exec(query, usuariID)
	if err != nil {
		return 0, fmt.Errorf("error al executar la consulta d'eliminació: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtenir el nombre de files afectades: %w", err)
	}
	return rowsAffected, nil
}
