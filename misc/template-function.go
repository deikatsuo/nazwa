package misc

import (
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
