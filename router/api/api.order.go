package api

import (
	"nazwa/dbquery"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ShowOrderList mengambil data/list order/penjualan
func ShowOrderList(db *sqlx.DB) gin.HandlerFunc {
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
		if t, err := dbquery.GetOrderTotalRow(db); err == nil {
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
			// Reverse urutan array user
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
			"lastid": lastid,
			"total":  total,
			"last":   last,
		})
	}
	return gin.HandlerFunc(fn)
}

// ShowOrderByID mengambil data order/penjualan berdasarkan ID
func ShowOrderByID(db *sqlx.DB) gin.HandlerFunc {
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
		if o, err := dbquery.GetOrderByID(db, oid); err == nil {
			order = o
		} else {
			httpStatus = http.StatusInternalServerError
			errMess = "Sepertinya telah terjadi kesalahan saat memuat data"
		}

		var total int
		if t, err := dbquery.GetOrderTotalRow(db); err == nil {
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
