package v1

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

func usernameValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.String {
		value := field.String()
		if n := utf8.RuneCountInString(value); n < 5 || n > 32 {
			return false
		}

		r, _ := utf8.DecodeRuneInString(value)
		if unicode.IsDigit(r) {
			return false
		}

		if strings.HasPrefix(value, "_") || strings.Contains(value, "__") || strings.HasSuffix(value, "_") {
			return false
		}

		for _, char := range value {
			if !unicode.IsLetter(char) && !unicode.IsDigit(char) && !strings.ContainsRune("_", char) {
				return false
			}
		}

		return true
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func passwordValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.String {
		value := field.String()
		if n := utf8.RuneCountInString(value); n < 8 || n > 64 {
			return false
		}

		var (
			hasUpper   bool
			hasLower   bool
			hasDigit   bool
			hasSpecial bool
		)

		for _, char := range value {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasDigit = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}

		return hasUpper && hasLower && hasDigit && hasSpecial
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}
