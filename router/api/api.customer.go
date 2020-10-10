package api

import (
	"nazwa/wrapper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// CustomerList menampilkan semua list customer
func CustomerList(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var customers wrapper.NullableUser
		c.JSON(http.StatusOK, gin.H{
			"error":     "",
			"customers": customers,
		})
	}
	return gin.HandlerFunc(fn)
}
