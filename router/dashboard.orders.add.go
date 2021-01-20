package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardOrdersAdd halaman tambah produk
func PageDashboardOrdersAdd(c *gin.Context) {
	gh := gin.H{
		"site_title": "Buat Order",
		"page":       "orders_add",
		/*"css": []string{
			"/assets/css/pikaday.css",
		},*/
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.orders.add.html", misc.Mete(df, gh))
}
