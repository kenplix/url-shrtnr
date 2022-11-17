package v1

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

func usernameValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		value := field.String()
		if n := utf8.RuneCountInString(value); n < 5 || n > 32 {
			return false
		}

		if unicode.IsDigit([]rune(value)[0]) {
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
	switch field.Kind() {
	case reflect.String:
		value := field.String()
		if n := utf8.RuneCountInString(value); n < 8 || n > 64 {
			return false
		}

		var hasUpper, hasLower, hasDigit, hasSpecial bool
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
