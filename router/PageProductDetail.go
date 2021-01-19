package router

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PageProductDetail Halaman Detail Produk
func PageProductDetail(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		httpStatus := http.StatusOK
		message := ""
		status := ""

		// Mengambil parameter id produk
		var pid int
		id, err := strconv.Atoi(c.Param("id"))
		if err == nil {
			pid = id
		} else {
			httpStatus = http.StatusBadRequest
			message = "Request tidak valid"
		}

		var product wrapper.Product
		if p, err := dbquery.ProductGetProductByID(db, pid); err == nil {
			product = p
		} else {
			httpStatus = http.StatusInternalServerError
			message = "Sepertinya telah terjadi kesalahan saat memuat data"
			status = "error"
		}

		// Ambil konfigurasi default
		df := c.MustGet("config").(wrapper.DefaultConfig).Site

		gh := gin.H{
			"title":   "Halaman Produk",
			"product": product,
			"message": message,
			"status":  status,
		}

		met := misc.Mete(df, gh)

		c.HTML(httpStatus, "product.detail.html", met)
	}
	return fn
}
