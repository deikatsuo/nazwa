package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

// RedirectWWW redirect www ke non-www
func RedirectWWW() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := location.Get(c)
		if host := strings.TrimPrefix(url.Host, "www."); host != url.Host {
			fmt.Println("Host", host)
			fmt.Println("URL Host", url.Host)
			url.Host = host
			c.Redirect(http.StatusMovedPermanently, url.String())
			c.Abort()
		}

		c.Next()
	}
}
