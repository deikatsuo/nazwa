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
						if !mon.Done {
							orders[oi].Undone += 1
						}
						// Kwitansi yang harus di print hari ini
						if !orders[oi].SuggestPrint {
							if (misc.IsLastMonth(orders[oi].OrderInfo.CreditDetail.LastPaid) || misc.IsThisMonth(orders[oi].OrderInfo.CreditDetail.LastPaid)) && ord.Undone < 2 {
								if !mon.Printed {
									mon.Print = true
									orders[oi].SuggestPrint = true
									var tmpBillAddr string
									var tmpItems string

									if orders[oi].OrderInfo.ImportedAddress == "" {
										tmpBillAddr = orders[oi].OrderInfo.BillingAddress
									} else {
										tmpBillAddr = orders[oi].OrderInfo.ImportedAddress
									}

									if orders[oi].OrderInfo.ImportedItems == "" {
										tmpItems = misc.ItemsToString(orders[oi].OrderInfo.Items)
									} else {
										tmpItems = orders[oi].OrderInfo.ImportedItems
									}
									checked = append(checked, wrapper.InstalmentPrintReceipt{
										ID:             mon.ID,
										Nth:            mon.Nth,
										DueDate:        mon.DueDate,
										Promise:        mon.Promise,
										PrintDate:      mon.PrintDate,
										Code:           mon.Code,
										CreditCode:     orders[oi].OrderInfo.CreditDetail.CreditCode,
										Customer:       fmt.Sprintf("%s (%s)", orders[oi].OrderInfo.Customer.Name, orders[oi].OrderInfo.Customer.Code),
										BillingAddress: tmpBillAddr,
										Deposit:        orders[oi].OrderInfo.Deposit,
										Monthly:        orders[oi].OrderInfo.CreditDetail.Monthly,
										Items:          tmpItems,
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
						var tmpBillAddr string
						var tmpItems string

						if orderInfo.ImportedAddress == "" {
							tmpBillAddr = orderInfo.BillingAddress
						} else {
							tmpBillAddr = orderInfo.ImportedAddress
						}

						if orderInfo.ImportedItems == "" {
							tmpItems = misc.ItemsToString(orderInfo.Items)
						} else {
							tmpItems = orderInfo.ImportedItems
						}

						checked = append(checked, wrapper.InstalmentPrintReceipt{
							ID:             mon.ID,
							Nth:            mon.Nth,
							DueDate:        mon.DueDate,
							Promise:        mon.Promise,
							PrintDate:      mon.PrintDate,
							Code:           mon.Code,
							CreditCode:     orderInfo.CreditDetail.CreditCode,
							Customer:       fmt.Sprintf("%s (%s)", orderInfo.Customer.Name, orderInfo.Customer.Code),
							BillingAddress: tmpBillAddr,
							Deposit:        orderInfo.Deposit,
							Monthly:        orderInfo.CreditDetail.Monthly,
							Items:          tmpItems,
							Total:          orderInfo.CreditDetail.Total,
							Collector:      orderInfo.Collector.Name,
						})
					}
				}

				var Undone int
				if !mon.Done {
					Undone += 1
				}

				// push data order
				orders = append(orders, wrapper.InstalmentOrderReceipt{
					OrderID:      mon.OrderID,
					SuggestPrint: suggestPrint,
					Undone:       Undone,
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

// InstalmentMoneyIn uang tagihan masuk
func InstalmentMoneyIn(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	// id order
	oid, err := strconv.Atoi(c.Param("oid"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var moneyIn wrapper.InstalmentMoneyIn
	if err := c.ShouldBindJSON(&moneyIn); err != nil {
		log.Warn("Gagal unmarshal json")
		log.Error(err)

		message = "Request tidak valid"
		status = "error"
		next = false
	}

	money, _ := strconv.Atoi(moneyIn.MoneyIn)

	// Masukan angsuran
	if next {
		if err := dbquery.InstalmentsMoneyIn(oid, moneyIn.Receiver, money, moneyIn.Notes, moneyIn.Mode); err != nil {
			message = "Gagal membayar angsuran"
			status = "error"
			next = false
		} else {
			message = "Angsuran berhasil dibayarkan"
			status = "success"
			next = true
			httpStatus = http.StatusOK
		}
	}

	c.JSON(httpStatus, gin.H{
		"message": message,
		"status":  status,
	})
}

// InstalmentMoneyOut uang tagihan batal masuk
func InstalmentMoneyOut(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	// id order
	oid, err := strconv.Atoi(c.Param("oid"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var moneyIn wrapper.InstalmentMoneyOut
	if err := c.ShouldBindJSON(&moneyIn); err != nil {
		log.Warn("Gagal unmarshal json")
		log.Error(err)

		message = "Request tidak valid"
		status = "error"
		next = false
	}

	// Batal masuk angsuran
	if next {
		if err := dbquery.InstalmentsMoneyOut(oid, moneyIn.PaymentId); err != nil {
			message = "Gagal menghapus angsuran"
			status = "error"
			next = false
		} else {
			message = "Angsuran berhasil dibatalkan"
			status = "success"
			next = true
			httpStatus = http.StatusOK
		}
	}

	c.JSON(httpStatus, gin.H{
		"message": message,
		"status":  status,
	})
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

	// Update status print
	if next {
		if err := dbquery.InstalmentsReceiptPrintedStatus(rid, *printed.Printed); err != nil {
			log.Warn("api.instalments.go InstalmentUpdateReceiptPrintStatus() Gagal mengubah status print")
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

// InstalmentUpdateReceiptNotes update notes
func InstalmentUpdateReceiptNotes(c *gin.Context) {
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

	notes := c.Query("set")

	// Update notes
	if next {
		if err := dbquery.InstalmentsReceiptUpdateNotes(rid, notes); err != nil {
			log.Warn("api.instalments.go InstalmentUpdateReceiptNotes() Gagal mengubah notes")
			log.Error(err)
			message = "Gagal mengubah notes"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Disimpan"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}
