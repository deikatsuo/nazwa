package api

import (
	"io/ioutil"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"nazwa/wrapper"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type updateContact struct {
	Username    string `json:"username" binding:"alphanum,min=4,max=25"`
	Password    string `json:"password" binding:"omitempty,alphanumunicode,min=8,max=25"`
	Repassword  string `json:"repassword" binding:"eqfield=Password"`
	Oldpassword string `json:"oldpassword" binding:"required_with=Password"`
}

////////////
/* UPDATE */
////////////

// UserUpdateContact api untuk mengupdate kontak user
func UserUpdateContact(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		// Periksa jika id peng request tidak sama dengan
		// id yang tersimpan di session
		if next {
			if nowID != uid {
				errMess = "User tidak memiliki ijin untuk mengubah data akun ini"
				next = false
			}
		}

		// Check ketersediaan username
		if next {
			if dbquery.UsernameExist(db, update.Username) {
				errMess = "Username tidak tersedia"
				next = false
			}
		}

		// Update username
		if next {
			if err := dbquery.UpdateUsername(db, nowID.(int), update.Username); err != nil {
				errMess = "Gagal mengubah username"
				next = false
			}
		}

		// Cocokan password lama
		if next {
			if update.Password != "" {
				if !dbquery.MatchPassword(db, nowID.(int), update.Oldpassword) {
					errMess = "Kata sandi lama salah"
					next = false
				}
			}
		}

		// Update password
		if next {
			if update.Password != "" {
				if err := dbquery.UpdatePassword(db, nowID.(int), update.Password); err != nil {
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
	return gin.HandlerFunc(fn)
}

////////////
/* DELETE */
////////////

// UserDeletePhone api untuk menghapus nomor HP
func UserDeletePhone(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		// Periksa jika id peng request tidak sama dengan
		// id yang tersimpan di session
		if next {
			if nowID.(int) != uid {
				errMess = "User tidak memiliki ijin untuk menghapus nomor ini"
				next = false
			}
		}

		// Delete nomor HP
		if next {
			if err := dbquery.DeletePhone(db, id, nowID.(int)); err != nil {
				errMess = "Gagal menghapus nomor HP"
				next = false
			}
		}

		// Ambil email sisa jika masih ada
		var phones []wrapper.UserPhone
		if next {
			ph, err := dbquery.GetPhone(db, nowID.(int))
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
	return gin.HandlerFunc(fn)
}

// UserDeleteEmail api untuk menghapus email
func UserDeleteEmail(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		// Periksa jika id peng request tidak sama dengan
		// id yang tersimpan di session
		if next {
			if nowID.(int) != uid {
				errMess = "User tidak memiliki ijin untuk menghapus email ini"
				next = false
			}
		}

		// Delete email
		if next {
			if err := dbquery.DeleteEmail(db, id, nowID.(int)); err != nil {
				errMess = "Gagal menghapus email"
				next = false
			}
		}

		// Ambil email sisa jika masih ada
		var emails []wrapper.UserEmail
		if next {
			em, err := dbquery.GetEmail(db, nowID.(int))
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
	return gin.HandlerFunc(fn)
}

// UserDeleteAddress menghapus alamat
func UserDeleteAddress(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		// Periksa jika id peng request tidak sama dengan
		// id yang tersimpan di session
		if next {
			if nowID.(int) != uid {
				errMess = "User tidak memiliki ijin untuk menghapus alamat ini"
				next = false
			}
		}

		// Hapus alamat
		if next {
			if err := dbquery.DeleteAddress(db, id, nowID.(int)); err != nil {
				errMess = "Gagal menghapus alamat"
				next = false
			}
		}

		// Ambil alamat sisa jika masih ada
		var addresses []wrapper.UserAddress
		if next {
			em, err := dbquery.GetAddress(db, nowID.(int))
			if err != nil {
				errMess = "Gagal memuat alamat/semua alamat sudah dihapus"
			} else {
				addresses = em
				httpStatus = http.StatusOK
				success = "Alamat berhasil dihapus"
			}
		}
		c.JSON(httpStatus, gin.H{
			"error":     errMess,
			"success":   success,
			"addresses": addresses,
		})
	}
	return gin.HandlerFunc(fn)
}

/////////
/* ADD */
/////////

type tmpMail struct {
	Email string `json:"email" binding:"email"`
}

// UserAddEmail api untuk menambah email
func UserAddEmail(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		// Periksa jika id peng request tidak sama dengan
		// id yang tersimpan di session
		if next {
			if nowID.(int) != uid {
				errMess = "User tidak memiliki ijin untuk menambahkan email ke akun ini"
				next = false
			}
		}

		// Periksa jika email sudah digunakan
		if next {
			if dbquery.EmailExist(db, newEmail.Email) {
				errMess = "Email ini sudah digunakan"
				next = false
			}
		}

		// Tambah email
		if next {
			if err := dbquery.AddEmail(db, newEmail.Email, nowID.(int)); err != nil {
				errMess = "Gagal menambahkan email"
				next = false
			}
		}

		// Ambil email sisa jika masih ada
		var emails []wrapper.UserEmail
		if next {
			em, err := dbquery.GetEmail(db, nowID.(int))
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
	return gin.HandlerFunc(fn)
}

type tmpPhone struct {
	Phone string `json:"phone" binding:"numeric,min=6,max=15"`
}

// UserAddPhone api untuk menambah nomor HP
func UserAddPhone(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		// Periksa jika id peng request tidak sama dengan
		// id yang tersimpan di session
		if next {
			if nowID.(int) != uid {
				errMess = "User tidak memiliki ijin untuk menambahkan nomor HP ke akun ini"
				next = false
			}
		}

		// Periksa jika nomor sudah digunakan
		if next {
			if dbquery.PhoneExist(db, newPhone.Phone) {
				errMess = "Nomor ini sudah digunakan"
				next = false
			}
		}

		// Tambahkan nomor HP
		if next {
			if err := dbquery.AddPhone(db, newPhone.Phone, nowID.(int)); err != nil {
				errMess = "Gagal menambahkan nomor HP"
				next = false
			}
		}

		// Ambil nomor HP dari database
		var phones []wrapper.UserPhone
		if next {
			ph, err := dbquery.GetPhone(db, nowID.(int))
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
	return gin.HandlerFunc(fn)
}

// UserAddAddress api untuk menambahkan address
func UserAddAddress(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		var newAddress wrapper.UserAddress
		if err := c.ShouldBindJSON(&newAddress); err != nil {
			simpleErr = validation.SimpleValErrMap(err)
			next = false
		}

		// Periksa jika id peng request tidak sama dengan
		// id yang tersimpan di session
		if next {
			if nowID.(int) != uid {
				errMess = "User tidak memiliki ijin untuk menambahkan alamat ke akun ini"
				next = false
			}
		}

		// Tambahkan alamat
		newAddress.UserID = nowID.(int)
		if next {
			if err := dbquery.AddAddress(db, newAddress); err != nil {
				errMess = "Gagal menambahkan alamat"
				next = false
			}
		}

		// Ambil data alamat dari database
		var addresses []wrapper.UserAddress
		if next {
			ph, err := dbquery.GetAddress(db, nowID.(int))
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
	return gin.HandlerFunc(fn)
}

///////////
//  GET  //
///////////

// ShowUserList mengambil data/list pengguna
func ShowUserList(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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
		if t, err := dbquery.GetUserTotalRow(db); err == nil {
			total = t
		}

		var users []wrapper.User

		if next {
			u.LastID(lastid)

			usr, err := u.Show(db)
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

// ShowUserByID mengambil data pengguna berdasarkan ID
func ShowUserByID(db *sqlx.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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
		if u, err := dbquery.GetUserByID(db, uid2); err == nil {
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
	return fn
}
