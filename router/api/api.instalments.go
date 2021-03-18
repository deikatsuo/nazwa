package api

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
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
	var checked []wrapper.InstalmentPrintReceipt

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
							if misc.IsLastMonth(orders[oi].OrderInfo.CreditDetail.LastPaid) && len(orders[oi].Monthly) <= 2 {
								if !mon.Printed {
									mon.Print = true
									orders[oi].SuggestPrint = true
									checked = append(checked, wrapper.InstalmentPrintReceipt{
										ID:             mon.ID,
										Nth:            mon.Nth,
										DueDate:        mon.DueDate,
										Promise:        mon.Promise,
										PrintDate:      mon.PrintDate,
										Code:           mon.Code,
										CreditCode:     orders[oi].OrderInfo.CreditDetail.CreditCode,
										Customer:       fmt.Sprintf("%s (%s)", orders[oi].OrderInfo.Customer.Name, orders[oi].OrderInfo.Customer.Code),
										BillingAddress: orders[oi].OrderInfo.BillingAddress,
										Deposit:        orders[oi].OrderInfo.Deposit,
										Monthly:        orders[oi].OrderInfo.CreditDetail.Monthly,
										Items:          misc.ItemsToString(orders[oi].OrderInfo.Items),
										Total:          orders[oi].OrderInfo.CreditDetail.Total,
										Collector:      orders[oi].OrderInfo.Collector.Name,
									})
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
						checked = append(checked, wrapper.InstalmentPrintReceipt{
							ID:             mon.ID,
							Nth:            mon.Nth,
							DueDate:        mon.DueDate,
							Promise:        mon.Promise,
							PrintDate:      mon.PrintDate,
							Code:           mon.Code,
							CreditCode:     orderInfo.CreditDetail.CreditCode,
							Customer:       fmt.Sprintf("%s (%s)", orderInfo.Customer.Name, orderInfo.Customer.Code),
							BillingAddress: orderInfo.BillingAddress,
							Deposit:        orderInfo.Deposit,
							Monthly:        orderInfo.CreditDetail.Monthly,
							Items:          misc.ItemsToString(orderInfo.Items),
							Total:          orderInfo.CreditDetail.Total,
							Collector:      orderInfo.Collector.Name,
						})
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
		"checked": checked,
	}

	c.JSON(httpStatus, gh)
}

// InstalmentUpdateReceiptPrintStatus update printed status
func InstalmentUpdateReceiptPrintStatus(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	rid, err := strconv.Atoi(c.Param("rid"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var printed wrapper.InstalmentUpdateReceiptPrintStatus
	if err := c.ShouldBindQuery(&printed); err != nil {
		next = false
		message = "Request tidak valid"
		status = "error"
	}

	// Update arah
	if next {
		if err := dbquery.InstalmentsPrintedStatus(rid, *printed.Printed); err != nil {
			log.Warn("api.instalments.go InstalmentUpdateReceiptPrintStatus() Gagal mengubah nama arah")
			log.Error(err)
			message = "Gagal mengubah status di print"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Status print berhasil dirubah"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}
