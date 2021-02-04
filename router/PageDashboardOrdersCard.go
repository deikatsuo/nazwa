package router

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PageDashboardOrdersCard kartu angsuran
func PageDashboardOrdersCard(c *gin.Context) {
	// id order
	oid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Page404(c)
		return
	}

	var order wrapper.Order
	if o, err := dbquery.OrderGetOrderByID(oid); err == nil {
		order = o
	}

	gh := gin.H{
		"site_title": "Kode " + order.Code,
		"page":       "orders_card",
		"order":      order,
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.orders.card.html", misc.Mete(df, gh))
}
