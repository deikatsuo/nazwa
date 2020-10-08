package wrapper

import "database/sql"

// DefaultConfig - untuk menyimpan konfigurasi bawaan
type DefaultConfig struct {
	Site map[string]interface{}
}

// NullableUser struk user dari database
type NullableUser struct {
	ID        int
	Firstname string         `db:"first_name"`
	Lastname  sql.NullString `db:"last_name"`
	Username  sql.NullString `db:"username"`
	Avatar    string         `db:"avatar"`
	Gender    string         `db:"gender"`
	CreatedAt string         `db:"created_at"`
	Balance   []uint8        `db:"balance"`
	Password  sql.NullString `db:"password"`
	Phone     sql.NullString `db:"phone"`
	Email     sql.NullString `db:"email"`
	Role      string         `db:"role"`
	Emails    []UserEmail
	Phones    []UserPhone
}

// Fullname menampilkan nama lengkap
func (u NullableUser) Fullname() string {
	return u.Firstname + " " + u.Lastname.String
}

// UserEmail menampung email user
type UserEmail struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Verified bool   `db:"verified"`
}

// UserPhone menampung email user
type UserPhone struct {
	ID       int    `db:"id"`
	Phone    string `db:"phone"`
	Verified bool   `db:"verified"`
}
