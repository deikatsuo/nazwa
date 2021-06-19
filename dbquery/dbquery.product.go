package dbquery

import (
	"database/sql"
	"fmt"
	"nazwa/wrapper"
	"strings"
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
func (p *GetProducts) Show() ([]wrapper.Product, error) {
	db := DB
	var product []wrapper.NullableProduct
	var parse []wrapper.Product
	limit := 0
	if p.limit > 0 {
		limit = p.limit
	}

	// Where logic
	where := p.where

	if limit != 0 {
		where = fmt.Sprintf("%s LIMIT %d", where, limit)
	}

	// query pengambilan data produk
	query := `SELECT
		id,
		name,
		slug,
		stock,
		brand,
		category,
		base_price,
		price,
		thumbnail,
		created_by,
		TO_CHAR(created_at, 'DD-MM-YYYY HH12:MI:SS AM') AS created_at
		FROM "product"
		%s`

	query = fmt.Sprintf(query, where)

	err := db.Select(&product, query)
	if err != nil {
		log.Warn("dbquery.product.go (p *GetProducts) Show() Select all product")
		return []wrapper.Product{}, err
	}

	// mapping data produk
	for _, p := range product {
		var creditPrice []wrapper.ProductCreditPriceSelect
		if pr, err := ProductGetProductCreditPrice(p.ID); err == nil {
			creditPrice = pr
		} else {
			log.Println(err)
		}

		parse = append(parse, wrapper.Product{
			ID:          p.ID,
			Name:        strings.Title(p.Name),
			Slug:        p.Slug,
			Stock:       p.Stock,
			Brand:       p.Brand.String,
			Category:    p.Category.String,
			CreatedAt:   p.CreatedAt,
			CreatedBy:   p.CreatedBy,
			BasePrice:   p.BasePrice,
			Price:       p.Price,
			Thumbnail:   p.Thumbnail.String,
			CreditPrice: creditPrice,
		})
	}

	return parse, nil
}

// ProductGetProductTotalRow menghitung jumlah row pada tabel user
func ProductGetProductTotalRow() (int, error) {
	db := DB
	var total int
	query := `SELECT COUNT(id) FROM "product"`
	err := db.Get(&total, query)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// ProductGetProductByID mengambil data produk berdasarkan ID produk
func ProductGetProductByID(pid int) (wrapper.Product, error) {
	db := DB
	var product wrapper.Product
	var p wrapper.NullableProduct
	query := `SELECT
		id,
		name,
		stock,
		base_price,
		price,
		TO_CHAR(created_at, 'DD-MM-YYYY HH12:MI:SS AM') AS created_at,
		category,
		brand,
		description
		FROM "product"
		WHERE id=$1
		LIMIT 1`

	err := db.Get(&p, query, pid)
	if err != nil {
		log.Warn("dbquery.product.go ProductGetProductByID() Select product berdasarkan ID")
		return wrapper.Product{}, err
	}

	var photos []wrapper.ProductPhotoListSelect

	if pp, err := ProductGetProductPhoto(p.ID); err == nil {
		photos = pp
	}

	var creditPrice []wrapper.ProductCreditPriceSelect
	if pr, err := ProductGetProductCreditPrice(p.ID); err == nil {
		creditPrice = pr
	}

	product = wrapper.Product{
		ID:          p.ID,
		Name:        strings.Title(p.Name),
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		BasePrice:   p.BasePrice,
		Price:       p.Price,
		Category:    strings.Title(p.Category.String),
		Brand:       strings.Title(p.Brand.String),
		Description: p.Description.String,
		Photos:      photos,
		CreditPrice: creditPrice,
	}

	return product, nil
}

// ProductGetProductBySlug mengambil data produk berdasarkan ID produk
func ProductGetProductBySlug(ps string) (wrapper.Product, error) {
	db := DB
	var product wrapper.Product
	var p wrapper.NullableProduct
	query := `SELECT
		id,
		name,
		stock,
		description,
		base_price,
		price,
		TO_CHAR(created_at, 'DD-MM-YYYY HH12:MI:SS AM') AS created_at,
		category,
		brand
		FROM "product"
		WHERE slug=$1
		LIMIT 1`

	err := db.Get(&p, query, ps)
	if err != nil {
		log.Warn("dbquery.product.go ProductGetProductBySlug() Select product berdasarkan slug")
		log.Error(err)
		return wrapper.Product{}, err
	}

	var photos []wrapper.ProductPhotoListSelect

	if pp, err := ProductGetProductPhoto(p.ID); err == nil {
		photos = pp
	}

	var creditPrice []wrapper.ProductCreditPriceSelect
	if pr, err := ProductGetProductCreditPrice(p.ID); err == nil {
		creditPrice = pr
	}

	product = wrapper.Product{
		ID:          p.ID,
		Name:        strings.Title(p.Name),
		Stock:       p.Stock,
		Description: p.Description.String,
		CreatedAt:   p.CreatedAt,
		Slug:        ps,
		BasePrice:   p.BasePrice,
		Price:       p.Price,
		Category:    strings.Title(p.Category.String),
		Brand:       strings.Title(p.Brand.String),
		Photos:      photos,
		CreditPrice: creditPrice,
	}

	return product, nil
}

// ProductGetProductCreditPrice mengambil data harga kredit
func ProductGetProductCreditPrice(pid int) ([]wrapper.ProductCreditPriceSelect, error) {
	db := DB
	var prices []wrapper.ProductCreditPriceSelect
	query := `SELECT id, duration, price
	FROM "product_credit_price"
	WHERE product_id=$1`
	err := db.Select(&prices, query, pid)
	if err != nil {
		return []wrapper.ProductCreditPriceSelect{}, err
	}
	return prices, err
}

// ProductGetProductPrice mengambil harga barang
func ProductGetProductPrice(pid int) (int, error) {
	db := DB
	var price int
	query := `SELECT
		price
		FROM "product"
		WHERE id=$1
		LIMIT 1`
	err := db.Get(&price, query, pid)
	if err != nil {
		log.Warn("dbquery.product.go ProductGetProductPrice() error saat mengambil data harga jual")
		return 0, err
	}

	return price, nil
}

// ProductGetProductBasePrice mengambil harga barang
func ProductGetProductBasePrice(pid int) (int, error) {
	db := DB
	var price int
	query := `SELECT
		base_price
		FROM "product"
		WHERE id=$1
		LIMIT 1`
	err := db.Get(&price, query, pid)
	if err != nil {
		log.Warn("dbquery.product.go ProductGetProductBasePrice() error saat mengambil data harga beli")
		return 0, err
	}

	return price, nil
}

// ProductGetProductPhoto mengambil data photo produk
func ProductGetProductPhoto(pid int) ([]wrapper.ProductPhotoListSelect, error) {
	db := DB
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

// ProductGetProductPhotoName mengambil nama foto
func ProductGetProductPhotoName(photoID int) (string, error) {
	db := DB
	var photo string
	query := `SELECT photo
	FROM "product_photo"
	WHERE id=$1`
	err := db.Get(&photo, query, photoID)
	if err != nil {
		return photo, err
	}
	return photo, err
}

// ProductGetProductThumbName mengambil thumb foto
func ProductGetProductThumbName(pid int) (string, error) {
	db := DB
	var thumb sql.NullString
	query := `SELECT thumbnail
	FROM "product"
	WHERE id=$1`
	err := db.Get(&thumb, query, pid)
	if err != nil {
		return thumb.String, err
	}
	return thumb.String, err
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

// SetSlug slug url
func (c *CreateProduct) SetSlug(p string) *CreateProduct {
	if p != "" {
		c.Slug = p
		c.into["slug"] = ":slug"
	}
	return c
}

// SetStock Stok produk
func (c *CreateProduct) SetStock(p int) *CreateProduct {
	if p >= 0 {
		c.Stock = p
		c.into["stock"] = ":stock"
	}
	return c
}

// SetBasePrice Harga beli produk
func (c *CreateProduct) SetBasePrice(p int) *CreateProduct {
	c.BasePrice = p
	c.into["base_price"] = ":base_price"
	return c
}

// SetDescription deskripsi produk
func (c *CreateProduct) SetDescription(p string) *CreateProduct {
	c.Description = p
	c.into["description"] = ":description"
	return c
}

// SetPrice Harga jual produk (kontan/cash)
func (c *CreateProduct) SetPrice(p int) *CreateProduct {
	c.Price = p
	c.into["price"] = ":price"
	return c
}

// SetCategory kategori produk
func (c *CreateProduct) SetCategory(p string) *CreateProduct {
	c.Category = p
	c.into["category"] = ":category"
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
func (c *CreateProduct) Save() error {
	db := DB

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
			// Set foto produk
			if _, err := tx.Exec(`INSERT INTO "product_photo" (product_id, photo) VALUES ($1, $2)`, tempReturnID, s); err != nil {
				log.Warn("dbquery.product.go (c *CreateProduct) Save() Insert photo ID: ", id)
				return err
			}
		}

		if _, err := tx.Exec(`UPDATE "product" SET thumbnail=$1	WHERE id=$2`, c.photos[0], tempReturnID); err != nil {
			log.Warn("dbquery.product.go (c *CreateProduct) Save() Update thumbnail")
			return err
		}
	}

	if len(c.creditPrice) > 0 {
		for _, cp := range c.creditPrice {
			// Set harga kredit untuk produk
			if _, err := tx.Exec(`INSERT INTO "product_credit_price" (product_id, duration, price) VALUES ($1, $2, $3)`, tempReturnID, cp.Duration, cp.Price); err != nil {
				log.Warn("dbquery.product.go (c *CreateProduct) Save() Menambahkan harga produk")
				return err
			}
		}
	}

	// Komit
	err := tx.Commit()
	return err
}

// ProductInsertCreditPrice menambahkan harga kredit produk ke database
func ProductInsertCreditPrice(cps []wrapper.ProductCreditPriceInsert) error {
	db := DB
	query := `INSERT INTO "product_credit_price" (product_id, duration, price) VALUES (:product_id, :duration, :price)`
	if _, err := db.NamedQuery(query, cps); err != nil {
		log.Warn("dbquery.product.go ProductInsertCreditPrice() Gagal menambahkan harga kredit")
		return err
	}

	return nil
}

// ProductDeleteCreditPrice menghapus harga kredit barang
func ProductDeleteCreditPrice(pcpid int64, pid int) error {
	db := DB
	query := `DELETE FROM "product_credit_price"
	WHERE id=$1 AND product_id=$2`
	_, err := db.Exec(query, pcpid, pid)
	return err
}

// ProductDeletePhoto menghapus foto produk
func ProductDeletePhoto(foto int64, pid int) error {
	db := DB
	query := `DELETE FROM "product_photo"
	WHERE id=$1 AND product_id=$2`
	_, err := db.Exec(query, foto, pid)
	return err
}

///////////
// CHECK //
///////////

// ProductCreditDurationExist chek apakah durasi kredit sudah ada
func ProductCreditDurationExist(pid int, dur int) bool {
	db := DB
	var indb int
	query := `SELECT duration FROM "product_credit_price" WHERE product_id=$1 AND duration=$2 LIMIT 1`
	err := db.Get(&indb, query, pid, dur)
	if err == nil {
		if indb != 0 {
			return true
		}

	}

	return false
}

// ProductCheckStock cek stok produk
func ProductCheckStock(pid int) (int, error) {
	db := DB
	var stock int

	query := `SELECT stock FROM "product" WHERE id=$1 LIMIT 1`
	err := db.Get(&stock, query, pid)

	return stock, err
}

// ProductUpdateStock mengubah jumlah stok produk
func ProductUpdateStock(pid, stock int) error {
	db := DB
	query := `UPDATE "product"
	SET stock=$1
	WHERE id=$2`
	_, err := db.Exec(query, stock, pid)

	return err
}

// ProductUpdateName mengubah nama produk
func ProductUpdateName(pid int, name string) error {
	db := DB
	query := `UPDATE "product"
	SET name=$1
	WHERE id=$2`
	_, err := db.Exec(query, name, pid)

	return err
}

// ProductUpdateBrand mengubah brand produk
func ProductUpdateBrand(pid int, brand string) error {
	db := DB
	query := `UPDATE "product"
	SET brand=$1
	WHERE id=$2`
	_, err := db.Exec(query, brand, pid)

	return err
}

// ProductUpdateCategory mengubah kategori produk
func ProductUpdateCategory(pid int, category string) error {
	db := DB
	query := `UPDATE "product"
	SET category=$1
	WHERE id=$2`
	_, err := db.Exec(query, category, pid)

	return err
}

// ProductUpdateDescription mengubah deskripsi produk
func ProductUpdateDescription(pid int, desc string) error {
	db := DB
	query := `UPDATE "product"
	SET description=$1
	WHERE id=$2`
	_, err := db.Exec(query, desc, pid)

	return err
}

// ProductUpdatePriceBuy ubah harga beli
func ProductUpdatePriceBuy(pid int, bp int) error {
	db := DB
	query := `UPDATE "product"
	SET base_price=$1
	WHERE id=$2`
	_, err := db.Exec(query, bp, pid)

	return err
}

// ProductUpdatePriceSell ubah harga jual
func ProductUpdatePriceSell(pid int, p int) error {
	db := DB
	query := `UPDATE "product"
	SET price=$1
	WHERE id=$2`
	_, err := db.Exec(query, p, pid)

	return err
}

// ProductUpdateThumb ubah thumbnail
func ProductUpdateThumb(pid int, thumb sql.NullString) error {
	db := DB
	query := `UPDATE "product"
	SET thumbnail=$1
	WHERE id=$2`
	_, err := db.Exec(query, thumb, pid)

	return err
}
