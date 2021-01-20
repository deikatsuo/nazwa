package router

import (
	"nazwa/misc"
	"nazwa/wrapper"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PageCreateAccount ...
// Halaman buat akun
func PageCreateAccount(c *gin.Context) {
	// Larang akses untuk user login
	session := sessions.Default(c)
	userid := session.Get("userid")
	if userid != nil {
		Page403(c)
		return
	}

	gh := gin.H{
		"site_title": "Membuat akun baru",

		// Formulir
		"l_c_form_title":        "Buat akun",
		"l_c_form_phone":        "Nomor Hp",
		"l_c_form_firstname":    "Nama depan",
		"l_c_form_lastname":     "Nama belakang",
		"l_c_form_gender":       "Jenis kelamin",
		"l_c_form_gender_m":     "Laki-laki",
		"l_c_form_gender_f":     "Perempuan",
		"l_c_form_password":     "Kata sandi",
		"l_c_form_repassword":   "Ulangi kata sandi",
		"l_c_form_agree":        "Saya setuju dengan",
		"l_c_form_privacy_link": "kebijakan privasi",
		"l_c_form_create":       "Buat akun",

		// Link sudah punya akun
		"l_c_have": "Sudah punya akun? Masuk",
	}

	df := c.MustGet("config").(wrapper.DefaultConfig).Info
	c.HTML(200, "create-account.html", misc.Mete(df, gh))
}
