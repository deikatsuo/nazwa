package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboard Halaman dashboard utama
func PageDashboard(c *gin.Context) {
	gh := gin.H{
		"site_title": "Dashboard",
		"page":       "dashboard",
		"chart":      true,
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.html", misc.Mete(df, gh))
}
