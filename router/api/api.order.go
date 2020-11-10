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
		pts := dbquery.GetProducts{}
		next := true

		// Mengambil parameter limit
		lim, err := strconv.Atoi(c.Param("limit"))
		if err == nil {
			limit = lim
			pts.Limit(limit)
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
			pts.Direction(direction)
		} else {
			pts.Direction("next")
		}

		var total int
		if t, err := dbquery.GetOrderTotalRow(db); err == nil {
			total = t
		}

		var products []wrapper.Order

		if next {
			if lastid != 0 {
				pts.LastID(lastid)
				p, err := pts.Show(db)
				if err != nil {
					errMess = err.Error()
					httpStatus = http.StatusInternalServerError
				}
				products = p
			} else {
				p, err := pts.Show(db)
				if err != nil {
					errMess = err.Error()
					httpStatus = http.StatusInternalServerError
				}
				products = p
			}
		}

		// Cek id terakhir
		if len(products) > 0 {
			lastid = products[len(products)-1].ID
			if len(products) < limit {
				lastid = (products[0].ID - 1) + limit
				// Periksa apakah ini data terakhir atau bukan
				last = true
			}
		}

		c.JSON(httpStatus, gin.H{
			"products": products,
			"error":    errMess,
			"lastid":   lastid,
			"total":    total,
			"last":     last,
		})
	}
	return gin.HandlerFunc(fn)
}