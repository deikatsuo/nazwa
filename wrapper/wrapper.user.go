package wrapper

import "database/sql"

// Definisikan ROLE
const (
	// UserRoleDev Tingkatan teratas
	UserRoleDev int8 = 0
	// UserRoleAdmin memiliki akses penuh
	// sebagai admin
	UserRoleAdmin int8 = 1
	// UserRoleCollector sebagai seorang penagih
	UserRoleCollector int8 = 2
	// UserRoleDriver sebagai supir
	UserRoleDriver int8 = 3
	// UserRoleSurveyor adalah penyurvey
	UserRoleSurveyor int8 = 4
	// UserRoleSales bagian marketing atau pemasaran
	UserRoleSales int8 = 5
	// UserRoleCustomer pelanggan
	UserRoleCustomer int8 = 6
	// UserRoleSubstitute pendamping/pengganti
	UserRoleSubstitute int8 = 7
)

// UserInsert base struk
type UserInsert struct {
	Firstname  string  `db:"first_name"`
	Lastname   string  `db:"last_name"`
	Username   string  `db:"username"`
	Avatar     string  `db:"avatar"`
	Gender     string  `db:"gender"`
	Occupation string  `db:"occupation"`
	CreatedAt  string  `db:"created_at"`
	Balance    []uint8 `db:"balance"`
	Password   string  `db:"password"`
	Role       string  `db:"role"`
	RIC        string  `db:"ric"`
	CreatedBy  int     `db:"created_by"`
}

// UserForm formulir buat/create user
type UserForm struct {
	FC         string `json:"fc" binding:"omitempty,numeric,min=16,max=16"`
	RIC        string `json:"ric" binding:"numeric,min=16,max=16"`
	Phone      string `json:"phone" binding:"omitempty,numeric,min=6,max=15"`
	Firstname  string `json:"firstname" binding:"required,min=3,max=25"`
	Lastname   string `json:"lastname" binding:"omitempty,min=1,max=25"`
	Gender     string `json:"gender" binding:"required,oneof=m f"`
	Occupation string `json:"occupation" binding:"omitempty,min=4,max=25"`
	Password   string `json:"password" binding:"omitempty,alphanumunicode,min=8,max=25"`
	Repassword string `json:"repassword" binding:"eqfield=Password"`
	Photo      string `json:"photo" binding:"omitempty,base64"`
	PhotoType  string `json:"photo_type"`
}

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

// UserAddressForm form data alamat user
type UserAddressForm struct {
	Name        string `json:"name" binding:"required,min=5,max=50"`
	Description string `json:"description" binding:"omitempty,max=80"`
	One         string `json:"one" binding:"required,min=5,max=80"`
	Two         string `json:"two" binding:"omitempty,max=80"`
	Zip         string `json:"zip" binding:"omitempty,numeric,min=5,max=5"`
	Province    string `json:"province" binding:"required,numeric"`
	City        string `json:"city" binding:"required,numeric"`
	District    string `json:"district" binding:"required,numeric"`
	Village     string `json:"village" binding:"required,numeric"`
}

// UserAddressInsert insert data alamat user
type UserAddressInsert struct {
	UserID      int    `db:"user_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	One         string `db:"one"`
	Two         string `db:"two"`
	Zip         string `db:"zip"`
	Province    string `db:"province_id"`
	City        string `db:"city_id"`
	District    string `db:"district_id"`
	Village     string `db:"village_id"`
}

// UserAddressSelect select data alamat user
type UserAddressSelect struct {
	ID           int            `db:"id"`
	UserID       int            `db:"user_id"`
	Name         string         `db:"name"`
	Description  sql.NullString `db:"description"`
	One          string         `db:"one"`
	Two          sql.NullString `db:"two"`
	Zip          sql.NullString `db:"zip"`
	Province     string         `db:"province_id"`
	City         string         `db:"city_id"`
	District     string         `db:"district_id"`
	Village      string         `db:"village_id"`
	ProvinceName string         `db:"province_name"`
	CityName     string         `db:"city_name"`
	DistrictName string         `db:"district_name"`
	VillageName  string         `db:"village_name"`
}

// UserAddress menampung data alamat user
type UserAddress struct {
	ID           int
	UserID       int
	Name         string
	Description  string
	One          string
	Two          string
	Zip          string
	Province     string
	City         string
	District     string
	Village      string
	ProvinceName string
	CityName     string
	DistrictName string
	VillageName  string
}
