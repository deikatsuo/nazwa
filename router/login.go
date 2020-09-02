package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// Login ...
// Halaman login
func Login(c *gin.Context) {
	gh := gin.H{
		"site_title": "Masuk akun",

		"l_login_title":  "Masuk",
		"l_login_btn":    "Masuk",
		"l_login_forgot": "Lupa password?",
		"l_login_create": "Buat akun",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "login.html", misc.Mete(df, gh))
}
