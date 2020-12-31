package middleware

import (
	"log"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/router"
	"nazwa/wrapper"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// NewDefaultConfig - mengambil konfigurasi default dari .env
func NewDefaultConfig() gin.HandlerFunc {
	config := map[string]interface{}{
		"site_url":   misc.GetEnv("SITE_URL", ""),
		"site_name":  misc.GetEnv("SITE_NAME", ""),
		"site_title": misc.GetEnv("SITE_TITLE", ""),
	}

	return func(c *gin.Context) {

		c.Set("config", wrapper.DefaultConfig{Site: config})

		c.Next()
	}
}

// NewDashboardDefaultConfig - konfigurasi default halaman
// dashboard
func NewDashboardDefaultConfig(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userid := session.Get("userid")
		var user wrapper.NullableUser

		// Periksa apakah session nil
		// guna menghindari error saat konversi nil ke int
		if userid != nil {
			if userid.(int) > 0 {
				abort := false
				if u, err := dbquery.UserGetNullableUserByID(db, userid.(int)); err != nil {
					log.Print("ERROR: default-config.go NewDashboardDefaultConfig() Gagal mengambil user by ID")
					log.Print(err)
					abort = true
				} else {
					user = u
				}

				if abort {
					router.Page500(c)
					c.Abort()
					return
				}

				user.ID = userid.(int)
				// Ambil data email
				email, err := dbquery.UserGetEmail(db, userid.(int))
				if err != nil {
					log.Print("User tidak memiliki email ", err)
				}
				user.Emails = email

				// Ambil nomor HP
				phone, err := dbquery.UserGetPhone(db, userid.(int))
				if err != nil {
					log.Print("User tidak memiliki nomor HP ", err)
				}
				user.Phones = phone

			}
		}

		dashboard := map[string]interface{}{
			"user":            user,
			"chart":           false,
			"l_modal_header":  "Modal",
			"l_modal_content": "Kosong",
			"l_modal_btn_one": "Batal",
			"l_modal_btn_two": "Oke",

			"l_admin_create": "Buat pelanggan",
		}

		siteDefault := c.MustGet("config").(wrapper.DefaultConfig).Site

		met := misc.Mete(dashboard, siteDefault)

		c.Set("dashboard", met)

		c.Next()
	}
}
