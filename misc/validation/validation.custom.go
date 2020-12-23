package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// CustomValidationDate validator tanggal
func CustomValidationDate() validator.Func {
	return func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(string)
		if ok {
			_, err := time.Parse("2006-01-02", date)
			if err != nil {
				return false
			}
		}

		return true
	}
}
