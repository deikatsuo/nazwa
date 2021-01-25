package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardLocationsZones tampilkan halaman zona wilayah
func PageDashboardLocationsZones(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman Zona Wilayah",
		"page":       "locations_zones",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.locations.zones.html", misc.Mete(df, gh))
}
