package misc

import (
	"html/template"
	"strings"
)

// karena balance yang di ambil dari database
// formatnya adalah []uint8 maka harus dirubah dulu
// ke string agar menjadi angka balance yang sebenarnya
func balance(u []uint8) string {
	return string(u)
}

func title(s string) string {
	return strings.Title(s)
}

// RegTmplFunc - mendaftarkan fungsi ke template
func RegTmplFunc() template.FuncMap {
	return template.FuncMap{
		"balance": balance,
		"title":   title,
	}
}
