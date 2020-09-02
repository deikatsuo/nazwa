package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// NoRoute ...
// Jika route gak ada yang cocok
// Maka jalankan ini
// Menggantikan default "halaman tidak ditemukan"
func NoRoute(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman tidak ditemukan",
		"l_back":     "Kembali",
		"l_missing":  "Oops, halaman tidak ditemukan.",
		"l_maybe":    "Mungkin salah mengetikan alamat, atau halaman tersebut sudah dipindahkan.",
	}
	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(404, "404.html", misc.Mete(df, gh))
}
