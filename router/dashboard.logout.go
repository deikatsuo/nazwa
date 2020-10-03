package router

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PageDashboardLogout menghapus session userid
func PageDashboardLogout(c *gin.Context) {
	// Ambil session
	session := sessions.Default(c)
	userid := session.Get("userid")

	// Cek session
	if userid == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token session salah",
		})
		return
	}

	// Menghapus session
	session.Delete("userid")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus session",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}