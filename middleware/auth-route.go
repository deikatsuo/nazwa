package middleware

import (
	"log"
	"nazwa/dbquery"
	"nazwa/router"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// RoutePermission memeriksa ijin akses route
func RoutePermission(db *sqlx.DB, e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userid := session.Get("userid")

		// Role default adalah guest
		role := "guest"
		// Periksa jika userid kosong
		if userid != nil {
			if ur, err := dbquery.UserGetRole(db, userid.(int)); err == nil {
				role = ur
			} else {
				log.Print(err)
			}
		}

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
