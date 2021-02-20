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

// ZoneGetList Tampilkan data zona
func ZoneGetList(c *gin.Context) {
	httpStatus := http.StatusOK
	message := ""
	status := "success"

	var zones []wrapper.LocationZone

	if z, err := dbquery.ZoneShowAll(); err == nil {
		zones = z
	} else {
		log.Warn("Terjadi kesalahan saat memuat data zona")
		log.Error(err)
		httpStatus = http.StatusInternalServerError
		message = "Tidak dapat memuat zona/Tidak ada zona"
		status = "error"
	}

	c.JSON(httpStatus, gin.H{
		"zones":   zones,
		"message": message,
		"status":  status,
	})
}

// ZoneUpdateCollector api untuk mengupdate kolektor pada zone
func ZoneUpdateCollector(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	zid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	newCollector := c.Query("set")
	newCollectorID, err := strconv.Atoi(newCollector)
	if err != nil {
		message = "Request tidak valid"
		status = "error"
		next = false
	}

	// Update collector
	if next {
		if err := dbquery.ZoneUpdateCollector(zid, newCollectorID); err != nil {
			log.Warn("api.zone.go ZoneUpdateCollector() Gagal mengubah kolektor")
			log.Error(err)
			message = "Gagal mengubah kolektor"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Kolektor berhasil dirubah"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}

// ZoneDeleteCollector api untuk menghapus/mengosongkan kolektor pada zone
func ZoneDeleteCollector(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	zid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	// Delete collector
	if next {
		if err := dbquery.ZoneDeleteCollector(zid); err != nil {
			log.Warn("api.zone.go ZoneDeleteCollector() Gagal menghapus Kolektor dari Zona")
			log.Error(err)
			message = "Gagal mengosongkan kolektor"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Kolektor berhasil dikosongkan"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}

// ZoneDeleteList hapus list dari zone
func ZoneDeleteList(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	zid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	deleteList := c.Query("lid")
	deleteListID, err := strconv.Atoi(deleteList)
	if err != nil {
		message = "Request tidak valid"
		status = "error"
		next = false
	}

	// Delete list
	if next {
		if err := dbquery.ZoneDeleteList(zid, deleteListID); err != nil {
			log.Warn("api.zone.go ZoneDeleteList() Gagal menghapus List dari Zona")
			log.Error(err)
			message = "Gagal menghapus wilayah dari zona"
			status = "error"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		message = "Wilayah berhasil dikeluarkan dari zona"
		status = "success"
	}

	gh := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, gh)
}

// ZoneDelete hapus zona
func ZoneDelete(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	zid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	// Delete collector
	if next {
		if err := dbquery.ZoneDelete(zid); err != nil {
			log.Warn("api.zone.go ZoneDelete() Gagal menghapus Zona")
			log.Error(err)
			message = "Gagal menghapus zona"
			status = "error"
			next = false
		} else {
			httpStatus = http.StatusOK
			message = "Zona telah dihapus"
			status = "success"
		}
	}

	// Ambil zone
	var zones []wrapper.LocationZone
	if next {
		if z, err := dbquery.ZoneShowAll(); err == nil {
			zones = z
		} else {
			log.Warn("Terjadi kesalahan saat memuat data zona")
			log.Error(err)
			httpStatus = http.StatusInternalServerError
			message = "Tidak dapat memuat zona/Tidak ada zona"
			status = "error"
		}

	}

	gh := gin.H{
		"message": message,
		"status":  status,
		"zones":   zones,
	}

	c.JSON(httpStatus, gh)
}

// ZoneAddList api untuk menambah list ke zona
func ZoneAddList(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	// id zone
	zid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	var lists wrapper.LocationZoneAddListForm
	if err := c.ShouldBindJSON(&lists); err != nil {
		log.Warn("Gagal unmarshal json")
		log.Error(err)

		message = "Request tidak valid"
		status = "error"
		next = false
	}

	var ret wrapper.NameIDNameID
	if next {
		for _, lid := range lists.Lists {
			if dbquery.ZoneListExistsAndRet(lid, &ret) {
				message = fmt.Sprintf("%s Sudah terdaftar di %s", ret.NameOne, ret.NameTwo)
				status = "error"
				next = false
				break
			}
		}
	}

	// Tambahkan list
	if next {
		if err := dbquery.ZoneAddList(zid, lists); err != nil {
			message = "Gagal menambahkan wilayah kedalam zona"
			status = "error"
			next = false
		} else {
			message = "Wilayah ditambahkan kedalam zona"
			status = "success"
			next = true
			httpStatus = http.StatusOK
		}
	}

	// Ambil list wilayah
	var newLists []wrapper.LocationZoneListsSelect
	if zl, err := dbquery.ZoneShowZoneList(zid); err == nil {
		newLists = zl
	} else {
		log.Warn("api.zone.go ZoneAddList() ambil data list zona")
		log.Error(err)
	}

	c.JSON(httpStatus, gin.H{
		"message": message,
		"status":  status,
		"lists":   newLists,
	})
}

// ZoneNewZone buat zona baru
func ZoneNewZone(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	var newZone wrapper.LocationZoneNewForm
	if err := c.ShouldBindJSON(&newZone); err != nil {
		simpleErrMap = validation.SimpleValErrMap(err)
		next = false
	}

	// Buat Zone
	if next {
		if err := dbquery.ZoneNew(newZone.Zone, nowID.(int)); err != nil {
			message = "Gagal membuat zona baru"
			status = "error"
			next = false
		} else {
			message = fmt.Sprintf("Zona %s berhasil dibuat", newZone.Zone)
			status = "success"
			next = true
		}
	}

	// Ambil zone
	var zones []wrapper.LocationZone
	if next {
		if z, err := dbquery.ZoneShowAll(); err == nil {
			zones = z
		} else {
			log.Warn("Terjadi kesalahan saat memuat data zona")
			log.Error(err)
			httpStatus = http.StatusInternalServerError
			message = "Tidak dapat memuat zona/Tidak ada zona"
			status = "error"
		}
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"zones":   zones,
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}
