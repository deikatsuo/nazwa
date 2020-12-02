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
	into        map[string]string
	returnID    bool
	returnIDTO  *int
	creditPrice []wrapper.ProductCreditPriceForm
	photos      []string
}

// NewProduct membuat user baru
// mengembalikan struct User {}
func NewProduct() *CreateProduct {
	return &CreateProduct{
		into: make(map[string]string),
	}
}

// SetName Nama produk
func (c *CreateProduct) SetName(p string) *CreateProduct {
	c.Name = strings.ToLower(p)
	c.into["name"] = ":name"
	return c
}

// SetCode Kode produk
func (c *CreateProduct) SetCode(p string) *CreateProduct {
	c.Code = strings.ToLower(p)
	c.into["code"] = ":code"
	return c
}

// SetBasePrice Harga beli produk
func (c *CreateProduct) SetBasePrice(p string) *CreateProduct {
	c.BasePrice = p
	c.into["base_price"] = ":base_price"
	return c
}

// SetPrice Harga jual produk (kontan/cash)
func (c *CreateProduct) SetPrice(p string) *CreateProduct {
	c.Price = p
	c.into["price"] = ":price"
	return c
}

// SetType Tipe atau model produk
func (c *CreateProduct) SetType(p string) *CreateProduct {
	c.Type = p
	c.into["type"] = ":type"
	return c
}

// SetBrand Brand produk
func (c *CreateProduct) SetBrand(p string) *CreateProduct {
	c.Brand = strings.ToLower(p)
	c.into["brand"] = ":brand"
	return c
}

// SetCreditPrice harga kredit barang
func (c *CreateProduct) SetCreditPrice(p []wrapper.ProductCreditPriceForm) *CreateProduct {
	c.creditPrice = p
	return c
}

// SetPhotos Brand produk
func (c *CreateProduct) SetPhotos(p []string) *CreateProduct {
	c.photos = p
	return c
}

// SetCreatedBy user yang menambahkan produk
func (c *CreateProduct) SetCreatedBy(p int) *CreateProduct {
	c.CreatedBy = p
	c.into["created_by"] = ":created_by"
	return c
}

// ReturnID Mengembalikan ID produk terakhir
func (c *CreateProduct) ReturnID(id *int) *CreateProduct {
	c.returnID = true
	c.returnIDTO = id
	return c
}

// Insert query berdasarka data yang diisi
func (c CreateProduct) generateInsertQuery() string {
	iq := c.into
	var kk []string
	var kv []string
	for k, v := range iq {
		kk = append(kk, k)
		kv = append(kv, v)
	}
	result := fmt.Sprintf("(%s) VALUES (%s) RETURNING id", strings.Join(kk, ","), strings.Join(kv, ","))

	return result
}

// Save Simpan produk
func (c *CreateProduct) Save(db *sqlx.DB) error {
	// Mulai transaksi
	tx := db.MustBegin()
	var tempReturnID int
	productInsertQuery := fmt.Sprintf(`INSERT INTO "product" %s`, c.generateInsertQuery())
	if rows, err := tx.NamedQuery(productInsertQuery, c); err == nil {
		// Ambil id dari transaksi terakhir
		if rows.Next() {
			rows.Scan(&tempReturnID)
		}

		if c.returnID && tempReturnID != 0 {
			*c.returnIDTO = tempReturnID
		}

		if err := rows.Close(); err != nil {
			return err
		}
	} else {
		tx.Rollback()
		return err
	}

	if len(c.photos) > 0 {
		for id, s := range c.photos {
			// Set role user
			if _, err := tx.Exec(`INSERT INTO "product_photo" (product_id, photo) VALUES ($1, $2)`, tempReturnID, s); err != nil {
				log.Println("ERROR: product.go Save() Insert photo ID: ", id)
				return err
			}
		}
	}

	// Komit
	err := tx.Commit()
	return err
}
