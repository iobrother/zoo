package validate

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	phone = "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,1,3,5-8])|(18[0-9])|166|191|198|199|(147))\\d{8}$"
)

var (
	rxPhone = regexp.MustCompile(phone)
)

func ValidateMobile(fl validator.FieldLevel) bool {
	return rxPhone.MatchString(fl.Field().String())
}
