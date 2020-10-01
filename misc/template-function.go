package misc

import "html/template"

// karena balance yang di ambil dari database
// formatnya adalah []uint8 maka harus dirubah dulu
// ke string agar menjadi angka balance yang sebenarnya
func balance(u []uint8) string {
	return string(u)
}

// RegTmplFunc - mendaftarkan fungsi ke template
func RegTmplFunc() template.FuncMap {
	return template.FuncMap{
		"balance": balance,
	}
}
