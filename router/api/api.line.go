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

// LineNew arah baru
func LineNew(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	var newLine wrapper.LocationLineNewForm
	if err := c.ShouldBindJSON(&newLine); err != nil {
		simpleErrMap = validation.SimpleValErrMap(err)
		next = false
	}

	if next {
		if dbquery.LineCodeExist(newLine.Code) {
			next = false
			message = fmt.Sprintf("Kode %s sudah terdaftar", newLine.Code)
			httpStatus = http.StatusConflict
			status = "error"
		}
	}

	// Buat Line
	if next {
		if err := dbquery.LineNew(newLine); err != nil {
			message = "Gagal membuat arah baru"
			status = "error"
			next = false
		} else {
			message = fmt.Sprintf("Arah %s (%s) berhasil dibuat", newLine.Name, newLine.Code)
			status = "success"
			next = true
			httpStatus = http.StatusOK
		}
	}

	// Ambil lines
	var lines []wrapper.LocationLine
	if next {
		if l, err := dbquery.LineShowAll(); err == nil {
			lines = l
		} else {
			log.Warn("Terjadi kesalahan saat memuat data arah")
			log.Error(err)
			httpStatus = http.StatusInternalServerError
			message = "Tidak dapat memuat arah/Tidak ada arah"
			status = "error"
		}
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"lines":   lines,
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}

// LineShowAll ambil semua data arah
func LineShowAll(c *gin.Context) {
	httpStatus := http.StatusOK
	message := ""
	status := ""
	// Ambil line
	var lines []wrapper.LocationLine

	if l, err := dbquery.LineShowAll(); err == nil {
		lines = l
	} else {
		log.Warn("Terjadi kesalahan saat memuat data arah")
		log.Error(err)
		httpStatus = http.StatusInternalServerError
		message = "Tidak dapat memuat arah/Tidak ada arah"
		status = "error"
	}

	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"lines":   lines,
	})
}

// LineDelete hapus line
func LineDelete(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	lid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	// Delete arah
	if next {
		if err := dbquery.LineDelete(lid); err != nil {
			log.Warn("api.line.go LineDelete() Gagal menghapus Arah")
			log.Error(err)
			message = "Gagal menghapus arah"
			status = "error"
			next = false
		} else {
			httpStatus = http.StatusOK
			message = "Arah telah dihapus"
			status = "success"
		}
	}

	// Ambil line
	var lines []wrapper.LocationLine
	if next {
		if l, err := dbquery.LineShowAll(); err == nil {
			lines = l
		} else {
			log.Warn("Terjadi kesalahan saat memuat data arah")
			log.Error(err)
			httpStatus = http.StatusInternalServerError
			message = "Tidak dapat memuat arah/Tidak ada arah"
			status = "error"
		}
	}

	gh := gin.H{
		"message": message,
		"status":  status,
		"lines":   lines,
	}

	c.JSON(httpStatus, gh)
}
