package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Login ...
// Struct untuk menyimpan data login
type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// APIUserLogin ...
// API untuk login user
func APIUserLogin(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		var message string
		if strings.Contains(err.Error(), "Email") {
			message = "Email harus diisi \n"
		}
		if strings.Contains(err.Error(), "Password") {
			message = message + "Password harus diisi \n"
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": message,
			"error":   true,
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
			"picture":  "..assets/img/test/deri.jpg",
		},
	}

	for i, v := range users {
		if json.Email == i && json.Password == v["password"] {
			break
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User tidak ditemukan",
				"fail":    true,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil masuk",
		"success": true,
	})
}
