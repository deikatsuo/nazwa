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
	Policy     bool   `json:"policy" binding:"required"`
}

// APIUserCreate ...
// API untuk membuat user baru
func APIUserCreate(c *gin.Context) {
	var json FormUser

	var status string
	var httpStatus int
	message := ""
	errmPhone := ""
	errmFirstname := ""
	errmPassword := ""
	errmRepassword := ""
	errmPolicy := ""
	if err := c.ShouldBindJSON(&json); err != nil {
		if strings.Contains(err.Error(), "Phone") {
			errmPhone = "Phone harus diisi \n"
		}
		if strings.Contains(err.Error(), "Firstname") {
			errmFirstname = "Nama depan harus diisi \n"
		}
		if strings.Contains(err.Error(), "Password") {
			errmPassword = "Password harus diisi \n"
		}
		if strings.Contains(err.Error(), "Repassword") {
			errmRepassword = "Ulangi password harus diisi \n"
		}
		if strings.Contains(err.Error(), "Policy") {
			errmPolicy = "Kebijakan privasi harus dicheck \n"
		}
		httpStatus = http.StatusBadRequest
		status = "fail"
		message = "Data tidak lengkap"
	}

	c.JSON(httpStatus, gin.H{
		"message":         message,
		"status":          status,
		"errm_phone":      errmPhone,
		"errm_firstname":  errmFirstname,
		"errm_password":   errmPassword,
		"errm_repassword": errmRepassword,
		"errm_policy":     errmPolicy,
	})
}
