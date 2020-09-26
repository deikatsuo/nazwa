package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type FormUser struct {
	Phone      string `json:"phone" binding:"required"`
	Firstname  string `json:"firstname" binding:"required"`
	Lastname   string `json:"lastname"`
	Password   string `json:"password" binding:"required"`
	Repassword string `json:"repassword" binding:"required"`
	Policy     string `json:"policy" binding:"required"`
}

// APIUserCreate ...
// API untuk membuat user baru
func APIUserCreate(c *gin.Context) {
	var json FormUser
	message := ""
	if err := c.ShouldBindJSON(&json); err != nil {
		if strings.Contains(err.Error(), "Phone") {
			message = message + "Phone harus diisi \n"
		}
		if strings.Contains(err.Error(), "Firstname") {
			message = message + "Nama depan harus diisi \n"
		}
		if strings.Contains(err.Error(), "Password") {
			message = message + "Password harus diisi \n"
		}
		if strings.Contains(err.Error(), "Repassword") {
			message = message + "Ulangi password harus diisi \n"
		}
		if strings.Contains(err.Error(), "Policy") {
			message = message + "Kebijakan privasi harus dicheck \n"
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"success": true,
	})
}
