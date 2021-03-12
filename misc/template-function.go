package misc

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// Balance - karena balance yang di ambil dari database
// formatnya adalah []uint8 maka harus dirubah dulu
// ke string agar menjadi angka balance yang sebenarnya
func Balance(u []uint8) string {
	return string(u)
}

// Title Membuat huruf pertama pada string menjadi kapital
func Title(s string) string {
	return strings.Title(s)
}

// MayNull memeriksa jika data dari database berupa null
// dan merubahnya ke string
func MayNull(s interface{}) string {
	var val string
	switch v := s.(type) {
	case sql.NullBool:
		val = strconv.FormatBool(v.Bool)
	case sql.NullFloat64:
		val = fmt.Sprintf("%f", v.Float64)
	case sql.NullInt32:
		val = string(v.Int32)
	case sql.NullInt64:
		val = fmt.Sprint(v.Int64)
	case sql.NullString:
		val = v.String
	case sql.NullTime:
		val = v.Time.String()
	default:
	}
	return val
}
