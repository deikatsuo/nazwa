package misc

import (
	"os"

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

// NewAdminDefaultConfig ...
// Ambil konfigurasi admin default
func NewAdminDefaultConfig() gin.HandlerFunc {
	config := map[string]interface{}{
		"l_admin_create_customer": "Buat pelanggan",
	}
	return func(c *gin.Context) {
		siteDefault := c.MustGet("config").(DefaultConfig).Site
		mergo.Map(&siteDefault, config, mergo.WithOverride)
		c.Set("config", DefaultConfig{Site: siteDefault})
		c.Next()
	}
}

// Cari konfigurasi dari .env
func getEnv(k string, df interface{}) interface{} {
	if v, e := os.LookupEnv(k); e {
		return v
	}
	return df
}
