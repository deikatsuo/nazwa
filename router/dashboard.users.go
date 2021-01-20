package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardUsers tampilkan halaman user
func PageDashboardUsers(c *gin.Context) {
	gh := gin.H{
		"site_title": "Halaman User",
		"page":       "users",

		// Detail tambah alamat
		"l_u_address_add":         "Tambah",
		"l_u_address_remove":      "Batal",
		"l_u_address_name":        "Panggilan",
		"l_u_address_description": "Deskripsi",
		"l_u_address_one":         "Alamat",
		"l_u_address_two_desc":    "Bagian ini boleh dikosongkan",
		"l_u_address_zip":         "Kode pos",
		"l_u_address_province":    "Provinsi",
		"l_u_address_city":        "Kota/Kabupaten",
		"l_u_address_district":    "Distrik/Kecamatan",
		"l_u_address_village":     "Kelurahan/Desa",
		"l_u_address_add_btn":     "Tambah",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.users.html", misc.Mete(df, gh))
}
