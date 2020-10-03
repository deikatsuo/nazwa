package dbquery

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// User base struk
type User struct {
	Firstname string  `db:"first_name"`
	Lastname  string  `db:"last_name"`
	Username  string  `db:"username"`
	Avatar    string  `db:"avatar"`
	Gender    string  `db:"gender"`
	CreatedAt string  `db:"created_at"`
	Balance   []uint8 `db:"balance"`
	Password  string  `db:"password"`
	Phone     string  `db:"phone"`
	Email     string  `db:"email"`
	Role      string  `db:"role"`
}

// CreateUser struk buat user baru
// Struct data user
type CreateUser struct {
	User
	tempPassword string
	into         map[string]string
	returnID     bool
	returnIDTO   *int
	role         int8
	phone        string
	email        string
}

const (
	// RoleAdmin memiliki akses penuh
	// sebagai admin
	RoleAdmin int8 = 1
	// RoleSurveyor adalah penyurvey
	RoleSurveyor int8 = 2
	// RoleCollector sebagai seorang penagih
	RoleCollector int8 = 3
	// RoleSales bagian marketing atau pemasaran
	RoleSales int8 = 4
	// RoleCustomer pelanggan
	RoleCustomer int8 = 5
)

// NewUser ...
// Membuat user baru
// mengembalikan struct User {}
func NewUser() *CreateUser {
	return &CreateUser{
		into: make(map[string]string),
		role: RoleCustomer,
	}
}

// SetFirstName ...
// Set nama depan
func (u *CreateUser) SetFirstName(p string) *CreateUser {
	u.Firstname = p
	u.into["first_name"] = ":first_name"
	return u
}

// SetLastName ...
// Set nama belakang
func (u *CreateUser) SetLastName(p string) *CreateUser {
	u.Lastname = p
	u.into["last_name"] = ":last_name"
	return u
}

// SetUserName ...
// Set nilai username
func (u *CreateUser) SetUserName(p string) *CreateUser {
	u.Username = p
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
	u.tempPassword = p
	return u
}

// SetGender meng set jenis kelamin user
// Set gender/jenis kelamin
func (u *CreateUser) SetGender(p string) *CreateUser {
	u.Gender = p
	u.into["gender"] = ":gender"
	return u
}

// SetRole tentukan peran/role
func (u *CreateUser) SetRole(p int8) *CreateUser {
	if p >= 1 && p <= 5 {
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
	u.email = p
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
	userInsertQuery := fmt.Sprintf(`INSERT INTO "user" %s`, u.generateInsertQuery())
	rows, err := tx.NamedQuery(userInsertQuery, u)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Ambil id dari transaksi terakhir
	if u.returnID {
		if rows.Next() {
			rows.Scan(u.returnIDTO)
		}
		err = rows.Close()
		if err != nil {
			return err
		}
	}

	// Set role user
	_, err = tx.Exec(`INSERT INTO "user_role" (role_id, user_id) VALUES ($1, $2)`, u.role, *u.returnIDTO)
	if err != nil {
		return err
	}

	// Simpan nomor hp user jika tersedia
	if len(u.phone) >= 6 && len(u.phone) <= 15 {
		_, err = tx.Exec(`INSERT INTO "phone" (user_id, phone) VALUES ($1, $2)`, *u.returnIDTO, u.phone)
		if err != nil {
			return err
		}
	}

	// Simpan email user jika ada
	if u.email != "" {
		_, err = tx.Exec(`INSERT INTO "email" (user_id, email) VALUES ($1, $2)`, *u.returnIDTO, u.email)
		if err != nil {
			return err
		}
	}

	// Komit
	err = tx.Commit()
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
	var result string
	if u.returnID {
		result = fmt.Sprintf("(%s) VALUES (%s) RETURNING id", strings.Join(kk, ","), strings.Join(kv, ","))
	} else {
		result = fmt.Sprintf("(%s) VALUES (%s)", strings.Join(kk, ","), strings.Join(kv, ","))
	}
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

// Fullname mengambil nama lengkap
func (u User) Fullname() string {
	return u.Firstname + " " + u.Lastname
}

// GetUser - mengambil data user berdasarkan ID
func GetUser(db *sqlx.DB, userid int) (User, error) {
	var user User
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
	JOIN "phone" p ON p.user_id=u.id
	JOIN "email" e ON e.user_id=u.id
	JOIN "user_role" ur ON ur.user_id=u.id
	JOIN "role" r ON r.id=ur.role_id
	WHERE u.id=$1
	GROUP BY u.first_name, u.last_name, u.username, u.avatar, u.gender, u.created_at, u.balance, r.name`
	err := db.Get(&user, query, userid)
	if err != nil {
		return User{}, err
	}

	return user, err
}
