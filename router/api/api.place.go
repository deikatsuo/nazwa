package api

import (
	"log"
	"nazwa/dbquery"
	"nazwa/router"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PlaceProvinces mengambil data provinsi
func PlaceProvinces(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		p, err := dbquery.GetProvinces(db)
		if err != nil {
			router.Page500(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"province": p,
		})
	}
	return gin.HandlerFunc(fn)
}

// PlaceCities mengambil data kota/kabupaten
func PlaceCities(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// ID provinsi
		pid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			router.Page404(c)
			return
		}

		p, err := dbquery.GetCities(db, pid)
		if err != nil {
			log.Print(err)
			router.Page500(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"cities": p,
		})
	}
	return gin.HandlerFunc(fn)
}

// PlaceDistricts mengambil data distrik/kecamatan
func PlaceDistricts(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// ID kota
		pid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			router.Page404(c)
			return
		}

		p, err := dbquery.GetDistricts(db, pid)
		if err != nil {
			log.Print(err)
			router.Page500(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"districts": p,
		})
	}
	return gin.HandlerFunc(fn)
}

// PlaceVillages mengambil data kelurahan/desa
func PlaceVillages(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// ID kota
		pid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			router.Page404(c)
			return
		}

		p, err := dbquery.GetVillages(db, pid)
		if err != nil {
			log.Print(err)
			router.Page500(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"villages": p,
		})
	}
	return gin.HandlerFunc(fn)
}
