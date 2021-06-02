package api

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/wrapper"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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

	var imLine int

	// File
	file, err := c.FormFile("file")

	if err != nil {
		log.Warn("api.developer.go DeveloperImportUpload() File tidak valid")
		log.Error(err)
		message = "Tidak ada file, atau format file tidak valid"
		status = "error"
		next = false
	}

	lineName := file.Filename
	lineName = strings.ToLower(strings.Split(lineName, ".")[0])
	if lid, err := dbquery.LineGetId(lineName); err == nil {
		imLine = lid
	} else {
		log.Warn("Select line id by code, fail")
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
	customers, _ := f.GetRows("Pelanggan")
	for rid, row := range customers {
		// Baca mulai dari baris ke 5
		if rid >= 4 && row[2] != "" && row[29] != "Lunas" {
			var lineCode string
			gender := "m"
			if row[2][0:3] == "Ibu" {
				gender = "f"
			}
			lineCode = strings.ToLower(row[1])

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
			var imSubstitutes []wrapper.OrderUserSubstituteForm

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
					imSubstitutes = append(imSubstitutes, wrapper.OrderUserSubstituteForm{
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
			var code string
			var imLineMax int
			var imItems string
			var imDeposit int
			var imDuration int
			var imMonthly int
			var imAddress string
			var imShippingDate string
			var imSales string
			var imSurveyor string
			var imDue int
			var imNotes string

			user := dbquery.UserNew()
			err := user.SetFirstName(firstname).
				SetLastName(lastname).
				SetPhone(phone).
				SetGender(gender).
				SetRole(wrapper.UserRoleCustomer).
				ReturnID(&uid).
				Save()
			if err != nil {
				log.Warn("ERROR: api.developer.go DeveloperImportUpload() Gagal membuat user baru")
				log.Error(err)
			} else {

				uname, err := dbquery.UserGetUsername(uid)
				if err == nil {
					if uname[:3] == "NZ-" || uname[:3] == "NE-" {
						if len(uname) >= 7 {
							uname = uname[3:]
						}
					}
					tm := time.Now()
					dt := strings.ReplaceAll(tm.Format("01-02-2006"), "-", "")
					dy := tm.Format("Mon")
					uq := tm.Format(".000")[1:]
					code = strings.ToUpper(fmt.Sprintf("%s%s-%s%s-%s%s", uname[4:], dy, uq, dt[4:], dt[:4], uname[:4]))
				} else {
					log.Warn("ERROR: api.developer.go DeveloperImportUpload() Gagal membuat kode")
					log.Error(err)
				}

				imLineMax, _ = strconv.Atoi(strings.TrimPrefix(strings.ReplaceAll(lineCode, lineName, ""), "0"))

				if row[4] != "" {
					imItems = row[4]
				}
				if row[6] != "" {
					imDeposit, _ = strconv.Atoi(row[6])
				}
				if row[8] != "" {
					imDuration, _ = strconv.Atoi(row[8])
				}
				if row[10] != "" {
					imMonthly, _ = strconv.Atoi(row[10])
				}
				if row[11] != "" {
					imAddress = row[11]
				}
				if row[13] != "" {
					imShippingDate = row[13]
					imShippingDate = strings.ReplaceAll(imShippingDate, " ", "")
					imShippingDate = strings.ReplaceAll(imShippingDate, "/", "-")
					imShippingDate = strings.ReplaceAll(imShippingDate, ".", "-")
					imSDS := strings.Split(imShippingDate, "-")
					if len(imSDS) == 3 {
						imShippingDate = fmt.Sprintf("%s-%s-%s", imSDS[2], imSDS[1], imSDS[0])
					} else {
						log.Warn("Tanggal tidak diformat ulang")
					}
				}
				if row[14] != "" {
					imSales = row[14]
				}
				if row[15] != "" {
					imSurveyor = row[15]
				}
				if row[16] != "" {
					imDue, _ = strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(row[16], " ", ""), "[", ""), "]", ""))
				}
				if row[17] != "" {
					imDue, _ = strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(row[17], " ", ""), "{", ""), "}", ""))
				}
				if row[21] != "" {
					imNotes = row[21]
				}

				// Input order
				order := dbquery.NewOrder()
				errs := order.SetCustomer(uid).
					SetImportMode(true).
					SetImportedSales(imSales).
					SetImportedSurveyor(imSurveyor).
					SetImportedAddress(imAddress).
					SetImportedItems(imItems).
					SetCredit(true).
					SetDeposit(imDeposit).
					SetDuration(imDuration).
					SetDue(imDue).
					SetLine(imLine).
					SetLineCodeMaxNumber(imLineMax).
					SetNotes(imNotes).
					SetCode(code).
					SetOrderDate(imShippingDate).
					SetShippingDate(imShippingDate).
					SetSubstitutes(imSubstitutes).
					SetImportedMonthly(imMonthly).
					Save()

				if errs != nil {
					log.Warn("ERROR: api.developer.go DeveloperImportUpload() Gagal membuat order baru")
					log.Error(errs)
				}
			}

			// fmt.Println("NAMA FILE: ", lineName)
			// fmt.Println("ID DATABASE", uid)
			// fmt.Println("Kode: ", lineCode)
			// fmt.Println("Jenis Kelamin: ", gender)
			// fmt.Println("Nama depan: ", firstname)
			// fmt.Println("Nama belakang: ", lastname)
			// fmt.Println(imSubstitutes)
			// fmt.Println("Nomor HP: ", phone)
			// fmt.Println("ORDER")
			// fmt.Println("Kode: ", code)
			// fmt.Println("Arah: ", imLine)
			// fmt.Println("Sales: ", imSales)
			// fmt.Println("Survey: ", imSurveyor)
			// fmt.Println("Deposit: ", imDeposit)
			// fmt.Println("Durasi: ", imDuration)
			// fmt.Println("Alamat: ", imAddress)
			// fmt.Println("Tanggal:", imShippingDate)
			// fmt.Println("Jatuh tempo: ", imDue)
			// fmt.Println("Catatan: ", imNotes)
			// fmt.Println("Bulanan: ", imMonthly)
			// fmt.Println()
		}
	}

	importPayment := func(oid int, agDate string, agAmountOne string, agAmountTwo string, agMethod string, agReceiver string) wrapper.OrderPaymentInsert {
		var pLineDate string
		var amount int
		cash := true

		dirtyDate := agDate
		dirtyDate = strings.ReplaceAll(dirtyDate, "/", "-")
		dirtyDate = strings.ReplaceAll(dirtyDate, " ", "")
		ddSplit := strings.Split(dirtyDate, "-")
		if len(ddSplit) == 3 {
			pLineDate = fmt.Sprintf("%02s-%02s-%s", ddSplit[2], ddSplit[1], ddSplit[0])
		} else {
			log.Warn("Tanggal tidak diformat ulang")
		}

		dirtyAmount := agAmountOne
		if agAmountTwo != "" {
			dirtyAmount = agAmountTwo
		}

		amount, _ = strconv.Atoi(dirtyAmount)

		// metode pembayaran
		if agMethod != "Tunai" {
			cash = false
		}

		receiver := agReceiver

		_, err := time.Parse(`2006-01-02`, pLineDate)
		if err != nil {
			pLineDate = "2016-01-01"
		}

		if len(receiver) > 20 {
			receiver = "ERROR PANJANG"
		}

		paymentData := wrapper.OrderPaymentInsert{
			OrderID:          oid,
			ImportedReceiver: receiver,
			PaymentDate:      pLineDate,
			Cash:             cash,
			Amount:           amount,
		}

		// fmt.Println("P Oid: ", oid)
		// fmt.Println("P Line date: ", pLineDate)
		// fmt.Println("P Penerima: ", receiver)
		// fmt.Println("P Cash: ", cash)
		// fmt.Println("P Amount: ", amount)

		return paymentData
	}
	// Import data pembayaran masuk
	paids, _ := f.GetRows("Pembayaran")
	var paymentInsertData []wrapper.OrderPaymentInsert

	var plastCode string
	var plastDate string
	for rid, row := range paids {
		if rid >= 3 && row[1] != "" {
			record := true

			if plastCode == row[1] {
				if plastDate == row[2] {
					record = false
				}
			}

			// Check apakah Order dengan kode kredit ini terdaftar
			oid, err := dbquery.OrderGetIDByCode(strings.ToUpper(row[1]))
			if err != nil {
				record = false
			}

			// Data untuk disimpan
			if record {
				paymentInsertData = append(paymentInsertData, importPayment(oid, row[2], row[4], row[6], row[5], row[7]))
			}

			plastCode = row[1]
			plastDate = row[2]
		}
	}

	// Simpan data pembayaran
	err = dbquery.OrderCreditAddPayment(paymentInsertData)
	if err != nil {
		fmt.Println("Gagal insert data pembayaran")
		log.Error(err)
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
