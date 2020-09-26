package router

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PageLogout ...
// Halaman logout akun
func PageLogout(c *gin.Context) {
	// Ambil session
	session := sessions.Default(c)
	loginid := session.Get("loginid")
	picture := session.Get("picture")

	// Cek session
	if loginid == nil || picture == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token session salah",
		})
		return
	}

	// Menghapus session
	session.Delete("loginid")
	session.Delete("picture")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus session",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}
