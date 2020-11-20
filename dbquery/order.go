package dbquery

import (
	"fmt"
	"log"
	"nazwa/wrapper"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// GetOrders mengambil list order/penjualan
type GetOrders struct {
	limit     int
	lastid    int
	direction string
}

// Limit set limit
func (p *GetOrders) Limit(limit int) *GetOrders {
	p.limit = limit
	return p
}

// LastID set lastid
func (p *GetOrders) LastID(lastid int) *GetOrders {
	p.lastid = lastid
	return p
}

// Direction untuk backward/forward
// @direction "back","next"
func (p *GetOrders) Direction(direction string) *GetOrders {
	p.direction = direction
	return p
}

// Show tampilkan data
func (p *GetOrders) Show(db *sqlx.DB) ([]wrapper.Order, error) {
	var order []wrapper.NullableOrder
	var parse []wrapper.Order
	limit := 10
	if p.limit > 0 {
		limit = p.limit
	}

	where := ""
	if p.direction == "next" {
		where = "WHERE id > " + strconv.Itoa(p.lastid) + "ORDER BY id ASC"
	} else if p.direction == "back" {
		where = "WHERE id < " + strconv.Itoa(p.lastid) + " ORDER BY id DESC"
	}

	query := `SELECT
		id,
		code,
		status,
		credit,
		TO_CHAR(order_date, 'MM/DD/YYYY HH12:MI:SS AM') AS order_date
		FROM "order"
		%s
		LIMIT $1`

	query = fmt.Sprintf(query, where)

	err := db.Select(&order, query, limit)
	if err != nil {
		log.Println("Error: order.go Select all product")
		log.Println(err)
		return []wrapper.Order{}, err
	}

	for _, p := range order {
		parse = append(parse, wrapper.Order{
			ID:        p.ID,
			OrderDate: p.OrderDate,
			Credit:    p.Credit,
			Code:      p.Code,
			Status:    strings.Title(p.Status),
		})
	}

	return parse, nil
}

// GetOrderTotalRow menghitung jumlah row pada tabel user
func GetOrderTotalRow(db *sqlx.DB) (int, error) {
	var total int
	query := `SELECT COUNT(id) FROM "order"`
	err := db.Get(&total, query)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetOrderByID mengambil data order berdasarkan ID order
func GetOrderByID(db *sqlx.DB, oid int) (wrapper.Order, error) {
	var order wrapper.Order
	var o wrapper.NullableOrder
	query := `SELECT
		o.id,
		o.customer_id,
		c.first_name || ' ' || c.last_name as customer_name,
		o.sales_id,
		sa.first_name || ' ' || sa.last_name as sales_name,
		o.surveyor_id,
		su.first_name || ' ' || su.last_name as surveyor_name,
		o.collector_id,
		co.first_name || ' ' || co.last_name as collector_name,
		o.shipping_address_id,
		sad.one || ' ' || sad.two as shipping_address_name,
		o.billing_address_id,
		bad.one || ' ' || bad.two as billing_address_name,
		o.status,
		o.credit,
		o.notes,
		TO_CHAR(o.order_date, 'MM/DD/YYYY HH12:MI:SS AM') AS order_date,
		TO_CHAR(o.shipping_date, 'MM/DD/YYYY HH12:MI:SS AM') AS shipping_date,
		o.first_time,
		o.code
		FROM "order" o
		LEFT JOIN "user" c ON c.id=o.customer_id
		LEFT JOIN "user" sa ON sa.id=o.sales_id
		LEFT JOIN "user" su ON su.id=o.surveyor_id
		LEFT JOIN "user" co ON co.id=o.collector_id
		LEFT JOIN "address" sad ON sad.id=o.shipping_address_id
		LEFT JOIN "address" bad ON bad.id=o.billing_address_id
		WHERE o.id=$1
		LIMIT 1`

	err := db.Get(&o, query, oid)
	if err != nil {
		log.Println("order.go Select order berdasarkan ID")
		log.Println(err)
		return wrapper.Order{}, err
	}

	// Mengambil list item dari transaksi
	var items []wrapper.OrderItem
	if oi, err := GetOrderItem(db, o.ID); err == nil {
		items = oi
	}

	order = wrapper.Order{
		ID: o.ID,
		Customer: wrapper.NameID{
			ID:   o.CustomerID,
			Name: o.CustomerName,
		},
		Sales: wrapper.NameID{
			ID:   int(o.SalesID.Int64),
			Name: o.SalesName.String,
		},
		Surveyor: wrapper.NameID{
			ID:   int(o.SurveyorID.Int64),
			Name: o.SurveyorName.String,
		},
		Collector: wrapper.NameID{
			ID:   int(o.CollectorID.Int64),
			Name: o.SurveyorName.String,
		},
		ShippingAddress: wrapper.NameID{
			ID:   o.ShippingAddressID,
			Name: o.ShippingAddressName,
		},
		BillingAddress: wrapper.NameID{
			ID:   int(o.BillingAddressID.Int64),
			Name: string(o.BillingAddressName.String),
		},
		Status:       strings.Title(o.Status),
		Code:         o.Code,
		Credit:       o.Credit,
		FirstTime:    o.FirstTime,
		Notes:        o.Notes.String,
		OrderDate:    o.OrderDate,
		ShippingDate: string(o.ShippingDate.String),
		Items:        items,
	}

	return order, nil
}

// GetOrderItem mengambil data barang transaksi berdasarkan id order
func GetOrderItem(db *sqlx.DB, oid int) ([]wrapper.OrderItem, error) {
	var items []wrapper.NullableOrderItem
	var parse []wrapper.OrderItem
	query := `SELECT oi.id, oi.product_id, oi.quantity, oi.notes, p.name, p.code
	FROM "order_item" oi
	LEFT JOIN "product" p ON p.id=oi.product_id
	WHERE order_id=$1`
	err := db.Select(&items, query, oid)
	if err != nil {
		return []wrapper.OrderItem{}, err
	}

	for _, i := range items {
		parse = append(parse, wrapper.OrderItem{
			ID: i.ID,
			Product: wrapper.NameIDCode{
				ID:   i.ProductID,
				Name: i.ProductName,
				Code: i.ProductCode,
			},
			Quantity: i.Quantity,
			Notes:    string(i.Notes.String),
		})
	}
	return parse, err
}
