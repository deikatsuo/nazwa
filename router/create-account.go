package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageCreateAccount ...
// Halaman buat akun
func PageCreateAccount(c *gin.Context) {
	gh := gin.H{
		"site_title": "Membuat akun baru",

		// Formulir
		"l_c_form_title":        "Buat akun",
		"l_c_form_password":     "Kata sandi",
		"l_c_form_repassword":   "Ulangi kata sandi",
		"l_c_form_agree":        "Saya setuju dengan",
		"l_c_form_privacy_link": "kebijakan privasi",
		"l_c_form_create":       "Buat akun",

		// Link sudah punya akun
		"l_c_have": "Sudah punya akun? Masuk",
	}

	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(200, "create-account.html", misc.Mete(df, gh))
}
