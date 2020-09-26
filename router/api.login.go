package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Login ...
// Struct untuk menyimpan data login
type Login struct {
	Loginid  string `json:"loginid" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// APIUserLogin ...
// API untuk login user
func APIUserLogin(c *gin.Context) {
	// Mulai session
	session := sessions.Default(c)

	// Nilai awal
	errLoginid := false
	errmLoginid := ""
	errPassword := false
	errmPassword := ""
	var status string
	statusMessage := ""
	var httpStatus int

	// Variabel untuk di bind ke Login
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		if strings.Contains(err.Error(), "Loginid") {
			errLoginid = true
			errmLoginid = "ID login harus diisi \n"
		}
		if strings.Contains(err.Error(), "Password") {
			errPassword = true
			errmPassword = "Password harus diisi \n"
		}
		httpStatus = http.StatusBadRequest
	}

	// Generate user
	// Sementara belum ada data dari database
	users := map[string]map[string]string{
		"rika@nazwa": {
			"password": "deri",
			"picture":  "../assets/img/test/teteh.jpeg",
		},
		"deri@deri": {
			"password": "rika",
			"picture":  "../assets/img/test/deri.jpeg",
		},
	}

	// Iterate users dan cocokan dengan input dari user
	var picture string
	userExist := false
	for i, v := range users {
		if json.Loginid == i && json.Password == v["password"] {
			picture = v["picture"]
			userExist = true
			break
		}
	}

	// Check apakah user ditemukan
	// Atau tidak
	if userExist == false {
		httpStatus = http.StatusUnauthorized
		statusMessage = "User tidak ditemukan"
		status = "fail"
	} else {
		// Simpan user ke session
		session.Set("loginid", json.Loginid)
		session.Set("picture", picture)
		if err := session.Save(); err != nil {
			httpStatus = http.StatusInternalServerError
			statusMessage = "Gagal membuat session"
			status = "error"
		} else {
			// User ditemukan
			// Dan berhasil diverifikasi
			httpStatus = http.StatusOK
			statusMessage = "Berhasil masuk"
			status = "success"
		}
	}

	// Kirim data ke browser klien
	c.JSON(httpStatus, gin.H{
		"message":       statusMessage,
		"status":        status,
		"err_loginid":   errLoginid,
		"errm_loginid":  errmLoginid,
		"err_password":  errPassword,
		"errm_password": errmPassword,
	})
}
