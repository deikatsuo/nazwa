package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardLocations tampilkan halaman zona wilayah
func PageDashboardLocations(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Lokasi",
		"page":       "locations",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.locations.html", misc.Mete(df, gh))
}
