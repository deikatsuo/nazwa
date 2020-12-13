package wrapper

// DefaultConfig - untuk menyimpan konfigurasi bawaan
type DefaultConfig struct {
	Site map[string]interface{}
}

// NameID menampilkan nama dan id
type NameID struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// NameIDCode menampilkan nama, id, dan kode
type NameIDCode struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Code string `db:"code"`
}
