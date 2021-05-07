package router

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PageDashboardInstalmentsReceipt halaman kwitansi berdasarkan ID
func PageDashboardInstalmentsReceipt(c *gin.Context) {
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
		if ord, err := dbquery.OrderGetOrderByID(oid); err == nil {
			order = ord
		} else {
			httpStatus = http.StatusInternalServerError
			message = "Sepertinya telah terjadi kesalahan saat memuat data order"
			status = "error"
			next = false
		}
	}

	var monthly []wrapper.OrderMonthlyCredit
	if next {
		if mon, err := dbquery.OrderGetMonthlyCredit(oid); err == nil {
			monthly = mon
		} else {
			httpStatus = http.StatusInternalServerError
			message = "Sepertinya telah terjadi kesalahan saat memuat data kredit"
			status = "error"
			next = false
		}
	}

	/*
		for i, mon := range monthly {
			// Buat QR
			var png []byte
			png, err := qrcode.Encode(fmt.Sprintf("%s/check/receipt/%s", misc.GetEnv("SITE_URL", "").(string), mon.Code), qrcode.Medium, 100)

			// Tanggal sekarang
			current := time.Now()

			if err == nil {
				monthly[i].QR = png
				if mon.PrintDate == "" {
					monthly[i].PrintDate = current.Format("02-01-2006")
				}
			}
		}
	*/

	if next {
		status = "success"
	}

	gh := gin.H{
		"site_title": "Kwitansi",
		"monthly":    monthly,
		"order":      order,
		"code":       order.Code,
		"message":    message,
		"status":     status,
		"page":       "instalments_receipt",
		"date":       time.Now().Format("02-01-2006"),
		"css": []string{
			"/assets/css/print.css",
			"/assets/css/loading.css",
		},
		"js": []string{
			"/assets/js/terbilang.js",
			"/assets/js/jspm/zip.js",
			"/assets/js/jspm/zip-ext.js",
			"/assets/js/jspm/deflate.js",
			"/assets/js/jspm/JSPrintManager.js",
			"/assets/js/html2canvas.min.js",
			"/assets/js/xlsx-populate.min.js",
			"/assets/js/print.js",
		},
	}

	// Ambil konfigurasi default dashboard
	df := c.MustGet("dashboard").(map[string]interface{})

	c.HTML(httpStatus, "dashboard.instalments.receipt.html", misc.Mete(df, gh))
}
