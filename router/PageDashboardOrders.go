package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardOrders tampilkan halaman order
func PageDashboardOrders(c *gin.Context) {
	gh := gin.H{
		"site_title": "Penjualan",
		"page":       "orders",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.orders.html", misc.Mete(df, gh))
}
