package dbquery

import (
	"fmt"
	"log"
	"nazwa/wrapper"
	"strconv"

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
	if p.direction == "next" && p.lastid > 0 {
		where = "WHERE id > " + strconv.Itoa(p.lastid)
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
			Status:    p.Status,
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
