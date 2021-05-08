package dbquery

import (
	"errors"
	"fmt"
	"nazwa/misc"
	"nazwa/wrapper"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
)

// CreateOrder struct untuk menyimpan data order yang akan di insert
type CreateOrder struct {
	wrapper.OrderInsert
	into            map[string]string
	returnID        bool
	returnIDTO      *int
	orderItems      []wrapper.OrderItemForm
	substitutes     []wrapper.OrderUserSubstituteForm
	due             int
	line            int
	lineMax         int
	duration        int
	importMode      bool
	importedMonthly int
}

// NewOrder membuat order baru
func NewOrder() *CreateOrder {
	return &CreateOrder{
		into:       make(map[string]string),
		importMode: false,
	}
}

// SetCustomer ID kostumer
func (c *CreateOrder) SetCustomer(o int) *CreateOrder {
	if o > 0 {
		c.CustomerID = o
		c.into["customer_id"] = ":customer_id"
	}
	return c
}

// SetSales ID sales
func (c *CreateOrder) SetSales(o int) *CreateOrder {
	if o > 0 {
		c.SalesID = o
		c.into["sales_id"] = ":sales_id"
	}
	return c
}

// SetSurveyor ID surveyor
func (c *CreateOrder) SetSurveyor(o int) *CreateOrder {
	if o > 0 {
		c.SurveyorID = o
		c.into["surveyor_id"] = ":surveyor_id"
	}
	return c
}

// SetDriver ID supir
func (c *CreateOrder) SetDriver(o int) *CreateOrder {
	if o > 0 {
		c.DriverID = o
		c.into["driver_id"] = ":driver_id"
	}
	return c
}

// SetShipping ID alamat pengiriman
func (c *CreateOrder) SetShipping(o int) *CreateOrder {
	if o > 0 {
		c.ShippingAddressID = o
		c.into["shipping_address_id"] = ":shipping_address_id"
	}
	return c
}

// SetBilling ID alamat penagihan
func (c *CreateOrder) SetBilling(o int) *CreateOrder {
	if o > 0 {
		c.BillingAddressID = o
		c.into["billing_address_id"] = ":billing_address_id"
	}
	return c
}

// SetCode Tentukan kode transaksi secara manual
func (c *CreateOrder) SetCode(o string) *CreateOrder {
	if o != "" {
		c.Code = o
		c.into["code"] = ":code"
	}
	return c
}

// SetStatus Status order
// default 'pending'
func (c *CreateOrder) SetStatus(o string) *CreateOrder {
	if o != "" {
		c.Status = o
		c.into["status"] = ":status"
	}
	return c
}

// SetCredit true | false
func (c *CreateOrder) SetCredit(o bool) *CreateOrder {
	c.Credit = o
	c.into["credit"] = ":credit"
	return c
}

// SetDeposit uang muka
func (c *CreateOrder) SetDeposit(o int) *CreateOrder {
	if o > 0 {
		c.Deposit = o
		c.into["deposit"] = ":deposit"
	}
	return c
}

// SetDuration Durasi kredit barang
func (c *CreateOrder) SetDuration(o int) *CreateOrder {
	if o > 0 {
		c.duration = o
	}
	return c
}

// SetDue tenggang/tanggal waktu pembayaran
func (c *CreateOrder) SetDue(o int) *CreateOrder {
	if o > 0 {
		c.due = o
	}
	return c
}

// SetLine tentukan arah tagih
func (c *CreateOrder) SetLine(o int) *CreateOrder {
	if o > 0 {
		c.line = o
	}
	return c
}

// SetLineCodeMaxNumber tentukan max number manual
func (c *CreateOrder) SetLineCodeMaxNumber(o int) *CreateOrder {
	if o > 0 {
		c.lineMax = o
	}
	return c
}

// SetNotes Catatan order
func (c *CreateOrder) SetNotes(o string) *CreateOrder {
	c.Notes = o
	c.into["notes"] = ":notes"
	return c
}

// SetOrderDate Tanggal order
func (c *CreateOrder) SetOrderDate(o string) *CreateOrder {
	if o != "" {
		c.OrderDate = o
		c.into["order_date"] = ":order_date"
	}
	return c
}

// SetShippingDate Tanggal pengiriman barang yang diorder
func (c *CreateOrder) SetShippingDate(o string) *CreateOrder {
	if o != "" {
		c.ShippingDate = o
		c.into["shipping_date"] = ":shipping_date"
	}
	return c
}

// SetCreatedBy ID admin yang membuat order
func (c *CreateOrder) SetCreatedBy(o int) *CreateOrder {
	if o > 0 {
		c.CreatedBy = o
		c.into["created_by"] = ":created_by"
	}
	return c
}

// SetOrderItems items
func (c *CreateOrder) SetOrderItems(o []wrapper.OrderItemForm) *CreateOrder {
	if len(o) > 0 {
		c.orderItems = o
	}
	return c
}

// SetSubstitutes menambahkan pengganti user/konsumen
func (c *CreateOrder) SetSubstitutes(o []wrapper.OrderUserSubstituteForm) *CreateOrder {
	if len(o) > 0 {
		c.substitutes = o
	}
	return c
}

// SetImportMode mode import
func (c *CreateOrder) SetImportMode(o bool) *CreateOrder {
	c.importMode = o
	return c
}

// SetImportedItems import sales
func (c *CreateOrder) SetImportedItems(o string) *CreateOrder {
	if o != "" {
		c.ImportedItems = o
		c.into["imported_items"] = ":imported_items"
	}
	return c
}

// SetImportedAddress import alamat
func (c *CreateOrder) SetImportedAddress(o string) *CreateOrder {
	if o != "" {
		c.ImportedAddress = o
		c.into["imported_address"] = ":imported_address"
	}
	return c
}

// SetImportedSurveyor import survey
func (c *CreateOrder) SetImportedSurveyor(o string) *CreateOrder {
	if o != "" {
		c.ImportedSurveyor = o
		c.into["imported_surveyor"] = ":imported_surveyor"
	}
	return c
}

// SetImportedSales import sales
func (c *CreateOrder) SetImportedSales(o string) *CreateOrder {
	if o != "" {
		c.ImportedSales = o
		c.into["imported_sales"] = ":imported_sales"
	}
	return c
}

// SetImportedMonthly import harga bulanan
func (c *CreateOrder) SetImportedMonthly(monthly int) *CreateOrder {
	c.importedMonthly = monthly
	return c
}

// ReturnID Mengembalikan ID produk terakhir
func (c *CreateOrder) ReturnID(id *int) *CreateOrder {
	c.returnID = true
	c.returnIDTO = id
	return c
}

// Insert query berdasarka data yang diisi
func (c CreateOrder) generateInsertQuery() string {
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
func (c *CreateOrder) Save() error {
	db := DB
	// Jika bukan mode import
	if !c.importMode {
		// Jika tidak ada barang yang di order
		if len(c.orderItems) == 0 {
			return errors.New("dbuery.order.go (CreateOrder) Save() Item order kosong")
		}
	}

	// Jika tanggal pengiriman kosong
	// Maka tanggal pengiriman disamakan dengan tanggal pemesanan
	if c.ShippingDate == "" {
		c.ShippingDate = c.OrderDate
		c.into["shipping_date"] = ":shipping_date"
	}

	if !c.importMode {
		// Jika alamat penagihan kosong
		// Maka alamat penagihan = alamat pengiriman
		if c.BillingAddressID == 0 {
			c.BillingAddressID = c.ShippingAddressID
			c.into["billing_address_id"] = ":billing_address_id"
		}
	}

	// Total keseluruhan tagihan
	var priceTotal int
	// Total keseluruhan harga awal barang (harga beli) sebelum profit
	var basePriceTotal int

	var prices []int
	var basePrices []int

	if !c.importMode {
		// Periksa apakah pembelian kredit atau cash
		// Lalu kalkulasikan
		for _, item := range c.orderItems {
			if c.Credit {
				// Temporary credit price
				var tmpcp int
				p, err := ProductGetProductCreditPrice(item.ProductID)
				if err == nil {
					for _, ps := range p {
						if ps.Duration == c.duration {
							tmpcp = ps.Price
						}
					}
				}
				if tmpcp <= 0 {
					log.Warn("dbquery.order.go (CreateOrder) Save() Harga kredit tidak ada")
					return fmt.Errorf("tidak ditemukan harga kredit untuk durasi %d bulan", c.duration)
				}

				if item.Discount > 0 {
					priceTotal += ((item.Discount * item.Quantity) * c.duration)
				} else {
					priceTotal += ((tmpcp * item.Quantity) * c.duration)
				}

				prices = append(prices, tmpcp)
			} else {
				p, err := ProductGetProductPrice(item.ProductID)
				if err != nil {
					log.Warn("dbquery.order.go (CreateOrder) Save() Get item price")
					return err
				}

				// Jika menggunakan harga diskon
				if item.Discount > 0 {
					priceTotal += item.Discount * item.Quantity
				} else {
					priceTotal += p * item.Quantity
				}
				prices = append(prices, p)
			}

			bp, err := ProductGetProductBasePrice(item.ProductID)
			if err != nil {
				log.Warn("dbquery.order.go (CreateOrder) Save() Get item base price")
				return err
			}
			basePriceTotal += bp * item.Quantity
			basePrices = append(basePrices, bp)
		}
	} else {
		priceTotal = c.importedMonthly * c.duration
	}

	remaining := priceTotal
	total := priceTotal
	var monthly int
	var luckyDiscount int

	if c.Credit {
		monthly = priceTotal / c.duration

		if c.Deposit > 0 {
			//remaining = priceTotal - c.Deposit
			priceTotal += c.Deposit
			//total = priceTotal - c.Deposit
			//monthly = (priceTotal - c.Deposit) / c.duration
			//monthly = int(math.Floor(float64(monthly)/1000)) * 1000
			//luckyDiscount = (priceTotal - c.Deposit) - (monthly * c.duration)
		}

		if c.Deposit > priceTotal {
			return errors.New("deposit lebih besar dari total harga barang, lebih baik beli cash")
		}
	}

	if priceTotal != 0 {
		c.PriceTotal = priceTotal
		c.into["price_total"] = ":price_total"
	}

	if basePriceTotal != 0 {
		c.BasePriceTotal = basePriceTotal
		c.into["base_price_total"] = ":base_price_total"
	}

	// Mulai transaksi
	tx := db.MustBegin()
	var tempReturnID int
	orderInsertQuery := fmt.Sprintf(`INSERT INTO "order" %s`, c.generateInsertQuery())
	if rows, err := tx.NamedQuery(orderInsertQuery, c); err == nil {
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
		log.Warn("dbquery.order.go (c *CreateOrder) Save(db *sqlx.DB) Gagal insert order")
		tx.Rollback()
		return err
	}

	if !c.importMode {
		// Item yang akan di insert
		var itemInsert []wrapper.OrderItemInsert

		itemInsertQuery := `INSERT INTO "order_item" (order_id, product_id, quantity, notes, base_price, price, discount) VALUES (:order_id, :product_id, :quantity, :notes, :base_price, :price, :discount)`
		for n, i := range c.orderItems {
			if stock, err := ProductCheckStock(i.ProductID); err == nil {
				remainingStock := stock - i.Quantity
				if err := ProductUpdateStock(i.ProductID, remainingStock); err != nil {
					log.Warn("dbquery.order.go Save() mengubah stok")
					return err
				}
			}
			itemInsert = append(itemInsert, wrapper.OrderItemInsert{
				OrderID:   tempReturnID,
				ProductID: i.ProductID,
				Quantity:  i.Quantity,
				Notes:     i.Notes,
				Price:     prices[n],
				BasePrice: basePrices[n],
				Discount:  i.Discount,
			})
		}

		if rows, err := tx.NamedQuery(itemInsertQuery, itemInsert); err == nil {
			if err := rows.Close(); err != nil {
				log.Warn("dbquery.order.go Save() Insert item closing row")
				return err
			}
		} else {
			log.Warn("dbquery.order.go (c *CreateOrder) Save(db *sqlx.DB) Gagal insert item ")
			tx.Rollback()
			return err
		}
	}

	// Simpan data substitutes
	if len(c.substitutes) > 0 {
		type uos struct {
			RIC          string `db:"ric"`
			Firstname    string `db:"first_name"`
			Lastname     string `db:"last_name"`
			Gender       string `db:"gender"`
			SubstituteTo int    `db:"substitute_to"`
			CreatedBy    int    `db:"created_by"`
		}
		for _, s := range c.substitutes {

			// Simpan data pendamping
			into := map[string]string{}
			uosd := uos{}
			if s.RIC != "" {
				into["ric"] = ":ric"
				uosd.RIC = s.RIC
			}
			if s.Firstname != "" {
				into["first_name"] = ":first_name"
				uosd.Firstname = s.Firstname
			}
			if s.Lastname != "" {
				into["last_name"] = ":last_name"
				uosd.Lastname = s.Lastname
			}
			if s.Gender != "" {
				into["gender"] = ":gender"
				uosd.Gender = s.Gender
			}
			if tempReturnID > 0 {
				into["substitute_to"] = ":substitute_to"
				uosd.SubstituteTo = tempReturnID
			}
			if c.CreatedBy > 0 {
				into["created_by"] = ":created_by"
				uosd.CreatedBy = c.CreatedBy
			}
			gev := fmt.Sprintf(`INSERT INTO "order_user_substitute" %s`, misc.GenerateSimpleInsertValues(into))

			if rows, err := tx.NamedQuery(gev, uosd); err == nil {
				if err := rows.Close(); err != nil {
					log.Warn("dbquery.order.go Save() Insert substitute closing row")
					return err
				}
			} else {
				tx.Rollback()
				return err
			}

		}
	}

	// Simpan credit detail
	if c.Credit {
		var lc string
		var lmax int

		if c, err := LineGetCodeByID(c.line); err == nil {
			lc = c
		} else {
			return err
		}

		if mx, err := LineGetCodeMaxNumber(lc); err == nil {
			lmax = mx
		} else {
			log.Warn("Error sepertinya belum ada data, skip error")
		}
		if c.lineMax > 0 {
			lmax = c.lineMax
		} else {
			// tambah satu
			lmax++
		}

		if _, err := tx.Exec(`INSERT INTO "zone_line_list" (zone_line_code, number) VALUES ($1, $2)`, lc, lmax); err != nil {
			log.Warn("dbquery.order.go Save() insert nomor arah total")
			return err
		}

		creditCode := fmt.Sprintf("%s%04d", strings.ToUpper(lc), lmax)

		if _, err := tx.Exec(`INSERT INTO "order_credit_detail" (order_id, zone_line_id, credit_code, monthly, duration, due, total, remaining, lucky_discount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, tempReturnID, c.line, creditCode, monthly, c.duration, c.due, total, remaining, luckyDiscount); err != nil {
			log.Warn("dbquery.order.go Save() Insert product detail")
			return err
		}
	}

	// Buat kwitansi nth
	if c.Credit {
		var monthlyCredit []wrapper.OrderMonthlyCreditQuery
		var tm time.Time
		if t, err := time.Parse(time.RFC3339, c.ShippingDate+"T00:00:00.0000Z"); err == nil {
			tm = t
		}
		tm2 := time.Date(tm.Year(), tm.Month(), c.due, 0, 0, 0, 0, tm.Location())
		tm2 = tm2.AddDate(0, 1, 0)
		if c.Deposit > 0 {
			monthlyCredit = append(monthlyCredit, wrapper.OrderMonthlyCreditQuery{
				OrderID: tempReturnID,
				Code:    xid.New().String(),
				Nth:     0,
				DueDate: c.ShippingDate,
			})
		}
		var dueDate string
		for i := 0; i < c.duration; i++ {

			if i == 0 {
				if c.Deposit == 0 {
					dueDate = c.ShippingDate
				} else {
					dueDate = fmt.Sprintf("%d-%d-%d", tm2.Year(), tm2.Month(), tm2.Day())
				}
			} else {
				if c.Deposit > 0 {
					tm2 = tm2.AddDate(0, 1, 0)
				} else {
					if i > 1 {
						tm2 = tm2.AddDate(0, 1, 0)
					}
				}
				dueDate = fmt.Sprintf("%d-%d-%d", tm2.Year(), tm2.Month(), tm2.Day())
			}

			monthlyCredit = append(monthlyCredit, wrapper.OrderMonthlyCreditQuery{
				OrderID: tempReturnID,
				Code:    xid.New().String(),
				Nth:     i + 1,
				DueDate: dueDate,
			})
		}

		monthlyQuery := `INSERT INTO "order_monthly_credit" (order_id, code, nth, due_date) VALUES (:order_id, :code, :nth, :due_date)`

		if rows, err := tx.NamedQuery(monthlyQuery, monthlyCredit); err == nil {
			if err := rows.Close(); err != nil {
				log.Warn("dbquery.order.go Save() Monthly credit closing row")
				return err
			}
		} else {
			log.Warn("dbquery.order.go (c *CreateOrder) Save(db *sqlx.DB) Gagal membuat kwitansi")
			tx.Rollback()
			return err
		}
	}

	// Komit
	err := tx.Commit()
	return err
}

///////////
/// GET ///
///////////

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
func (p *GetOrders) Show() ([]wrapper.Order, error) {
	db := DB
	var order []wrapper.NullableOrder
	var parse []wrapper.Order
	limit := 10
	if p.limit > 0 {
		limit = p.limit
	}

	where := ""
	if p.direction == "next" {
		where = "WHERE o.id > " + strconv.Itoa(p.lastid) + "ORDER BY o.id ASC"
	} else if p.direction == "back" {
		where = "WHERE o.id < " + strconv.Itoa(p.lastid) + " ORDER BY o.id DESC"
	}

	query := `SELECT
		o.id,
		o.code,
		o.status,
		o.credit,
		TO_CHAR(o.order_date, 'DD-MM-YYYY HH12:MI:SS AM') AS order_date,
		TO_CHAR(o.shipping_date, 'DD-MM-YYYY HH12:MI:SS AM') AS shipping_date,
		o.customer_id,
		concat_ws(' ', c.first_name, c.last_name) as customer_name,
		c.avatar as customer_thumb,
		o.imported_items
		FROM "order" o
		LEFT JOIN "user" c ON c.id=customer_id
		%s
		LIMIT $1`

	query = fmt.Sprintf(query, where)

	err := db.Select(&order, query, limit)
	if err != nil {
		log.Warn("Error: order.go Select all product")
		return []wrapper.Order{}, err
	}

	for _, p := range order {
		// Mengambil list item dari transaksi
		var items []wrapper.OrderItem
		if oi, err := OrderGetOrderItem(p.ID); err == nil {
			items = oi
		}

		parse = append(parse, wrapper.Order{
			ID:           p.ID,
			OrderDate:    p.OrderDate,
			ShippingDate: p.ShippingDate,
			Credit:       p.Credit,
			Code:         p.Code,
			Status:       strings.Title(p.Status),
			Customer: wrapper.NameIDCode{
				ID:        p.CustomerID,
				Name:      p.CustomerName,
				Thumbnail: p.CustomerThumb,
			},
			Items:         items,
			ImportedItems: p.ImportedItems.String,
		})
	}

	return parse, nil
}

// OrderGetOrderTotalRow menghitung jumlah row pada tabel user
func OrderGetOrderTotalRow() (int, error) {
	db := DB
	var total int
	query := `SELECT COUNT(id) FROM "order"`
	err := db.Get(&total, query)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// OrderGetOrderByID mengambil data order berdasarkan ID order
func OrderGetOrderByID(oid int) (wrapper.Order, error) {
	db := DB
	var order wrapper.Order
	var o wrapper.NullableOrder
	query := `SELECT
		o.id,
		o.customer_id,
		c.username as customer_code,
		concat_ws(' ', c.first_name, c.last_name) as customer_name,
		c.avatar as customer_thumb,
		o.sales_id,
		concat_ws(' ', sa.first_name, sa.last_name) as sales_name,
		sa.avatar as sales_thumb,
		o.surveyor_id,
		concat_ws(' ', su.first_name, su.last_name) as surveyor_name,
		su.avatar as surveyor_thumb,
		concat_ws(' ', co.first_name, co.last_name) as collector_name,
		co.avatar as collector_thumb,
		o.created_by as created_by_id,
		concat_ws(' ', cb.first_name, cb.last_name) as created_by_name,
		cb.avatar as created_by_thumb,
		o.shipping_address_id,
		o.billing_address_id,
		o.status,
		o.credit,
		o.notes,
		TO_CHAR(o.order_date, 'DD-MM-YYYY') AS order_date,
		TO_CHAR(o.shipping_date, 'DD-MM-YYYY') AS shipping_date,
		TO_CHAR(o.created_at, 'DD-MM-YYYY HH12:MI:SS AM') AS created_at,
		o.code,
		o.deposit,
		o.price_total,
		o.base_price_total,
		o.imported_items,
		o.imported_address,
		o.imported_sales,
		o.imported_surveyor
		FROM "order" o
		LEFT JOIN "order_credit_detail" ocd ON order_id=o.id
		LEFT JOIN "zone_list" zlt ON zlt.zone_line_id=ocd.zone_line_id
		LEFT JOIN "zone" z ON z.id=zlt.zone_id
		LEFT JOIN "user" c ON c.id=o.customer_id
		LEFT JOIN "user" sa ON sa.id=o.sales_id
		LEFT JOIN "user" su ON su.id=o.surveyor_id
		LEFT JOIN "user" co ON co.id=z.collector_id
		LEFT JOIN "user" cb ON cb.id=o.created_by
		LEFT JOIN "address" sad ON sad.id=o.shipping_address_id
		LEFT JOIN "address" bad ON bad.id=o.billing_address_id
		
		WHERE o.id=$1
		LIMIT 1`

	err := db.Get(&o, query, oid)
	if err != nil {
		log.Warn("dbquery.order.go OrderGetOrderByID() Select order berdasarkan ID")
		log.Error(err)
		return wrapper.Order{}, err
	}

	// Mengambil list item dari transaksi
	var items []wrapper.OrderItem
	if oi, err := OrderGetOrderItem(o.ID); err == nil {
		items = oi
	}

	// Detail kredit
	var creditDetail wrapper.OrderCreditDetail
	if cd, err := OrderGetCreditInfo(o.ID); err == nil {
		creditDetail = cd
	} else {
		log.Warn("dbquery.order.go OrderGetOrderByID() select credit detail")
		log.Error(err)
	}

	var shipping string
	var billing string

	if int(o.ShippingAddressID.Int64) > 0 {
		if ship, err := AddressGetByID(int(o.ShippingAddressID.Int64)); err == nil {
			shipping = ship.String()
		}
	}

	if int(o.BillingAddressID.Int64) > 0 {
		if bill, err := AddressGetByID(int(o.BillingAddressID.Int64)); err == nil {
			billing = bill.String()
		}
	}

	order = wrapper.Order{
		ID: o.ID,
		Customer: wrapper.NameIDCode{
			ID:        o.CustomerID,
			Code:      o.CustomerCode.String,
			Name:      o.CustomerName,
			Thumbnail: o.CustomerThumb,
		},
		Sales: wrapper.NameID{
			ID:        int(o.SalesID.Int64),
			Name:      o.SalesName.String,
			Thumbnail: o.SalesThumb.String,
		},
		Surveyor: wrapper.NameID{
			ID:        int(o.SurveyorID.Int64),
			Name:      o.SurveyorName.String,
			Thumbnail: o.SurveyorThumb.String,
		},
		Collector: wrapper.NameID{
			ID:        int(o.CollectorID.Int64),
			Name:      o.CollectorName.String,
			Thumbnail: o.CollectorThumb.String,
		},
		CreatedBy: wrapper.NameID{
			ID:        int(o.CreatedByID.Int64),
			Name:      o.CreatedByName.String,
			Thumbnail: o.CreatedByThumb.String,
		},
		ShippingAddress:  shipping,
		BillingAddress:   billing,
		Status:           strings.Title(o.Status),
		Code:             o.Code,
		Credit:           o.Credit,
		Notes:            o.Notes.String,
		OrderDate:        o.OrderDate,
		ShippingDate:     o.ShippingDate,
		CreatedAt:        o.CreatedAt,
		Deposit:          o.Deposit,
		PriceTotal:       o.PriceTotal,
		BasePriceTotal:   int(o.BasePriceTotal.Int64),
		Items:            items,
		CreditDetail:     creditDetail,
		ImportedItems:    o.ImportedItems.String,
		ImportedAddress:  o.ImportedAddress.String,
		ImportedSales:    o.ImportedSales.String,
		ImportedSurveyor: o.ImportedSurveyor.String,
	}

	return order, nil
}

// OrderGetSimpleOrderByID mengambil data order sederhana berdasarkan ID order
func OrderGetSimpleOrderByID(oid int) (wrapper.OrderSimple, error) {
	db := DB
	var order wrapper.OrderSimple
	var o wrapper.NullableOrder
	query := `SELECT
		o.id,
		o.customer_id,
		c.username as customer_code,
		concat_ws(' ', c.first_name, c.last_name) as customer_name,
		c.avatar as customer_thumb,
		o.sales_id,
		concat_ws(' ', sa.first_name, sa.last_name) as sales_name,
		sa.avatar as sales_thumb,
		o.surveyor_id,
		concat_ws(' ', su.first_name, su.last_name) as surveyor_name,
		su.avatar as surveyor_thumb,
		concat_ws(' ', co.first_name, co.last_name) as collector_name,
		co.avatar as collector_thumb,
		o.billing_address_id,
		o.credit,
		TO_CHAR(o.shipping_date, 'DD-MM-YYYY') AS shipping_date,
		o.code,
		o.deposit,
		o.price_total,
		o.base_price_total
		FROM "order" o
		LEFT JOIN "order_credit_detail" ocd ON order_id=o.id
		LEFT JOIN "zone_list" zlt ON zlt.zone_line_id=ocd.zone_line_id
		LEFT JOIN "zone" z ON z.id=zlt.zone_id
		LEFT JOIN "user" c ON c.id=o.customer_id
		LEFT JOIN "user" sa ON sa.id=o.sales_id
		LEFT JOIN "user" su ON su.id=o.surveyor_id
		LEFT JOIN "user" co ON co.id=z.collector_id
		LEFT JOIN "user" cb ON cb.id=o.created_by
		LEFT JOIN "address" sad ON sad.id=o.shipping_address_id
		LEFT JOIN "address" bad ON bad.id=o.billing_address_id
		
		WHERE o.id=$1
		LIMIT 1`

	err := db.Get(&o, query, oid)
	if err != nil {
		log.Warn("dbquery.order.go OrderGetSimpleOrderByID() Select order berdasarkan ID")
		log.Error(err)
		return wrapper.OrderSimple{}, err
	}

	// Mengambil list item dari transaksi
	var items []wrapper.OrderItem
	if oi, err := OrderGetOrderItem(o.ID); err == nil {
		items = oi
	}

	// Detail kredit
	var creditDetail wrapper.OrderCreditDetail
	if cd, err := OrderGetCreditInfo(o.ID); err == nil {
		creditDetail = cd
	} else {
		log.Warn("dbquery.order.go OrderGetSimpleOrderByID() select credit detail")
		log.Error(err)
	}

	var billing string

	if bill, err := AddressGetByID(int(o.BillingAddressID.Int64)); err == nil {
		billing = bill.String()
	}

	order = wrapper.OrderSimple{
		ID: o.ID,
		Customer: wrapper.NameIDCode{
			ID:        o.CustomerID,
			Code:      o.CustomerCode.String,
			Name:      o.CustomerName,
			Thumbnail: o.CustomerThumb,
		},
		Sales: wrapper.NameID{
			ID:        int(o.SalesID.Int64),
			Name:      o.SalesName.String,
			Thumbnail: o.SalesThumb.String,
		},
		Surveyor: wrapper.NameID{
			ID:        int(o.SurveyorID.Int64),
			Name:      o.SurveyorName.String,
			Thumbnail: o.SurveyorThumb.String,
		},
		Collector: wrapper.NameID{
			ID:        int(o.CollectorID.Int64),
			Name:      o.CollectorName.String,
			Thumbnail: o.CollectorThumb.String,
		},
		BillingAddress: billing,
		Code:           o.Code,
		Credit:         o.Credit,
		ShippingDate:   o.ShippingDate,
		Deposit:        o.Deposit,
		PriceTotal:     o.PriceTotal,
		BasePriceTotal: int(o.BasePriceTotal.Int64),
		Items:          items,
		CreditDetail:   creditDetail,
	}

	return order, nil
}

// OrderGetOrderByCode ambil order berdasarkan kode
func OrderGetOrderByCode(code string) (wrapper.Order, error) {
	db := DB
	var order wrapper.Order
	var o wrapper.NullableOrder
	query := `SELECT
		o.id,
		o.customer_id,
		c.username as customer_code,
		concat_ws(' ', c.first_name, c.last_name) as customer_name,
		c.avatar as customer_thumb,
		o.sales_id,
		concat_ws(' ', sa.first_name, sa.last_name) as sales_name,
		sa.avatar as sales_thumb,
		o.surveyor_id,
		concat_ws(' ', su.first_name, su.last_name) as surveyor_name,
		su.avatar as surveyor_thumb,
		o.collector_id,
		concat_ws(' ', co.first_name, co.last_name) as collector_name,
		co.avatar as collector_thumb,
		o.created_by as created_by_id,
		concat_ws(' ', cb.first_name, cb.last_name) as created_by_name,
		cb.avatar as created_by_thumb,
		o.shipping_address_id,
		o.billing_address_id,
		o.status,
		o.credit,
		o.notes,
		TO_CHAR(o.order_date, 'DD-MM-YYYY HH12:MI:SS AM') AS order_date,
		TO_CHAR(o.shipping_date, 'DD-MM-YYYY HH12:MI:SS AM') AS shipping_date,
		TO_CHAR(o.created_at, 'DD-MM-YYYY HH12:MI:SS AM') AS created_at,
		o.code,
		o.deposit,
		o.price_total,
		o.base_price_total
		FROM "order" o
		LEFT JOIN "user" c ON c.id=o.customer_id
		LEFT JOIN "user" sa ON sa.id=o.sales_id
		LEFT JOIN "user" su ON su.id=o.surveyor_id
		LEFT JOIN "user" co ON co.id=o.collector_id
		LEFT JOIN "user" cb ON cb.id=o.created_by
		LEFT JOIN "address" sad ON sad.id=o.shipping_address_id
		LEFT JOIN "address" bad ON bad.id=o.billing_address_id
		WHERE o.code=$1
		LIMIT 1`

	err := db.Get(&o, query, code)
	if err != nil {
		log.Warn("dbquery.order.go OrderGetOrderByCode() Select order berdasarkan code")
		log.Error(err)
		return wrapper.Order{}, err
	}

	// Mengambil list item dari transaksi
	var items []wrapper.OrderItem
	if oi, err := OrderGetOrderItem(o.ID); err == nil {
		items = oi
	}

	// Detail kredit
	var creditDetail wrapper.OrderCreditDetail
	if cd, err := OrderGetCreditInfo(o.ID); err == nil {
		creditDetail = cd
	} else {
		log.Warn("dbquery.order.go OrderGetOrderByID() select credit detail")
		log.Error(err)
	}

	var shipping string
	var billing string

	if ship, err := AddressGetByID(int(o.ShippingAddressID.Int64)); err == nil {
		shipping = ship.String()
	}

	if bill, err := AddressGetByID(int(o.BillingAddressID.Int64)); err == nil {
		billing = bill.String()
	}

	order = wrapper.Order{
		ID: o.ID,
		Customer: wrapper.NameIDCode{
			ID:        o.CustomerID,
			Code:      o.CustomerCode.String,
			Name:      o.CustomerName,
			Thumbnail: o.CustomerThumb,
		},
		Sales: wrapper.NameID{
			ID:        int(o.SalesID.Int64),
			Name:      o.SalesName.String,
			Thumbnail: o.SalesThumb.String,
		},
		Surveyor: wrapper.NameID{
			ID:        int(o.SurveyorID.Int64),
			Name:      o.SurveyorName.String,
			Thumbnail: o.SurveyorThumb.String,
		},
		Collector: wrapper.NameID{
			ID:        int(o.CollectorID.Int64),
			Name:      o.CollectorName.String,
			Thumbnail: o.CollectorThumb.String,
		},
		CreatedBy: wrapper.NameID{
			ID:        int(o.CreatedByID.Int64),
			Name:      o.CreatedByName.String,
			Thumbnail: o.CreatedByThumb.String,
		},
		ShippingAddress: shipping,
		BillingAddress:  billing,
		Status:          strings.Title(o.Status),
		Code:            o.Code,
		Credit:          o.Credit,
		Notes:           o.Notes.String,
		OrderDate:       o.OrderDate,
		ShippingDate:    o.ShippingDate,
		CreatedAt:       o.CreatedAt,
		Deposit:         o.Deposit,
		PriceTotal:      o.PriceTotal,
		BasePriceTotal:  int(o.BasePriceTotal.Int64),
		Items:           items,
		CreditDetail:    creditDetail,
	}

	return order, nil
}

// OrderGetOrderItem mengambil data barang transaksi berdasarkan id order
func OrderGetOrderItem(oid int) ([]wrapper.OrderItem, error) {
	db := DB
	var items []wrapper.NullableOrderItem
	var parse []wrapper.OrderItem
	query := `SELECT oi.id, oi.product_id, oi.quantity, oi.notes, oi.discount, oi.base_price, oi.price, p.name, p.code, p.thumbnail
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
				ID:        i.ProductID,
				Name:      i.ProductName,
				Code:      i.ProductCode,
				Thumbnail: i.Thumbnail,
			},
			Quantity:  i.Quantity,
			Notes:     string(i.Notes.String),
			Price:     i.Price,
			BasePrice: i.BasePrice,
			Discount:  i.Discount,
		})
	}
	return parse, err
}

// OrderGetSubstituteByRic ambil data substitute berdasarkan NIK
func OrderGetSubstituteByRic(ric string) ([]wrapper.NameID, error) {
	db := DB
	var substitute []wrapper.NameID
	query := `SELECT
	id,
	concat_ws(' ', first_name, last_name) as name
	FROM "order_user_substitute"
	WHERE ric=$1`

	err := db.Select(&substitute, query, ric)
	if err != nil {
		return []wrapper.NameID{}, err
	}

	keys := make(map[string]bool)
	new := []wrapper.NameID{}

	// Hilangkan data dengan nama yang sama
	for _, entry := range substitute {
		if _, value := keys[entry.Name]; !value {
			keys[entry.Name] = true
			new = append(new, entry)
		}
	}

	return new, nil
}

// OrderGetCreditInfo mengambil informasi kredit
func OrderGetCreditInfo(oid int) (wrapper.OrderCreditDetail, error) {
	db := DB
	var cd wrapper.OrderCreditDetailSelect
	query := `SELECT ocd.id, ocd.zone_line_id, ocd.credit_code, ocd.monthly, ocd.duration, ocd.due, ocd.total, ocd.remaining, ocd.lucky_discount, TO_CHAR(ocd.last_paid, 'DD-MM-YYYY') AS last_paid, zl.name as zone_line_name, zl.code as zone_line_code
	FROM "order_credit_detail" ocd
	LEFT JOIN "zone_line" zl ON zl.id=ocd.zone_line_id
	WHERE ocd.order_id=$1`

	err := db.Get(&cd, query, oid)
	if err != nil {
		return wrapper.OrderCreditDetail{}, err
	}

	return wrapper.OrderCreditDetail{
		ID:      cd.ID,
		OrderID: oid,
		ZoneLine: wrapper.NameIDCode{
			ID:   int(cd.ZoneLineID.Int32),
			Code: cd.ZoneLineCode.String,
			Name: cd.ZoneLineName.String,
		},
		CreditCode:    cd.CreditCode,
		Monthly:       cd.Monthly,
		Duration:      cd.Duration,
		Due:           cd.Due,
		Total:         cd.Total,
		Remaining:     cd.Remaining,
		LuckyDiscount: cd.LuckyDiscount,
		LastPaid:      cd.LastPaid.String,
	}, nil
}

// OrderDeleteByID delete order
func OrderDeleteByID(oid int) error {
	db := DB
	query := `DELETE FROM "order"
	WHERE id=$1`
	_, err := db.Exec(query, oid)

	return err
}

// OrderGetMonthlyCredit ambil data kredit bulanan
func OrderGetMonthlyCredit(oid int) ([]wrapper.OrderMonthlyCredit, error) {
	db := DB
	var monthly []wrapper.OrderMonthlyCredit
	var monthlyQ []wrapper.OrderMonthlyCreditQuery

	query := `SELECT *, TO_CHAR(due_date, 'DD-MM-YYYY') AS due_date
	FROM "order_monthly_credit"
	WHERE order_id=$1 ORDER BY nth`

	err := db.Select(&monthlyQ, query, oid)
	if err != nil {
		return []wrapper.OrderMonthlyCredit{}, err
	}

	for _, mon := range monthlyQ {

		monthly = append(monthly, wrapper.OrderMonthlyCredit{
			ID:        mon.ID,
			OrderID:   mon.OrderID,
			Code:      mon.Code,
			Nth:       mon.Nth,
			DueDate:   mon.DueDate,
			PrintDate: mon.PrintDate.String,
			Promise:   mon.Promise.String,
			Paid:      mon.Paid,
			Notes:     mon.Notes.String,
			Position:  mon.Position,
			Printed:   mon.Printed,
			Done:      mon.Done,
		})
	}

	return monthly, nil
}

// OrderGetMonthlyCreditByDate ambil data kredit bulanan
func OrderGetMonthlyCreditByDate(zid int, date string) ([]wrapper.OrderMonthlyCredit, error) {
	db := DB
	var monthly []wrapper.OrderMonthlyCredit
	var monthlyQ []wrapper.OrderMonthlyCreditQuery

	var tm time.Time
	if t, err := time.Parse(time.RFC3339, date+"T00:00:00.0000Z"); err == nil {
		tm = t
	}

	query := `SELECT omc.*, TO_CHAR(omc.due_date, 'DD-MM-YYYY') AS due_date
	FROM "order_monthly_credit" omc
	LEFT JOIN "order_credit_detail" ocd ON ocd.order_id=omc.order_id
	LEFT JOIN "zone_list" zl ON zl.zone_line_id=ocd.zone_line_id
	WHERE TO_CHAR(omc.due_date, 'YYYY-MM-DD')<=$1 AND omc.done=false AND ocd.due=$2 AND zl.zone_id=$3 ORDER BY omc.nth`

	err := db.Select(&monthlyQ, query, date, tm.Day(), zid)
	if err != nil {
		log.Warn("dbquery.order.go OrderGetMonthlyCreditByDate() select")
		log.Error(err)
		return []wrapper.OrderMonthlyCredit{}, err
	}

	for _, mon := range monthlyQ {

		monthly = append(monthly, wrapper.OrderMonthlyCredit{
			ID:        mon.ID,
			OrderID:   mon.OrderID,
			Code:      mon.Code,
			Nth:       mon.Nth,
			DueDate:   mon.DueDate,
			PrintDate: mon.PrintDate.String,
			Promise:   mon.Promise.String,
			Paid:      mon.Paid,
			Notes:     mon.Notes.String,
			Position:  mon.Position,
			Printed:   mon.Printed,
			Done:      mon.Done,
		})
	}

	return monthly, nil
}

// OrderGetCodeByID ambil kode berdasarkan ID
func OrderGetCodeByID(oid int) (string, error) {
	db := DB
	var code string
	query := `SELECT code FROM "order" WHERE id=$1`
	err := db.Get(&code, query, oid)
	if err != nil {
		return code, err
	}

	return code, nil
}
