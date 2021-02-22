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
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// OrderShowList mengambil data/list order/penjualan
func OrderShowList(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	// Tolak jika yang request bukan user terdaftar
	uid := session.Get("userid")
	if uid == nil {
		router.Page404(c)
		return
	}

	var lastid int
	last := false
	limit := 10
	var direction string
	httpStatus := http.StatusOK
	errMess := ""
	o := dbquery.GetOrders{}
	next := true

	// Mengambil parameter limit
	lim, err := strconv.Atoi(c.Param("limit"))
	if err == nil {
		limit = lim
		o.Limit(limit)
	} else {
		errMess = "Limit tidak valid"
		httpStatus = http.StatusBadRequest
		next = false
	}

	// Ambil query id terakhir
	lst, err := strconv.Atoi(c.Query("lastid"))
	if err == nil {
		lastid = lst
	}

	// Forward/backward
	direction = c.Query("direction")
	if direction == "back" {
		o.Direction(direction)
	} else {
		o.Direction("next")
	}

	var total int
	if t, err := dbquery.OrderGetOrderTotalRow(); err == nil {
		total = t
	}

	var orders []wrapper.Order

	if next {
		o.LastID(lastid)
		or, err := o.Show()
		if err != nil {
			errMess = err.Error()
			httpStatus = http.StatusInternalServerError
		}
		orders = or
	}

	if len(orders) > 0 && direction == "back" {
		// Reverse urutan array order
		temp := make([]wrapper.Order, len(orders))
		in := 0
		for i := len(orders) - 1; i >= 0; i-- {
			temp[in] = orders[i]
			in++
		}
		orders = temp
	}

	// Cek id terakhir
	if len(orders) > 0 && len(orders) < limit {
		// Periksa apakah ini data terakhir atau bukan
		last = true
	}

	c.JSON(httpStatus, gin.H{
		"orders": orders,
		"error":  errMess,
		"total":  total,
		"last":   last,
	})
}

// OrderShowByID mengambil data order/penjualan berdasarkan ID
func OrderShowByID(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	// Tolak jika yang request bukan user terdaftar
	uid := session.Get("userid")
	if uid == nil {
		router.Page404(c)
		return
	}
	httpStatus := http.StatusOK
	errMess := ""

	// Mengambil parameter id order
	var oid int
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		oid = id
	} else {
		httpStatus = http.StatusBadRequest
		errMess = "Request tidak valid"
	}

	var order wrapper.Order
	if o, err := dbquery.OrderGetOrderByID(oid); err == nil {
		order = o
	} else {
		httpStatus = http.StatusInternalServerError
		errMess = "Sepertinya telah terjadi kesalahan saat memuat data"
	}

	c.JSON(httpStatus, gin.H{
		"order": order,
		"error": errMess,
	})
}

// OrderSubstituteByRicShow menampilkan informasi ric substitute
func OrderSubstituteByRicShow(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	// Tolak jika yang request bukan user terdaftar
	uid := session.Get("userid")
	if uid == nil {
		router.Page404(c)
		return
	}
	httpStatus := http.StatusOK
	errMess := ""

	ric := c.Query("number")

	var substitutes []wrapper.NameID
	if s, err := dbquery.OrderGetSubstituteByRic(ric); err == nil {
		substitutes = s
	} else {
		httpStatus = http.StatusInternalServerError
		errMess = "Tidak ditemukan data pendamping berdasarkan NIK ini"
	}

	c.JSON(httpStatus, gin.H{
		"substitutes": substitutes,
		"error":       errMess,
	})
}

// OrderCreate API untuk menambahkan produk baru
func OrderCreate(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	// Tolak jika yang request bukan user terdaftar
	uid := session.Get("userid")
	if uid == nil {
		router.Page404(c)
		return
	}

	var json wrapper.OrderForm

	status := "success"
	var httpStatus int
	message := ""
	var simpleErrMap = make(map[string]interface{})
	save := true

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println("ERROR: api.order.go OrderCreate() bind json")
		log.Println(err)
		if fmt.Sprintf("%T", err) == "validator.ValidationErrors" {
			simpleErrMap = validation.SimpleValErrMap(err)
		}
		httpStatus = http.StatusBadRequest
		status = "fail"
		save = false
	}

	if len(json.OrderItems) == 0 {
		simpleErrMap["orderitems"] = "Item tidak boleh kosong"
		httpStatus = http.StatusBadRequest
		status = "error"
		message = "Harap pilih setidaknya satu produk"
		save = false
	}

	var code string
	if save {
		uname, err := dbquery.UserGetUsername(json.Customer)
		if err == nil {
			if uname[:3] == "NZ-" || uname[:3] == "NE-" {
				if len(uname) >= 7 {
					uname = uname[3:]
				}
			}
			tm := time.Now()
			dt := strings.ReplaceAll(tm.Format("01-02-2006"), "-", "")
			dy := tm.Format("Mon")
			uq := tm.Format(".000")[1:]
			code = strings.ToUpper(fmt.Sprintf("%s%s-%s%s-%s%s", uname[4:], dy, uq, dt[4:], dt[:4], uname[:4]))
		} else {
			message = "Telah terjadi kesalahan saat memeriksa data konsumen"
			status = "error"
			save = false
			httpStatus = http.StatusInternalServerError
		}
	}

	var retOrder wrapper.Order
	var retID int

	if save {
		for _, oi := range json.OrderItems {
			if stock, err := dbquery.ProductCheckStock(oi.ProductID); err == nil {
				if stock < oi.Quantity {
					message = fmt.Sprintf("stok:%d-%d", oi.ProductID, stock)
					status = "warning"
					save = false
					httpStatus = http.StatusBadRequest
				}
			}
		}
	}

	if save {
		deposit, _ := strconv.Atoi(json.Deposit)
		order := dbquery.NewOrder()
		err := order.SetCustomer(json.Customer).
			SetSales(json.Sales).
			SetCollector(json.Collector).
			SetSurveyor(json.Surveyor).
			SetShipping(json.ShippingAddress).
			SetBilling(json.BillingAddress).
			SetCredit(*json.Credit).
			SetDeposit(deposit).
			SetDuration(json.Duration).
			SetDue(json.Due).
			SetLine(json.Line).
			SetNotes(json.Notes).
			SetCode(code).
			SetOrderDate(json.OrderDate).
			SetShippingDate(json.ShippingDate).
			SetOrderItems(json.OrderItems).
			SetSubstitutes(json.Substitutes).
			SetCreatedBy(uid.(int)).
			ReturnID(&retID).
			Save()

		if err != nil {
			log.Warn("ERROR: api.order.go OrderCreate() Gagal membuat order baru")
			log.Error(err)
			status = "error"
			message = err.Error()
			httpStatus = http.StatusBadRequest
		} else {
			status = "success"
			message = "Berhasil membuat order"

			if o, err := dbquery.OrderGetOrderByID(retID); err == nil {
				retOrder = o
			} else {
				httpStatus = http.StatusInternalServerError
				message = "Sepertinya telah terjadi kesalahan saat memuat data"
				status = "error"
			}
		}
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"order":   retOrder,
	}
	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))

}

// OrderDeleteByID delete order
func OrderDeleteByID(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// id order
	oid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	// Delete order
	if next {
		if err := dbquery.OrderDeleteByID(oid); err != nil {
			message = "Gagal menghapus order"
			status = "error"
		} else {
			httpStatus = http.StatusOK
			message = "Order berhasil dihapus"
			status = "success"
		}
	}

	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
	})
}
