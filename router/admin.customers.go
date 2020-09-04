package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageAdminCustomers ...
// Halaman admin pelanggan
func PageAdminCustomers(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Pelanggan Admin",
		"page":       "customers",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "admin.customers.html", misc.Mete(df, gh))
}
