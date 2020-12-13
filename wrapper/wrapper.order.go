package wrapper

import "database/sql"

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
	Collector       NameID
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
	ID          int            `db:"id"`
	ProductID   int            `db:"product_id"`
	ProductName string         `db:"name"`
	ProductCode string         `db:"code"`
	Quantity    int            `db:"quantity"`
	Notes       sql.NullString `db:"notes"`
}

// OrderItem item/produk yang di order
type OrderItem struct {
	ID       int
	Product  NameIDCode
	Quantity int
	Notes    string
}
