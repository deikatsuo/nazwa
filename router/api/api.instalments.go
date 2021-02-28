package api

import (
	"nazwa/dbquery"
	"nazwa/misc/validation"
	"nazwa/wrapper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// InstalmentShowByDate list tagihan berdasarkan tanggal
func InstalmentShowByDate(c *gin.Context) {
	httpStatus := http.StatusOK
	message := ""
	status := ""
	next := true

	var simpleErr map[string]interface{}

	var byDate wrapper.InstalmentShowByDate
	if err := c.ShouldBindUri(&byDate); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
		message = simpleErr["date"].(string)
		status = "error"
		httpStatus = http.StatusBadRequest
	}

	var monthly []wrapper.OrderMonthlyCredit
	if next {
		if mon, err := dbquery.OrderGetMonthlyCreditByDate(byDate.Date); err == nil {
			monthly = mon
		} else {
			httpStatus = http.StatusInternalServerError
			message = "Sepertinya telah terjadi kesalahan saat memuat data kredit"
			status = "error"
			next = false
		}
	}

	var orders []wrapper.InstalmentOrderReceipt
	if next {
		for i, mon := range monthly {
			// Tanggal sekarang
			current := time.Now()

			if mon.PrintDate == "" {
				monthly[i].PrintDate = current.Format("02/01/2006")
			}

			if len(orders) > 1 {
				for oi, ord := range orders {
					if ord.OrderID == mon.OrderID {
						orders[oi].Monthly = append(orders[oi].Monthly, mon)
					}
				}
			} else {
				orders = append(orders, wrapper.InstalmentOrderReceipt{
					OrderID: mon.OrderID,
					Monthly: []wrapper.OrderMonthlyCredit{mon},
				})
			}
		}
	}

	if next {
		status = "success"
	}

	gh := gin.H{
		"orders":  orders,
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}
