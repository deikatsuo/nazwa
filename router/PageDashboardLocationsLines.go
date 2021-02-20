package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardLocationsLines halaman lines/arah
func PageDashboardLocationsLines(c *gin.Context) {
	gh := gin.H{
		"site_title": "Arah Tagihan",
		"page":       "locations_lines",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.locations.lines.html", misc.Mete(df, gh))
}
