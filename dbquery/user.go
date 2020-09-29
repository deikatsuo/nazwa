package dbquery

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// User ...
// Struct data user
type User struct {
	ID        string
	Firstname string `db:"first_name"`
	Lastname  string `db:"last_name"`
	Username  string `db:"username"`
	Avatar    string `db:"avatar"`
	Password  string `db:"password"`
	Gender    string `db:"gender"`
	CreatedAt string `db:"created_at"`
	Balance   int    `db:"balance"`

	tempPassword string
	into         map[string]string
	returnid     bool
	returnidto   *int
}

// NewUser ...
// Membuat user baru
// mengembalikan struct User {}
func NewUser() *User {
	return &User{
		into: make(map[string]string),
	}
}

// SetFirstName ...
// Set nama depan
func (u *User) SetFirstName(p string) *User {
	u.Firstname = p
	u.into["first_name"] = ":first_name"
	return u
}

// SetLastName ...
// Set nama belakang
func (u *User) SetLastName(p string) *User {
	u.Lastname = p
	u.into["last_name"] = ":last_name"
	return u
}

// SetUserName ...
// Set nilai username
func (u *User) SetUserName(p string) *User {
	u.Username = p
	u.into["username"] = ":username"
	return u
}

// SetAvatar ...
// Set avatar
func (u *User) SetAvatar(p string) *User {
	u.Avatar = p
	u.into["avatar"] = ":avatar"
	return u
}

// SetPassword ...
// Set password
func (u *User) SetPassword(p string) *User {
	u.tempPassword = p
	return u
}

// Hash password dan secret/salt
// jika password telah di tetapkan
func (u *User) hashPassword() error {
	var hashedpwd []byte
	var err error
	if len(u.tempPassword) > 0 {
		hashedpwd, err = bcrypt.GenerateFromPassword([]byte(u.tempPassword), 10)
		if err == nil {
			u.Password = string(hashedpwd)
			u.into["password"] = ":password"
		}
	}

	return err
}

// Cocokan password
func matchPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println(err)
	return err == nil
}

// SetGender ...
// Set gender/jenis kelamin
func (u *User) SetGender(p string) *User {
	u.Gender = p
	u.into["gender"] = ":gender"
	return u
}

// Save ...
// Simpan user ke database
func (u *User) Save(db *sqlx.DB) error {
	if u.hashPassword() != nil {
		return errors.New("Gagal meng enkripsi password")
	}
	tx := db.MustBegin()
	userInsertQuery := fmt.Sprintf(`INSERT INTO "user" %s`, u.generateInsertQuery())
	rows, err := tx.NamedQuery(userInsertQuery, u)
	if err != nil {
		return err
	}
	if u.returnid {
		if rows.Next() {
			rows.Scan(u.returnidto)
		}
	}
	err = tx.Commit()
	return err
}

// Insert query berdasarka data yang diisi
func (u User) generateInsertQuery() string {
	iq := u.into
	var kk []string
	var kv []string
	for k, v := range iq {
		kk = append(kk, k)
		kv = append(kv, v)
	}
	var result string
	if u.returnid {
		result = fmt.Sprintf("(%s) VALUES (%s) RETURNING id", strings.Join(kk, ","), strings.Join(kv, ","))
	} else {
		result = fmt.Sprintf("(%s) VALUES (%s)", strings.Join(kk, ","), strings.Join(kv, ","))
	}
	return result
}

// ReturnID ...
// Mengembalikan ID user terakhir
func (u *User) ReturnID(id *int) *User {
	u.returnid = true
	u.returnidto = id
	return u
}
