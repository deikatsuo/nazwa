package router

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// PageCheckCard kartu angsuran
func PageCheckCard(c *gin.Context) {
	// kode order
	code := c.Param("code")
	if code == "" {
		Page404(c)
		return
	}

	message := ""

	var order wrapper.Order
	if o, err := dbquery.OrderGetOrderByCode(code); err == nil {
		// Generate kode QR
		var png []byte
		png, err := qrcode.Encode(fmt.Sprintf("%s/check/card/%s", misc.GetEnv("SITE_URL", "").(string), o.Code), qrcode.Medium, 100)
		if err == nil {
			o.QR = png
		}
		order = o
	} else {
		Page404(c)
		return
	}

	gh := gin.H{
		"site_title": "Kartu Tagihan",
		"order":      order,
		"message":    message,
	}

	c.HTML(200, "check.card.html", gh)
}
