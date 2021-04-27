package router

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PageProductDetail Halaman Detail Produk
func PageProductDetail(c *gin.Context) {
	httpStatus := http.StatusOK

	// Mengambil parameter slug
	slug := c.Param("slug")

	var product wrapper.Product
	if p, err := dbquery.ProductGetProductBySlug(slug); err == nil {
		product = p
	} else {
		Page404(c)
		return
	}

	// Ambil konfigurasi default
	df := c.MustGet("config").(wrapper.DefaultConfig).Info

	gh := gin.H{
		"title":   product.Name,
		"product": product,
	}

	met := misc.Mete(df, gh)

	c.HTML(httpStatus, "product.detail.html", met)
}
