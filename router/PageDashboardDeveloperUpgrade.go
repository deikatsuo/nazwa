package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardDeveloperUpgrade halaman upgrade system
func PageDashboardDeveloperUpgrade(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Peningkatan Sistem",
		"page":       "developer_upgrade",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.developer.upgrade.html", misc.Mete(df, gh))
}
