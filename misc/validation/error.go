package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// SimpleValErr memparse error dari `validator` menjadi lebih simple
//
// Dikembalikan error versi sederhana
func SimpleValErr(verr interface{}) string {
	ve := verr.(validator.ValidationErrors)

	/*
		for _, v := range ve {
			fmt.Println("Namespace ", v.Namespace())
			fmt.Println("Field ", v.Field())
			fmt.Println("Struct Namespace ", v.StructNamespace())
			fmt.Println("Struct Field ", v.StructField())
			fmt.Println("Tag ", v.Tag())
			fmt.Println("Actual Tag ", v.ActualTag())
			fmt.Println("Kind ", v.Kind())
			fmt.Println("Type ", v.Type())
			fmt.Println("Value ", v.Value())
			fmt.Println("Param ", v.Param())
		}
	*/

	var erbar string
	for _, v := range ve {
		for mi, mv := range ValidationErrorsMask {
			if v.Tag() == mi {
				tmperbar := erbar + v.Field() + " " + mv + " \n"
				if v.Tag() == "min" {
					tmperbar = erbar + v.Field() + " " + fmt.Sprintf(mv, v.Param()) + "\n"
				}
				if v.Tag() == "max" {
					tmperbar = erbar + v.Field() + " " + fmt.Sprintf(mv, v.Param()) + "\n"
				}
				if v.Tag() == "oneof" {
					tmperbar = erbar + v.Field() + " " + fmt.Sprintf(mv, strings.Replace(v.Param(), " ", "' atau '", 1)) + "\n"
				}
				erbar = tmperbar
			}
		}
	}

	return erbar
}

// SimpleValErrMap menyederhanakan error dari validator
// dan mengembalikannya dalam bentuk map
func SimpleValErrMap(verr interface{}) map[string]interface{} {
	ve := verr.(validator.ValidationErrors)
	errmap := make(map[string]interface{})
	for _, v := range ve {
		for mi, mv := range ValidationErrorsMask {
			if v.Tag() == mi {
				errmap[strings.ToLower(v.StructField())] = mv
				if v.Tag() == "min" || v.Tag() == "max" {
					errmap[strings.ToLower(v.StructField())] = fmt.Sprintf(mv, v.Param())
				}
				if v.Tag() == "oneof" {
					errmap[strings.ToLower(v.StructField())] = fmt.Sprintf(mv, strings.Replace(v.Param(), " ", "' atau '", 1))
				}
			}
		}
	}
	return errmap
}

// ValidationErrorsMask mengubah pesan error default dari validator
var ValidationErrorsMask = map[string]string{
	"alpha":           "Harus menggunakan huruf alphabet",
	"alphanum":        "Harus huruf alphabet atau nomor",
	"alphanumunicode": "Hanya boleh menggunakan huruf alphabet, nomor dan simbol",
	"required":        "Tidak boleh kosong",
	"min":             "Tidak boleh kurang dari %s",
	"max":             "Tidak boleh lebih dari %s",
	"oneof":           "Hanya boleh di isi oleh '%s'",
	"numeric":         "Hanya boleh di isi dengan nomor",
	"email":           "Format email salah",
	"eqfield":         "Kolom ini harus sama",
	"required_with":   "Kolom ini perlu diisi",
	"base64":          "Harus string base64",
}
