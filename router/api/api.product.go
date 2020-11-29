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

		var json wrapper.FormProduct

		status := "success"
		var httpStatus int
		message := ""
		var simpleErrMap = make(map[string]interface{})
		save := true
		var files []string

		if err := c.ShouldBindJSON(&json); err != nil {
			log.Println("ERROR: api.create-account.go UserCreate() bind json")
			log.Println(err)
			if fmt.Sprintf("%T", err) == "validator.ValidationErrors" {
				simpleErrMap = validation.SimpleValErrMap(err)
			}
			httpStatus = http.StatusBadRequest
			status = "fail"
			save = false
		}

		if len(json.Photo) > 0 {
			for _, p := range json.Photo {
				if p.PhotoType != "" && p.Photo != "" {
					if f, err := misc.Base64ToFileWithData("./upload/product/", p.Photo, p.PhotoType); err == nil {
						files = append(files, f)
					} else {
						log.Println("ERROR: api.product.go ProductCreate() Konversi base64 ke dalam bentuk file")
						message = err.Error()
					}
				}
			}
		}

		var retProduct wrapper.Product
		var pid int
		if save {
			user := dbquery.NewUser()
			err := user.SetFirstName(json.Firstname).
				SetLastName(json.Lastname).
				SetFamilyCard(json.FC).
				SetRIC(json.RIC).
				SetPhone(json.Phone).
				SetAvatar(file).
				SetPassword(json.Password).
				SetGender(json.Gender).
				SetOccupation(json.Occupation).
				SetRole(dbquery.RoleCustomer).
				ReturnID(&pid).
				Save(db)
			if err != nil {
				log.Println("ERROR: api.product.go ProductCreate() Gagal menambahkan produk baru")
				log.Print(err)
				if err := os.Remove("./upload/product/" + file); err != nil {
					log.Println("ERROR: api.product.go ProductCreate() Gagal menghapus file")
					log.Println(err)
				}
			} else {
				httpStatus = http.StatusOK
				status = "success"
				message = "Berhasil menambahkan produk"

				if p, err := dbquery.GetUserByID(db, pid); err == nil {
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

// ShowProductList mengambil data/list produk
func ShowProductList(db *sqlx.DB) gin.HandlerFunc {
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

		var total int
		if t, err := dbquery.GetProductTotalRow(db); err == nil {
			total = t
		}

		var products []wrapper.Product

		if next {
			pts.LastID(lastid)
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

// ShowProductByID mengambil data produk berdasarkan ID
func ShowProductByID(db *sqlx.DB) gin.HandlerFunc {
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
		if p, err := dbquery.GetProductByID(db, pid); err == nil {
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
