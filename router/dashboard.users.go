package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PageDashboardUsers tampilkan halaman user
func PageDashboardUsers(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		gh := gin.H{
			"site_title": "Halaman User",
			"page":       "users",
		}

		df := c.MustGet("dashboard").(map[string]interface{})
		c.HTML(200, "dashboard.users.html", misc.Mete(df, gh))
	}
	return gin.HandlerFunc(fn)
}
