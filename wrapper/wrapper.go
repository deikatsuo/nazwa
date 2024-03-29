package wrapper

import "database/sql"

// DefaultConfig - untuk menyimpan konfigurasi bawaan
type DefaultConfig struct {
	Info map[string]interface{}
}

// IntString kadang struk gini dibutuhkan juga hhha
type IntString struct {
	Integer sql.NullInt32  `db:"is_int"`
	String  sql.NullString `db:"is_string"`
}

// NameID menampilkan nama dan id
type NameID struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Thumbnail string `db:"thumbnail"`
}

// NameIDCode menampilkan nama, id, dan kode
type NameIDCode struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Code      string `db:"code"`
	Thumbnail string `db:"thumbnail"`
}

// NameIDNameID membandingkan dua nama
type NameIDNameID struct {
	NameOne   string `db:"name_one"`
	NameOneID int    `db:"name_one_id"`
	NameTwo   string `db:"name_two"`
	NameTwoID int    `db:"name_two_id"`
}
