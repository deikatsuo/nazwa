package api

import (
	"fmt"
	"io/ioutil"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

////////////
// CREATE //
////////////

// UserCreate API untuk membuat user baru
func UserCreate(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")

	var json wrapper.UserForm

	status := "success"
	var httpStatus int
	message := ""
	var simpleErrMap = make(map[string]interface{})
	save := true
	var avatarExt string
	var avatar string
	var file string

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println("ERROR: api.user.go UserCreate() bind json")
		log.Println(err)
		if fmt.Sprintf("%T", err) == "validator.ValidationErrors" {
			simpleErrMap = validation.SimpleValErrMap(err)
		}
		httpStatus = http.StatusBadRequest
		status = "fail"
		save = false
	} else {
		if dbquery.UserRICExist(json.RIC) {
			simpleErrMap["ric"] = "Nomor KTP sudah terdaftar"
			status = "fail"
			save = false
		}
		if dbquery.UserPhoneExist(json.Phone) {
			simpleErrMap["phone"] = "Nomor ini sudah terdaftar"
			status = "fail"
			save = false
		}
	}

	if json.PhotoType != "" && json.Photo != "" {
		avatarExt = json.PhotoType
		avatar = json.Photo
		if f, err := misc.FileBase64ToFileWithData("../data/upload/profile/", avatar, avatarExt); err == nil {
			file = f
		} else {
			log.Println("ERROR: api.user.go UserCreate() Konversi base64 ke file gambar")
			message = err.Error()
		}
	}

	if file != "" {
		err := misc.FileGenerateThumb(file, "../data/upload/profile/")
		if err != nil {
			message = err.Error()
		}
	}

	var uid int
	var retUser wrapper.User
	if save {
		user := dbquery.UserNew()
		err := user.SetFirstName(json.Firstname).
			SetLastName(json.Lastname).
			SetFamilyCard(json.FC).
			SetRIC(json.RIC).
			SetPhone(json.Phone).
			SetAvatar(file).
			SetPassword(json.Password).
			SetGender(json.Gender).
			SetOccupation(json.Occupation).
			SetRole(wrapper.UserRoleCustomer).
			SetCreatedBy(nowID.(int)).
			ReturnID(&uid).
			Save()
		if err != nil {
			log.Println("ERROR: api.user.go UserCreate() Gagal membuat user baru")
			log.Print(err)
			if err := os.Remove("../data/upload/profile/" + file); err != nil {
				log.Println("ERROR: api.user.go UserCreate() Gagal menghapus photo")
				log.Println(err)
			}

			if err := os.Remove("../data/upload/profile/thumbnail/" + file); err != nil {
				log.Println("ERROR: api.user.go UserCreate() Gagal menghapus thumbnail")
				log.Println(err)
			}
		} else {
			httpStatus = http.StatusOK
			status = "success"
			message = "Berhasil membuat user baru"

			if u, err := dbquery.UserGetByID(uid); err == nil {
				retUser = u
			} else {
				httpStatus = http.StatusInternalServerError
				message = "Sepertinya telah terjadi kesalahan saat memuat data"
			}
		}
	} else {
		httpStatus = http.StatusBadRequest
		status = "error"
		message = "Gagal membuat pelanggan, silahkan perbaiki formulir"
	}

	m := gin.H{
		"message": message,
		"status":  status,
		"user":    retUser,
	}
	c.JSON(httpStatus, misc.Mete(m, simpleErrMap))
}

////////////
/* UPDATE */
////////////

// TODO: Pindahkan ini
type updateContact struct {
	Username    string `json:"username" binding:"alphanum,min=4,max=25"`
	Password    string `json:"password" binding:"omitempty,alphanumunicode,min=8,max=25"`
	Repassword  string `json:"repassword" binding:"eqfield=Password"`
	Oldpassword string `json:"oldpassword" binding:"required_with=Password"`
}

// UserUpdateContact api untuk mengupdate kontak user
func UserUpdateContact(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""
	var simpleErr map[string]interface{}

	var update updateContact
	if err := c.ShouldBindJSON(&update); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
	}

	// Check ketersediaan username
	if next {
		if dbquery.UsernameExist(update.Username) {
			errMess = "Username tidak tersedia"
			next = false
		}
	}

	// Update username
	if next {
		if err := dbquery.UserUpdateUsername(uid, update.Username); err != nil {
			errMess = "Gagal mengubah username"
			next = false
		}
	}

	// Cocokan password lama
	if next {
		if update.Password != "" {
			if !dbquery.UserMatchPassword(uid, update.Oldpassword) {
				errMess = "Kata sandi lama salah"
				next = false
			}
		}
	}

	// Update password
	if next {
		if update.Password != "" {
			if err := dbquery.UserUpdatePassword(uid, update.Password); err != nil {
				errMess = "Gagal merubah kata sandi"
				next = false
			}
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		success = "Data berhasil disimpan"
	}

	gh := gin.H{
		"error":   errMess,
		"success": success,
	}
	c.JSON(httpStatus, misc.Mete(gh, simpleErr))
}

// UserUpdateRole api untuk mengupdate kontak user
func UserUpdateRole(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// Role user yang melakukan request
	nowRole := session.Get("role")
	// ID user yang role nya akan di ubah
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""
	var simpleErr map[string]interface{}

	// Role baru
	newRole := c.Query("set")
	newRoleID, err := strconv.Atoi(newRole)
	if err != nil || nowRole == nil {
		errMess = "Request tidak valid"
		next = false
	}

	if nowID == uid {
		errMess = "Tidak bisa merubah peran sendiri"
		next = false
	}

	// Jika request set role ke dev (0)
	if next {
		if newRoleID == 0 && nowRole != 0 {
			errMess = "Harus dev untuk membuat dev baru"
			next = false
		}
	}

	// Update role
	if next {
		if err := dbquery.UserUpdateRole(uid, newRoleID); err != nil {
			errMess = "Gagal mengubah role"
			next = false
		}
	}

	// Berhasil update data
	if next {
		httpStatus = http.StatusOK
		success = "Data berhasil disimpan"
	}

	gh := gin.H{
		"error":   errMess,
		"success": success,
	}

	c.JSON(httpStatus, misc.Mete(gh, simpleErr))
}

////////////
/* DELETE */
////////////

// UserDeletePhone api untuk menghapus nomor HP
func UserDeletePhone(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errMess = "Data tidak benar"
		next = false
	}

	id, err := jsonparser.GetInt(body, "id")
	if err != nil {
		errMess = "Request tidak valid"
	}

	// Delete nomor HP
	if next {
		if err := dbquery.UserDeletePhone(id, uid); err != nil {
			errMess = "Gagal menghapus nomor HP"
			next = false
		}
	}

	// Ambil nomor sisa jika masih ada
	var phones []wrapper.UserPhone
	if next {
		ph, err := dbquery.UserGetPhone(uid)
		if err != nil {
			errMess = "Gagal memuat nomor HP/semua nomor sudah dihapus"
		} else {
			phones = ph
			httpStatus = http.StatusOK
			success = "Nomor berhasil dihapus"
		}
	}
	c.JSON(httpStatus, gin.H{
		"error":   errMess,
		"success": success,
		"phones":  phones,
	})
}

// UserDeleteEmail api untuk menghapus email
func UserDeleteEmail(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errMess = "Data tidak benar"
		next = false
	}

	id, err := jsonparser.GetInt(body, "id")
	if err != nil {
		errMess = "Request tidak valid"
	}

	// Delete email
	if next {
		if err := dbquery.UserDeleteEmail(id, uid); err != nil {
			errMess = "Gagal menghapus email"
			next = false
		}
	}

	// Ambil email sisa jika masih ada
	var emails []wrapper.UserEmail
	if next {
		em, err := dbquery.UserGetEmail(uid)
		if err != nil {
			errMess = "Gagal memuat email/semua email sudah dihapus"
		} else {
			emails = em
			httpStatus = http.StatusOK
			success = "Email berhasil dihapus"
		}
	}
	c.JSON(httpStatus, gin.H{
		"error":   errMess,
		"success": success,
		"emails":  emails,
	})
}

// UserDeleteAddress menghapus alamat
func UserDeleteAddress(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	message := ""
	next := true
	httpStatus := http.StatusBadRequest
	status := ""
	errAlert := ""

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		message = "Data tidak benar"
		status = "error"
		next = false
	}

	id, err := jsonparser.GetInt(body, "id")
	if err != nil {
		message = "Request tidak valid"
		status = "error"
	}

	// Hapus alamat
	if next {
		if err := dbquery.UserDeleteAddress(id, uid); err != nil {
			errAlert = err.Error()
			message = "Gagal menghapus alamat, kemungkinan karena ada beberapa order yang menggunakan alamat ini"
			status = "error"
			next = false
		} else {
			httpStatus = http.StatusOK
			message = "Alamat berhasil dihapus"
			status = "success"
			next = true
		}
	}

	// Ambil alamat sisa jika masih ada
	var addresses []wrapper.Address
	if next {
		em, err := dbquery.AddressGetByUserID(uid)
		if err == nil {
			addresses = em
		}
	}
	c.JSON(httpStatus, gin.H{
		"status":    status,
		"message":   message,
		"addresses": addresses,
		"error":     errAlert,
	})
}

/////////
/* ADD */
/////////

type tmpMail struct {
	Email string `json:"email" binding:"email"`
}

// UserAddEmail api untuk menambah email
func UserAddEmail(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""

	var newEmail tmpMail
	if err := c.ShouldBindJSON(&newEmail); err != nil {
		simpleErr := validation.SimpleValErrMap(err)
		errMess = simpleErr["email"].(string)
		next = false
	}

	// Periksa jika email sudah digunakan
	if next {
		if dbquery.UserEmailExist(newEmail.Email) {
			errMess = "Email ini sudah digunakan"
			next = false
		}
	}

	// Tambah email
	if next {
		if err := dbquery.UserAddEmail(newEmail.Email, uid); err != nil {
			errMess = "Gagal menambahkan email"
			next = false
		}
	}

	// Ambil email sisa jika masih ada
	var emails []wrapper.UserEmail
	if next {
		em, err := dbquery.UserGetEmail(uid)
		if err != nil {
			errMess = "Gagal memuat email/semua email sudah dihapus"
		} else {
			emails = em
			httpStatus = http.StatusOK
			success = "Email berhasil ditambahkan"
		}
	}
	c.JSON(httpStatus, gin.H{
		"error":   errMess,
		"success": success,
		"emails":  emails,
	})
}

type tmpPhone struct {
	Phone string `json:"phone" binding:"numeric,min=6,max=15"`
}

// UserAddPhone api untuk menambah nomor HP
func UserAddPhone(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""

	var newPhone tmpPhone
	if err := c.ShouldBindJSON(&newPhone); err != nil {
		simpleErr := validation.SimpleValErrMap(err)
		errMess = simpleErr["phone"].(string)
		next = false
	}

	// Periksa jika nomor sudah digunakan
	if next {
		if dbquery.UserPhoneExist(newPhone.Phone) {
			errMess = "Nomor ini sudah digunakan"
			next = false
		}
	}

	// Tambahkan nomor HP
	if next {
		if err := dbquery.UserAddPhone(newPhone.Phone, uid); err != nil {
			errMess = "Gagal menambahkan nomor HP"
			next = false
		}
	}

	// Ambil nomor HP dari database
	var phones []wrapper.UserPhone
	if next {
		ph, err := dbquery.UserGetPhone(uid)
		if err != nil {
			errMess = "Gagal memuat nomor HP"
		} else {
			phones = ph
			httpStatus = http.StatusOK
			success = "Nomor HP berhasil ditambahkan"
		}
	}
	c.JSON(httpStatus, gin.H{
		"error":   errMess,
		"success": success,
		"phones":  phones,
	})
}

// UserAddAddress api untuk menambahkan address
func UserAddAddress(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	nowID := session.Get("userid")
	// User id yang merequest
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || nowID == nil {
		router.Page404(c)
		return
	}

	errMess := ""
	next := true
	httpStatus := http.StatusBadRequest
	success := ""
	var simpleErr map[string]interface{}

	var newAddress wrapper.AddressForm
	if err := c.ShouldBindJSON(&newAddress); err != nil {
		simpleErr = validation.SimpleValErrMap(err)
		next = false
	}

	// Tambahkan alamat
	addressInsert := wrapper.AddressInsert{
		UserID:      uid,
		Description: newAddress.Description,
		One:         newAddress.One,
		Two:         newAddress.Two,
		Zip:         newAddress.Zip,
		Province:    newAddress.Province,
		City:        newAddress.City,
		District:    newAddress.District,
		Village:     newAddress.Village,
	}

	if next {
		if err := dbquery.AddressAdd(addressInsert); err != nil {
			log.Warning(err)
			errMess = "Gagal menambahkan alamat"
			next = false
		}
	}

	// Ambil data alamat dari database
	var addresses []wrapper.Address
	if next {
		ph, err := dbquery.AddressGetByUserID(uid)
		if err != nil {
			errMess = "Gagal memuat alamat"
		} else {
			addresses = ph
			httpStatus = http.StatusOK
			success = "Alamat berhasil ditambahkan"
		}
	}
	gh := gin.H{
		"error":     errMess,
		"success":   success,
		"addresses": addresses,
	}
	c.JSON(httpStatus, misc.Mete(gh, simpleErr))
}

///////////
//  GET  //
///////////

// UserShowList mengambil data/list pengguna
func UserShowList(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	// Tolak jika yang request bukan user terdaftar
	uid := session.Get("userid")
	if uid == nil {
		router.Page404(c)
		return
	}

	lastid := 1
	last := false
	limit := 10
	var direction string
	httpStatus := http.StatusOK
	errMess := ""
	u := dbquery.GetUsers{}
	next := true

	// Mengambil parameter limit
	lim, err := strconv.Atoi(c.Param("limit"))
	if err == nil {
		limit = lim
		u.Limit(limit)
	} else {
		errMess = "Limit tidak valid"
		httpStatus = http.StatusBadRequest
		next = false
	}

	// Ambil query id terakhir
	lst, err := strconv.Atoi(c.Query("lastid"))
	if err == nil {
		lastid = lst
	}

	// Forward/backward
	direction = c.Query("direction")
	if direction == "back" {
		u.Direction(direction)
	} else {
		u.Direction("next")
	}

	var total int
	if t, err := dbquery.UserGetUserTotalRow(); err == nil {
		total = t
	}

	var users []wrapper.User

	if next {
		u.LastID(lastid)
		// Maju/Mundur
		if direction == "next" {
			u.Where("WHERE u.id > " + strconv.Itoa(lastid) + " ORDER BY u.id ASC")
		} else if direction == "back" {
			u.Where("WHERE u.id < " + strconv.Itoa(lastid) + " ORDER BY u.id DESC")
		}

		usr, err := u.Show()
		if err != nil {
			errMess = err.Error()
			httpStatus = http.StatusInternalServerError
		}
		users = usr
	}

	if len(users) > 0 && direction == "back" {
		// Reverse urutan array user
		temp := make([]wrapper.User, len(users))
		in := 0
		for i := len(users) - 1; i >= 0; i-- {
			temp[in] = users[i]
			in++
		}
		users = temp
	}

	// Cek id terakhir
	if len(users) > 0 && len(users) < limit {
		// Periksa apakah ini data terakhir atau bukan
		last = true
	}

	c.JSON(httpStatus, gin.H{
		"users": users,
		"error": errMess,
		"total": total,
		"last":  last,
	})
}

// UserShowByID mengambil data pengguna berdasarkan ID
func UserShowByID(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	// Tolak jika yang request bukan user terdaftar
	uid := session.Get("userid")
	if uid == nil {
		router.Page404(c)
		return
	}
	httpStatus := http.StatusOK
	errMess := ""

	// Mengambil parameter id user
	var uid2 int
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		uid2 = id
	} else {
		httpStatus = http.StatusBadRequest
		errMess = "Request tidak valid"
	}

	var user wrapper.User
	if u, err := dbquery.UserGetByID(uid2); err == nil {
		user = u
	} else {
		httpStatus = http.StatusInternalServerError
		errMess = "Sepertinya telah terjadi kesalahan saat memuat data"
	}

	c.JSON(httpStatus, gin.H{
		"user":  user,
		"error": errMess,
	})
}

// UserShowAddressByUserID mengambil data alamat pengguna berdasarkan ID pengguna
func UserShowAddressByUserID(c *gin.Context) {
	httpStatus := http.StatusOK
	errMess := ""

	// Mengambil parameter id user
	var uid2 int
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		uid2 = id
	} else {
		httpStatus = http.StatusBadRequest
		errMess = "Request tidak valid"
	}

	var address []wrapper.Address
	if addr, err := dbquery.AddressGetByUserID(uid2); err == nil {
		address = addr
	} else {
		httpStatus = http.StatusInternalServerError
		errMess = "Sepertinya telah terjadi kesalahan saat memuat data"
	}

	c.JSON(httpStatus, gin.H{
		"addresses": address,
		"error":     errMess,
	})
}

// UserSearchByNIK cari user berdasarkan Nomor induk kependudukan
func UserSearchByNIK(c *gin.Context) {
	session := sessions.Default(c)
	// User session saat ini
	// Tolak jika yang request bukan user terdaftar
	uid := session.Get("userid")
	if uid == nil {
		router.Page404(c)
		return
	}

	search := ""
	lastid := 1
	last := false
	limit := 10
	var direction string
	httpStatus := http.StatusOK
	errMess := ""
	u := dbquery.GetUsers{}
	next := true

	// Mengambil parameter limit
	lim, err := strconv.Atoi(c.Param("limit"))
	if err == nil {
		limit = lim
		u.Limit(limit)
	} else {
		errMess = "Limit tidak valid"
		httpStatus = http.StatusBadRequest
		next = false
	}

	// Ambil query pencarian
	search = c.Query("search")

	// Ambil query id terakhir
	lst, err := strconv.Atoi(c.Query("lastid"))
	if err == nil {
		lastid = lst
	}

	// Forward/backward
	direction = c.Query("direction")
	if direction == "back" {
		u.Direction(direction)
	} else {
		u.Direction("next")
	}

	var users []wrapper.User

	if next {
		u.Where("WHERE u.ric LIKE '" + search + "%' ORDER BY u.id ASC")
		fmt.Println(search)
		u.LastID(lastid)

		usr, err := u.Show()
		if err != nil {
			errMess = err.Error()
			httpStatus = http.StatusInternalServerError
		}
		users = usr
	}

	if len(users) > 0 && direction == "back" {
		// Reverse urutan array user
		temp := make([]wrapper.User, len(users))
		in := 0
		for i := len(users) - 1; i >= 0; i-- {
			temp[in] = users[i]
			in++
		}
		users = temp
	}

	// Cek id terakhir
	if len(users) > 0 && len(users) < limit {
		// Periksa apakah ini data terakhir atau bukan
		last = true
	}

	c.JSON(httpStatus, gin.H{
		"users": users,
		"error": errMess,
		"last":  last,
	})
}

// UserSearchByNameType cari sales berdasarkan nama
func UserSearchByNameType(roleid string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		session := sessions.Default(c)
		// User session saat ini
		// Tolak jika yang request bukan user terdaftar
		uid := session.Get("userid")
		if uid == nil {
			router.Page404(c)
			return
		}

		search := ""
		lastid := 1
		last := false
		limit := 10
		var direction string
		httpStatus := http.StatusOK
		errMess := ""
		u := dbquery.GetUsers{}
		next := true

		// Mengambil parameter limit
		lim, err := strconv.Atoi(c.Param("limit"))
		if err == nil {
			limit = lim
			u.Limit(limit)
		} else {
			errMess = "Limit tidak valid"
			httpStatus = http.StatusBadRequest
			next = false
		}

		// Ambil query pencarian
		search = strings.ToLower(c.Query("search"))

		// Ambil query id terakhir
		lst, err := strconv.Atoi(c.Query("lastid"))
		if err == nil {
			lastid = lst
		}

		// Forward/backward
		direction = c.Query("direction")
		if direction == "back" {
			u.Direction(direction)
		} else {
			u.Direction("next")
		}

		var total int
		if t, err := dbquery.UserGetUserTotalRow(); err == nil {
			total = t
		}

		var users []wrapper.User

		if next {
			u.Where("WHERE r.id=" + roleid + " AND concat(u.first_name, ' ', u.last_name) LIKE '" + search + "%' ORDER BY u.id ASC")
			u.LastID(lastid)

			usr, err := u.Show()
			if err != nil {
				errMess = err.Error()
				httpStatus = http.StatusInternalServerError
			}
			users = usr
		}

		if len(users) > 0 && direction == "back" {
			// Reverse urutan array user
			temp := make([]wrapper.User, len(users))
			in := 0
			for i := len(users) - 1; i >= 0; i-- {
				temp[in] = users[i]
				in++
			}
			users = temp
		}

		// Cek id terakhir
		if len(users) > 0 && len(users) < limit {
			// Periksa apakah ini data terakhir atau bukan
			last = true
		}

		c.JSON(httpStatus, gin.H{
			"users": users,
			"error": errMess,
			"total": total,
			"last":  last,
		})
	}
	return gin.HandlerFunc(fn)
}
