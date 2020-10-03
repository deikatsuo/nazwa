package api

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Login struk formulir login
// menyimpan data login
type Login struct {
	Loginid  string `json:"loginid" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// APIUserLogin gerbang user login
// API untuk login user
func APIUserLogin(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Mulai session
		session := sessions.Default(c)

		// Nilai awal
		var status string
		statusMessage := ""
		var httpStatus int
		var simpleErrMap map[string]interface{}
		var userid int
		var userExist bool
		// Variabel untuk di bind ke Login
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			// Sederhanakan pesan error
			simpleErrMap = validation.SimpleValErrMap(err)
			httpStatus = http.StatusBadRequest
		} else {
			// Periksa user di database
			if id, err := dbquery.Login(db, json.Loginid, json.Password); err == nil {
				userid = id
				userExist = true
			}
		}

		// Check apakah user ditemukan
		// Atau tidak
		if userExist == false {
			httpStatus = http.StatusUnauthorized
			statusMessage = "User tidak ditemukan"
			status = "fail"
		} else {
			// Simpan userid ke session
			session.Set("userid", userid)
			if err := session.Save(); err != nil {
				httpStatus = http.StatusInternalServerError
				statusMessage = "Gagal membuat session"
				status = "error"
			} else {
				// userid berhasil disimpan ke session
				httpStatus = http.StatusOK
				statusMessage = "Berhasil masuk"
				status = "success"
			}
		}

		m := gin.H{
			"message": statusMessage,
			"status":  status,
		}

		// Kirim data ke browser klien
		c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
	}
	return gin.HandlerFunc(fn)
}
