package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PageDashboardUsersAdd tampilkan halaman tambah user
func PageDashboardUsersAdd(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		gh := gin.H{
			"site_title": "Tambah Users",
			"page":       "users_add",
		}

		df := c.MustGet("dashboard").(map[string]interface{})
		c.HTML(200, "dashboard.users.add.html", misc.Mete(df, gh))
	}
	return gin.HandlerFunc(fn)
}
