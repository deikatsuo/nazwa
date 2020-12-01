package dbquery

import (
	"fmt"
	"log"
	"nazwa/wrapper"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

//////////
// GET ///
//////////

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
	if p.direction == "next" {
		where = "WHERE id > " + strconv.Itoa(p.lastid) + " ORDER BY id ASC"
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
			Code:      p.Code,
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

	var photos []wrapper.ListProductPhoto

	if pp, err := GetProductPhoto(db, p.ID); err == nil {
		photos = pp
	}

	product = wrapper.Product{
		ID:        p.ID,
		Name:      strings.Title(p.Name),
		CreatedAt: p.CreatedAt,
		BasePrice: string(p.BasePrice),
		Price:     string(p.Price),
		Code:      p.Code,
		Type:      strings.Title(p.Type.String),
		Brand:     strings.Title(p.Brand.String),
		Photos:    photos,
	}

	return product, nil
}

// GetProductPhoto mengambil data photo produk
func GetProductPhoto(db *sqlx.DB, pid int) ([]wrapper.ListProductPhoto, error) {
	var photos []wrapper.ListProductPhoto
	query := `SELECT id, photo
	FROM "product_photo"
	WHERE product_id=$1`
	err := db.Select(&photos, query, pid)
	if err != nil {
		return []wrapper.ListProductPhoto{}, err
	}
	return photos, err
}

//////////
// POST //
//////////

// CreateProduct membuat produk baru
type CreateProduct struct {
	wrapper.ProductInsert
	into       map[string]string
	returnID   bool
	returnIDTO *int
}

// NewProduct membuat user baru
// mengembalikan struct User {}
func NewProduct() *CreateProduct {
	return &CreateProduct{
		into: make(map[string]string),
	}
}

// SetName Nama produk
func (u *CreateProduct) SetName(p string) *CreateProduct {
	u.Name = strings.ToLower(p)
	u.into["name"] = ":name"
	return u
}

// SetCode Kode produk
func (u *CreateProduct) SetCode(p string) *CreateProduct {
	u.Code = strings.ToLower(p)
	u.into["code"] = ":code"
	return u
}

// SetBasePrice Harga beli produk
func (u *CreateProduct) SetBasePrice(p string) *CreateProduct {
	u.BasePrice = p
	u.into["base_price"] = ":base_price"
	return u
}

// SetPrice Harga jual produk (kontan/cash)
func (u *CreateProduct) SetPrice(p string) *CreateProduct {
	u.Price = p
	u.into["price"] = ":price"
	return u
}

// SetType Tipe atau model produk
func (u *CreateProduct) SetType(p string) *CreateProduct {
	u.Type = p
	u.into["type"] = ":type"
	return u
}

// SetBrand Brand produk
func (u *CreateProduct) SetBrand(p string) *CreateProduct {
	u.Brand = strings.ToLower(p)
	u.into["brand"] = ":brand"
	return u
}

// SetCreatedBy user yang menambahkan produk
func (u *CreateProduct) SetCreatedBy(p int) *CreateProduct {
	u.CreatedBy = p
	u.into["created_by"] = ":created_by"
	return u
}

// ReturnID Mengembalikan ID produk terakhir
func (u *CreateProduct) ReturnID(id *int) *CreateProduct {
	u.returnID = true
	u.returnIDTO = id
	return u
}

// Save Simpan produk
func (u *CreateProduct) Save(db *sqlx.DB) error {
	return nil
}
