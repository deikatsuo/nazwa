package api

import (
	"fmt"
	"io/ioutil"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

//////////////////////////////////// [POST] //////////////////////////////////////////

// ProductCreate API untuk menambahkan produk baru
func ProductCreate(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	// Tolak jika yang request bukan user terdaftar
	uid := session.Get("userid")
	if uid == nil {
		router.Page404(c)
		return
	}

	var json wrapper.ProductForm

	status := "success"
	var httpStatus int
	message := ""
	var simpleErrMap = make(map[string]interface{})
	save := true
	var files []string

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Warn("api.product.go ProductCreate() bind json")
		log.Error(err)
		simpleErrMap = validation.SimpleValErrMap(err)
		httpStatus = http.StatusBadRequest
		status = "fail"
		save = false
	}

	if dbquery.ProductSkuExist(json.Code) {
		simpleErrMap["code"] = "SKU atau Kode produk sudah terdaftar"
		status = "fail"
		save = false
	}

	if save {
		if len(json.Photo) > 0 {
			for _, p := range json.Photo {
				if p.PhotoType != "" && p.Photo != "" {
					if f, err := misc.FileBase64ToFileWithData("../data/upload/product/", p.Photo, p.PhotoType); err == nil {
						files = append(files, f)
					} else {
						log.Warn("api.product.go ProductCreate() Konversi base64 ke dalam bentuk file")
						log.Error(err)
					}
				}
			}
		}

		// Buat thumbnail
		if len(files) > 0 {
			err := misc.FileGenerateThumb(files[0], "../data/upload/product/")
			if err != nil {
				log.Warn("Gagal membuat thumbnail")
				log.Error(err)
			}
		}
	}

	slugURL := slug.Make(json.Name)

	base, _ := strconv.Atoi(json.BasePrice)
	price, _ := strconv.Atoi(json.Price)
	stock, _ := strconv.Atoi(json.Stock)
	var retProduct wrapper.Product
	var pid int
	if save {
		user := dbquery.NewProduct()
		err := user.SetName(json.Name).
			SetSlug(slugURL).
			SetCode(json.Code).
			SetType(json.Type).
			SetBrand(json.Brand).
			SetStock(stock).
			SetBasePrice(base).
			SetPrice(price).
			SetCreatedBy(uid.(int)).
			SetPhotos(files).
			SetCreditPrice(json.CreditPrice).
			ReturnID(&pid).
			Save()
		if err != nil {
			log.Warn("api.product.go ProductCreate() Gagal menambahkan produk baru")
			log.Error(err)
			status = "error"
			message = "Gagal menambahkan produk baru"

			if len(files) > 0 {
				for _, s := range files {
					if err := os.Remove("../data/upload/product/" + s); err != nil {
						log.Warn("api.product.go ProductCreate() Gagal menghapus file")
						log.Error(err)
					}
				}
				// Hapus thumbnail
				if err := os.Remove("../data/upload/product/thumbnail/" + files[0]); err != nil {
					log.Warn("api.product.go ProductCreate() Gagal menghapus file thumbnail")
					log.Error(err)
				}
			}
		} else {
			httpStatus = http.StatusOK
			status = "success"
			message = "Berhasil menambahkan produk"

			if p, err := dbquery.ProductGetProductByID(pid); err == nil {
				retProduct = p
			} else {
				log.Warn("Gagal mengambil data produk by ID")
				log.Error(err)
				httpStatus = http.StatusInternalServerError
				message = "Sepertinya telah terjadi kesalahan saat memuat data"
			}
		}
	} else {
		httpStatus = http.StatusBadRequest
		status = "error"
		message = "Gagal menambahkan produk, silahkan perbaiki formulir"
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"product": retProduct,
	}
	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}

// ProductAddCreditPrice menambahkan harga kredit barang
func ProductAddCreditPrice(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		log.Warn("api.product.go ProductAddCreditPrice() parameter id tidak valid")
		log.Error(err)

		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""

	var creditPrice wrapper.ProductCreditPriceForm
	if err := c.ShouldBindJSON(&creditPrice); err != nil {
		log.Warn("api.product.go ProductAddCreditPrice() gagal bind json")
		log.Error(err)
		next = false
	}

	if creditPrice.Duration <= 1 {
		errMess = "Durasi harus lebih dari satu bulan"
		next = false
	}

	if next {
		if dbquery.ProductCreditDurationExist(pid, creditPrice.Duration) {
			errMess = fmt.Sprintf("Produk dengan durasi %d sudah ada", creditPrice.Duration)
			next = false
		}
	}

	var insertCreditPrice []wrapper.ProductCreditPriceInsert
	if next {
		insertCreditPrice = append(insertCreditPrice, wrapper.ProductCreditPriceInsert{
			ProductID: pid,
			Duration:  creditPrice.Duration,
			Price:     creditPrice.Price,
		})
		if err := dbquery.ProductInsertCreditPrice(insertCreditPrice); err != nil {
			log.Warn("api.product.go ProductAddCreditPrice() Gagal menyimpan harga kredit")
			log.Error(err)

			errMess = "Gagal menambahkan harga kredit"
			next = false
		}

	}

	// Ambil harga kredit barang dari database
	var retCreditPrices []wrapper.ProductCreditPriceSelect
	if next {
		pp, err := dbquery.ProductGetProductCreditPrice(pid)
		if err != nil {
			log.Warn("api.product.go ProductAddCreditPrice() gagal mengambil harga barang dari database")
			log.Error(err)

			errMess = "Gagal memuat harga kredit barang"
		} else {
			retCreditPrices = pp
			httpStatus = http.StatusOK
			success = "Harga kredit berhasil ditambahkan"
		}
	}
	c.JSON(httpStatus, gin.H{
		"error":        errMess,
		"success":      success,
		"credit_price": retCreditPrices,
	})
}

/////////////////////////////////////// DELETE /////////////////////////////////////////

// ProductDeleteCreditPrice hapus harga kredit barang
func ProductDeleteCreditPrice(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		log.Warn("api.product.go ProductDeleteCreditPrice() parameter id tidak valid")
		log.Error(err)

		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Warn("api.product.go ProductDeleteCreditPrice() gagal membaca body")
		log.Error(err)

		errMess = "Data tidak benar"
		next = false
	}

	pcpid, err := jsonparser.GetInt(body, "id")
	if err != nil {
		errMess = "Request tidak valid"
	}

	// Delete harga kredit
	if next {
		if err := dbquery.ProductDeleteCreditPrice(pcpid, pid); err != nil {
			errMess = "Gagal menghapus harga kredit"
			next = false
		}
	}

	// Ambil harga kredit sisa jika masih ada
	var retCreditPrices []wrapper.ProductCreditPriceSelect
	if next {
		rcp, err := dbquery.ProductGetProductCreditPrice(pid)
		if err != nil {
			errMess = "Gagal memuat harga kredit/semua harga kredit sudah dihapus"
		} else {
			retCreditPrices = rcp
			httpStatus = http.StatusOK
			success = "Harga kredit berhasil dihapus"
		}
	}
	c.JSON(httpStatus, gin.H{
		"error":        errMess,
		"success":      success,
		"credit_price": retCreditPrices,
	})
}

/////////////////////////////////////// GET ////////////////////////////////////////////

// ProductShowList mengambil data/list produk
func ProductShowList(c *gin.Context) {
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
	pts := dbquery.GetProducts{}
	next := true

	// Mengambil parameter limit
	lim, err := strconv.Atoi(c.Param("limit"))
	if err == nil {
		limit = lim
		pts.Limit(limit)
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
		pts.Direction(direction)
	} else {
		pts.Direction("next")
	}

	// Total produk
	var total int
	if t, err := dbquery.ProductGetProductTotalRow(); err == nil {
		total = t
	}

	var products []wrapper.Product

	if next {
		pts.LastID(lastid)
		// Maju/Mundur
		if direction == "next" {
			pts.Where("WHERE id > " + strconv.Itoa(lastid) + " ORDER BY id ASC")
		} else if direction == "back" {
			pts.Where("WHERE id < " + strconv.Itoa(lastid) + " ORDER BY id DESC")
		}
		p, err := pts.Show()
		if err != nil {
			errMess = err.Error()
			httpStatus = http.StatusInternalServerError
		}
		products = p

	}

	if len(products) > 0 && direction == "back" {
		// Reverse urutan array produk
		temp := make([]wrapper.Product, len(products))
		in := 0
		for i := len(products) - 1; i >= 0; i-- {
			temp[in] = products[i]
			in++
		}
		products = temp
	}

	// Cek id terakhir
	if len(products) > 0 && len(products) < limit {
		// Periksa apakah ini data terakhir atau bukan
		last = true
	}

	c.JSON(httpStatus, gin.H{
		"products": products,
		"error":    errMess,
		"total":    total,
		"last":     last,
	})
}

// ProductShowAll mengambil semua data/list produk
func ProductShowAll(c *gin.Context) {
	var direction string
	httpStatus := http.StatusOK
	errMess := ""
	pts := dbquery.GetProducts{}

	// Forward/backward
	direction = c.Query("direction")
	if direction == "back" {
		pts.Direction(direction)
	} else {
		pts.Direction("next")
	}

	var products []wrapper.Product

	// Maju/Mundur
	if direction == "next" {
		pts.Where("ORDER BY name ASC")
	} else if direction == "back" {
		pts.Where("ORDER BY name DESC")
	}
	p, err := pts.Show()
	if err != nil {
		errMess = err.Error()
		httpStatus = http.StatusInternalServerError
	}
	products = p

	if len(products) > 0 && direction == "back" {
		// Reverse urutan array produk
		temp := make([]wrapper.Product, len(products))
		in := 0
		for i := len(products) - 1; i >= 0; i-- {
			temp[in] = products[i]
			in++
		}
		products = temp
	}

	c.JSON(httpStatus, gin.H{
		"products": products,
		"error":    errMess,
	})
}

// ProductShowByID mengambil data produk berdasarkan ID
func ProductShowByID(c *gin.Context) {
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

	// Mengambil parameter id produk
	var pid int
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		pid = id
	} else {
		httpStatus = http.StatusBadRequest
		errMess = "Request tidak valid"
	}

	var product wrapper.Product
	if p, err := dbquery.ProductGetProductByID(pid); err == nil {
		product = p
	} else {
		httpStatus = http.StatusInternalServerError
		errMess = "Sepertinya telah terjadi kesalahan saat memuat data"
	}

	c.JSON(httpStatus, gin.H{
		"product": product,
		"error":   errMess,
	})
}

// ProductSearchByName cari produk berdasarkan nama
func ProductSearchByName(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	// Tolak jika yang request bukan user terdaftar
	uid := session.Get("userid")
	if uid == nil {
		router.Page404(c)
		return
	}

	search := ""
	lastid := 1
	last := false
	limit := 10
	var direction string
	httpStatus := http.StatusOK
	errMess := ""
	p := dbquery.GetProducts{}
	next := true

	// Mengambil parameter limit
	lim, err := strconv.Atoi(c.Param("limit"))
	if err == nil {
		limit = lim
		p.Limit(limit)
	} else {
		errMess = "Limit tidak valid"
		httpStatus = http.StatusBadRequest
		next = false
	}

	// Ambil query pencarian
	search = c.Query("search")

	// Ambil query id terakhir
	lst, err := strconv.Atoi(c.Query("lastid"))
	if err == nil {
		lastid = lst
	}

	// Forward/backward
	direction = c.Query("direction")
	if direction == "back" {
		p.Direction(direction)
	} else {
		p.Direction("next")
	}

	var products []wrapper.Product

	if next {
		p.Where("WHERE name LIKE '" + search + "%' ORDER BY id ASC")
		p.LastID(lastid)

		prod, err := p.Show()
		if err != nil {
			errMess = err.Error()
			httpStatus = http.StatusInternalServerError
		}
		products = prod
	}

	if len(products) > 0 && direction == "back" {
		// Reverse urutan array user
		temp := make([]wrapper.Product, len(products))
		in := 0
		for i := len(products) - 1; i >= 0; i-- {
			temp[in] = products[i]
			in++
		}
		products = temp
	}

	// Cek id terakhir
	if len(products) > 0 && len(products) < limit {
		// Periksa apakah ini data terakhir atau bukan
		last = true
	}

	c.JSON(httpStatus, gin.H{
		"products": products,
		"error":    errMess,
		"last":     last,
	})
}

// ProductUpdateStock Tambah/Kurangi stok
func ProductUpdateStock(c *gin.Context) {
	// ID produk yang akan di update stoknya
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := "error"
	var simpleErr map[string]interface{}

	// Stock
	setStock := strings.ReplaceAll(c.Query("set"), " ", "")
	var newStock int

	if stock, err := dbquery.ProductCheckStock(pid); err == nil {
		// Tambah/Kurang
		if setStock[:1] == "-" {
			if tmpStock, err := strconv.Atoi(setStock[1:]); err == nil {
				if tmpStock > stock {
					message = "Tidak bisa mengurangi stok melebihi nilai stok saat ini"
					next = false
				} else {
					newStock = stock - tmpStock
				}
			} else {
				message = "Format pengurangan stok tidak valid"
				next = false
			}
		} else {
			if tmpStock, err := strconv.Atoi(setStock); err == nil {
				newStock = stock + tmpStock
			} else {
				message = "Format penambahan stok tidak valid"
				next = false
			}
		}

	}

	if next {
		if err := dbquery.ProductUpdateStock(pid, newStock); err != nil {
			log.Warn("api.product.go ProductUpdateStock() menambah/mengurangi stock stok")
			log.Error(err)
			message = "Gagal menambah/mengurangi stok"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Stok berhasil diperbarui"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, misc.Mete(gh, simpleErr))
}
