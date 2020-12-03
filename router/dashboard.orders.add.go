package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PageDashboardOrdersAdd halaman tambah produk
func PageDashboardOrdersAdd(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		gh := gin.H{
			"site_title": "Buat Order",
			"page":       "orders_add",
		}

		df := c.MustGet("dashboard").(map[string]interface{})
		c.HTML(200, "dashboard.orders.add.html", misc.Mete(df, gh))
	}
	return gin.HandlerFunc(fn)
}
