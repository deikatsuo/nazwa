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
			if lastid != 0 {
				o.LastID(lastid)
				or, err := o.Show(db)
				if err != nil {
					errMess = err.Error()
					httpStatus = http.StatusInternalServerError
				}
				orders = or
			} else {
				or, err := o.Show(db)
				if err != nil {
					errMess = err.Error()
					httpStatus = http.StatusInternalServerError
				}
				orders = or
			}
		}

		// Cek id terakhir
		if len(orders) > 0 {
			lastid = orders[len(orders)-1].ID
			if len(orders) < limit {
				lastid = (orders[0].ID - 1) + limit
				// Periksa apakah ini data terakhir atau bukan
				last = true
			}
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
