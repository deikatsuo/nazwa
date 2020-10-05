package middleware

import (
	"html/template"
	"nazwa/misc"
)

// RegTmplFunc - mendaftarkan fungsi ke template
func RegTmplFunc() template.FuncMap {
	return template.FuncMap{
		"balance": misc.Balance,
		"title":   misc.Title,
	}
}
