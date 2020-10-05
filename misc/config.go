package misc

import (
	"fmt"
	"os"
)

// Cari konfigurasi dari .env, dengan nilai default
// Jika belum ditentukan atau kosong
func getEnv(k string, df interface{}) interface{} {
	if v, e := os.LookupEnv(k); e {
		return v
	}
	return df
}

// GetEnv buat ngewrap getEnv
func GetEnv(k string, df interface{}) interface{} {
	return getEnv(k, df)
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
