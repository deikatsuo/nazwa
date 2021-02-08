package router

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PageDashboardOrdersCard kartu angsuran
func PageDashboardOrdersCard(c *gin.Context) {
	// id order
	//oid, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	Page404(c)
	//	return
	//}

	qid := c.QueryArray("ids")
	var failsParse []string
	var failsFetch []string

	var orders []wrapper.Order
	for _, soid := range qid {
		if oid, err := strconv.Atoi(soid); err == nil {
			if o, err := dbquery.OrderGetOrderByID(oid); err == nil {
				orders = append(orders, o)
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
