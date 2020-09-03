package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageAdminBlank ...
// Contoh halaman admin kosong
func PageAdminBlank(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Kosong Admin",
		"page":       "blank",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "admin.blank.html", misc.Mete(df, gh))
}
