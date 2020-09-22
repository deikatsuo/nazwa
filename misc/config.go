package misc

import (
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
)

// DefaultConfig ...
// Struct untuk menyimpan konfigurasi bawaan
type DefaultConfig struct {
	Site map[string]interface{}
}

// NewDefaultConfig ...
// Atur konfigurasi bawaan
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

// NewDashboardDefaultConfig ...
// Ambil konfigurasi default untuk
// Halaman dashboard
func NewDashboardDefaultConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		email := session.Get("email")
		picture := session.Get("picture")

		config := map[string]interface{}{
			"email":   email,
			"picture": picture,

			"l_admin_create_customer": "Buat pelanggan",
		}

		siteDefault := c.MustGet("config").(DefaultConfig).Site
		mergo.Map(&siteDefault, config, mergo.WithOverride)
		c.Set("config", DefaultConfig{Site: siteDefault})
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

// Tipe database
func SetupDBType() string {
	return getEnv("DB_TYPE", "").(string)
}

// Setup db
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
