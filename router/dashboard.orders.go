package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PageDashboardOrders tampilkan halaman order
func PageDashboardOrders(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		gh := gin.H{
			"site_title": "Halaman Penjualan",
			"page":       "orders",
		}

		df := c.MustGet("dashboard").(map[string]interface{})
		c.HTML(200, "dashboard.orders.html", misc.Mete(df, gh))
	}
	return gin.HandlerFunc(fn)
}
