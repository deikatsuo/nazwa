package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PageDashboardCustomers tampilkan halaman customers
// Halaman dashboard data pelanggan
func PageDashboardCustomers(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		gh := gin.H{
			"site_title": "Halaman Pelanggan",
			"page":       "customers",
		}

		df := c.MustGet("dashboard").(map[string]interface{})
		c.HTML(200, "dashboard.customers.html", misc.Mete(df, gh))
	}
	return gin.HandlerFunc(fn)
}
