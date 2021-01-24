package wrapper

// DefaultConfig - untuk menyimpan konfigurasi bawaan
type DefaultConfig struct {
	Info map[string]interface{}
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

// NameName membandingkan dua nama
type NameName struct {
	NameOne string `db:"name_one"`
	NameTwo string `db:"name_two"`
}
