package router

import (
	"nazwa/misc"
	"nazwa/misc/validation"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Login struk formulir login
// menyimpan data login
type Login struct {
	Loginid  string `json:"loginid" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// APIUserLogin gerbang user login
// API untuk login user
func APIUserLogin(c *gin.Context) {
	// Mulai session
	session := sessions.Default(c)

	// Nilai awal
	var status string
	statusMessage := ""
	var httpStatus int
	var simpleErrMap map[string]interface{}

	// Variabel untuk di bind ke Login
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		// Sederhanakan pesan error
		simpleErrMap = validation.SimpleValErrMap(err)
		httpStatus = http.StatusBadRequest
	}

	// Periksa user di database

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

	m := gin.H{
		"message": statusMessage,
		"status":  status,
	}

	// Kirim data ke browser klien
	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}
