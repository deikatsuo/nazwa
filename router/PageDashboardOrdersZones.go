package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardOrdersZones tampilkan halaman zona wilayah
func PageDashboardOrdersZones(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Zona Wilayah",
		"page":       "orders_zones",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.orders.zones.html", misc.Mete(df, gh))
}
