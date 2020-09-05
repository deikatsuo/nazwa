package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardCustomers ...
// Halaman dashboard pelanggan
func PageDashboardCustomers(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Pelanggan",
		"page":       "customers",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "dashboard.customers.html", misc.Mete(df, gh))
}
