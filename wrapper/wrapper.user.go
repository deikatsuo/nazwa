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
	Addresses        []Address
}

// User buat nge map data user
type User struct {
	ID               int
	Firstname        string
	Lastname         string
	Password         bool
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
	Addresses        []Address
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
