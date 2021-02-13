package router

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// PageDashboardOrdersReceipt halaman kwitansi berdasarkan ID
func PageDashboardOrdersReceipt(c *gin.Context) {
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

	// Kode order
	var code string
	if next {
		if co, err := dbquery.OrderGetCodeByID(oid); err == nil {
			code = co
		} else {
			next = false
			message = "Sepertinya telah terjadi masalah saat mencoba mengambil kode order"
			httpStatus = http.StatusInternalServerError
		}
	}

	var monthly []wrapper.OrderMonthlyCredit
	if next {
		if mon, err := dbquery.OrderGetMonthlyCredit(oid); err == nil {
			monthly = mon
		} else {
			httpStatus = http.StatusInternalServerError
			message = "Sepertinya telah terjadi kesalahan saat memuat data"
			status = "error"
			next = false
		}
	}

	for i, mon := range monthly {
		// Buat QR
		var png []byte
		png, err := qrcode.Encode(fmt.Sprintf("%s/check/receipt/%s", misc.GetEnv("SITE_URL", "").(string), mon.Code), qrcode.Medium, 100)
		if err == nil {
			monthly[i].QR = png
		}
	}

	if next {
		status = "success"
	}

	gh := gin.H{
		"site_title": "Kwitansi",
		"monthly":    monthly,
		"code":       code,
		"message":    message,
		"status":     status,
		"page":       "orders_receipt",
		"css": []string{
			"/assets/css/print.css",
		},
	}

	// Ambil konfigurasi default dashboard
	df := c.MustGet("dashboard").(map[string]interface{})

	c.HTML(httpStatus, "dashboard.orders.receipt.html", misc.Mete(df, gh))
}
