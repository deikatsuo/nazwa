package api

import (
	"nazwa/dbquery"
	"nazwa/router"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ProductList mengambil data provinsi
func ProductList(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		p, err := dbquery.GetAllProduct(db)
		if err != nil {
			router.Page500(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"product": p,
		})
	}
	return gin.HandlerFunc(fn)
}
