package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// Page404 halaman tidak ditemukan
// Jika route gak ada yang cocok
func Page404(c *gin.Context) {
	gh := gin.H{
		"title":         "Halaman tidak ditemukan",
		"type":          404,
		"picture":       "404.png",
		"l_reason":      "Oops! halaman tidak ditemukan.",
		"l_description": "Mohon maaf, halaman yang anda tuju tidak tersedia. Mungkin salah mengetikan alamat, atau halaman tersebut sudah dipindahkan/dihapus",
	}
	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(404, "error.html", misc.Mete(df, gh))
}

// Page403 forbidden page
// tidak memiliki akses untuk membuka halaman
func Page403(c *gin.Context) {
	gh := gin.H{
		"title":         "Akses dilarang",
		"type":          403,
		"picture":       "403.png",
		"l_reason":      "Oops! Akses dilarang.",
		"l_description": "Mohon maaf, anda tidak memiliki ijin untuk mengakses halaman ini",
	}
	df := c.MustGet("config").(misc.DefaultConfig).Site
	c.HTML(403, "error.html", misc.Mete(df, gh))
}
