package api

import (
	"log"
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

		// Id terakhir
		var lastid int
		// Total data yang sudah diload
		var loaded int
		limit := 10
		if next {
			lim, err := strconv.Atoi(c.Param("limit"))
			if err == nil {
				limit = lim
			}

			// Ambil id terakhir
			lst, err := strconv.Atoi(c.Query("lastid"))
			if err == nil {
				lastid = lst
			}

			// Total yang sudah diload
			lod, err := strconv.Atoi(c.Query("loaded"))
			if err == nil {
				loaded = lod
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

		// ID terakhir yang diambil database
		if len(users) > 0 {
			log.Print(len(users) - 1)
			lastid = users[len(users)-1].ID
		}

		if limit < len(users) {
			// Total user yang sudah di load
			loaded = loaded + limit
		} else {
			loaded = len(users)
		}

		c.JSON(httpStatus, gin.H{
			"error":  errMess,
			"users":  users,
			"total":  total,
			"lastid": lastid,
			"loaded": loaded,
		})
	}
	return gin.HandlerFunc(fn)
}
