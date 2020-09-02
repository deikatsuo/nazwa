package misc

import (
	"os"

	"github.com/gin-gonic/gin"
)

// Struct untuk menyimpan konfigurasi bawaan
type DefaultConfig struct {
	Site map[string]interface{}
}

// Atur konfigurasi bawaan
func NewDefaultConfig() gin.HandlerFunc {
	config := map[string]interface{}{
		"site_url":   getEnv("SITE_URL", ""),
		"site_name":  getEnv("SITE_NAME", ""),
		"site_title": getEnv("SITE_TITLE", ""),
	}

	return func(c *gin.Context) {
		c.Set("config", DefaultConfig{config})
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
