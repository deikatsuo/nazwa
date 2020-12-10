package middleware

import (
	"log"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/router"
	"nazwa/wrapper"
	"sync"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var mut = sync.RWMutex{}

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
				if u, err := dbquery.GetNullableUserByID(db, userid.(int)); err != nil {
					log.Print("ERR-500")
					log.Print(err)
					router.Page500(c)
					c.Abort()
				} else {
					user = u
					user.ID = userid.(int)
					// Ambil data email
					email, err := dbquery.GetEmail(db, userid.(int))
					if err != nil {
						log.Print("User tidak memiliki email ", err)
					}
					user.Emails = email

					// Ambil nomor HP
					phone, err := dbquery.GetPhone(db, userid.(int))
					if err != nil {
						log.Print("User tidak memiliki nomor HP ", err)
					}
					user.Phones = phone
				}
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
