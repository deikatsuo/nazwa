package wrapper

import "database/sql"

// ListProductPhoto mengambil photo produk dari database
type ListProductPhoto struct {
	ID    int    `db:"id"`
	Photo string `db:"photo"`
}

// ProductPhoto list photo produk
type ProductPhoto struct {
	PhotoType string `json:"photo_type" binding:"omitempty,base64"`
	Photo     string `json:"photo"`
}

// FormProduct menyimpan input produk
type FormProduct struct {
	Name      string         `json:"name" binding:"alphanumeric,min=4,max=100"`
	Brand     string         `json:"brand" binding:"omitempty,alphanumeric,min=2,max=25"`
	Type      string         `json:"type" binding:"omitempty,alphanumeric,min=2,max=25"`
	BasePrice string         `json:"base_price" binding:"required,numeric,min=1,max=15"`
	Price     string         `json:"price" binding:"required,numeric,min=1,max=15"`
	Photo     []ProductPhoto `json:"photo" binding:"omitempty"`
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

// Product map data produk
type Product struct {
	ID        int
	Name      string
	Code      string
	BasePrice string
	Price     string
	Type      string
	Brand     string
	CreatedAt string
	CreatedBy int
	Photos    []ListProductPhoto
}

// ProductCreditPrice list harga kredit produk
type ProductCreditPrice struct {
	ID       int    `db:"id"`
	Duration int    `db:"duration"`
	Price    string `db:"price"`
}
