package api

import (
	"io/ioutil"
	"nazwa/dbquery"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// UserDeleteEmail api untuk menghapus email
func UserDeleteEmail(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			errMess = "Data tidak benar"
			next = false
		}
		id, err := jsonparser.GetInt(body, "id")
		if err != nil {
			errMess = "Request tidak valid"
		}

		// Periksa jika id peng request tidak sama dengan
		// id yang tersimpan di session
		if next {
			if nowID != uid {
				errMess = "User tidak memiliki ijin untuk menghapus email ini"
				next = false
			}
		}

		// Delete email
		if next {
			if err := dbquery.DeleteEmail(db, id, nowID.(int)); err != nil {
				errMess = "Gagal menghapus email"
				next = false
			}
		}

		// Ambil email sisa jika masih ada
		var emails []wrapper.UserEmail
		if next {
			em, err := dbquery.GetEmail(db, nowID.(int))
			if err != nil {
				errMess = "Gagal memuat email/semua email sudah dihapus"
			} else {
				emails = em
				httpStatus = http.StatusOK
				success = "Email berhasil dihapus"
			}
		}
		c.JSON(httpStatus, gin.H{
			"error":   errMess,
			"success": success,
			"emails":  emails,
		})
	}
	return gin.HandlerFunc(fn)
}
