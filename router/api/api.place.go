package api

import (
	"nazwa/dbquery"
	"nazwa/router"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PlaceProvinces mengambil data provinsi
func PlaceProvinces(c *gin.Context) {
	p, err := dbquery.PlaceGetProvinces()
	if err != nil {
		router.Page500(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"province": p,
	})
}

// PlaceCities mengambil data kota/kabupaten
func PlaceCities(c *gin.Context) {
	// ID provinsi
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	p, err := dbquery.PlaceGetCities(pid)
	if err != nil {
		log.Print(err)
		router.Page500(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"cities": p,
	})
}

// PlaceDistricts mengambil data distrik/kecamatan
func PlaceDistricts(c *gin.Context) {
	// ID kota
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	p, err := dbquery.PlaceGetDistricts(pid)
	if err != nil {
		log.Print(err)
		router.Page500(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"districts": p,
	})
}

// PlaceVillages mengambil data kelurahan/desa
func PlaceVillages(c *gin.Context) {
	// ID kota
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.Page404(c)
		return
	}

	p, err := dbquery.PlaceGetVillages(pid)
	if err != nil {
		log.Print(err)
		router.Page500(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"villages": p,
	})
}
