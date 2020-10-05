package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardCustomers tampilkan halaman customers
// Halaman dashboard data pelanggan
func PageDashboardCustomers(c *gin.Context) {

	gh := gin.H{
		"site_title": "Halaman Pelanggan",
		"page":       "customers",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.customers.html", misc.Mete(df, gh))
}
