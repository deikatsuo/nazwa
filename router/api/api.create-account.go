package api

import (
	"log"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// FormUser menyimpan input pendaftaran user
type FormUser struct {
	Phone      string `json:"phone" binding:"numeric,min=6,max=15"`
	Firstname  string `json:"firstname" binding:"required,alpha,min=3,max=25"`
	Lastname   string `json:"lastname"`
	Gender     string `json:"gender" binding:"required,oneof=m f"`
	Password   string `json:"password" binding:"alphanumunicode,min=8,max=25"`
	Repassword string `json:"repassword" binding:"eqfield=Password"`
	Policy     bool   `json:"policy" binding:"required"`
}

// UserCreate API untuk membuat user baru
func UserCreate(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var json FormUser

		var status string
		var httpStatus int
		message := ""
		var simpleErrMap = make(map[string]interface{})
		var save bool

		if err := c.ShouldBindJSON(&json); err != nil {
			simpleErrMap = validation.SimpleValErrMap(err)
			httpStatus = http.StatusBadRequest
			status = "fail"
			message = "Data tidak lengkap"
		} else {
			if !dbquery.PhoneExist(db, json.Phone) {
				save = true
			} else {
				simpleErrMap["phone"] = "Nomor ini sudah terdaftar"
				httpStatus = http.StatusBadRequest
				status = "error"
				message = "Gagal mendaftar! anda atau seseorang selain anda telah mendaftar dengan nomor ini. \n Silahkan hubungi Administrator jika anda tidak pernah mendaftarkan nomor ini."
			}
		}

		if save {
			user := dbquery.NewUser()
			err := user.SetFirstName(json.Firstname).
				SetLastName(json.Lastname).
				SetPhone(json.Phone).
				SetPassword(json.Password).
				SetGender(json.Gender).
				SetRole(dbquery.RoleCustomer).
				Save(db)
			if err != nil {
				log.Print(err)
			}
			httpStatus = http.StatusOK
			status = "success"
			message = "Pendaftaran berhasil"
		}

		m := gin.H{
			"message": message,
			"status":  status,
		}
		c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
	}
	return gin.HandlerFunc(fn)
}
