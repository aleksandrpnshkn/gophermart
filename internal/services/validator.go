package services

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidate(uni *AppUni) *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	uni.RegisterValidationTranslations(validate)

	// использовать имя поля из тега json вместо имени поля структуры
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("luhn", func(fl validator.FieldLevel) bool {
		return IsValidLuhnNumber(fl.Field().String())
	})

	return validate
}
