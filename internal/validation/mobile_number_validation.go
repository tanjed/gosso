package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^(?:\+88|88)?01[3-9]\d{8}$`)
    return re.MatchString(fl.Field().String())
}