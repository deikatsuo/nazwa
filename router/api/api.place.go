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

// PlaceCountryAddProvince tambah provinsi
func PlaceCountryAddProvince(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	// ID country
	countryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "ID negara tidak valid"
		status = "error"
		next = false
	}

	var newProvince wrapper.PlaceNewProvince
	if err := c.ShouldBindJSON(&newProvince); err != nil {
		simpleErrMap = validation.SimpleValErrMap(err)
		next = false
	}

	// Provinsi baru
	if next {
		if err := dbquery.PlaceNewProvince(countryID, newProvince.Province, nowID.(int)); err != nil {
			message = "Gagal menambahkan provinsi"
			status = "error"
			httpStatus = http.StatusInternalServerError
			next = false
		} else {
			message = fmt.Sprintf("Provinsi %s berhasil ditambahkan", newProvince.Province)
			status = "success"
			httpStatus = http.StatusOK
			next = true
		}
	}

	// Ambil provinsi
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

	m := gin.H{
		"message": message,
		"status":  status,
		"provinces": map[string]interface{}{
			"original": po,
			"manual":   provincesNotOri,
		},
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}

// PlaceProvinceAddCity API untuk menambahkan kota
func PlaceProvinceAddCity(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	// ID provinsi
	provinceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "ID provinsi tidak valid"
		status = "error"
		next = false
	}

	var newCity wrapper.PlaceNewCity
	if err := c.ShouldBindJSON(&newCity); err != nil {
		simpleErrMap = validation.SimpleValErrMap(err)
		next = false
	}

	// Kota baru
	if next {
		if err := dbquery.PlaceNewCity(provinceID, newCity.City, nowID.(int)); err != nil {
			message = "Gagal menambahkan kota/kabupaten"
			status = "error"
			httpStatus = http.StatusInternalServerError
			next = false
		} else {
			message = fmt.Sprintf("Kota/Kabupaten %s berhasil ditambahkan", newCity.City)
			status = "success"
			httpStatus = http.StatusOK
			next = true
		}
	}

	// Ambil data kota/kabupaten
	co, err := dbquery.PlaceGetCities(provinceID, true)
	if err != nil {
		message = "Gagal mengambil data kota/kabupaten original"
		status = "error"
		httpStatus = http.StatusBadRequest
		next = false
	}

	var citiesNotOri []wrapper.Place
	if next {
		cno, err := dbquery.PlaceGetCities(provinceID, false)
		if err != nil {
			message = "Gagal mengambil data kota/kabupaten manual"
			status = "error"
			httpStatus = http.StatusBadRequest
		} else {
			citiesNotOri = cno
		}
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"cities": map[string]interface{}{
			"original": co,
			"manual":   citiesNotOri,
		},
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}

// PlaceCityAddDistrict API untuk menambahkan distrik
func PlaceCityAddDistrict(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	// ID kota
	cityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "ID kota/kabupaten tidak valid"
		status = "error"
		next = false
	}

	var newDistrict wrapper.PlaceNewDistrict
	if err := c.ShouldBindJSON(&newDistrict); err != nil {
		simpleErrMap = validation.SimpleValErrMap(err)
		next = false
	}

	// Distrik baru
	if next {
		if err := dbquery.PlaceNewDistrict(cityID, newDistrict.District, nowID.(int)); err != nil {
			message = "Gagal menambahkan distrik/kecamatan"
			status = "error"
			httpStatus = http.StatusInternalServerError
			next = false
		} else {
			message = fmt.Sprintf("Distrik/Kecamatan %s berhasil ditambahkan", newDistrict.District)
			status = "success"
			httpStatus = http.StatusOK
			next = true
		}
	}

	// Ambil data distrik/kecamatan
	do, err := dbquery.PlaceGetDistricts(cityID, true)
	if err != nil {
		message = "Gagal mengambil data distrik/kecamatan original"
		status = "error"
		httpStatus = http.StatusBadRequest
		next = false
	}

	var districtsNotOri []wrapper.Place
	if next {
		dno, err := dbquery.PlaceGetDistricts(cityID, false)
		if err != nil {
			message = "Gagal mengambil data distrik/kecamatan manual"
			status = "error"
			httpStatus = http.StatusBadRequest
		} else {
			districtsNotOri = dno
		}
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"districts": map[string]interface{}{
			"original": do,
			"manual":   districtsNotOri,
		},
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}

// PlaceDistrictAddVillage API untuk menambahkan kelurahan/desa
func PlaceDistrictAddVillage(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	// ID distrik
	districtID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "ID distrik/kecamatan tidak valid"
		status = "error"
		next = false
	}

	var newVillage wrapper.PlaceNewVillage
	if err := c.ShouldBindJSON(&newVillage); err != nil {
		simpleErrMap = validation.SimpleValErrMap(err)
		next = false
	}

	// Kelurahan/desa baru
	if next {
		if err := dbquery.PlaceNewVillage(districtID, newVillage.Village, nowID.(int)); err != nil {
			message = "Gagal menambahkan kelurahan/desa"
			status = "error"
			httpStatus = http.StatusInternalServerError
			next = false
		} else {
			message = fmt.Sprintf("Kelurahan/Desa %s berhasil ditambahkan", newVillage.Village)
			status = "success"
			httpStatus = http.StatusOK
			next = true
		}
	}

	// Ambil data kelurahan/desa
	vo, err := dbquery.PlaceGetVillages(districtID, true)
	if err != nil {
		message = "Gagal mengambil data kelurahan/desa original"
		status = "error"
		httpStatus = http.StatusBadRequest
		next = false
	}

	var villagesNotOri []wrapper.Place
	if next {
		vno, err := dbquery.PlaceGetVillages(districtID, false)
		if err != nil {
			message = "Gagal mengambil data kelurahan/desa manual"
			status = "error"
			httpStatus = http.StatusBadRequest
		} else {
			villagesNotOri = vno
		}
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"villages": map[string]interface{}{
			"original": vo,
			"manual":   villagesNotOri,
		},
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}

// PlaceCountryDeleteProvinceByID API untuk menghapus provinsi
func PlaceCountryDeleteProvinceByID(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	// id negara
	countryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "ID negara tidak valid"
		status = "error"
		next = false
	}

	// id provinsi
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		message = "Data tidak valid, tidak bisa menghapus provinsi"
		status = "error"
		next = false
	}

	// Delete provinsi
	if next {
		if err := dbquery.PlaceProvinceDeleteByID(countryID, pid); err != nil {
			message = "Gagal menghapus provinsi"
			status = "error"
		} else {
			httpStatus = http.StatusOK
			message = "Provinsi berhasil dihapus"
			status = "success"
		}
	}

	var provincesNotOri []wrapper.Place
	pno, err := dbquery.PlaceGetProvinces(false)
	if err != nil {
		message = "Gagal mengambil data provinsi manual"
		status = "error"
		httpStatus = http.StatusBadRequest
	} else {
		provincesNotOri = pno
	}

	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"provinces": map[string]interface{}{
			"manual": provincesNotOri,
		},
	})
}

// PlaceProvinceDeleteCityByID API untuk menghapus provinsi
func PlaceProvinceDeleteCityByID(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	// id provinsi
	provinceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "ID provinsi tidak valid"
		status = "error"
		next = false
	}

	// id kota
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		message = "Data tidak valid, tidak bisa menghapus kota/kabupaten"
		status = "error"
		next = false
	}

	// Delete kota
	if next {
		if err := dbquery.PlaceCityDeleteByID(provinceID, cid); err != nil {
			message = "Gagal menghapus kota/kabupaten"
			status = "error"
		} else {
			httpStatus = http.StatusOK
			message = "Kota/Kabupaten berhasil dihapus"
			status = "success"
		}
	}

	var citiesNotOri []wrapper.Place
	cno, err := dbquery.PlaceGetCities(provinceID, false)
	if err != nil {
		message = "Gagal mengambil data kota/kabupaten manual"
		status = "error"
		httpStatus = http.StatusBadRequest
	} else {
		citiesNotOri = cno
	}

	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"cities": map[string]interface{}{
			"manual": citiesNotOri,
		},
	})
}

// PlaceCityDeleteDistrictByID API delete distrik
func PlaceCityDeleteDistrictByID(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	// id kota
	cityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "ID kota tidak valid"
		status = "error"
		next = false
	}

	// id distrik
	did, err := strconv.Atoi(c.Param("did"))
	if err != nil {
		message = "Data tidak valid, tidak bisa menghapus distrik/kecamatan"
		status = "error"
		next = false
	}

	// Delete distrik
	if next {
		if err := dbquery.PlaceDistrictDeleteByID(cityID, did); err != nil {
			message = "Gagal menghapus distrik/kecamatan"
			status = "error"
		} else {
			httpStatus = http.StatusOK
			message = "Distrik/Kecamatan berhasil dihapus"
			status = "success"
		}
	}

	var districtsNotOri []wrapper.Place
	dno, err := dbquery.PlaceGetDistricts(cityID, false)
	if err != nil {
		message = "Gagal mengambil data distrik/kecamatan manual"
		status = "error"
		httpStatus = http.StatusBadRequest
	} else {
		districtsNotOri = dno
	}

	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"districts": map[string]interface{}{
			"manual": districtsNotOri,
		},
	})
}

// PlaceDistrictDeleteVillageByID API delete vilage
func PlaceDistrictDeleteVillageByID(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	// id distrik
	districtID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message = "ID distrik tidak valid"
		status = "error"
		next = false
	}

	// id kelurahan
	vid, err := strconv.Atoi(c.Param("vid"))
	if err != nil {
		message = "Data tidak valid, tidak bisa menghapus kelurahan/desa"
		status = "error"
		next = false
	}

	// Delete kelurahan
	if next {
		if err := dbquery.PlaceVillageDeleteByID(districtID, vid); err != nil {
			message = "Gagal menghapus kelurahan/desa"
			status = "error"
		} else {
			httpStatus = http.StatusOK
			message = "Kelurahan/Desa berhasil dihapus"
			status = "success"
		}
	}

	var villagesNotOri []wrapper.Place
	vno, err := dbquery.PlaceGetVillages(districtID, false)
	if err != nil {
		message = "Gagal mengambil data kelurahan/desa manual"
		status = "error"
		httpStatus = http.StatusBadRequest
	} else {
		villagesNotOri = vno
	}

	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"villages": map[string]interface{}{
			"manual": villagesNotOri,
		},
	})
}
