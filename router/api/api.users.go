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
		httpStatus := http.StatusBadRequest

		// Id terakhir
		var lastid int
		// Total data yang sudah diload
		limit := 10
		var direction string
		userRole := 0

		lim, err := strconv.Atoi(c.Param("limit"))
		if err == nil {
			limit = lim
		}

		// Ambil id terakhir
		lst, err := strconv.Atoi(c.Query("lastid"))
		if err == nil {
			lastid = lst
		}

		// Filter berdasarkan tipe user
		uty, err := strconv.Atoi(c.Query("role"))
		if err == nil {
			userRole = uty
		}

		// Forward/backward
		direction = c.Query("direction")

		// Melihat total user di database
		var total int
		if t, err := dbquery.GetUserTotalRow(db); err == nil {
			total = t
		}

		var users []wrapper.User
		// Gunakan offset jika tersedia
		if lastid != 0 {
			if direction == "back" {
				u, err := dbquery.GetAllUser(db, false, userRole, limit, lastid-limit+1)
				if err != nil {
					errMess = "Tidak bisa mengambil data"
				}
				users = u
				// Reverse urutan array user
				tempUsers := make([]wrapper.User, len(users))
				in := 0
				for i := len(users) - 1; i >= 0; i-- {
					tempUsers[in] = users[i]
					in++
				}
				users = tempUsers
			} else {
				u, err := dbquery.GetAllUser(db, true, userRole, limit, lastid)
				if err != nil {
					errMess = "Tidak bisa mengambil data"
				}
				users = u
			}
			httpStatus = http.StatusOK
		} else {
			u, err := dbquery.GetAllUser(db, true, userRole, limit)
			if err != nil {
				errMess = "Tidak bisa mengambil data"
			}
			users = u
			httpStatus = http.StatusOK
		}

		// ID terakhir yang diambil database
		last := false
		if len(users) > 0 {
			lastid = users[len(users)-1].ID
			if len(users) < limit {
				lastid = (users[0].ID - 1) + limit
				last = true
			}
		}

		c.JSON(httpStatus, gin.H{
			"error":  errMess,
			"users":  users,
			"total":  total,
			"lastid": lastid,
			"last":   last,
		})
	}
	return gin.HandlerFunc(fn)
}
