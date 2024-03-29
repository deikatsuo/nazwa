package router

import (
	"nazwa/misc"
	"nazwa/wrapper"

	"github.com/gin-gonic/gin"
)

// PageLogin Halaman login
func PageLogin(c *gin.Context) {
	df := c.MustGet("config").(wrapper.DefaultConfig).Info

	// Larang akses untuk user login
	if df["login"] != false {
		Page403(c)
		return
	}

	gh := gin.H{
		"site_title": "Masuk akun",

		"l_login_title":    "Masuk",
		"l_login_loginid":  "ID login",
		"l_login_password": "Kata sandi",
		"l_login_btn":      "Masuk",
		"l_login_forgot":   "Lupa password?",
		"l_login_create":   "Buat akun",
	}

	c.HTML(200, "login.html", misc.Mete(df, gh))
}
