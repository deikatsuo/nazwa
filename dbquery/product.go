package dbquery

import (
	"log"
	"nazwa/wrapper"
	"strings"

	"github.com/jmoiron/sqlx"
)

// GetAllProduct menampilkan semua product
func GetAllProduct(db *sqlx.DB) ([]wrapper.Product, error) {
	var product []wrapper.NullableProduct
	var parse []wrapper.Product

	query := `SELECT
		id,
		name,
		TO_CHAR(base_price,'Rp999G999G999G999G999') AS base_price,
		TO_CHAR(price,'Rp999G999G999G999G999') AS price,
		code,
		TO_CHAR(created_at, 'MM/DD/YYYY HH12:MI:SS AM') AS created_at
		FROM "product"`

	err := db.Select(&product, query)
	if err != nil {
		log.Println("product.go Select all product")
		log.Println(err)
		return []wrapper.Product{}, err
	}

	for _, p := range product {
		parse = append(parse, wrapper.Product{
			ID:        p.ID,
			Name:      strings.Title(p.Name),
			CreatedAt: p.CreatedAt,
			BasePrice: string(p.BasePrice),
			Price:     string(p.Price),
			Code:      string(p.Code.String),
		})
	}

	return parse, nil
}
