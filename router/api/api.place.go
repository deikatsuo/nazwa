package api

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//////////////////////////////// GET //

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

///////////////////////////////// ADD //

// PlaceAddProvince tambah provinsi
func PlaceAddProvince(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	var newProvince wrapper.PlaceNewProvince
	if err := c.ShouldBindJSON(&newProvince); err != nil {
		simpleErrMap = validation.SimpleValErrMap(err)
		next = false
	}

	// Provinsi baru
	if next {
		if err := dbquery.PlaceNewProvince(newProvince.Province, nowID.(int)); err != nil {
			message = "Gagal menambahkan provinsi"
			status = "error"
			next = false
		} else {
			message = fmt.Sprintf("Provinsi %s berhasil ditambahkan", newProvince.Province)
			status = "success"
			next = true
		}
	}

	// Ambil provinsi
	var provinces []wrapper.Place
	if next {
		p, err := dbquery.PlaceGetProvinces()
		if err != nil {
			message = "Gagal mengambil data provinsi"
			status = "error"
		} else {
			provinces = p
		}
	}

	m := gin.H{
		"message":   message,
		"status":    status,
		"provinces": provinces,
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}
