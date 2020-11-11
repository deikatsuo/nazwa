package wrapper

import "database/sql"

// DefaultConfig - untuk menyimpan konfigurasi bawaan
type DefaultConfig struct {
	Site map[string]interface{}
}

// NullableUser struk user dari database
type NullableUser struct {
	ID        int            `db:"id"`
	Firstname string         `db:"first_name"`
	Lastname  sql.NullString `db:"last_name"`
	Username  sql.NullString `db:"username"`
	Avatar    string         `db:"avatar"`
	Gender    string         `db:"gender"`
	CreatedAt string         `db:"created_at"`
	Balance   []uint8        `db:"balance"`
	Password  sql.NullString `db:"password"`
	Phone     sql.NullString `db:"phone"`
	Email     sql.NullString `db:"email"`
	Role      string         `db:"role"`
	Emails    []UserEmail
	Phones    []UserPhone
	Addresses []UserAddress
}

// Fullname menampilkan nama lengkap
func (u NullableUser) Fullname() string {
	return u.Firstname + " " + u.Lastname.String
}

// User buat nge map data user
type User struct {
	ID        int
	Firstname string
	Lastname  string
	Username  string
	Avatar    string
	Gender    string
	CreatedAt string
	Balance   string
	Role      string
	Emails    []UserEmail
	Phones    []UserPhone
	Addresses []UserAddress
}

// UserEmail menampung email user
type UserEmail struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Verified bool   `db:"verified"`
}

// UserPhone menampung email user
type UserPhone struct {
	ID       int    `db:"id"`
	Phone    string `db:"phone"`
	Verified bool   `db:"verified"`
}

// UserAddress menampung data alamat user
type UserAddress struct {
	ID           int    `db:"id"`
	UserID       int    `db:"user_id"`
	Name         string `db:"name" json:"name" binding:"required,min=5,max=50"`
	Description  string `db:"description" json:"description" binding:"omitempty,max=80"`
	One          string `db:"one" json:"one" binding:"required,min=5,max=80"`
	Two          string `db:"two" json:"two" binding:"omitempty,max=80"`
	Zip          string `db:"zip" json:"zip" binding:"required,numeric,min=5,max=5"`
	Province     string `db:"province_id" json:"province" binding:"required,numeric"`
	City         string `db:"city_id" json:"city" binding:"required,numeric"`
	District     string `db:"district_id" json:"district" binding:"required,numeric"`
	Village      string `db:"village_id" json:"village" binding:"required,numeric"`
	ProvinceName string `db:"province_name"`
	CityName     string `db:"city_name"`
	DistrictName string `db:"district_name"`
	VillageName  string `db:"village_name"`
}

// NullableProduct menampilkan list produk
type NullableProduct struct {
	ID            int            `db:"id"`
	Name          string         `db:"name"`
	Code          sql.NullString `db:"code"`
	BasePrice     []uint8        `db:"base_price"`
	Price         []uint8        `db:"price"`
	CreditSix     []uint8        `db:"credit_six"`
	CreditEight   []uint8        `db:"credit_eight"`
	CreditTen     []uint8        `db:"credit_ten"`
	CreditTwelve  []uint8        `db:"credit_twelve"`
	CreditFifteen []uint8        `db:"credit_fifteen"`
	Type          sql.NullString `db:"type"`
	Brand         sql.NullString `db:"brand"`
	CreatedAt     string         `db:"created_at"`
}

// Product map data produk
type Product struct {
	ID            int
	Name          string
	Code          string
	BasePrice     string
	Price         string
	CreditSix     string
	CreditEight   string
	CreditTen     string
	CreditTwelve  string
	CreditFifteen string
	Type          string
	Brand         string
	CreatedAt     string
	Photos        []ProductPhoto
}

// ProductPhoto mengambil photo produk dari database
type ProductPhoto struct {
	ID    int    `db:"id"`
	Photo string `db:"photo"`
}

// NullableOrder menampilkan data order
type NullableOrder struct {
	ID                  int            `db:"id"`
	CustomerID          int            `db:"customer_id"`
	CustomerName        string         `db:"customer_name"`
	SalesID             sql.NullInt64  `db:"sales_id"`
	SalesName           sql.NullString `db:"sales_name"`
	SurveyorID          sql.NullInt64  `db:"surveyor_id"`
	SurveyorName        sql.NullString `db:"surveyor_name"`
	ShippingAddressID   int            `db:"shipping_address_id"`
	ShippingAddressName string         `db:"shipping_address_name"`
	BillingAddressID    sql.NullInt64  `db:"billing_address_id"`
	BillingAddressName  sql.NullString `db:"billing_address_name"`
	Status              string         `db:"status"`
	Code                string         `db:"code"`
	Credit              bool           `db:"credit"`
	FirstTime           bool           `db:"first_time"`
	Notes               sql.NullString `db:"notes"`
	OrderDate           string         `db:"order_date"`
	ShippingDate        sql.NullString `db:"shipping_date"`
}

// Order mapping data order
type Order struct {
	ID              int
	Customer        NameID
	Sales           NameID
	Surveyor        NameID
	ShippingAddress NameID
	BillingAddress  NameID
	Status          string
	Code            string
	Credit          bool
	FirstTime       bool
	Notes           string
	OrderDate       string
	ShippingDate    string
	Items           []OrderItem
}

// NullableOrderItem item/produk yang di order
type NullableOrderItem struct {
	ID        int            `db:"id"`
	ProductID int            `db:"product_id"`
	Quantity  int            `db:"quantity"`
	Notes     sql.NullString `db:"notes"`
}

// OrderItem item/produk yang di order
type OrderItem struct {
	ID        int
	ProductID int
	Quantity  int
	Notes     string
}

// NameID menampilkan nama dan id
type NameID struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
