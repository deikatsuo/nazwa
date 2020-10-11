package api

import (
	"nazwa/dbquery"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"strconv"

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

		/*
			body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				errMess = "Data tidak benar"
				next = false
			}
		*/

		var lastid int
		limit := 10
		if next {
			lim, err := strconv.Atoi(c.Param("limit"))
			if err == nil {
				limit = lim
			}

			lst, err := strconv.Atoi(c.Query("lastid"))
			if err == nil {
				lastid = lst
			}
		}

		var total int
		if t, err := dbquery.GetTotalRow(db); err == nil {
			total = t
		}
		var users []wrapper.User
		// Gunakan offset jika tersedia
		if lastid != 0 {
			u, err := dbquery.GetAllUser(db, limit, lastid)
			if err != nil {
				errMess = "Tidak bisa mengambil data"
			}
			users = u
			httpStatus = http.StatusOK
		} else {
			u, err := dbquery.GetAllUser(db, limit)
			if err != nil {
				errMess = "Tidak bisa mengambil data"
			}
			users = u
			httpStatus = http.StatusOK
		}
		nextid := users[len(users)-1].ID

		c.JSON(httpStatus, gin.H{
			"error":  errMess,
			"users":  users,
			"total":  total,
			"lastid": nextid,
		})
	}
	return gin.HandlerFunc(fn)
}
