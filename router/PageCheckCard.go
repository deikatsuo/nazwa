package router

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"strings"

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
	df := c.MustGet("config").(wrapper.DefaultConfig).Info

	var order wrapper.Order
	if o, err := dbquery.OrderGetOrderByCode(code); err == nil {
		// Generate kode QR
		var png []byte
		png, err := qrcode.Encode(fmt.Sprintf("%s/check/card/%s", df["site_url"].(string), o.Code), qrcode.Medium, 100)
		if err == nil {
			o.QR = png
		}
		order = o
	} else {
		Page404(c)
		return
	}

	gh := gin.H{
		"site_title": "Kartu Tagihan " + strings.Title(order.Customer.Name),
		"order":      order,
		"message":    message,
	}

	c.HTML(200, "check.card.html", misc.Mete(df, gh))
}
