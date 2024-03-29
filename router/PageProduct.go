package router

import (
	"nazwa/misc"
	"nazwa/wrapper"

	"github.com/gin-gonic/gin"
)

// PageProduct Halaman Produk
func PageProduct(c *gin.Context) {

	// Ambil konfigurasi default
	df := c.MustGet("config").(wrapper.DefaultConfig).Info

	gh := gin.H{
		"title": "Daftar Produk",
	}

	met := misc.Mete(df, gh)

	c.HTML(200, "product.html", met)
}
