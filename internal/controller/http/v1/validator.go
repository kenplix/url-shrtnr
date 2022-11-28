package v1

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/ru"
	"github.com/pkg/errors"
)

func configureValidator(validate *validator.Validate, universalTranslator *ut.UniversalTranslator) error {
	vt := newValidationsTranslator(validate, universalTranslator)

	if err := vt.registerDefaultTranslations(); err != nil {
		return errors.Wrapf(err, "failed to register default translations")
	}

	if err := vt.overrideDefaultTranslations(); err != nil {
		return errors.Wrapf(err, "failed to override default translations")
	}

	if err := vt.registerCustomValidations(); err != nil {
		return errors.Wrapf(err, "failed to register custom validations")
	}

	if err := vt.registerAliases(); err != nil {
		return errors.Wrapf(err, "failed to register aliases")
	}

	return nil
}

type validationsTranslator struct {
	validator           *validator.Validate
	universalTranslator *ut.UniversalTranslator
}

func newValidationsTranslator(validate *validator.Validate, universalTranslator *ut.UniversalTranslator) *validationsTranslator {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	vt := validationsTranslator{
		validator:           validate,
		universalTranslator: universalTranslator,
	}

	return &vt
}

type Translation struct {
	translation     string
	customRegisFunc validator.RegisterTranslationsFunc
	customTransFunc validator.TranslationFunc
	override        bool
}

type Translations map[string]Translation

type TranslatorNotFoundError struct {
	Locale string
}

func (e *TranslatorNotFoundError) Error() string {
	return fmt.Sprintf("translator for %q locale not found", e.Locale)
}

func (vt *validationsTranslator) registerDefaultTranslations() error {
	for locale, recorder := range map[string]func(validate *validator.Validate, translator ut.Translator) error{
		"en": en.RegisterDefaultTranslations,
		"ru": ru.RegisterDefaultTranslations,
	} {
		log.Printf("debug: registering %q locale default translations", locale)

		translator, found := vt.universalTranslator.GetTranslator(locale)
		if !found {
			return &TranslatorNotFoundError{Locale: locale}
		}

		if err := recorder(vt.validator, translator); err != nil {
			return errors.Wrapf(err, "failed to register %q locale default translations", locale)
		}
	}

	return nil
}

func (vt *validationsTranslator) overrideDefaultTranslations() error {
	for tag, translations := range map[string]Translations{} {
		log.Printf("debug: overriding %q tag translations", tag)

		if err := vt.registerTranslations(tag, translations); err != nil {
			return errors.Wrapf(err, "failed to override %q tag translations", tag)
		}
	}

	return nil
}

func (vt *validationsTranslator) registerCustomValidations() error {
	for tag, st := range map[string]struct {
		validationFn validator.Func
		translations Translations
	}{
		"username": {
			validationFn: usernameValidation,
			translations: Translations{
				"en": {
					translation: "{0} must begin with a letter and contain only letters, underscores, numbers and has length 5 to 32 characters",
					override:    false,
				},
				"ru": {
					translation: "{0} должно начинаться с буквы и содержать только буквы, нижние подчеркивания, цифры и иметь длину от 5 до 32 символов",
					override:    false,
				},
			},
		},
		"password": {
			validationFn: passwordValidation,
			translations: Translations{
				"en": {
					translation: "{0} must contain uppercase and lowercase letters, digits, special characters and has length 8 to 64 characters",
					override:    false,
				},
				"ru": {
					translation: "{0} должен содержать заглавные и строчные буквы, цифры, специальные символы и иметь длину от 8 до 64 символов",
					override:    false,
				},
			},
		},
	} {
		log.Printf("debug: registering custom %q tag validation", tag)

		if err := vt.validator.RegisterValidation(tag, st.validationFn); err != nil {
			return errors.Wrapf(err, "failed to register custom %q tag validation", tag)
		}

		log.Printf("debug: registering custom %q tag translations", tag)

		if err := vt.registerTranslations(tag, st.translations); err != nil {
			return errors.Wrapf(err, "failed to register custom %q tag translations", tag)
		}
	}

	return nil
}

func (vt *validationsTranslator) registerAliases() error {
	for alias, st := range map[string]struct {
		tags         string
		translations Translations
	}{
		"login": {
			tags: "username|email",
			translations: Translations{
				"en": {
					translation: "{0} must be a valid username or email address",
					override:    false,
				},
				"ru": {
					translation: "{0} должен быть действительным именем пользователя или адресом электронной почты",
					override:    false,
				},
			},
		},
	} {
		log.Printf("debug: registering %q alias for tags %q", alias, st.tags)

		vt.validator.RegisterAlias(alias, st.tags)

		log.Printf("debug: registering %q alias translations", alias)

		if err := vt.registerTranslations(alias, st.translations); err != nil {
			return errors.Wrapf(err, "failed to register %q alias translations", alias)
		}
	}

	return nil
}

func (vt *validationsTranslator) registerTranslations(tag string, translations Translations) error {
	for locale, t := range translations {
		translator, found := vt.universalTranslator.GetTranslator(locale)
		if !found {
			return &TranslatorNotFoundError{Locale: locale}
		}

		regisFunc := t.customRegisFunc
		if t.customRegisFunc == nil {
			regisFunc = registrationFunc(tag, t.translation, t.override)
		}

		transFunc := t.customTransFunc
		if t.customTransFunc == nil {
			transFunc = translateFunc
		}

		err := vt.validator.RegisterTranslation(tag, translator, regisFunc, transFunc)
		if err != nil {
			return errors.Wrapf(err, "failed to register %q tag translation in %q locale", tag, locale)
		}
	}

	return nil
}

func registrationFunc(tag, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(translator ut.Translator) error {
		return translator.Add(tag, translation, override)
	}
}

func translateFunc(translator ut.Translator, fieldError validator.FieldError) string {
	translation, err := translator.T(fieldError.Tag(), fieldError.Field())
	if err != nil {
		log.Printf("warning: error translating FieldError: %#v", fieldError)
		return fieldError.(error).Error()
	}

	return translation
}
