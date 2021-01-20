package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardProducts tampilkan halaman produk
func PageDashboardProducts(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Produk",
		"page":       "products",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.products.html", misc.Mete(df, gh))
}
