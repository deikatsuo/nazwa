package api

import (
	"io/ioutil"
	"nazwa/dbquery"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// UsersList menampilkan semua list customer
func UsersList(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		session := sessions.Default(c)
		// User session saat ini
		uid := session.Get("userid")
		if uid == nil {
			router.Page404(c)
			return
		}

		errMess := ""
		next := true
		httpStatus := http.StatusBadRequest

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			errMess = "Data tidak benar"
			next = false
		}

		var offset int64
		var limit int64 = 10
		if next {
			v, err := jsonparser.GetInt(body, "limit")
			if err != nil {
				limit = v
			}
			v, err = jsonparser.GetInt(body, "offset")
			if err != nil {
				offset = v
			}
		}

		var users []wrapper.User
		// Gunakan offset jika tersedia
		if offset != 0 {
			u, err := dbquery.GetAllUser(db, limit, offset)
			if err != nil {
				errMess = "Tidak bisa mengambil data"
			}
			users = u
		} else {
			u, err := dbquery.GetAllUser(db, limit)
			if err != nil {
				errMess = "Tidak bisa mengambil data"
			}
			users = u
		}

		c.JSON(httpStatus, gin.H{
			"error": errMess,
			"users": users,
		})
	}
	return gin.HandlerFunc(fn)
}
