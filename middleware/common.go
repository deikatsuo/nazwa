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
			u := c.Request.URL
			u.Host = host

			// Cek apakah absolut atau tidak
			if c.Request.URL.IsAbs() == false {
				u.Scheme = "http"
			}

			c.Redirect(http.StatusMovedPermanently, u.String())
			c.Abort()
			return
		}

		c.Next()
	}
}
