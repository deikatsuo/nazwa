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
	message := ""
	status := ""

	// Mengambil parameter slug
	slug := c.Param("slug")
	if len(slug) < 4 || len(slug) > 100 {
		httpStatus = http.StatusBadRequest
		message = "Request tidak valid"
	}

	var product wrapper.Product
	if p, err := dbquery.ProductGetProductBySlug(slug); err == nil {
		product = p
	} else {
		httpStatus = http.StatusInternalServerError
		message = "Sepertinya telah terjadi kesalahan saat memuat data"
		status = "error"
	}

	// Ambil konfigurasi default
	df := c.MustGet("config").(wrapper.DefaultConfig).Info

	gh := gin.H{
		"title":   product.Name,
		"product": product,
		"message": message,
		"status":  status,
	}

	met := misc.Mete(df, gh)

	c.HTML(httpStatus, "product.detail.html", met)
}
