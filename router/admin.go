package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

func PageAdmin(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Admin",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "admin.html", misc.Mete(df, gh))
}
