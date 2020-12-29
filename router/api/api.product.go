package api

import (
	"fmt"
	"log"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ProductCreate API untuk menambahkan produk baru
func ProductCreate(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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
			log.Println("ERROR: api.product.go ProductCreate() bind json")
			log.Println(err)
			if fmt.Sprintf("%T", err) == "validator.ValidationErrors" {
				simpleErrMap = validation.SimpleValErrMap(err)
			}
			httpStatus = http.StatusBadRequest
			status = "fail"
			save = false
		}

		if dbquery.ProductSkuExist(db, json.Code) {
			simpleErrMap["code"] = "SKU atau Kode produk sudah terdaftar"
			status = "fail"
			save = false
		}

		if len(json.Photo) > 0 {
			for _, p := range json.Photo {
				if p.PhotoType != "" && p.Photo != "" {
					if f, err := misc.Base64ToFileWithData("../data/upload/product/", p.Photo, p.PhotoType); err == nil {
						files = append(files, f)
					} else {
						log.Println("ERROR: api.product.go ProductCreate() Konversi base64 ke dalam bentuk file")
						message = err.Error()
					}
				}
			}
		}

		// Buat thumbnail
		if len(files) > 0 {
			err := misc.FileGenerateThumb(files[0], "../data/upload/product/")
			if err != nil {
				message = err.Error()
			}
		}

		var retProduct wrapper.Product
		var pid int
		if save {
			user := dbquery.NewProduct()
			err := user.SetName(json.Name).
				SetCode(json.Code).
				SetType(json.Type).
				SetBrand(json.Brand).
				SetBasePrice(json.BasePrice).
				SetPrice(json.Price).
				SetCreatedBy(uid.(int)).
				SetPhotos(files).
				SetCreditPrice(json.CreditPrice).
				ReturnID(&pid).
				Save(db)
			if err != nil {
				log.Println("ERROR: api.product.go ProductCreate() Gagal menambahkan produk baru")
				log.Print(err)
				status = "error"
				message = "Gagal menambahkan produk baru"

				if len(files) > 0 {
					for _, s := range files {
						if err := os.Remove("../data/upload/product/" + s); err != nil {
							log.Println("ERROR: api.product.go ProductCreate() Gagal menghapus file")
							log.Println(err)
						}
					}
					// Hapus thumbnail
					if err := os.Remove("../data/upload/product/thumbnail/" + files[0]); err != nil {
						log.Println("ERROR: api.product.go ProductCreate() Gagal menghapus file thumbnail")
						log.Println(err)
					}
				}
			} else {
				httpStatus = http.StatusOK
				status = "success"
				message = "Berhasil menambahkan produk"

				if p, err := dbquery.ProductGetProductByID(db, pid); err == nil {
					retProduct = p
				} else {
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
	return gin.HandlerFunc(fn)
}

// ProductShowList mengambil data/list produk
func ProductShowList(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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
		if t, err := dbquery.ProductGetProductTotalRow(db); err == nil {
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
			p, err := pts.Show(db)
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
	return gin.HandlerFunc(fn)
}

// ProductShowAll mengambil semua data/list produk
func ProductShowAll(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		session := sessions.Default(c)
		// User session saat ini
		// Tolak jika yang request bukan user terdaftar
		uid := session.Get("userid")
		if uid == nil {
			router.Page404(c)
			return
		}

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
			pts.Where("ORDER BY id ASC")
		} else if direction == "back" {
			pts.Where("ORDER BY id DESC")
		}
		p, err := pts.Show(db)
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
	return gin.HandlerFunc(fn)
}

// ProductShowByID mengambil data produk berdasarkan ID
func ProductShowByID(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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
		if p, err := dbquery.ProductGetProductByID(db, pid); err == nil {
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
	return fn
}

////////////
// SEARCH //
////////////

// ProductSearchByName cari user berdasarkan Nomor induk kependudukan
func ProductSearchByName(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

			prod, err := p.Show(db)
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
	return gin.HandlerFunc(fn)
}
