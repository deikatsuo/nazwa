package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PageLogin ...
// Halaman login
func PageLogin(c *gin.Context) {
	// Larang akses untuk user login
	session := sessions.Default(c)
	userid := session.Get("userid")
	if userid != nil {
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

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "login.html", misc.Mete(df, gh))
}
