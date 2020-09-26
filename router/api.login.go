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
	errLoginid := false
	errmLoginid := ""
	errPassword := false
	errmPassword := ""
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
		c.JSON(http.StatusBadRequest, gin.H{
			"err_loginid":   errLoginid,
			"errm_loginid":  errmLoginid,
			"err_password":  errPassword,
			"errm_password": errmPassword,
		})
		return
	}

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

	var picture string
	for i, v := range users {
		if json.Loginid == i && json.Password == v["password"] {
			picture = v["picture"]
			break
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User tidak ditemukan",
				"fail":    true,
			})
			return
		}
	}

	// Simpan user ke session
	session.Set("loginid", json.Loginid)
	session.Set("picture", picture)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal membuat session",
			"fail":    true,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil masuk",
		"success": true,
	})
}
