package api

import (
	"fmt"
	"log"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// OrderShowList mengambil data/list order/penjualan
func OrderShowList(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		session := sessions.Default(c)
		// User session saat ini
		// Tolak jika yang request bukan user terdaftar
		uid := session.Get("userid")
		if uid == nil {
			router.Page404(c)
			return
		}

		var lastid int
		last := false
		limit := 10
		var direction string
		httpStatus := http.StatusOK
		errMess := ""
		o := dbquery.GetOrders{}
		next := true

		// Mengambil parameter limit
		lim, err := strconv.Atoi(c.Param("limit"))
		if err == nil {
			limit = lim
			o.Limit(limit)
		} else {
			errMess = "Limit tidak valid"
			httpStatus = http.StatusBadRequest
			next = false
		}

		// Ambil query id terakhir
		lst, err := strconv.Atoi(c.Query("lastid"))
		if err == nil {
			lastid = lst
		}

		// Forward/backward
		direction = c.Query("direction")
		if direction == "back" {
			o.Direction(direction)
		} else {
			o.Direction("next")
		}

		var total int
		if t, err := dbquery.OrderGetOrderTotalRow(db); err == nil {
			total = t
		}

		var orders []wrapper.Order

		if next {
			o.LastID(lastid)
			or, err := o.Show(db)
			if err != nil {
				errMess = err.Error()
				httpStatus = http.StatusInternalServerError
			}
			orders = or
		}

		if len(orders) > 0 && direction == "back" {
			// Reverse urutan array order
			temp := make([]wrapper.Order, len(orders))
			in := 0
			for i := len(orders) - 1; i >= 0; i-- {
				temp[in] = orders[i]
				in++
			}
			orders = temp
		}

		// Cek id terakhir
		if len(orders) > 0 && len(orders) < limit {
			// Periksa apakah ini data terakhir atau bukan
			last = true
		}

		c.JSON(httpStatus, gin.H{
			"orders": orders,
			"error":  errMess,
			"total":  total,
			"last":   last,
		})
	}
	return gin.HandlerFunc(fn)
}

// OrderShowByID mengambil data order/penjualan berdasarkan ID
func OrderShowByID(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		session := sessions.Default(c)
		// User session saat ini
		// Tolak jika yang request bukan user terdaftar
		uid := session.Get("userid")
		if uid == nil {
			router.Page404(c)
			return
		}
		httpStatus := http.StatusOK
		errMess := ""

		// Mengambil parameter id order
		var oid int
		id, err := strconv.Atoi(c.Param("id"))
		if err == nil {
			oid = id
		} else {
			httpStatus = http.StatusBadRequest
			errMess = "Request tidak valid"
		}

		var order wrapper.Order
		if o, err := dbquery.OrderGetOrderByID(db, oid); err == nil {
			order = o
		} else {
			httpStatus = http.StatusInternalServerError
			errMess = "Sepertinya telah terjadi kesalahan saat memuat data"
		}

		var total int
		if t, err := dbquery.OrderGetOrderTotalRow(db); err == nil {
			total = t
		}

		c.JSON(httpStatus, gin.H{
			"order": order,
			"error": errMess,
			"total": total,
		})
	}
	return fn
}

// OrderCreate API untuk menambahkan produk baru
func OrderCreate(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		session := sessions.Default(c)
		// User session saat ini
		// Tolak jika yang request bukan user terdaftar
		uid := session.Get("userid")
		if uid == nil {
			router.Page404(c)
			return
		}

		var json wrapper.OrderForm

		status := "success"
		var httpStatus int
		message := ""
		var simpleErrMap = make(map[string]interface{})
		save := true

		if err := c.ShouldBindJSON(&json); err != nil {
			log.Println("ERROR: api.order.go OrderCreate() bind json")
			log.Println(err)
			if fmt.Sprintf("%T", err) == "validator.ValidationErrors" {
				simpleErrMap = validation.SimpleValErrMap(err)
			}
			httpStatus = http.StatusBadRequest
			status = "fail"
			save = false
		}

		if len(json.OrderItems) == 0 {
			simpleErrMap["orderitems"] = "Item tidak boleh kosong"
			httpStatus = http.StatusBadRequest
			status = "error"
			message = "Harap pilih setidaknya satu produk"
			save = false
		}

		var code string
		if save {
			uname, err := dbquery.UserGetUsername(db, json.Customer)
			if err == nil {
				uname = uname[3:]
				tm := time.Now()
				dt := strings.ReplaceAll(tm.Format("01-02-2006"), "-", "")
				dy := tm.Format("Mon")
				uq := tm.Format(".000")[1:]
				code = strings.ToUpper(fmt.Sprintf("%s%s-%s%s-%s%s", uname[4:], dy, uq, dt[4:], dt[:4], uname[:4]))
			} else {
				message = "Telah terjadi kesalahan saat memeriksa data konsumen"
				status = "error"
				save = false
				httpStatus = http.StatusInternalServerError
			}
		}

		if save {
			order := dbquery.NewOrder()
			err := order.SetCustomer(json.Customer).
				SetSales(json.Sales).
				SetCollector(json.Collector).
				SetSurveyor(json.Surveyor).
				SetShipping(json.ShippingAddress).
				SetBilling(json.BillingAddress).
				SetCredit(*json.Credit).
				SetNotes(json.Notes).
				SetCode(code).
				SetOrderDate(json.OrderDate).
				SetShippingDate(json.ShippingDate).
				SetOrderItems(json.OrderItems).
				SetCreatedBy(uid.(int)).
				Save(db)

			if err != nil {
				log.Println("ERROR: api.order.go OrderCreate() Gagal membuat order baru")
				log.Print(err)
				status = "error"
				message = "Gagal membuat order"
			} else {
				status = "success"
				message = "di save"
			}
		}

		m := gin.H{
			"message": message,
			"status":  status,
		}
		c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
	}
	return gin.HandlerFunc(fn)
}
