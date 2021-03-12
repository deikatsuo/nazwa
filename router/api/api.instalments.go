package api

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/wrapper"
	"net/http"

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
		if simpleErr["date"] != nil {
			message = simpleErr["date"].(string)
		}
		status = "error"
		httpStatus = http.StatusBadRequest
	}

	var monthly []wrapper.OrderMonthlyCredit
	if next {
		if mon, err := dbquery.OrderGetMonthlyCreditByDate(byDate.ZoneID, byDate.Date); err == nil {
			monthly = mon
		} else {
			httpStatus = http.StatusInternalServerError
			message = "Sepertinya telah terjadi kesalahan saat memuat data kredit"
			status = "error"
			next = false
		}
	}

	var orders []wrapper.InstalmentOrderReceipt
	var tlines []wrapper.LocationLine
	var lines []wrapper.LocationLine

	// Kumpulkan data order
	if next {
		for _, mon := range monthly {
			// Tanggal cetak

			if mon.PrintDate == "" {
				mon.PrintDate = "-"
			}

			exist := false
			if len(orders) > 0 {
				for oi, ord := range orders {
					if ord.OrderID == mon.OrderID {

						// Kwitansi yang harus di print hari ini
						if !orders[oi].SuggestPrint {
							if misc.IsLastMonth(orders[oi].OrderInfo.CreditDetail.LastPaid) {
								if !mon.Printed {
									mon.Print = true
									orders[oi].SuggestPrint = true
								}
							}
						}
						orders[oi].Monthly = append(orders[oi].Monthly, mon)

						exist = true
					}
				}
			}

			if !exist {
				// order info
				var orderInfo wrapper.OrderSimple

				if ord, err := dbquery.OrderGetSimpleOrderByID(mon.OrderID); err == nil {
					orderInfo = ord
				} else {
					httpStatus = http.StatusInternalServerError
					message = "Sepertinya telah terjadi kesalahan saat memuat info order"
					status = "error"
					next = false
					break
				}

				// kwitansi yang harus di print hari ini
				var suggestPrint bool
				if misc.IsLastMonth(orderInfo.CreditDetail.LastPaid) {
					if !mon.Printed {
						mon.Print = true
						suggestPrint = true
					}
				}

				// push data order
				orders = append(orders, wrapper.InstalmentOrderReceipt{
					OrderID:      mon.OrderID,
					SuggestPrint: suggestPrint,
					OrderInfo:    orderInfo,
					Monthly:      []wrapper.OrderMonthlyCredit{mon},
				})

				// push data arah
				tlines = append(tlines, wrapper.LocationLine{
					ID:   orderInfo.CreditDetail.ZoneLine.ID,
					Code: orderInfo.CreditDetail.ZoneLine.Code,
					Name: orderInfo.CreditDetail.ZoneLine.Name,
				})
			}
		}
	}

	// kumpulkan data arah
	if next {
		for _, tln := range tlines {
			exist := false
			if len(lines) > 0 {
				for li, ln := range lines {
					if ln.ID == tln.ID {
						lines[li].Count++
						exist = true
					}
				}
			}

			if !exist {
				lines = append(lines, wrapper.LocationLine{
					ID:    tln.ID,
					Code:  tln.Code,
					Name:  tln.Name,
					Count: 1,
				})
			}
		}
	}

	if next {
		if len(orders) == 0 {
			message = "Tidak ada tagihan untuk hari ini"
			status = "error"
			next = false
		}
	}

	if next {
		status = "success"
	}

	gh := gin.H{
		"orders":  orders,
		"lines":   lines,
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}
