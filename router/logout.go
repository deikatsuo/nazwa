package router

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PageLogout menghapus session userid
func PageLogout(c *gin.Context) {
	// Ambil session
	session := sessions.Default(c)
	userid := session.Get("userid")

	// Cek session
	if userid == nil {
		Page403(c)
		return
	}

	// Menghapus session
	session.Delete("userid")
	session.Delete("role")
	session.Delete("user")

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus session",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}
