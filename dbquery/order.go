package dbquery

import "github.com/jmoiron/sqlx"

// GetOrderTotalRow menghitung jumlah row pada tabel user
func GetOrderTotalRow(db *sqlx.DB) (int, error) {
	var total int
	query := `SELECT COUNT(id) FROM "order"`
	err := db.Get(&total, query)
	if err != nil {
		return 0, err
	}
	return total, nil
}
