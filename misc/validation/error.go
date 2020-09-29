package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// SimpleVErr memparse error dari `validator` menjadi lebih simple
//
// Dikembalikan error versi sederhana
func SimpleVErr(verr interface{}) string {
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

// ValidationErrorsMask mengubah pesan error default dari validator
var ValidationErrorsMask = map[string]string{
	"alpha":           "harus menggunakan huruf alphabet",
	"alphanum":        "harus huruf alphabet atau nomor",
	"alphanumunicode": "hanya boleh menggunakan huruf alphabet, nomor dan simbol",
	"required":        "tidak boleh kosong",
	"min":             "tidak boleh kurang dari %s",
	"oneof":           "hanya boleh di isi oleh '%s'",
	"numeric":         "hanya boleh di isi dengan nomor",
	"email":           "format email salah",
}
