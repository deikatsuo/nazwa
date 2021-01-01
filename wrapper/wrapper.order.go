package wrapper

import (
	"database/sql"
)

// OrderInsert base struk
type OrderInsert struct {
	CustomerID        int    `db:"customer_id"`
	SalesID           int    `db:"sales_id"`
	CollectorID       int    `db:"collector_id"`
	SurveyorID        int    `db:"surveyor_id"`
	DriverID          int    `db:"driver_id"`
	ShippingAddressID int    `db:"shipping_address_id"`
	BillingAddressID  int    `db:"billing_address_id"`
	Code              string `db:"code"`
	Status            string `db:"status"`
	Credit            bool   `db:"credit"`

	Notes        string `db:"notes"`
	OrderDate    string `db:"order_date"`
	ShippingDate string `db:"shipping_date"`

	CreatedAt string `db:"created_at"`
	CreatedBy int    `db:"created_by"`

	Deposit        int `db:"deposit"`
	PriceTotal     int `db:"price_total"`
	BasePriceTotal int `db:"base_price_total"`
}

// OrderItemInsert insert item
type OrderItemInsert struct {
	OrderID   int    `db:"order_id"`
	ProductID int    `db:"product_id"`
	Quantity  int    `db:"quantity"`
	Notes     string `db:"notes"`
	BasePrice int    `db:"base_price"`
	Price     int    `db:"price"`
	Discount  int    `db:"discount"`
}

// OrderForm formulir buat order
type OrderForm struct {
	Customer        int             `json:"customer" binding:"required,numeric"`
	Sales           int             `json:"sales" binding:"omitempty,numeric"`
	Surveyor        int             `json:"surveyor" binding:"omitempty,numeric"`
	Collector       int             `json:"collector" binding:"omitempty,numeric"`
	Credit          *bool           `json:"credit" binding:"required"`
	Duration        int             `json:"duration" binding:"omitempty,numeric"`
	Notes           string          `json:"notes" binding:"omitempty"`
	OrderDate       string          `json:"order_date" binding:"required,date"`
	ShippingDate    string          `json:"shipping_date" binding:"omitempty,date"`
	BillingAddress  int             `json:"billing_address" binding:"omitempty,numeric"`
	ShippingAddress int             `json:"shipping_address" binding:"required,numeric"`
	Deposit         int             `json:"deposit" binding:"omitempty,numeric"`
	OrderItems      []OrderItemForm `json:"order_items" binding:"omitempty"`
}

// OrderItemForm form item
type OrderItemForm struct {
	ProductID int    `json:"id" binding:"required,numeric"`
	Quantity  int    `json:"quantity" binding:"required,numeric"`
	Notes     string `json:"notes" binding:"omitempty"`
	Discount  int    `json:"discount" binding:"omitempty,numeric"`
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
	CollectorID         sql.NullInt64  `db:"collector_id"`
	CollectorName       sql.NullString `db:"collector_name"`
	DriverID            sql.NullInt64  `db:"driver_id"`
	DriverName          sql.NullString `db:"driver_name"`
	ShippingAddressID   int            `db:"shipping_address_id"`
	ShippingAddressName string         `db:"shipping_address_name"`
	BillingAddressID    sql.NullInt64  `db:"billing_address_id"`
	BillingAddressName  sql.NullString `db:"billing_address_name"`
	Status              string         `db:"status"`
	Code                string         `db:"code"`
	Credit              bool           `db:"credit"`
	Notes               sql.NullString `db:"notes"`
	OrderDate           string         `db:"order_date"`
	ShippingDate        sql.NullString `db:"shipping_date"`
	CreatedAt           string         `db:"created_at"`
	CreatedBy           int            `db:"created_by"`
	Deposit             int            `db:"deposit"`
	PriceTotal          int            `db:"price_total"`
	BasePriceTotal      int            `db:"base_price_total"`
}

// Order mapping data order
type Order struct {
	ID              int
	Customer        NameID
	Sales           NameID
	Surveyor        NameID
	Collector       NameID
	Driver          NameID
	ShippingAddress NameID
	BillingAddress  NameID
	Status          string
	Code            string
	Credit          bool
	Notes           string
	OrderDate       string
	ShippingDate    string
	Items           []OrderItem
	CreatedAt       string
	CreatedBy       NameID
	Deposit         int
	PriceTotal      int
	BasePriceTotal  int
}

// NullableOrderItem item/produk yang di order
type NullableOrderItem struct {
	ID          int            `db:"id"`
	ProductID   int            `db:"product_id"`
	ProductName string         `db:"name"`
	ProductCode string         `db:"code"`
	Thumbnail   string         `db:"thumbnail"`
	Quantity    int            `db:"quantity"`
	Notes       sql.NullString `db:"notes"`
	BasePrice   int            `db:"base_price"`
	Price       int            `db:"price"`
	Discount    int            `db:"discount"`
}

// OrderItem item/produk yang di order
type OrderItem struct {
	ID        int
	Product   NameIDCode
	Quantity  int
	Notes     string
	BasePrice int
	Price     int
	Discount  int
}
