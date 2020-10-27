package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PageDashboardProducts tampilkan halaman user
func PageDashboardProducts(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		gh := gin.H{
			"site_title": "Halaman Produk",
			"page":       "products",
		}

		df := c.MustGet("dashboard").(map[string]interface{})
		c.HTML(200, "dashboard.products.html", misc.Mete(df, gh))
	}
	return gin.HandlerFunc(fn)
}
