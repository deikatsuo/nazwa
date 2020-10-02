package misc

import (
	"fmt"
	"nazwa/dbquery"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
	"github.com/jmoiron/sqlx"
)

// DefaultConfig - untuk menyimpan konfigurasi bawaan
type DefaultConfig struct {
	Site map[string]interface{}
}

// NewDefaultConfig - mengambil konfigurasi default dari .env
func NewDefaultConfig() gin.HandlerFunc {
	config := map[string]interface{}{
		"site_url":   getEnv("SITE_URL", ""),
		"site_name":  getEnv("SITE_NAME", ""),
		"site_title": getEnv("SITE_TITLE", ""),
	}

	return func(c *gin.Context) {
		c.Set("config", DefaultConfig{Site: config})
		c.Next()
	}
}

// NewDashboardDefaultConfig - konfigurasi default halaman
// dashboard
func NewDashboardDefaultConfig(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userid := session.Get("userid")
		var user dbquery.DashboardUser

		// Periksa apakah session nil
		// guna menghindari error saat konversi nil ke int
		if userid != nil {
			if userid.(int) > 0 {
				if u, err := dbquery.GetUser(db, userid.(int)); err != nil {
					fmt.Println(err)
				} else {
					user = u
				}
			}
		}

		dashboard := map[string]interface{}{
			"user": user,

			"l_admin_create_customer": "Buat pelanggan",
		}

		siteDefault := c.MustGet("config").(DefaultConfig).Site
		mergo.Map(&dashboard, siteDefault, mergo.WithOverride)
		c.Set("dashboard", dashboard)
		c.Next()
	}
}

// Cari konfigurasi dari .env, dengan nilai default
// Jika belum ditentukan atau kosong
func getEnv(k string, df interface{}) interface{} {
	if v, e := os.LookupEnv(k); e {
		return v
	}
	return df
}

// Ambil nilai dari konfigurasi dari .env tanpa default
func getEnvND(k string) string {
	if v, e := os.LookupEnv(k); e {
		return v
	}
	return ""
}

// Cek konfigurasi tidak kosong ""
func checkEnv(k string) bool {
	if v, e := os.LookupEnv(k); e {
		if v != "" {
			return true
		}
	}
	return false
}

// SetupDBType - mengambil value tipe database dari .env
func SetupDBType() string {
	return getEnv("DB_TYPE", "").(string)
}

// SetupDBSource - mengambil konfigurasi db dan mengubahnya menjadi
// source string
func SetupDBSource() string {
	var source = ""

	if checkEnv("DB_NAME") {
		source = "dbname=" + getEnvND("DB_NAME")
	}
	if checkEnv("DB_USER") {
		source = source + " user=" + getEnvND("DB_USER")
	}
	if checkEnv("DB_PASSWORD") {
		source = source + " password=" + getEnvND("DB_PASSWORD")
	}
	if checkEnv("DB_SSLMODE") {
		source = source + " sslmode=" + getEnvND("DB_SSLMODE")
	}

	return source
}

// SetupMigrationURL - mengambil URL migration
func SetupMigrationURL() string {
	var db string
	var user string
	var password string
	var host string
	var name string
	var ssl string

	if checkEnv("DB_TYPE") {
		db = getEnvND("DB_TYPE")
	}
	if checkEnv("DB_USER") {
		user = getEnvND("DB_USER")
	}
	if checkEnv("DB_PASSWORD") {
		password = getEnvND("DB_PASSWORD")
	}
	if checkEnv("DB_HOST") {
		host = getEnvND("DB_HOST")
	}
	if checkEnv("DB_NAME") {
		name = getEnvND("DB_NAME")
	}
	if checkEnv("DB_SSLMODE") {
		ssl = getEnvND("DB_SSLMODE")
	}
	url := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s", db, user, password, host, name, ssl)
	return url
}
