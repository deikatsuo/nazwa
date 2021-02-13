package wrapper

import "database/sql"

// AddressForm form data alamat
type AddressForm struct {
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

// AddressInsert insert data alamat
type AddressInsert struct {
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

// AddressSelect select data alamat
type AddressSelect struct {
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

// Address menampung data alamat
type Address struct {
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
