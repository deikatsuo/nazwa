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

// UUserAddress menampung alamat user
type UUserAddress struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	One      string `db:"one"`
	Two      string `db:"two"`
	Zip      string `db:"zip"`
	Province string `db:"province"`
	City     string `db:"city"`
	District string `db:"district"`
	Village  string `db:"village"`
}

// UserAddress menampung data alamat user
type UserAddress struct {
	ID       int    `db:"id"`
	UserID   int    `db:"user_id"`
	Name     string `db:"name" json:"name" binding:"required,min=5,max=50"`
	One      string `db:"one" json:"one" binding:"required,min=5,max=80"`
	Two      string `db:"two" json:"two" binding:"omitempty,min=5,max=80"`
	Zip      string `db:"zip" json:"zip" binding:"required,numeric,min=5,max=5"`
	Province string `db:"province_id" json:"province" binding:"required,numeric"`
	City     string `db:"city_id" json:"city" binding:"required,numeric"`
	District string `db:"district_id" json:"district" binding:"required,numeric"`
	Village  string `db:"village_id" json:"village" binding:"required,numeric"`
}
