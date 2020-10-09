package router

import (
	"nazwa/misc"

	"github.com/gin-gonic/gin"
)

// PageDashboardAccount halaman tempat user untuk mengubah data akun
func PageDashboardAccount(c *gin.Context) {
	gh := gin.H{
		"site_title":             "Akun",
		"page":                   "account",
		"l_account_header":       "Pengaturan Akun",
		"l_account_user_contact": "Informasi Pengguna",
		"l_account_user_address": "Alamat",

		"l_u_account_firstname":       "Nama depan",
		"l_u_account_lastname":        "Nama belakang",
		"l_u_account_gender":          "Jenis kelamin",
		"l_u_account_gender_m":        "Laki-laki",
		"l_u_account_gender_f":        "Perempuan",
		"l_u_account_add_phone":       "Tambah no. HP",
		"l_u_account_phone":           "No. HP",
		"l_u_account_add_email":       "Tambah email",
		"l_u_account_email":           "Email",
		"l_u_account_change_password": "Ubah kata sandi",
		"l_u_account_password":        "Kata sandi baru",
		"l_u_account_repassword":      "Ulangi sandi baru",
		"l_u_account_old_password":    "Kata sandi saat ini",
		"l_u_account_delete_btn":      "Hapus",
		"l_u_account_update_btn":      "Simpan Data",
		"l_u_account_verify_btn":      "Verifikasi",
	}

	df := c.MustGet("dashboard").(map[string]interface{})
	c.HTML(200, "dashboard.account.html", misc.Mete(df, gh))
}
