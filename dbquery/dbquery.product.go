package dbquery

import (
	"fmt"
	"log"
	"nazwa/wrapper"
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
	where     string
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

// Where kondisi
func (p *GetProducts) Where(where string) *GetProducts {
	p.where = where
	return p
}

// NoLimit tidak membatasi hasil query
func (p *GetProducts) NoLimit() *GetProducts {
	return p
}

// Show tampilkan data
func (p *GetProducts) Show(db *sqlx.DB) ([]wrapper.Product, error) {
	var product []wrapper.NullableProduct
	var parse []wrapper.Product
	limit := 0
	if p.limit > 0 {
		limit = p.limit
	}

	// Where logic
	where := p.where

	if limit == 0 {
		where = fmt.Sprintf("%s", where)
	} else {
		where = fmt.Sprintf("%s LIMIT %d", where, limit)
	}

	// query pengambilan data produk
	query := `SELECT
		id,
		name,
		TO_CHAR(base_price,'Rp999G999G999G999G999') AS base_price,
		TO_CHAR(price,'Rp999G999G999G999G999') AS price,
		code,
		thumbnail,
		TO_CHAR(created_at, 'MM/DD/YYYY HH12:MI:SS AM') AS created_at
		FROM "product"
		%s`

	query = fmt.Sprintf(query, where)

	err := db.Select(&product, query)
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
			Thumbnail: p.Thumbnail.String,
		})
	}

	return parse, nil
}

// ProductGetProductTotalRow menghitung jumlah row pada tabel user
func ProductGetProductTotalRow(db *sqlx.DB) (int, error) {
	var total int
	query := `SELECT COUNT(id) FROM "product"`
	err := db.Get(&total, query)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// ProductGetProductByID mengambil data produk berdasarkan ID produk
func ProductGetProductByID(db *sqlx.DB, pid int) (wrapper.Product, error) {
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

	var photos []wrapper.ProductPhotoListSelect

	if pp, err := ProductGetProductPhoto(db, p.ID); err == nil {
		photos = pp
	}

	var creditPrice []wrapper.ProductCreditPriceSelect
	if pr, err := ProductGetProductCreditPrice(db, p.ID); err == nil {
		creditPrice = pr
	}

	product = wrapper.Product{
		ID:          p.ID,
		Name:        strings.Title(p.Name),
		CreatedAt:   p.CreatedAt,
		BasePrice:   string(p.BasePrice),
		Price:       string(p.Price),
		Code:        p.Code,
		Type:        strings.Title(p.Type.String),
		Brand:       strings.Title(p.Brand.String),
		Photos:      photos,
		CreditPrice: creditPrice,
	}

	return product, nil
}

// ProductGetProductCreditPrice mengambil data harga kredit
func ProductGetProductCreditPrice(db *sqlx.DB, pid int) ([]wrapper.ProductCreditPriceSelect, error) {
	var prices []wrapper.ProductCreditPriceSelect
	query := `SELECT id, duration, TO_CHAR(price,'Rp999G999G999G999G999') AS price
	FROM "product_credit_price"
	WHERE product_id=$1`
	err := db.Select(&prices, query, pid)
	if err != nil {
		return []wrapper.ProductCreditPriceSelect{}, err
	}
	return prices, err
}

// ProductGetProductPhoto mengambil data photo produk
func ProductGetProductPhoto(db *sqlx.DB, pid int) ([]wrapper.ProductPhotoListSelect, error) {
	var photos []wrapper.ProductPhotoListSelect
	query := `SELECT id, photo
	FROM "product_photo"
	WHERE product_id=$1`
	err := db.Select(&photos, query, pid)
	if err != nil {
		return []wrapper.ProductPhotoListSelect{}, err
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
	if p != "" {
		c.Code = strings.ToLower(p)
		c.into["code"] = ":code"
	}
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

		if _, err := tx.Exec(`UPDATE "product" SET thumbnail=$1	WHERE id=$2`, c.photos[0], tempReturnID); err != nil {
			log.Println("ERROR: product.go Save() Update thumbnail")
			return err
		}
	}

	if len(c.creditPrice) > 0 {
		for _, cp := range c.creditPrice {
			// Set role user
			if _, err := tx.Exec(`INSERT INTO "product_credit_price" (product_id, duration, price) VALUES ($1, $2, $3)`, tempReturnID, cp.Duration, cp.Price); err != nil {
				log.Println("ERROR: product.go Save() Menambahkan harga produk")
				return err
			}
		}
	}

	// Komit
	err := tx.Commit()
	return err
}

///////////
// CHECK //
///////////

// ProductSkuExist kode produk sudah digunakan
func ProductSkuExist(db *sqlx.DB, sku string) bool {
	// Check bila sku sudah ada di database
	var indb string
	query := `SELECT code FROM "product" WHERE code=$1`
	err := db.Get(&indb, query, sku)
	if err == nil {
		if indb != "" {
			return true
		}
	}

	return false
}
