package middleware

import (
	"log"
	"nazwa/router"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// RoutePermission memeriksa ijin akses route
func RoutePermission(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//session := sessions.Default(c)
		//userid := session.Get("userid")

		// Role default adalah guest
		role := "guest"

		// Cek "ijin" untuk role
		res, err := e.Enforce("role::"+role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			log.Fatal(err)
		}

		// Jika tidak memiliki ijin untuk mengakses halaman
		if !res {
			router.Page403(c)
			c.Abort()
		}

		c.Next()
	}
}
