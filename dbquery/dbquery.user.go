package dbquery

import (
	"fmt"
	"nazwa/misc"
	"nazwa/wrapper"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/martinlindhe/base36"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers mengambil list user
type GetUsers struct {
	limit     int
	lastid    int
	direction string
	where     string
}

// Limit set limit
func (u *GetUsers) Limit(limit int) *GetUsers {
	u.limit = limit
	return u
}

// LastID set lastid
func (u *GetUsers) LastID(lastid int) *GetUsers {
	u.lastid = lastid
	return u
}

// Direction untuk backward/forward
// @direction "back","next"
func (u *GetUsers) Direction(direction string) *GetUsers {
	u.direction = direction
	return u
}

// Where kondisi
func (u *GetUsers) Where(where string) *GetUsers {
	u.where = where
	return u
}

// Show tampilkan data
func (u *GetUsers) Show(db *sqlx.DB) ([]wrapper.User, error) {
	var user []wrapper.NullableUser
	var parse []wrapper.User
	limit := 10
	if u.limit > 0 {
		limit = u.limit
	}

	// Where logic
	where := u.where

	// query pengambilan data user
	query := `SELECT
		u.id,
		u.first_name,
		u.last_name,
		u.username,
		u.avatar,
		u.gender,
		u.ric,
		TO_CHAR(u.created_at, 'MM/DD/YYYY HH12:MI:SS AM') AS created_at,
		u.balance,
		INITCAP(r.name) AS role
		FROM "user" u
		LEFT JOIN "user_role" ur ON ur.user_id=u.id
		LEFT JOIN "role" r ON r.id=ur.role_id
		%s
		LIMIT $1`

	query = fmt.Sprintf(query, where)

	err := db.Select(&user, query, limit)
	if err != nil {
		log.Println("Error: user.go Select all user")
		log.Println(err)
		return []wrapper.User{}, err
	}

	var addresses []wrapper.UserAddress
	// mapping data user
	for _, u := range user {

		if addrs, err := UserGetAddress(db, u.ID); err == nil {
			addresses = addrs
		}
		parse = append(parse, wrapper.User{
			ID:        u.ID,
			Firstname: strings.Title(u.Firstname),
			Lastname:  strings.Title(u.Lastname.String),
			CreatedAt: u.CreatedAt,
			Username:  string(u.Username.String),
			Balance:   string(u.Balance),
			Avatar:    u.Avatar,
			RIC:       u.RIC,
			Role:      u.Role,
			Addresses: addresses,
		})
	}

	return parse, nil
}

// CreateUser struk buat user baru
// Struct data user
type CreateUser struct {
	wrapper.UserInsert
	tempPassword string
	into         map[string]string
	returnID     bool
	returnIDTO   *int
	role         int8
	phone        string
	email        string
	familyCard   string
}

// UserNew membuat user baru
// mengembalikan struct User {}
func UserNew() *CreateUser {
	return &CreateUser{
		into: make(map[string]string),
		role: wrapper.UserRoleCustomer,
	}
}

// SetFirstName ...
// Set nama depan
func (u *CreateUser) SetFirstName(p string) *CreateUser {
	u.Firstname = strings.ToLower(p)
	u.into["first_name"] = ":first_name"
	return u
}

// SetLastName ...
// Set nama belakang
func (u *CreateUser) SetLastName(p string) *CreateUser {
	// Nama belakang adalah opsional, jadi gak divalidasi
	// maka perlu di cek sebelum di input ke database
	if len(p) > 0 {
		u.Lastname = strings.ToLower(p)
		u.into["last_name"] = ":last_name"
	}
	return u
}

// SetRIC menambahkan nomor ktp
func (u *CreateUser) SetRIC(p string) *CreateUser {
	u.RIC = p
	u.into["ric"] = ":ric"
	return u
}

// SetFamilyCard menambahkan nomor KK
func (u *CreateUser) SetFamilyCard(p string) *CreateUser {
	u.familyCard = p
	return u
}

// SetUserName ...
// Set nilai username
func (u *CreateUser) SetUserName(p string) *CreateUser {
	u.Username = strings.ToLower(p)
	u.into["username"] = ":username"
	return u
}

// SetAvatar tetapkan gambar avatar untuk user ini
func (u *CreateUser) SetAvatar(p string) *CreateUser {
	u.Avatar = p
	u.into["avatar"] = ":avatar"
	return u
}

// SetPassword ...
// Set password
func (u *CreateUser) SetPassword(p string) *CreateUser {
	if p != "" {
		u.tempPassword = p
	}
	return u
}

// SetGender meng set jenis kelamin user
// Set gender/jenis kelamin
func (u *CreateUser) SetGender(p string) *CreateUser {
	u.Gender = strings.ToLower(p)
	u.into["gender"] = ":gender"
	return u
}

// SetOccupation fungsi untuk menambahkan pekerjaan user
func (u *CreateUser) SetOccupation(p string) *CreateUser {
	if p != "" {
		u.Occupation = strings.ToLower(p)
		u.into["occupation"] = ":occupation"
	}
	return u
}

// SetRole tentukan peran/role
func (u *CreateUser) SetRole(p int8) *CreateUser {
	if p >= 0 && p <= 5 {
		u.role = p
	}
	return u
}

// SetPhone fungsi untuk menambahkan nomor
// hp untuk user baru
func (u *CreateUser) SetPhone(p string) *CreateUser {
	u.phone = p
	return u
}

// SetEmail fungsi untuk menambahkan email
func (u *CreateUser) SetEmail(p string) *CreateUser {
	u.email = strings.ToLower(p)
	return u
}

// SetCreatedBy tentukan admin/user yang membuat user ini
func (u *CreateUser) SetCreatedBy(p int) *CreateUser {
	if p > 0 {
		u.CreatedBy = p
		u.into["created_by"] = ":created_by"
	}
	return u
}

// Save ...
// Simpan user ke database
func (u *CreateUser) Save(db *sqlx.DB) error {
	if len(u.tempPassword) > 0 {
		hash, err := hashPassword(u.tempPassword)
		if err != nil {
			return err
		}
		u.Password = hash
		u.into["password"] = ":password"
	}

	// Set avatar default
	if u.Avatar == "" {
		if u.Gender == "m" {
			u.SetAvatar("male.png")
		} else if u.Gender == "f" {
			u.SetAvatar("female.png")
		}
	}

	// Mulai transaksi
	tx := db.MustBegin()
	var tempReturnID int
	userInsertQuery := fmt.Sprintf(`INSERT INTO "user" %s`, u.generateInsertQuery())
	if rows, err := tx.NamedQuery(userInsertQuery, u); err == nil {
		// Ambil id dari transaksi terakhir
		if rows.Next() {
			rows.Scan(&tempReturnID)
		}

		if u.returnID && tempReturnID != 0 {
			*u.returnIDTO = tempReturnID
		}

		if err := rows.Close(); err != nil {
			return err
		}
	} else {
		tx.Rollback()
		return err
	}

	// Lakukan pengecekan nomor kk/ insert jika belum ada
	var fcid int
	if len(u.familyCard) > 0 {
		if fe, fv := UserFamilyCardExist(db, u.familyCard); fe {
			fcid = fv
		} else {
			if rows, err := tx.Query(`INSERT INTO "family_card" (number) VALUES ($1) RETURNING id`, u.familyCard); err == nil {
				// Ambil id dari transaksi terakhir
				if rows.Next() {
					rows.Scan(&fcid)
				}
				if err := rows.Close(); err != nil {
					log.Println("ERROR: user.go Save() fcid closing rows")
					return err
				}
			} else {
				tx.Rollback()
				log.Println("ERROR: user.go Save() Insert family card rollback")
				return err
			}
		}
	}

	// Set family
	if fcid != 0 {
		if _, err := tx.Exec(`INSERT INTO "family" (family_card_id, user_id) VALUES ($1, $2)`, fcid, tempReturnID); err != nil {
			log.Println("ERROR: user.go Save() Insert set family")
			return err
		}
	}

	var username string
	if u.Username != "" {
		username = u.Username
	} else {
		// Set username/kode pelanggan
		num := fmt.Sprintf("% 10d", tempReturnID)

		if misc.CountDigits(tempReturnID) < 10 {
			// Generate angka acak
			if fill, err := misc.GenerateNumberFixedLength(strings.Count(num, " ")); err == nil {
				num = fmt.Sprintf("%s%s", fill, num)
			}
		}

		// Hapus semua spasi
		num = strings.ReplaceAll(num, " ", "")
		if iui, err := strconv.ParseUint(num, 10, 64); err == nil {
			nzne := time.Now().Hour()
			if nzne >= 12 {
				username = "NZ-" + base36.Encode(iui)
			} else {
				username = "NE-" + base36.Encode(iui)
			}
		} else {
			log.Println("ERROR: user.go Save() ParseUint username/kode pelanggan")
			return err
		}
	}

	// Simpan username
	if _, err := tx.Exec(`UPDATE "user"	SET username=$1	WHERE id=$2`, username, tempReturnID); err != nil {
		log.Println("ERROR: user.go Save() Insert username/kode pelanggan")
		return err
	}

	// Set role user
	if _, err := tx.Exec(`INSERT INTO "user_role" (role_id, user_id) VALUES ($1, $2)`, u.role, tempReturnID); err != nil {
		log.Println("ERROR: user.go Save() Insert user role")
		return err
	}

	// Simpan nomor hp user jika tersedia
	if len(u.phone) >= 6 && len(u.phone) <= 15 {
		if _, err := tx.Exec(`INSERT INTO "phone" (user_id, phone) VALUES ($1, $2)`, tempReturnID, u.phone); err != nil {
			log.Println("ERROR: user.go Save() Insert phone number")
			return err
		}
	}

	// Simpan email user jika ada
	if u.email != "" {
		if _, err := tx.Exec(`INSERT INTO "email" (user_id, email) VALUES ($1, $2)`, tempReturnID, u.email); err != nil {
			log.Println("ERROR: user.go Save() Insert email")
			return err
		}
	}

	// Komit
	err := tx.Commit()
	return err
}

// Insert query berdasarka data yang diisi
func (u CreateUser) generateInsertQuery() string {
	iq := u.into
	var kk []string
	var kv []string
	for k, v := range iq {
		kk = append(kk, k)
		kv = append(kv, v)
	}
	result := fmt.Sprintf("(%s) VALUES (%s) RETURNING id", strings.Join(kk, ","), strings.Join(kv, ","))

	return result
}

// ReturnID ...
// Mengembalikan ID user terakhir
func (u *CreateUser) ReturnID(id *int) *CreateUser {
	u.returnID = true
	u.returnIDTO = id
	return u
}

// Hash password dan secret/salt
// jika password telah di tetapkan
func hashPassword(key string) (string, error) {
	var hashedpwd []byte
	var result string
	var err error
	if len(key) > 0 {
		hashedpwd, err = bcrypt.GenerateFromPassword([]byte(key), 10)
		if err != nil {
			return "", err
		}
		result = string(hashedpwd)
	}
	return result, nil
}

// Cocokan password
func matchPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

// Login - User login
// loginid: ID login dari input user
// password: Password dari input user
func Login(db *sqlx.DB, loginid, password string) (int, error) {
	var userid int
	// Cari ID user berdasarkan login id
	var query = `SELECT id
	FROM "user"
	WHERE username=$1
	UNION
	SELECT user_id
	FROM "phone"
	WHERE phone=$1
	UNION
	SELECT user_id
	FROM "email"
	WHERE email=$1`
	err := db.Get(&userid, query, loginid)
	if err != nil {
		return 0, err
	}

	// Cocokan password user
	var pwd string
	query = `SELECT password 
	FROM "user"
	WHERE id=$1`
	err = db.Get(&pwd, query, userid)
	if err != nil {
		return 0, err
	}
	err = matchPassword(pwd, password)
	if err != nil {
		return 0, err
	}
	return userid, err
}

/////////
/* GET */
/////////

// UserGetNullableUserByID - mengambil data user berdasarkan ID
func UserGetNullableUserByID(db *sqlx.DB, userid int) (wrapper.NullableUser, error) {
	var user wrapper.NullableUser
	query := `SELECT
    u.first_name,
    u.last_name,
    u.username,
    u.avatar,
    u.gender,
    u.created_at,
    u.balance,
    string_agg(DISTINCT p.phone, ',' ORDER BY p.phone) AS phone,
    string_agg(DISTINCT e.email, ',' ORDER BY e.email) AS email,
    r.name AS role
	FROM "user" u
	LEFT JOIN "phone" p ON p.user_id=u.id
	LEFT JOIN "email" e ON e.user_id=u.id
	LEFT JOIN "user_role" ur ON ur.user_id=u.id
	LEFT JOIN "role" r ON r.id=ur.role_id
	WHERE u.id=$1
	GROUP BY u.first_name, u.last_name, u.username, u.avatar, u.gender, u.created_at, u.balance, r.name`
	err := db.Get(&user, query, userid)
	if err != nil {
		return wrapper.NullableUser{}, err
	}

	return user, err
}

// UserGetByID mengambil data pengguna berdasarkan ID produk
func UserGetByID(db *sqlx.DB, uid int) (wrapper.User, error) {
	var user wrapper.User
	var u wrapper.NullableUser
	query := `SELECT
		u.id,
		u.first_name,
		u.last_name,
		u.ric,
		u.username,
		u.avatar,
		u.gender,
		u.occupation,
		u.balance,
		fc.number as family_card_number,
		TO_CHAR(u.balance,'Rp999G999G999G999G999') AS balance,
		TO_CHAR(u.created_at, 'MM/DD/YYYY HH12:MI:SS AM') AS created_at,
		r.name AS role
		FROM "user" u
		LEFT JOIN "user_role" ur ON ur.user_id=u.id
		LEFT JOIN "role" r ON r.id=ur.role_id
		LEFT JOIN "family" f ON f.user_id=u.id
		LEFT JOIN "family_card" fc ON fc.id=f.family_card_id
		WHERE u.id=$1
		LIMIT 1`

	err := db.Get(&u, query, uid)
	if err != nil {
		log.Println("user.go GetUserByID() Select user berdasarkan ID")
		log.Println(err)
		return wrapper.User{}, err
	}

	var emails []wrapper.UserEmail
	if em, err := UserGetEmail(db, uid); err == nil {
		emails = em
	}

	var phones []wrapper.UserPhone
	if ph, err := UserGetPhone(db, uid); err == nil {
		phones = ph
	} else {
		log.Println(err)
	}

	var addresses []wrapper.UserAddress
	if addrs, err := UserGetAddress(db, uid); err == nil {
		addresses = addrs
	}

	user = wrapper.User{
		ID:               u.ID,
		Firstname:        strings.Title(u.Firstname),
		Lastname:         strings.Title(string(u.Lastname.String)),
		RIC:              u.RIC,
		FamilyCardNumber: string(u.FamilyCardNumber.String),
		Username:         string(u.Username.String),
		Avatar:           u.Avatar,
		Gender:           u.Gender,
		Occupation:       strings.Title(string(u.Occupation.String)),
		Role:             strings.Title(u.Role),
		CreatedAt:        u.CreatedAt,
		Balance:          string(u.Balance),
		Emails:           emails,
		Phones:           phones,
		Addresses:        addresses,
	}

	return user, nil
}

// UserGetUserTotalRow menghitung jumlah row pada tabel user
func UserGetUserTotalRow(db *sqlx.DB) (int, error) {
	var total int
	query := `SELECT COUNT(id) FROM "user"`
	err := db.Get(&total, query)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// UserGetPhone mengambil nomor HP berdasarkan ID
func UserGetPhone(db *sqlx.DB, userid int) ([]wrapper.UserPhone, error) {
	var phones []wrapper.UserPhone
	query := `SELECT id, phone, verified
	FROM "phone"
	WHERE user_id=$1`
	err := db.Select(&phones, query, userid)
	if err != nil {
		return []wrapper.UserPhone{}, err
	}
	return phones, err
}

// UserGetEmail mengambil email berdasarkan ID
func UserGetEmail(db *sqlx.DB, userid int) ([]wrapper.UserEmail, error) {
	var emails []wrapper.UserEmail
	query := `SELECT id, email, verified
	FROM "email"
	WHERE user_id=$1`
	err := db.Select(&emails, query, userid)
	if err != nil {
		return []wrapper.UserEmail{}, err
	}
	return emails, err
}

// UserGetAddress mengambil data alamat user
func UserGetAddress(db *sqlx.DB, userid int) ([]wrapper.UserAddress, error) {
	var addresses []wrapper.UserAddress
	query := `SELECT a.id, a.name, a.description, a.one, a.two, a.zip, a.village_id, a.district_id, a.city_id, a.province_id, INITCAP(p.name) AS province_name, INITCAP(c.name) AS city_name, INITCAP(d.name) AS district_name, INITCAP(v.name) AS village_name
	FROM "address" a
	JOIN "province" p ON p.id=a.province_id
	JOIN "city" c ON c.id=a.city_id
	JOIN "district" d ON d.id=a.district_id
	JOIN "village" v ON v.id=a.village_id
	WHERE user_id=$1`
	err := db.Select(&addresses, query, userid)
	if err != nil {
		log.Println("Warning: user.go GetAddress() Tidak ada alamat ditemukan")
		log.Println(err)
		return []wrapper.UserAddress{}, err
	}
	return addresses, err
}

// UserGetRole mengambil role berdasarkan id
func UserGetRole(db *sqlx.DB, userid int) (string, error) {
	var role string
	query := `SELECT
		r.name
		FROM "role" r
		JOIN user_role ur ON ur.role_id=r.id
		WHERE ur.user_id=$1`
	err := db.Get(&role, query, userid)
	if err != nil {
		return "", err
	}

	return role, err
}

// UserGetUsername Mengambil username user
// atau tidak
func UserGetUsername(db *sqlx.DB, uid int) (string, error) {
	// Ambil username
	var uname string
	query := `SELECT username FROM "user" WHERE id=$1`
	err := db.Get(&uname, query, uid)

	return uname, err
}

////////////
/* DELETE */
////////////

// UserDeleteEmail menghapus email
func UserDeleteEmail(db *sqlx.DB, id int64, uid int) error {
	query := `DELETE FROM "email"
	WHERE id=$1 AND user_id=$2`
	_, err := db.Exec(query, id, uid)
	return err
}

// UserDeletePhone menghapus nomor HP
func UserDeletePhone(db *sqlx.DB, id int64, uid int) error {
	query := `DELETE FROM "phone"
	WHERE id=$1 AND user_id=$2`
	_, err := db.Exec(query, id, uid)
	return err
}

// UserDeleteAddress menghapus alamat
func UserDeleteAddress(db *sqlx.DB, id int64, uid int) error {
	query := `DELETE FROM "address"
	WHERE id=$1 AND user_id=$2`
	_, err := db.Exec(query, id, uid)
	return err
}

/////////
/* ADD */
/////////

// UserAddEmail menambahkan email user
func UserAddEmail(db *sqlx.DB, email string, uid int) error {
	query := `INSERT INTO "email" (email, user_id) VALUES ($1, $2)`
	_, err := db.Exec(query, strings.ToLower(email), uid)
	return err
}

// UserAddPhone menambahkan nomor HP ke database
func UserAddPhone(db *sqlx.DB, phone string, uid int) error {
	query := `INSERT INTO "phone" (phone, user_id) VALUES ($1, $2)`
	_, err := db.Exec(query, phone, uid)
	return err
}

// UserAddAddress menambahkan alamat baru
func UserAddAddress(db *sqlx.DB, address wrapper.UserAddress) error {
	query := `INSERT INTO "address" (user_id, name, description, one, two, zip, province_id, city_id, district_id, village_id)
	VALUES (:user_id, :name, :description, :one, :two, :zip, :province_id, :city_id, :district_id, :village_id)`
	_, err := db.NamedExec(query, address)
	if err != nil {
		log.Println("ERROR: user.go AddAddress() Gagal menambahkan alamat")
		log.Println(err)
		return err
	}
	return nil
}

///////////
/* CHECK */
///////////

// UserFamilyCardExist check nomor KK
func UserFamilyCardExist(db *sqlx.DB, fc string) (bool, int) {
	var f int
	query := `SELECT id FROM "family_card" WHERE number=$1`
	err := db.Get(&f, query, fc)
	if err == nil {
		if f != 0 {
			return true, f
		}
	} else {
		log.Println("Warning: user.go FamilyCardExist() KK tidak ditemukan")
		log.Println(err)
	}
	return false, f
}

// UserRICExist check nomor KTP
func UserRICExist(db *sqlx.DB, ric string) bool {
	var r string
	query := `SELECT id FROM "user" WHERE ric=$1`
	err := db.Get(&r, query, ric)
	if err == nil {
		if r != "" {
			return true
		}
	} else {
		log.Println("Warning: NIK KTP tidak ditemukan")
		log.Println(err)
	}
	return false
}

// UserPhoneExist check nomor telepon
func UserPhoneExist(db *sqlx.DB, phone string) bool {
	var p string
	query := `SELECT id FROM "phone" WHERE phone=$1`
	err := db.Get(&p, query, phone)
	if err == nil {
		if p != "" {
			return true
		}
	} else {
		log.Println("Warning: Nomor hp tidak ditemukan")
		log.Println(err)
	}
	return false
}

// UserEmailExist check alamat email
func UserEmailExist(db *sqlx.DB, email string) bool {
	var p string
	query := `SELECT id FROM "email" WHERE email=$1`
	err := db.Get(&p, query, email)
	if err == nil {
		if p != "" {
			return true
		}
	} else {
		log.Println("Warning: Email tidak ditemukan")
		log.Println(err)
	}
	return false
}

// UsernameExist mengecek apakah username tersedia
// atau tidak
func UsernameExist(db *sqlx.DB, uname string) bool {
	// Check bila username sudah ada di database
	var indb string
	query := `SELECT username FROM "user" WHERE username=$1`
	err := db.Get(&indb, query, uname)
	if err == nil {
		if indb != "" {
			return true
		}
	}

	return false
}

////////////
/* UPDATE */
////////////

// UserUpdateRole mengubah role user
func UserUpdateRole(db *sqlx.DB, uid, roleid int) error {
	query := `UPDATE "user_role"
	SET role_id=$1
	WHERE user_id=$2`
	_, err := db.Exec(query, roleid, uid)

	return err
}

// UserUpdateUsername mengubah username user
func UserUpdateUsername(db *sqlx.DB, uid int, uname string) error {
	query := `UPDATE "user"
	SET username=$1
	WHERE id=$2`
	_, err := db.Exec(query, uname, uid)

	return err
}

// UserUpdatePassword mengubah password user
func UserUpdatePassword(db *sqlx.DB, uid int, pwd string) error {
	// Hash password
	pwd, err := hashPassword(pwd)
	if err != nil {
		return err
	}
	query := `UPDATE "user"
	SET password=$1
	WHERE id=$2`
	_, err = db.Exec(query, pwd, uid)

	return err
}

///////////
/* MATCH */
///////////

// UserMatchPassword mencocokan password input user
// dengan yang ada di database
func UserMatchPassword(db *sqlx.DB, uid int, pwd string) bool {
	// Cocokan password user
	var inpwd string
	var next bool
	query := `SELECT password 
	FROM "user"
	WHERE id=$1`
	err := db.Get(&inpwd, query, uid)
	if err == nil {
		if inpwd != "" {
			next = true
		}
	} else {
		log.Println("Warning: user.go MatchPassword() password tidak cocok")
		log.Println(err)
	}
	if next {
		err = matchPassword(inpwd, pwd)
		if err != nil {
			next = false
		}
	}
	return next
}
