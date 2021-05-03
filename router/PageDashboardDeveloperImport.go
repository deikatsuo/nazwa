package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardDeveloperImport halaman import data
func PageDashboardDeveloperImport(c *gin.Context) {
	gh := gin.H{
		"site_title": "Import Data File",
		"page":       "developer_import",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.developer.import.html", misc.Mete(df, gh))
}
