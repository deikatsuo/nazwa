package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// Forgot ...
// Halaman pemulihan password
func Forgot(c *gin.Context) {
	gh := gin.H{
		"site_title": "Lupa password",

		"l_forgot_title":       "Pulihkan password",
		"l_forgot_recover_btn": "Pulihkan",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "forgot-password.html", misc.Mete(df, gh))
}
