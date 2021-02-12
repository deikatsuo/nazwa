package router

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// PageDashboardOrdersCard kartu angsuran
func PageDashboardOrdersCard(c *gin.Context) {
	qid := c.QueryArray("ids")
	if len(qid) == 0 {
		Page404(c)
		return
	}

	var failsParse []string
	var failsFetch []string

	var orders []wrapper.Order
	for _, soid := range qid {
		if oid, err := strconv.Atoi(soid); err == nil {
			if o, err := dbquery.OrderGetOrderByID(oid); err == nil {
				if !o.Credit {
					failsFetch = append(failsFetch, fmt.Sprintf("ID %d bukan kredit", oid))
				} else {
					// Buat QR
					var png []byte
					png, err := qrcode.Encode(fmt.Sprintf("%s/check/card/%s", misc.GetEnv("SITE_URL", "").(string), o.Code), qrcode.Medium, 100)
					if err == nil {
						o.QR = png
					}
					orders = append(orders, o)
				}
			} else {
				failsFetch = append(failsFetch, fmt.Sprintf("ID %d %s", oid, err.Error()))
			}
		} else {
			failsParse = append(failsParse, fmt.Sprintf("ID %s %s", soid, err.Error()))
		}
	}

	gh := gin.H{
		"site_title": "Kartu Tagihan",
		"page":       "orders_card",
		"orders":     orders,
		"fails": map[string]interface{}{
			"parse": failsParse,
			"fetch": failsFetch,
		},
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.orders.card.html", misc.Mete(df, gh))
}
