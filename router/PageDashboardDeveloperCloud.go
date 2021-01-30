package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardDeveloperCloud halaman penyimpanan awan
func PageDashboardDeveloperCloud(c *gin.Context) {
	gh := gin.H{
		"site_title": "Penyimpanan Awan",
		"page":       "developer_cloud",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.developer.cloud.html", misc.Mete(df, gh))
}
