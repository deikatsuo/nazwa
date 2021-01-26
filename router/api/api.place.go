package api

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//////////////////////////////// GET //

// PlaceProvinces mengambil data provinsi
func PlaceProvinces(c *gin.Context) {
	message := ""
	status := ""
	httpStatus := http.StatusOK
	mode := c.Query("mode")
	next := true

	if mode == "split" {

		po, err := dbquery.PlaceGetProvinces(true)
		if err != nil {
			message = "Gagal mengambil data provinsi original"
			status = "error"
			httpStatus = http.StatusBadRequest
			next = false
		}

		var provincesNotOri []wrapper.Place
		if next {
			pno, err := dbquery.PlaceGetProvinces(false)
			if err != nil {
				message = "Gagal mengambil data provinsi manual"
				status = "error"
				httpStatus = http.StatusBadRequest
			} else {
				provincesNotOri = pno
			}
		}

		c.JSON(httpStatus, gin.H{
			"message": message,
			"status":  status,
			"provinces": map[string]interface{}{
				"original": po,
				"manual":   provincesNotOri,
			},
		})
	} else {
		p, err := dbquery.PlaceGetProvinces()
		if err != nil {
			message = "Tidak ada data provinsi"
			status = "error"
			httpStatus = http.StatusBadRequest
		}

		c.JSON(httpStatus, gin.H{
			"message":  message,
			"status":   status,
			"province": p,
		})
	}
}

// PlaceCities mengambil data kota/kabupaten
func PlaceCities(c *gin.Context) {
	message := ""
	status := ""
	httpStatus := http.StatusOK
	next := true

	// ID provinsi
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "Request tidak valid"
		status = "error"
		httpStatus = http.StatusBadRequest
		next = false
	}

	mode := c.Query("mode")

	if mode == "split" {

		co, err := dbquery.PlaceGetCities(pid, true)
		if err != nil {
			message = "Gagal mengambil data kota/kabupaten original"
			status = "error"
			httpStatus = http.StatusBadRequest
			next = false
		}

		var citiesNotOri []wrapper.Place
		if next {
			cno, err := dbquery.PlaceGetCities(pid, false)
			if err != nil {
				message = "Gagal mengambil data kota/kabupaten manual"
				status = "error"
				httpStatus = http.StatusBadRequest
			} else {
				citiesNotOri = cno
			}
		}

		c.JSON(httpStatus, gin.H{
			"message": message,
			"status":  status,
			"cities": map[string]interface{}{
				"original": co,
				"manual":   citiesNotOri,
			},
		})
	} else {
		var cities []wrapper.Place

		if next {
			p, err := dbquery.PlaceGetCities(pid)
			if err != nil {
				message = "Gagal mengambil data kota/kabupaten"
				status = "error"
				httpStatus = http.StatusBadRequest
				next = false
			} else {
				cities = p
			}
		}

		if len(cities) == 0 {
			message = "Tidak ada data kota/kabupaten"
			status = "error"
			next = false
		}

		c.JSON(httpStatus, gin.H{
			"message": message,
			"status":  status,
			"cities":  cities,
		})
	}
}

// PlaceDistricts mengambil data distrik/kecamatan
func PlaceDistricts(c *gin.Context) {
	message := ""
	status := ""
	httpStatus := http.StatusOK
	next := true

	// ID kota
	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "Request tidak valid"
		status = "error"
		httpStatus = http.StatusBadRequest
		next = false
	}

	mode := c.Query("mode")

	if mode == "split" {

		do, err := dbquery.PlaceGetDistricts(cid, true)
		if err != nil {
			message = "Gagal mengambil data distrik/kecamatan original"
			status = "error"
			httpStatus = http.StatusBadRequest
			next = false
		}

		var districtsNotOri []wrapper.Place
		if next {
			dno, err := dbquery.PlaceGetDistricts(cid, false)
			if err != nil {
				message = "Gagal mengambil data distrik/kecamatan manual"
				status = "error"
				httpStatus = http.StatusBadRequest
			} else {
				districtsNotOri = dno
			}
		}

		c.JSON(httpStatus, gin.H{
			"message": message,
			"status":  status,
			"districts": map[string]interface{}{
				"original": do,
				"manual":   districtsNotOri,
			},
		})
	} else {
		var districts []wrapper.Place
		if next {
			p, err := dbquery.PlaceGetDistricts(cid)
			if err != nil {
				message = "Gagal mengambil data distrik/kecamatan"
				status = "error"
				httpStatus = http.StatusBadRequest
				next = false
			} else {
				districts = p
			}
		}

		if len(districts) == 0 {
			message = "Tidak ada data distrik/kecamatan"
			status = "error"
			next = false
		}

		c.JSON(httpStatus, gin.H{
			"message":   message,
			"status":    status,
			"districts": districts,
		})
	}
}

// PlaceVillages mengambil data kelurahan/desa
func PlaceVillages(c *gin.Context) {
	message := ""
	status := ""
	httpStatus := http.StatusOK
	next := true

	// ID distrik/kecamatan
	did, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "Request tidak valid"
		status = "error"
		httpStatus = http.StatusBadRequest
		next = false
	}

	mode := c.Query("mode")

	if mode == "split" {

		vo, err := dbquery.PlaceGetVillages(did, true)
		if err != nil {
			message = "Gagal mengambil data kelurahan/desa original"
			status = "error"
			httpStatus = http.StatusBadRequest
			next = false
		}

		var villagesNotOri []wrapper.Place
		if next {
			vno, err := dbquery.PlaceGetVillages(did, false)
			if err != nil {
				message = "Gagal mengambil data kelurahan/desa manual"
				status = "error"
				httpStatus = http.StatusBadRequest
			} else {
				villagesNotOri = vno
			}
		}

		c.JSON(httpStatus, gin.H{
			"message": message,
			"status":  status,
			"villages": map[string]interface{}{
				"original": vo,
				"manual":   villagesNotOri,
			},
		})
	} else {
		var villages []wrapper.Place
		if next {
			p, err := dbquery.PlaceGetVillages(did)
			if err != nil {
				message = "Gagal mengambil data kelurahan/desa"
				status = "error"
				httpStatus = http.StatusBadRequest
				next = false
			} else {
				villages = p
			}
		}

		if len(villages) == 0 {
			message = "Tidak ada data kelurahan/desa"
			status = "error"
			next = false
		}

		c.JSON(httpStatus, gin.H{
			"message":  message,
			"status":   status,
			"villages": villages,
		})
	}
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
