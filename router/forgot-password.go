package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PageForgot ...
// Halaman pemulihan password
func PageForgot(c *gin.Context) {
	// Larang akses untuk user login
	session := sessions.Default(c)
	userid := session.Get("userid")
	if userid != nil {
		Page403(c)
		return
	}

	gh := gin.H{
		"site_title": "Lupa password",

		"l_forgot_title":       "Pulihkan password",
		"l_forgot_recover_btn": "Pulihkan",
		"l_forgot_login":       "Masuk",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "forgot-password.html", misc.Mete(df, gh))
}
