package api

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
)

// DeveloperUpgradeUpload API untuk upload file upgrade
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
		if err := os.Mkdir("../data/upgrade", 0755); err != nil {
			if os.IsExist(err) {
				log.Warn("api.developer.go DeveloperUpgradeUpload() direktori sudah ada")
				log.Error(err)
			} else {
				log.Warn("api.developer.go DeveloperUpgradeUpload()  gagal membuat direktori /upgrade")
				log.Error(err)
				message = "Terjadi kesalahan saat mencoba membuat direktori"
				status = "error"
				next = false
			}
		}

	}

	if next {
		// Simpan file
		path := "../data/upgrade/" + file.Filename

		if len(file.Filename) > 7 {
			if file.Filename[len(file.Filename)-7:] != ".tar.xz" {
				message = "Format file harus berupa .tar.xz"
				status = "error"
				next = false
			}
		} else {
			message = "File tidak valid"
			status = "error"
			next = false
		}

		if next {
			if err := c.SaveUploadedFile(file, path); err != nil {
				log.Warn("api.developer.go DeveloperUpgradeUpload() Gagal menyimpan file")
				log.Error(err)
				message = "Terjadi kesalahan saat mencoba menyimpan file"
				status = "error"
				next = false
			}
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

// DeveloperCloudUpload API untuk upload file ke cloud
func DeveloperCloudUpload(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	// File
	file, err := c.FormFile("file")

	if err != nil {
		log.Warn("api.developer.go DeveloperCloudUpload() File tidak valid")
		log.Error(err)
		message = "Tidak ada file, atau format file tidak valid"
		status = "error"
		next = false
	}

	if next {
		if err := os.Mkdir("../data/cloud", 0755); err != nil {
			if os.IsExist(err) {
				log.Warn("api.developer.go DeveloperCloudUpload() direktori sudah ada")
				log.Error(err)
			} else {
				log.Warn("api.developer.go DeveloperCloudUpload()  gagal membuat direktori /cloud")
				log.Error(err)
				message = "Terjadi kesalahan saat mencoba membuat direktori"
				status = "error"
				next = false
			}
		}

	}

	if next {
		// Simpan file
		path := "../data/cloud/" + file.Filename

		if err := c.SaveUploadedFile(file, path); err != nil {
			log.Warn("api.developer.go DeveloperCloudUpload() Gagal menyimpan file")
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

// DeveloperImportUpload API untuk upload file import
func DeveloperImportUpload(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	// File
	file, err := c.FormFile("file")

	if err != nil {
		log.Warn("api.developer.go DeveloperImportUpload() File tidak valid")
		log.Error(err)
		message = "Tidak ada file, atau format file tidak valid"
		status = "error"
		next = false
	}

	// Buka file
	fo, err := file.Open()
	if err != nil {
		log.Warn("api.developer.go DeveloperImportUpload() Gagal membuka file")
		log.Error(err)
		message = "Tidak dapat membuka file import"
		status = "error"
		next = false
	}
	// Tutup file diakhir fungsi
	defer fo.Close()

	// Baca file
	f, err := excelize.OpenReader(fo)
	if err != nil {
		log.Warn("api.developer.go DeveloperImportUpload() Gagal membaca file")
		log.Error(err)
		message = "Tidak dapat membaca file"
		status = "error"
		next = false
	}

	// Ambil semua baris di Sheet Pelanggan
	rows, _ := f.GetRows("Pelanggan")
	for rid, row := range rows {
		// Baca mulai dari baris ke 5
		if rid >= 4 && row[2] != "" {
			var kode string
			gender := "m"
			if row[2][0:3] == "Ibu" {
				gender = "f"
			}
			kode = row[1]

			// Pecah nama konsumen
			sname := strings.SplitN(row[2], "/", 2)

			name := strings.ReplaceAll(sname[0], "Bpk.", "")
			name = strings.ReplaceAll(name, "Bpk", "")
			name = strings.ReplaceAll(name, "Ibu.", "")
			name = strings.ReplaceAll(name, "Ibu", "")
			name = strings.TrimSpace(name)
			splitname := strings.SplitN(name, " ", 2)
			firstname := splitname[0]
			var lastname string
			if len(splitname) > 1 {
				lastname = splitname[1]
			}
			var substitutes []wrapper.OrderUserSubstituteForm

			if len(sname) > 1 {
				splitsubs := strings.Split(sname[1], "/")
				for _, sva := range splitsubs {
					sg := "f"
					sn := strings.TrimSpace(sva)
					sn = strings.ReplaceAll(sn, "Bpk.", "")
					sn = strings.ReplaceAll(sn, "Bpk", "")
					sn = strings.ReplaceAll(sn, "Ibu.", "")
					sn = strings.ReplaceAll(sn, "Ibu", "")
					sn = strings.TrimSpace(sn)

					splitsn := strings.SplitN(sn, " ", 2)
					sfn := splitsn[0]
					var sln string
					if len(splitsn) > 1 {
						sln = splitsn[1]
					}
					if len(sg) > 3 {
						if sg[0:3] == "Bpk" {
							gender = "m"
						}
					}
					substitutes = append(substitutes, wrapper.OrderUserSubstituteForm{
						Firstname: sfn,
						Lastname:  sln,
						Gender:    sg,
					})
				}
			}

			var phone string
			if row[12] != "0" && row[12] != "" {
				phone = strings.Split(strings.TrimSpace(row[12]), "/")[0]
				phone = strings.ReplaceAll(phone, "-", "")
			}

			// Simpan data user
			var uid int
			user := dbquery.UserNew()
			err := user.SetFirstName(firstname).
				SetLastName(lastname).
				SetPhone(phone).
				SetGender(gender).
				SetRole(wrapper.UserRoleCustomer).
				ReturnID(&uid).
				Save()
			if err != nil {
				log.Warn("ERROR: api.user.go UserCreate() Gagal membuat user baru")
				log.Error(err)
			}

			fmt.Println("ID DATABASE", uid)
			fmt.Println("Kode: ", kode)
			fmt.Println("Jenis Kelamin: ", gender)
			fmt.Println("Nama depan: ", firstname)
			fmt.Println("Nama belakang: ", lastname)
			fmt.Println(substitutes)
			fmt.Println("Nomor HP: ", phone)
			fmt.Println()
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

// DeveloperUpgradeListAvailable list file upgrade
func DeveloperUpgradeListAvailable(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	f, err := os.Open("../data/upgrade")
	if err != nil {
		log.Warn("api.developer.go DeveloperUpgradeListAvailable() Gagal membuka folder")
		log.Error(err)
		message = "Terjadi kesalahan saat membuka folder"
		status = "error"
		next = false
	}

	var listFile []map[string]interface{}

	if next {
		files, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			log.Warn("api.developer.go DeveloperUpgradeListAvailable() Gagal membuka membaca file")
			log.Error(err)
			message = "Terjadi kesalahan mencoba membaca file"
			status = "error"
			next = false
		}

		for _, file := range files {
			hour, minute, _ := file.ModTime().Clock()
			year, month, day := file.ModTime().Date()
			listFile = append(listFile, map[string]interface{}{
				"Name": file.Name(),
				"Size": misc.FileFormatSize(file),
				"Edit": fmt.Sprintf("%02d-%02d-%d", day, month, year) + " " + fmt.Sprintf("%02d:%02d", hour, minute),
			})
		}
	}

	if next {
		if len(listFile) > 0 {
			message = "Menampilkan list file upgrade"
			status = "success"
			httpStatus = http.StatusOK
		} else {
			message = "Tidak ada file upgrade"
			status = "error"
			httpStatus = http.StatusOK
		}
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"files":   listFile,
	}

	c.JSON(httpStatus, m)
}

// DeveloperCloudListAvailable list file di cloud
func DeveloperCloudListAvailable(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""

	f, err := os.Open("../data/cloud")
	if err != nil {
		log.Warn("api.developer.go DeveloperCloudListAvailable() Gagal membuka folder")
		log.Error(err)
		message = "Terjadi kesalahan saat membuka folder"
		status = "error"
		next = false
	}

	var listFile []map[string]interface{}

	if next {
		files, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			log.Warn("api.developer.go DeveloperCloudListAvailable() Gagal membuka membaca file")
			log.Error(err)
			message = "Terjadi kesalahan mencoba membaca file"
			status = "error"
			next = false
		}

		for _, file := range files {
			hour, minute, _ := file.ModTime().Clock()
			year, month, day := file.ModTime().Date()
			listFile = append(listFile, map[string]interface{}{
				"Name": file.Name(),
				"Size": misc.FileFormatSize(file),
				"Edit": fmt.Sprintf("%02d-%02d-%d", day, month, year) + " " + fmt.Sprintf("%02d:%02d", hour, minute),
			})
		}
	}

	if next {
		if len(listFile) > 0 {
			message = "Menampilkan list file dari awan"
			status = "success"
			httpStatus = http.StatusOK
		} else {
			message = "Tidak ada file di awan"
			status = "error"
			httpStatus = http.StatusOK
		}
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"files":   listFile,
	}

	c.JSON(httpStatus, m)
}

// DeveloperUpgradeRemove API untuk menghapus file upgrade
func DeveloperUpgradeRemove(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	// File
	file := c.Query("file")

	err := os.Remove("../data/upgrade/" + file)
	if err != nil {
		message = "Gagal menghapus file"
		status = "error"
	}

	if next {
		message = "File berhasil dihapus"
		status = "success"
		httpStatus = http.StatusOK
	}

	m := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}

// DeveloperCloudRemove API untuk menghapus file di cloud
func DeveloperCloudRemove(c *gin.Context) {
	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	var simpleErrMap map[string]interface{}

	// File
	file := c.Query("file")

	err := os.Remove("../data/cloud/" + file)
	if err != nil {
		message = "Gagal menghapus file"
		status = "error"
	}

	if next {
		message = "File berhasil dihapus"
		status = "success"
		httpStatus = http.StatusOK
	}

	m := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}

// DeveloperInstallUpgrade upgrade
func DeveloperInstallUpgrade(c *gin.Context) {
	message := ""
	//next := true
	httpStatus := http.StatusBadRequest
	status := ""

	if misc.GetEnv("REMOTE", "false") == "true" {
		cmd := exec.Command("systemctl", "restart", "cvnazwa")
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", out)
	}

	m := gin.H{
		"message": message,
		"status":  status,
	}

	c.JSON(httpStatus, m)
}
