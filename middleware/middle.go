package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RedirectWWW redirect www ke non-www
func RedirectWWW() gin.HandlerFunc {
	return func(c *gin.Context) {
		if host := strings.TrimPrefix(c.Request.Host, "www."); host != c.Request.Host {
			c.Redirect(http.StatusMovedPermanently, host+c.Request.RequestURI)
			c.Abort()
			return
		}

		c.Next()
	}
}
