package middleware

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/router"
	"nazwa/wrapper"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// NewDefaultConfig - mengambil konfigurasi default dari .env
func NewDefaultConfig(db *sqlx.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		session := sessions.Default(c)
		// User session saat ini
		userid := session.Get("userid")

		login := false

		if userid != nil {
			login = true
		} else {
			login = false
		}

		cur := session.Get("user")
		var user wrapper.User

		if cur == nil {
			// Periksa apakah session nil
			// guna menghindari error saat konversi nil ke int
			if userid != nil {
				if userid.(int) > 0 {
					abort := false
					if u, err := dbquery.UserGetByID(db, userid.(int)); err != nil {
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
					session.Set("user", user)
					if err := session.Save(); err != nil {
						log.Println("Data user tidak tersimpan di session")
						log.Println(err)
					} else {
						log.Println("Data user tersimpan di session")
					}
				}
			}
		} else {
			user = cur.(wrapper.User)
		}

		config := map[string]interface{}{
			"site_url":   misc.GetEnv("SITE_URL", ""),
			"site_name":  misc.GetEnv("SITE_NAME", ""),
			"site_title": misc.GetEnv("SITE_TITLE", ""),
			"login":      login,
			"user":       user,
		}

		c.Set("config", wrapper.DefaultConfig{Info: config})

		c.Next()
	}
}

// NewDashboardDefaultConfig konfigurasi default halaman dashboard
func NewDashboardDefaultConfig() gin.HandlerFunc {
	return func(c *gin.Context) {

		dashboard := map[string]interface{}{
			"chart":           false,
			"l_modal_header":  "Modal",
			"l_modal_content": "Kosong",
			"l_modal_btn_one": "Batal",
			"l_modal_btn_two": "Oke",

			"l_admin_create": "Buat pelanggan",
		}

		siteDefault := c.MustGet("config").(wrapper.DefaultConfig).Info

		met := misc.Mete(dashboard, siteDefault)

		c.Set("dashboard", met)

		c.Next()
	}
}
