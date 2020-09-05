package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardBlank ...
// Contoh halaman dashboard kosong
func PageDashboardBlank(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Kosong Dashboad",
		"page":       "blank",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "dashboard.blank.html", misc.Mete(df, gh))
}
