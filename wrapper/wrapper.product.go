package wrapper

import "database/sql"

/* ----------------------------------------------------- */
/* INSERT */

// ProductInsert base struk
type ProductInsert struct {
	Name      string `db:"name"`
	Code      string `db:"code"`
	BasePrice string `db:"base_price"`
	Price     string `db:"price"`
	Type      string `db:"type"`
	Brand     string `db:"brand"`
	CreatedAt string `db:"created_at"`
	CreatedBy int    `db:"created_by"`
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
	ID        int            `db:"id"`
	Name      string         `db:"name"`
	Code      string         `db:"code"`
	BasePrice []uint8        `db:"base_price"`
	Price     []uint8        `db:"price"`
	Type      sql.NullString `db:"type"`
	Brand     sql.NullString `db:"brand"`
	CreatedAt string         `db:"created_at"`
	CreatedBy int            `db:"created_by"`
}

// ProductCreditPriceSelect list harga kredit produk
type ProductCreditPriceSelect struct {
	ID       int    `db:"id"`
	Duration int    `db:"duration"`
	Price    string `db:"price"`
}

/* ----------------------------------------------------- */
/* FORM VALIDATION */

// ProductForm menyimpan input produk
type ProductForm struct {
	Name        string                   `json:"name" binding:"required,min=4,max=100"`
	Code        string                   `json:"code" binding:"omitempty,min=5,max=10"`
	Brand       string                   `json:"brand" binding:"omitempty,min=2,max=25"`
	Type        string                   `json:"type" binding:"omitempty,min=2,max=25"`
	BasePrice   string                   `json:"base_price" binding:"required,numeric,min=1,max=15"`
	Price       string                   `json:"price" binding:"required,numeric,min=1,max=15"`
	CreditPrice []ProductCreditPriceForm `json:"credit_price" binding:"omitempty"`
	Photo       []ProductPhotoForm       `json:"photos" binding:"omitempty"`
}

// ProductCreditPriceForm harga kredit barang
type ProductCreditPriceForm struct {
	Duration int `json:"duration"`
	Price    int `json:"price"`
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
	Code        string
	BasePrice   string
	Price       string
	Type        string
	Brand       string
	CreatedAt   string
	CreatedBy   int
	CreditPrice []ProductCreditPriceSelect
	Photos      []ProductPhotoListSelect
}
