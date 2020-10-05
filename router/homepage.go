package router

import (
	"nazwa/misc"
	"nazwa/wrapper"

	"github.com/gin-gonic/gin"
)

// PageHome ...
// Halaman homepage
func PageHome(c *gin.Context) {
	// Ambil konfigurasi default
	df := c.MustGet("config").(wrapper.DefaultConfig).Site

	gh := gin.H{
		"site_title": "Homepage",

		// Navbar
		"l_nav_home":  "Home",
		"l_nav_login": "Masuk",

		// Header
		"l_h_only":   "Apa barang impianmu?",
		"l_h_slogan": df["site_name"].(string) + " hanya menyediakan barang berkwalitas!",
		"l_h_desc":   "Ayo pesan barangmu sekarang, bisa kredit juga cash!",

		// Kenapa
		"l_k_title": "Kenapa " + df["site_name"].(string) + "?",

		// Bagian hubungi kami
		"l_call_act":      "Hubungi kami",
		"l_call_act_desc": "Siap melayani selama jam kerja",
		"l_call_act_btn":  "Hubungi!",

		// Bagian footer perusahaan
		"l_company":         "Perusahaan",
		"l_company_blog":    "Blog",
		"l_company_about":   "Tentang kami",
		"l_company_contact": "Hubungi kami",

		"l_social": "Sosial media",

		// Bagian footer hukum
		"l_legal":         "Hukum",
		"l_legal_term":    "Persyaratan",
		"l_legal_privacy": "Privasi",

		// Bagian footer link
		"l_link":         "Link",
		"l_link_help":    "Bantuan",
		"l_link_support": "Dukungan",
	}
	c.HTML(200, "index.html", misc.Mete(df, gh))
}
