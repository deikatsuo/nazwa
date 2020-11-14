package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PageDashboardProductsAdd halaman tambah produk
func PageDashboardProductsAdd(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		gh := gin.H{
			"site_title": "Tambah Produk",
			"page":       "products_add",
		}

		df := c.MustGet("dashboard").(map[string]interface{})
		c.HTML(200, "dashboard.products.add.html", misc.Mete(df, gh))
	}
	return gin.HandlerFunc(fn)
}
