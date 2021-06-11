package wrapper

import "database/sql"

/* ----------------------------------------------------- */
/* INSERT */

// ProductInsert base struk
type ProductInsert struct {
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Stock       int    `db:"stock"`
	BasePrice   int    `db:"base_price"`
	Price       int    `db:"price"`
	Category    string `db:"category"`
	Brand       string `db:"brand"`
	CreatedAt   string `db:"created_at"`
	CreatedBy   int    `db:"created_by"`
	Description string `db:"description"`
}

/* ----------------------------------------------------- */
/* SELECT */

// ProductPhotoListSelect mengambil photo produk dari database
type ProductPhotoListSelect struct {
	ID    int    `db:"id"`
	Photo string `db:"photo"`
}

// NullableProduct menampilkan list produk
type NullableProduct struct {
	ID          int            `db:"id"`
	Name        string         `db:"name"`
	Slug        string         `db:"slug"`
	Stock       int            `db:"stock"`
	BasePrice   int            `db:"base_price"`
	Price       int            `db:"price"`
	Category    sql.NullString `db:"category"`
	Brand       sql.NullString `db:"brand"`
	Thumbnail   sql.NullString `db:"thumbnail"`
	Description sql.NullString `db:"description"`
	CreatedAt   string         `db:"created_at"`
	CreatedBy   int            `db:"created_by"`
}

// ProductCreditPriceInsert tambah harga kredit barang
type ProductCreditPriceInsert struct {
	ProductID int `db:"product_id"`
	Duration  int `db:"duration"`
	Price     int `db:"price"`
}

// ProductCreditPriceForm harga kredit barang
type ProductCreditPriceForm struct {
	Duration int `json:"duration"`
	Price    int `json:"price"`
}

// ProductCreditPriceSelect list harga kredit produk
type ProductCreditPriceSelect struct {
	ID       int `db:"id"`
	Duration int `db:"duration"`
	Price    int `db:"price"`
}

/* ----------------------------------------------------- */
/* FORM VALIDATION */

// ProductForm menyimpan input produk
type ProductForm struct {
	Name        string                   `json:"name" binding:"required,min=4,max=100"`
	Stock       string                   `json:"stock" binding:"numeric,gte=0"`
	Brand       string                   `json:"brand" binding:"omitempty,min=2,max=25"`
	Category    string                   `json:"category" binding:"omitempty,min=2,max=25"`
	BasePrice   string                   `json:"base_price" binding:"required,numeric,min=1,max=15"`
	Price       string                   `json:"price" binding:"required,numeric,min=1,max=15"`
	CreditPrice []ProductCreditPriceForm `json:"credit_price" binding:"omitempty"`
	Photo       []ProductPhotoForm       `json:"photos" binding:"omitempty"`
	Description string                   `json:"description" binding:"omitempty"`
}

// ProductPhotoForm list photo produk
type ProductPhotoForm struct {
	PhotoType string `json:"photo_type" binding:"base64"`
	Photo     string `json:"photo"`
}

/* ----------------------------------------------------- */
/* SHOW/VIEW */

// Product map data produk
type Product struct {
	ID          int
	Name        string
	Slug        string
	Description string
	Stock       int
	BasePrice   int
	Price       int
	Category    string
	Brand       string
	Thumbnail   string
	CreatedAt   string
	CreatedBy   int
	CreditPrice []ProductCreditPriceSelect
	Photos      []ProductPhotoListSelect
}
