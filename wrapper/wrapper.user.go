package wrapper

import "database/sql"

// NullableUser struk user dari database
type NullableUser struct {
	ID               int            `db:"id"`
	Firstname        string         `db:"first_name"`
	Lastname         sql.NullString `db:"last_name"`
	Username         sql.NullString `db:"username"`
	RIC              string         `db:"ric"`
	Occupation       sql.NullString `db:"occupation"`
	Avatar           string         `db:"avatar"`
	Gender           string         `db:"gender"`
	CreatedAt        string         `db:"created_at"`
	Balance          []uint8        `db:"balance"`
	Password         sql.NullString `db:"password"`
	Phone            sql.NullString `db:"phone"`
	Email            sql.NullString `db:"email"`
	Role             string         `db:"role"`
	FamilyCardNumber sql.NullString `db:"family_card_number"`
	Emails           []UserEmail
	Phones           []UserPhone
	Addresses        []UserAddress
}

// User buat nge map data user
type User struct {
	ID               int
	Firstname        string
	Lastname         string
	RIC              string
	Username         string
	Occupation       string
	Avatar           string
	Gender           string
	CreatedAt        string
	Balance          string
	Role             string
	FamilyCardNumber string
	Emails           []UserEmail
	Phones           []UserPhone
	Addresses        []UserAddress
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

// UserAddress menampung data alamat user
type UserAddress struct {
	ID           int    `db:"id"`
	UserID       int    `db:"user_id"`
	Name         string `db:"name" json:"name" binding:"required,min=5,max=50"`
	Description  string `db:"description" json:"description" binding:"omitempty,max=80"`
	One          string `db:"one" json:"one" binding:"required,min=5,max=80"`
	Two          string `db:"two" json:"two" binding:"omitempty,max=80"`
	Zip          string `db:"zip" json:"zip" binding:"omitempty,numeric,min=5,max=5"`
	Province     string `db:"province_id" json:"province" binding:"required,numeric"`
	City         string `db:"city_id" json:"city" binding:"required,numeric"`
	District     string `db:"district_id" json:"district" binding:"required,numeric"`
	Village      string `db:"village_id" json:"village" binding:"required,numeric"`
	ProvinceName string `db:"province_name"`
	CityName     string `db:"city_name"`
	DistrictName string `db:"district_name"`
	VillageName  string `db:"village_name"`
}
