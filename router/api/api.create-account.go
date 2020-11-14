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
	RIC        string `json:"ric" binding:"numeric,min=16,max=16"`
	Phone      string `json:"phone" binding:"numeric,min=6,max=15"`
	Firstname  string `json:"firstname" binding:"required,min=3,max=25"`
	Lastname   string `json:"lastname" binding:"omitempty,min=1,max=25"`
	Gender     string `json:"gender" binding:"required,oneof=m f"`
	Password   string `json:"password" binding:"omitempty,alphanumunicode,min=8,max=25"`
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
		save := true

		if err := c.ShouldBindJSON(&json); err != nil {
			simpleErrMap = validation.SimpleValErrMap(err)
			httpStatus = http.StatusBadRequest
			status = "fail"
			message = "Data tidak lengkap"
			save = false
		} else {
			if dbquery.RICExist(db, json.RIC) {
				simpleErrMap["ric"] = "Nomor KTP sudah terdaftar"
				save = false
			}
			if dbquery.PhoneExist(db, json.Phone) {
				simpleErrMap["phone"] = "Nomor ini sudah terdaftar"
				save = false
			}
		}

		if save {
			user := dbquery.NewUser()
			err := user.SetFirstName(json.Firstname).
				SetLastName(json.Lastname).
				SetRIC(json.RIC).
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
		} else {
			httpStatus = http.StatusBadRequest
			status = "error"
			message = "Gagal membuat pelanggan, silahkan perbaiki formulir"
		}

		m := gin.H{
			"message": message,
			"status":  status,
		}
		c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
	}
	return gin.HandlerFunc(fn)
}
