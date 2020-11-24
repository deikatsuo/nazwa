package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PageDashboardUsersAdd tampilkan halaman tambah user
func PageDashboardUsersAdd(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		gh := gin.H{
			"site_title": "Tambah Users",
			"page":       "users_add",

			"l_account_user_contact": "Informasi Pengguna",
			"l_account_user_address": "Informasi Alamat",

			// Formulir buat akun
			"l_c_form_title":           "Buat akun",
			"l_c_form_fc":              "Nomor KK",
			"l_c_form_ric":             "Nomor KTP",
			"l_c_form_phone":           "Nomor Hp",
			"l_c_form_firstname":       "Nama depan",
			"l_c_form_lastname":        "Nama belakang",
			"l_c_form_gender":          "Jenis kelamin",
			"l_c_form_gender_m":        "Laki-laki",
			"l_c_form_gender_f":        "Perempuan",
			"l_c_form_occupation":      "Pekerjaan",
			"l_c_form_add_password":    "Tambahkan kata sandi",
			"l_c_form_remove_password": "Hapus kata sandi",
			"l_c_form_password":        "Kata sandi",
			"l_c_form_repassword":      "Ulangi kata sandi",
			"l_c_form_agree":           "Saya setuju dengan",
			"l_c_form_privacy_link":    "kebijakan privasi",
			"l_c_form_create":          "Buat user",

			// Detail alamat
			"l_u_address_delete_btn":    "Hapus",
			"l_u_address_show_one":      "Alamat 1",
			"l_u_address_show_two":      "Alamat 2",
			"l_u_address_show_zip":      "Kode Pos",
			"l_u_address_show_village":  "Kelurahan/Desa",
			"l_u_address_show_district": "Distrik/Kecamatan",
			"l_u_address_show_city":     "Kota/Kabupaten",
			"l_u_address_show_province": "Provinsi",

			// Detail tambah alamat
			"l_u_address_add":         "Buat alamat baru",
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
		c.HTML(200, "dashboard.users.add.html", misc.Mete(df, gh))
	}
	return gin.HandlerFunc(fn)
}
