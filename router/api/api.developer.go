package api

import (
	"nazwa/misc"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// DeveloperUpgradeUpload API untuk upload file server
func DeveloperUpgradeUpload(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	// File
	file, err := c.FormFile("file")

	if err != nil {
		log.Warn("api.developer.go DeveloperUpgradeUpload() File tidak valid")
		log.Error(err)
		message = "Tidak ada file, atau format file tidak valid"
		status = "error"
		next = false
	}

	if next {
		if err := os.Mkdir("../data/upload/upgrade", 0755); err != nil {
			if os.IsExist(err) {
				log.Warn("api.developer.go DeveloperUpgradeUpload() direktori sudah ada")
				log.Error(err)
			} else {
				log.Warn("api.developer.go DeveloperUpgradeUpload()  gagal membuat direktori upgrade")
				log.Error(err)
				message = "Terjadi kesalahan saat mencoba membuat direktori"
				status = "error"
				next = false
			}
		}

	}

	if next {
		// Simpan file
		path := "../data/upload/upgrade/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			log.Warn("api.developer.go DeveloperUpgradeUpload() Gagal menyimpan file")
			log.Error(err)
			message = "Terjadi kesalahan saat mencoba menyimpan file"
			status = "error"
			next = false
		}
	}

	if next {
		message = "File berhasil disimpan"
		status = "success"
		httpStatus = http.StatusOK
	}

	m := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}
