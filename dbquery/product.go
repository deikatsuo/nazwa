package dbquery

import (
	"fmt"
	"log"
	"nazwa/wrapper"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// GetProducts mengambil list produk
type GetProducts struct {
	limit     int
	lastid    int
	direction string
}

// Limit set limit
func (p *GetProducts) Limit(limit int) *GetProducts {
	p.limit = limit
	return p
}

// LastID set lastid
func (p *GetProducts) LastID(lastid int) *GetProducts {
	p.lastid = lastid
	return p
}

// Direction untuk backward/forward
// @direction "back","next"
func (p *GetProducts) Direction(direction string) *GetProducts {
	p.direction = direction
	return p
}

// Show tampilkan data
func (p *GetProducts) Show(db *sqlx.DB) ([]wrapper.Product, error) {
	var product []wrapper.NullableProduct
	var parse []wrapper.Product
	limit := 10
	if p.limit > 0 {
		limit = p.limit
	}

	// Where logic
	where := ""
	// Maju/Mundur
	if p.direction == "next" && p.lastid > 0 {
		where = "WHERE id > " + strconv.Itoa(p.lastid)
	} else if p.direction == "back" {
		where = "WHERE id < " + strconv.Itoa(p.lastid) + " ORDER BY id DESC"
	}

	// query pengambilan data produk
	query := `SELECT
		id,
		name,
		TO_CHAR(base_price,'Rp999G999G999G999G999') AS base_price,
		TO_CHAR(price,'Rp999G999G999G999G999') AS price,
		code,
		TO_CHAR(created_at, 'MM/DD/YYYY HH12:MI:SS AM') AS created_at
		FROM "product"
		%s
		LIMIT $1`

	query = fmt.Sprintf(query, where)

	err := db.Select(&product, query, limit)
	if err != nil {
		log.Println("Error: product.go Select all product")
		log.Println(err)
		return []wrapper.Product{}, err
	}

	// mapping data produk
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

// GetProductTotalRow menghitung jumlah row pada tabel user
func GetProductTotalRow(db *sqlx.DB) (int, error) {
	var total int
	query := `SELECT COUNT(id) FROM "product"`
	err := db.Get(&total, query)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetProductByID mengambil data produk berdasarkan ID produk
func GetProductByID(db *sqlx.DB, pid int) (wrapper.Product, error) {
	var product wrapper.Product
	var p wrapper.NullableProduct
	query := `SELECT
		id,
		name,
		TO_CHAR(base_price,'Rp999G999G999G999G999') AS base_price,
		TO_CHAR(price,'Rp999G999G999G999G999') AS price,
		code,
		TO_CHAR(created_at, 'MM/DD/YYYY HH12:MI:SS AM') AS created_at,
		type,
		brand
		FROM "product"
		WHERE id=$1
		LIMIT 1`

	err := db.Get(&p, query, pid)
	if err != nil {
		log.Println("product.go Select product berdasarkan ID")
		log.Println(err)
		return wrapper.Product{}, err
	}

	var photos []wrapper.ProductPhoto

	if pp, err := GetProductPhoto(db, p.ID); err == nil {
		photos = pp
	}

	product = wrapper.Product{
		ID:        p.ID,
		Name:      strings.Title(p.Name),
		CreatedAt: p.CreatedAt,
		BasePrice: string(p.BasePrice),
		Price:     string(p.Price),
		Code:      string(p.Code.String),
		Type:      p.Type.String,
		Brand:     p.Brand.String,
		Photos:    photos,
	}

	return product, nil
}

// GetProductPhoto mengambil data photo produk
func GetProductPhoto(db *sqlx.DB, pid int) ([]wrapper.ProductPhoto, error) {
	var photos []wrapper.ProductPhoto
	query := `SELECT id, photo
	FROM "product_photo"
	WHERE product_id=$1`
	err := db.Select(&photos, query, pid)
	if err != nil {
		return []wrapper.ProductPhoto{}, err
	}
	return photos, err
}
