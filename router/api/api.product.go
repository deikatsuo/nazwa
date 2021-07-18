package api

import (
	"database/sql"
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
	fmt.Println(json.Description)
	if save {
		user := dbquery.NewProduct()
		err := user.SetName(json.Name).
			SetSlug(slugURL).
			SetCategory(json.Category).
			SetBrand(json.Brand).
			SetStock(stock).
			SetDescription(json.Description).
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

// ProductAddPhotos tambah foto produk
func ProductAddPhotos(c *gin.Context) {
	// User id yang merequest
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Warn("api.product.go ProductAddPhotos() parameter id tidak valid")
		log.Error(err)

		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""

	var photos wrapper.ProductAddPhotosForm
	if err := c.ShouldBindJSON(&photos); err != nil {
		log.Warn("api.product.go ProductAddPhotos() gagal bind json")
		log.Error(err)
		next = false
	}

	var files []string
	if next {
		if len(photos.Photo) > 0 {
			for _, p := range photos.Photo {
				if p.PhotoType != "" && p.Photo != "" {
					if f, err := misc.FileBase64ToFileWithData("../data/upload/product/", p.Photo, p.PhotoType); err == nil {
						files = append(files, f)
					} else {
						next = false
						log.Warn("api.product.go ProductAddPhotos() Konversi base64 ke dalam bentuk file")
						log.Error(err)
					}
				}
			}
		}
	}

	if next {
		if err := dbquery.ProductAddPhotos(pid, files); err != nil {
			errMess = "Gagal menyimpan ke database"
			log.Warn("api.product.go ProductAddPhotos() Simpan file ke database")
			log.Error(err)

			for _, s := range files {
				if err := os.Remove("../data/upload/product/" + s); err != nil {
					log.Warn("api.product.go ProductCreate() Gagal menghapus file")
					log.Error(err)
				}
			}
		}
	}

	if next {
		success = "Foto berhasil ditambahkan"
		httpStatus = http.StatusOK
		errMess = ""
	}

	c.JSON(httpStatus, gin.H{
		"error":   errMess,
		"success": success,
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

// ProductDeletePhoto hapus foto produk
func ProductDeletePhoto(c *gin.Context) {
	// id produk
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Warn("api.product.go ProductDeletePhoto() id produk tidak valid")
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

	photoID, err := jsonparser.GetInt(body, "id")
	if err != nil {
		errMess = "Request tidak valid"
	}

	var photo string
	if next {
		if ph, err := dbquery.ProductGetProductPhotoName(int(photoID)); err != nil {
			log.Warn(err)
			errMess = "tidak ditemukan foto"
			next = false
		} else {
			photo = ph
		}
	}

	var thumb string
	if next {
		if tb, err := dbquery.ProductGetProductThumbName(pid); err != nil {
			log.Warn(err)
		} else {
			thumb = tb
		}
	}

	if thumb == photo {
		err := os.Remove("../data/upload/product/thumbnail/" + thumb)
		if err != nil {
			log.Error(err)
			errMess = "Gagal menghapus thumbnail dari direktori"
			next = false
		}
	}

	// Update thumb
	if next {
		var photos []wrapper.ProductPhotoListSelect

		if pp, err := dbquery.ProductGetProductPhoto(pid); err == nil {
			photos = pp
		}

		// Buat thumbnail
		if len(photos) > 0 {
			err := misc.FileGenerateThumb(photos[0].Photo, "../data/upload/product/")
			if err != nil {
				log.Warn("Gagal membuat thumbnail")
				log.Error(err)
			} else {
				if err := dbquery.ProductUpdateThumb(pid, sql.NullString{String: photos[0].Photo, Valid: true}); err != nil {
					errMess = "Gagal menghapus thumbnail"
					next = false
				} else {
					success = "Thumbnail telah dihapus"
				}
			}
		} else {

			if err := dbquery.ProductUpdateThumb(pid, sql.NullString{}); err != nil {
				errMess = "Gagal menghapus thumbnail"
				next = false
			} else {
				success = "Thumbnail telah dihapus"
			}
		}
	}

	// Hapus photo
	if next {
		if err := dbquery.ProductDeletePhoto(photoID, pid); err != nil {
			errMess = "Gagal menghapus foto"
			next = false
		} else {
			success = "Foto telah dihapus"
		}
	}

	if next {
		err := os.Remove("../data/upload/product/" + photo)
		if err != nil {
			errMess = "Gagal menghapus foto dari direktori"
			next = false
		}
	}

	c.JSON(httpStatus, gin.H{
		"error":   errMess,
		"success": success,
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

// ProductUpdateName ubah nama produk
func ProductUpdateName(c *gin.Context) {
	// ID produk yang akan di update
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var simpleErr map[string]interface{}

	var updateName wrapper.ProductUpdateName
	if err := c.ShouldBindQuery(&updateName); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
		message = simpleErr["name"].(string)
		status = "error"
	}

	// Update nama produk
	if next {
		if err := dbquery.ProductUpdateName(pid, updateName.Name); err != nil {
			log.Warn("api.product.go ProductUpdateName() Gagal mengubah nama produk")
			log.Error(err)
			message = "Gagal mengubah nama produk"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Nama produk telah dirubah"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}

// ProductUpdateBrand update brand
func ProductUpdateBrand(c *gin.Context) {
	// ID produk yang akan di update
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var simpleErr map[string]interface{}

	var updateBrand wrapper.ProductUpdateBrand
	if err := c.ShouldBindQuery(&updateBrand); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
		message = simpleErr["brand"].(string)
		status = "error"
	}

	// Update brand produk
	if next {
		if err := dbquery.ProductUpdateBrand(pid, updateBrand.Brand); err != nil {
			log.Warn("api.product.go ProductUpdateBrand() Gagal mengubah brand produk")
			log.Error(err)
			message = "Gagal mengubah brand produk"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Brand telah dirubah"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}

// ProductUpdateCategory update kategori
func ProductUpdateCategory(c *gin.Context) {
	// ID produk yang akan di update
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var simpleErr map[string]interface{}

	var updateCategory wrapper.ProductUpdateCategory
	if err := c.ShouldBindQuery(&updateCategory); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
		message = simpleErr["category"].(string)
		status = "error"
	}

	// Update kategori produk
	if next {
		if err := dbquery.ProductUpdateCategory(pid, updateCategory.Category); err != nil {
			log.Warn("api.product.go ProductUpdateCategory() Gagal mengubah kategori produk")
			log.Error(err)
			message = "Gagal mengubah kategori produk"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Kategori produk telah dirubah"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}

// ProductUpdateDescription update deskripsi produk
func ProductUpdateDescription(c *gin.Context) {
	// ID produk yang akan di update
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var simpleErr map[string]interface{}

	var updateDescription wrapper.ProductUpdateDescription
	if err := c.ShouldBindQuery(&updateDescription); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
		message = simpleErr["description"].(string)
		status = "error"
	}

	// Update deskripsi produk
	if next {
		if err := dbquery.ProductUpdateDescription(pid, updateDescription.Description); err != nil {
			log.Warn("api.product.go ProductUpdateDescription() Gagal mengubah deskripsi produk")
			log.Error(err)
			message = "Gagal mengubah deskripsi produk"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Deskripsi produk telah dirubah"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}

// ProductUpdatePriceBuy Update harga beli
func ProductUpdatePriceBuy(c *gin.Context) {
	// ID produk yang akan di update
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var simpleErr map[string]interface{}

	var updatePriceBuy wrapper.ProductUpdatePriceBuy
	if err := c.ShouldBindQuery(&updatePriceBuy); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
		message = simpleErr["base_price"].(string)
		status = "error"
	}

	price, err := strconv.Atoi(updatePriceBuy.BasePrice)
	if err != nil {
		next = false
		message = "Nilai input harga beli tidak benar"
		status = "error"
	}

	// Update harga beli produk
	if next {
		if err := dbquery.ProductUpdatePriceBuy(pid, price); err != nil {
			log.Warn("api.product.go ProductUpdatePriceBuy() Gagal mengubah harga beli produk")
			log.Error(err)
			message = "Gagal mengubah harga beli produk"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Harga beli produk telah dirubah"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}

// ProductUpdatePriceSell Update harga jual
func ProductUpdatePriceSell(c *gin.Context) {
	// ID produk yang akan di update
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var simpleErr map[string]interface{}

	var updatePriceSell wrapper.ProductUpdatePriceSell
	if err := c.ShouldBindQuery(&updatePriceSell); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
		message = simpleErr["price"].(string)
		status = "error"
	}

	price, err := strconv.Atoi(updatePriceSell.Price)
	if err != nil {
		next = false
		message = "Nilai input harga jual tidak benar"
		status = "error"
	}

	// Update harga jual produk
	if next {
		if err := dbquery.ProductUpdatePriceSell(pid, price); err != nil {
			log.Warn("api.product.go ProductUpdatePriceSell() Gagal mengubah harga jual produk")
			log.Error(err)
			message = "Gagal mengubah harga jual produk"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Harga jual produk telah dirubah"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}
