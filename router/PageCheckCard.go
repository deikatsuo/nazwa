package router

import (
	"nazwa/dbquery"
	"nazwa/wrapper"

	"github.com/gin-gonic/gin"
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
