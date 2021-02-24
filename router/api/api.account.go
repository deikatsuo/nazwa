package api

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// TODO: Pindahkan ini
type updateContact struct {
	Password    string `json:"password" binding:"omitempty,alphanumunicode,min=8,max=25"`
	Repassword  string `json:"repassword" binding:"eqfield=Password"`
	Oldpassword string `json:"oldpassword" binding:"required_with=Password"`
}

// AccountUpdateContact api untuk mengupdate kontak user
func AccountUpdateContact(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""
	var simpleErr map[string]interface{}

	var update updateContact
	if err := c.ShouldBindJSON(&update); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
	}

	// Cocokan password lama
	if next {
		if update.Password != "" {
			if !dbquery.UserMatchPassword(uid, update.Oldpassword) {
				errMess = "Kata sandi lama salah"
				next = false
			}
		}
	}

	// Update password
	if next {
		if update.Password != "" {
			if err := dbquery.UserUpdatePassword(uid, update.Password); err != nil {
				errMess = "Gagal merubah kata sandi"
				next = false
			}
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		success = "Data berhasil disimpan"
	}

	gh := gin.H{
		"error":   errMess,
		"success": success,
	}
	c.JSON(httpStatus, misc.Mete(gh, simpleErr))
}
