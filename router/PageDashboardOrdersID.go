package router

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PageDashboardOrdersID halaman order berdasarkan ID
func PageDashboardOrdersID(c *gin.Context) {
	httpStatus := http.StatusOK
	message := ""
	status := ""
	next := true

	// Mengambil parameter id produk
	var oid int
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		oid = id
	} else {
		httpStatus = http.StatusBadRequest
		message = "Request tidak valid"
		next = false
	}

	var order wrapper.Order
	if next {
		if o, err := dbquery.OrderGetOrderByIDFull(oid); err == nil {
			order = o
		} else {
			httpStatus = http.StatusInternalServerError
			message = "Sepertinya telah terjadi kesalahan saat memuat data"
			status = "error"
			next = false
		}
	}

	if next {
		status = "success"
	}

	gh := gin.H{
		"site_title": fmt.Sprintf("Order %s", order.Code),
		"order":      order,
		"message":    message,
		"status":     status,
		"page":       "orders_id",
	}

	// Ambil konfigurasi default dashboard
	df := c.MustGet("dashboard").(map[string]interface{})

	c.HTML(httpStatus, "dashboard.orders.id.html", misc.Mete(df, gh))
}
