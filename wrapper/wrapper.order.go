package wrapper

import (
	"database/sql"
)

// OrderInsert base struk
type OrderInsert struct {
	CustomerID        int    `db:"customer_id"`
	SalesID           int    `db:"sales_id"`
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

	ImportedItems    string `db:"imported_items"`
	ImportedAddress  string `db:"imported_address"`
	ImportedSales    string `db:"imported_sales"`
	ImportedSurveyor string `db:"imported_surveyor"`
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
	Customer        int                       `json:"customer" binding:"required,numeric"`
	Sales           int                       `json:"sales" binding:"omitempty,numeric"`
	Surveyor        int                       `json:"surveyor" binding:"omitempty,numeric"`
	Credit          *bool                     `json:"credit" binding:"required"`
	CreditMask      string                    `json:"credit_mask"`
	Duration        int                       `json:"duration" binding:"omitempty,numeric,gte=1,lte=24"`
	Due             int                       `json:"due" binding:"required_if=CreditMask credit,omitempty,numeric,gte=1,lte=28"`
	Line            int                       `json:"line" binding:"required_if=CreditMask credit,omitempty,numeric"`
	Notes           string                    `json:"notes" binding:"omitempty"`
	OrderDate       string                    `json:"order_date" binding:"required,date"`
	ShippingDate    string                    `json:"shipping_date" binding:"omitempty,date"`
	BillingAddress  int                       `json:"billing_address" binding:"omitempty,numeric"`
	ShippingAddress int                       `json:"shipping_address" binding:"required,numeric"`
	Deposit         string                    `json:"deposit" binding:"omitempty,numeric"`
	OrderItems      []OrderItemForm           `json:"order_items" binding:"omitempty"`
	Substitutes     []OrderUserSubstituteForm `json:"substitutes" binding:"omitempty"`
}

// OrderItemForm form item
type OrderItemForm struct {
	ProductID int    `json:"id" binding:"required,numeric"`
	Quantity  int    `json:"quantity" binding:"required,numeric,gte=0"`
	Notes     string `json:"notes" binding:"omitempty"`
	Discount  int    `json:"discount" binding:"omitempty,numeric"`
}

// OrderUserSubstituteForm penanggung jawab/pendamping user (konsumen)
type OrderUserSubstituteForm struct {
	RIC       string `json:"ric"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
}

// NullableOrder menampilkan data order
type NullableOrder struct {
	ID                int            `db:"id"`
	CustomerID        int            `db:"customer_id"`
	CustomerCode      sql.NullString `db:"customer_code"`
	CustomerName      string         `db:"customer_name"`
	CustomerThumb     string         `db:"customer_thumb"`
	SalesID           sql.NullInt64  `db:"sales_id"`
	SalesName         sql.NullString `db:"sales_name"`
	SalesThumb        sql.NullString `db:"sales_thumb"`
	SurveyorID        sql.NullInt64  `db:"surveyor_id"`
	SurveyorName      sql.NullString `db:"surveyor_name"`
	SurveyorThumb     sql.NullString `db:"surveyor_thumb"`
	CollectorID       sql.NullInt64  `db:"collector_id"`
	CollectorName     sql.NullString `db:"collector_name"`
	CollectorThumb    sql.NullString `db:"collector_thumb"`
	DriverID          sql.NullInt64  `db:"driver_id"`
	DriverName        sql.NullString `db:"driver_name"`
	DriverThumb       sql.NullString `db:"driver_thumb"`
	CreatedByID       sql.NullInt64  `db:"created_by_id"`
	CreatedByName     sql.NullString `db:"created_by_name"`
	CreatedByThumb    sql.NullString `db:"created_by_thumb"`
	ShippingAddressID sql.NullInt64  `db:"shipping_address_id"`
	BillingAddressID  sql.NullInt64  `db:"billing_address_id"`
	Status            string         `db:"status"`
	Code              string         `db:"code"`
	CreditCode        string         `db:"credit_code"`
	Credit            bool           `db:"credit"`
	Notes             sql.NullString `db:"notes"`
	OrderDate         string         `db:"order_date"`
	ShippingDate      string         `db:"shipping_date"`
	CreatedAt         string         `db:"created_at"`
	Deposit           int            `db:"deposit"`
	PriceTotal        int            `db:"price_total"`
	BasePriceTotal    sql.NullInt64  `db:"base_price_total"`

	ImportedItems    sql.NullString `db:"imported_items"`
	ImportedAddress  sql.NullString `db:"imported_address"`
	ImportedSales    sql.NullString `db:"imported_sales"`
	ImportedSurveyor sql.NullString `db:"imported_surveyor"`
}

// Order mapping data order
type Order struct {
	ID              int
	Customer        NameIDCode
	Sales           NameID
	Surveyor        NameID
	Collector       NameID
	Driver          NameID
	ShippingAddress string
	BillingAddress  string
	Status          string
	Code            string
	CreditCode      string
	QR              []byte
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
	CreditDetail    OrderCreditDetail
	MonthlyCredit   []OrderMonthlyCredit
	PaymentLogs     []OrderPayment

	ImportedItems    string
	ImportedAddress  string
	ImportedSales    string
	ImportedSurveyor string
}

// OrderSimple data order sederhana
type OrderSimple struct {
	ID             int
	Customer       NameIDCode
	Sales          NameID
	Surveyor       NameID
	Collector      NameID
	BillingAddress string
	Code           string
	Credit         bool
	ShippingDate   string
	Items          []OrderItem
	Deposit        int
	PriceTotal     int
	BasePriceTotal int
	CreditDetail   OrderCreditDetail

	ImportedItems    string
	ImportedAddress  string
	ImportedSales    string
	ImportedSurveyor string
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

// OrderCreditDetailInsert detail order kredit
type OrderCreditDetailInsert struct {
	OrderID       int `db:"order_id"`
	ZoneID        int `db:"zone_id"`
	Monthly       int `db:"monthly"`
	Duration      int `db:"duration"`
	Due           int `db:"due"`
	Total         int `db:"total"`
	Remaining     int `db:"remaining"`
	LuckyDiscount int `db:"lucky_discount"`
}

// OrderCreditDetailSelect detail order kredit
type OrderCreditDetailSelect struct {
	ID            int            `db:"id"`
	OrderID       int            `db:"order_id"`
	ZoneLineID    sql.NullInt32  `db:"zone_line_id"`
	ZoneLineCode  sql.NullString `db:"zone_line_code"`
	ZoneLineName  sql.NullString `db:"zone_line_name"`
	CreditCode    string         `db:"credit_code"`
	Monthly       int            `db:"monthly"`
	Duration      int            `db:"duration"`
	Due           int            `db:"due"`
	Total         int            `db:"total"`
	Remaining     int            `db:"remaining"`
	LuckyDiscount int            `db:"lucky_discount"`
	LastPaid      sql.NullString `db:"last_paid"`
	Active        bool           `db:"active"`
	Done          bool           `db:"done"`
}

// OrderCreditDetail detail order kredit
type OrderCreditDetail struct {
	ID            int
	OrderID       int
	ZoneLine      NameIDCode
	CreditCode    string
	Monthly       int
	Duration      int
	Due           int
	Total         int
	Remaining     int
	LuckyDiscount int
	LastPaid      string
	Active        bool
	Done          bool
}

// OrderMonthlyCreditSelect kredit bulanan
type OrderMonthlyCreditSelect struct {
	ID        int            `db:"id"`
	OrderID   int            `db:"order_id"`
	Code      string         `db:"code"`
	Nth       int            `db:"nth"`
	DueDate   string         `db:"due_date"`
	PrintDate sql.NullString `db:"print_date"`
	Promise   sql.NullString `db:"promise"`
	Paid      int            `db:"paid"`
	Notes     sql.NullString `db:"notes"`
	Position  string         `db:"position"`
	Printed   bool           `db:"printed"`
	Done      bool           `db:"done"`
}

// OrderMonthlyCredit kredit tagihan bulanan
type OrderMonthlyCredit struct {
	ID        int
	OrderID   int
	Code      string
	QR        []byte
	Nth       int
	DueDate   string
	PrintDate string
	Promise   string
	Paid      int
	Notes     string
	Position  string
	Printed   bool
	Print     bool
	// Kwitansi sudah selesai
	Done bool
}

// OrderPaymentInsert pembayaran
type OrderPaymentInsert struct {
	ID               int    `db:"id"`
	OrderID          int    `db:"order_id"`
	ReceiverID       int    `db:"receiver_id"`
	ImportedReceiver string `db:"imported_receiver"`
	PaymentDate      string `db:"payment_date"`
	Cash             bool   `db:"cash"`
	Notes            string `db:"notes"`
	Amount           int    `db:"amount"`
}

// OrderPaymentQuery data pembayaran
type OrderPaymentSelect struct {
	ID               int            `db:"id"`
	OrderID          int            `db:"order_id"`
	ReceiverID       int            `db:"receiver_id"`
	ReceiverName     string         `db:"receiver_name"`
	ImportedReceiver sql.NullString `db:"imported_receiver"`
	PaymentDate      sql.NullString `db:"payment_date"`
	Cash             bool           `db:"cash"`
	Notes            sql.NullString `db:"notes"`
	Amount           int            `db:"amount"`
}

// OrderPayment wrapper data pembayaran
type OrderPayment struct {
	ID               int
	OrderID          int
	ReceiverID       int
	ReceiverName     string
	ImportedReceiver string
	PaymentDate      string
	Cash             bool
	Notes            string
	Amount           int
}
