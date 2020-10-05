package middleware

import (
	"log"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/router"
	"nazwa/wrapper"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
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
		var user dbquery.User

		// Periksa apakah session nil
		// guna menghindari error saat konversi nil ke int
		if userid != nil {
			if userid.(int) > 0 {
				if u, err := dbquery.GetUserByID(db, userid.(int)); err != nil {
					log.Print("ERR-500")
					log.Print(err)
					router.Page500(c)
					c.Abort()
				} else {
					user = u
				}
			}
		}

		dashboard := map[string]interface{}{
			"user": user,

			"l_admin_create_customer": "Buat pelanggan",
		}

		siteDefault := c.MustGet("config").(wrapper.DefaultConfig).Site
		mergo.Map(&dashboard, siteDefault, mergo.WithOverride)
		c.Set("dashboard", dashboard)
		c.Next()
	}
}
