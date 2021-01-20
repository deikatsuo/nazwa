package api

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/wrapper"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserLogin gerbang user login
func UserLogin(c *gin.Context) {
	// Mulai session
	session := sessions.Default(c)

	// Nilai awal
	var status string
	statusMessage := ""
	var httpStatus int
	var simpleErrMap map[string]interface{}
	var userid int
	var role string
	var userExist bool
	// Variabel untuk di bind ke Login
	var json wrapper.Login
	if err := c.ShouldBindJSON(&json); err != nil {
		// Sederhanakan pesan error
		simpleErrMap = validation.SimpleValErrMap(err)
		httpStatus = http.StatusBadRequest
	} else {
		// Periksa user di database
		if id, err := dbquery.Login(json.Loginid, json.Password); err == nil {
			userid = id
			userExist = true
		}
	}

	// Check role
	if ur, err := dbquery.UserGetRole(userid); err == nil {
		role = ur
	} else {
		httpStatus = http.StatusInternalServerError
		statusMessage = "Gagal mengambil user role"
		status = "error"
	}

	// Check apakah user ditemukan
	// Atau tidak
	if userExist == false || role == "" {
		statusMessage = "User tidak ditemukan"
		status = "fail"
	} else {
		// Simpan userid dan role ke session
		session.Set("userid", userid)
		session.Set("role", role)
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
