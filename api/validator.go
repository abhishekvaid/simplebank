package api

import (
	"github.com/go-playground/validator/v10"
	"himavisoft.simple_bank/util"
)

var validCurrency = func(fl validator.FieldLevel) bool {
	if str, ok := fl.Field().Interface().(string); ok {
		return util.ValidateCurrency(str)
	}
	return false
}
