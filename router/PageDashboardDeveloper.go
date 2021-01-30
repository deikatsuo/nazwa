package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardDeveloper halaman developer
func PageDashboardDeveloper(c *gin.Context) {
	gh := gin.H{
		"site_title": "Pengembang",
		"page":       "developer",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.developer.html", misc.Mete(df, gh))
}
