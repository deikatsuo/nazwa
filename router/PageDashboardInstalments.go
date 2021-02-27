package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardInstalments halaman tagihan
func PageDashboardInstalments(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Tagihan",
		"page":       "instalments",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.instalments.html", misc.Mete(df, gh))
}
