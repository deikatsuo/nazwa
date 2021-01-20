package router

import (
	"nazwa/misc"
	"nazwa/wrapper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Page404 halaman tidak ditemukan
// Jika route gak ada yang cocok
func Page404(c *gin.Context) {
	gh := gin.H{
		"title":         "Halaman tidak ditemukan",
		"type":          http.StatusNotFound,
		"picture":       "404.png",
		"l_reason":      "Oops! halaman tidak ditemukan.",
		"l_description": "Mohon maaf, halaman yang anda tuju tidak tersedia. Mungkin salah mengetikan alamat, atau halaman tersebut sudah dipindahkan/dihapus",
	}

	df := c.MustGet("config").(wrapper.DefaultConfig).Info

	met := misc.Mete(df, gh)

	c.HTML(http.StatusNotFound, "error.html", met)

}

// Page403 halaman terlarang
// tidak memiliki akses untuk membuka halaman
func Page403(c *gin.Context) {
	gh := gin.H{
		"title":         "Akses dilarang",
		"type":          http.StatusForbidden,
		"picture":       "403.png",
		"l_reason":      "Oops! Akses dilarang.",
		"l_description": "Mohon maaf, anda tidak memiliki ijin untuk mengakses halaman ini",
	}

	df := c.MustGet("config").(wrapper.DefaultConfig).Info

	met := misc.Mete(df, gh)

	c.HTML(http.StatusForbidden, "error.html", met)
}

// Page500 halaman internal server error
// jika ada kesalahan pada server
func Page500(c *gin.Context) {
	gh := gin.H{
		"title":         "Terjadi kesalahan pada server",
		"type":          http.StatusInternalServerError,
		"picture":       "500.png",
		"l_reason":      "Oops! Error pada server",
		"l_description": "Kami mohon maaf, sepertinya telah terjadi kesalahan pada server kami. Mungkin data ",
	}

	df := c.MustGet("config").(wrapper.DefaultConfig).Info

	met := misc.Mete(df, gh)

	c.HTML(http.StatusInternalServerError, "error.html", met)
}
