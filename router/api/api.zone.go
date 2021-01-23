package api

import (
	"nazwa/dbquery"
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

	var zones []wrapper.Zone

	if z, err := dbquery.ZoneShowAll(); err == nil {
		zones = z
	} else {
		log.Warn("Terjadi kesalahan saat memuat data zona")
		log.Error(err)
		httpStatus = http.StatusInternalServerError
		message = "Sepertinya telah terjadi kesalahan saat memuat data"
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
		message = "Kolektor berhasil disimpan"
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
