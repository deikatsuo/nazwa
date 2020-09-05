package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboard ...
// Halaman dashboard utama
func PageDashboard(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Dashboard",
		"page":       "dashboard",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "dashboard.html", misc.Mete(df, gh))
}
