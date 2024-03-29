package router

import (
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"

	"github.com/gin-gonic/gin"
)

// PageDashboardInstalments halaman tagihan
func PageDashboardInstalments(c *gin.Context) {
	var zones []wrapper.LocationZone

	if z, err := dbquery.ZoneShowAll(); err == nil {
		zones = z
	} else {
		log.Warn("Terjadi kesalahan saat memuat data zona")
		log.Error(err)
	}

	gh := gin.H{
		"site_title": "Halaman Tagihan",
		"page":       "instalments",
		"zones":      zones,
		"css": []string{
			"/assets/css/print.css",
			"/assets/css/loading.css",
		},
		"js": []string{
			"/assets/js/dayjs.min.js",
			"/assets/js/dayjs.locale.id.js",
			"/assets/js/dayjs.customParseFormat.js",
			"/assets/js/dayjs.relativeTime.js",
			"/assets/js/terbilang.js",
			"/assets/js/jspm/zip.js",
			"/assets/js/jspm/zip-ext.js",
			"/assets/js/jspm/deflate.js",
			"/assets/js/jspm/JSPrintManager.js",
			"/assets/js/html2canvas.min.js",
			"/assets/js/xlsx-populate.min.js",
			"/assets/js/print.js",
		},
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.instalments.html", misc.Mete(df, gh))
}
