package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/ru"
	"github.com/pkg/errors"
)

type registrar struct {
	validator  *validator.Validate
	translator *ut.UniversalTranslator
}

func (r *registrar) configure(logger *zap.Logger) error {
	r.validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	if err := r.registerDefaultTranslations(logger); err != nil {
		return errors.Wrapf(err, "failed to register default translations")
	}

	if err := r.overrideDefaultTranslations(logger); err != nil {
		return errors.Wrapf(err, "failed to override default translations")
	}

	if err := r.registerCustomValidations(logger); err != nil {
		return errors.Wrapf(err, "failed to register custom validations")
	}

	if err := r.registerAliases(logger); err != nil {
		return errors.Wrapf(err, "failed to register aliases")
	}

	return nil
}

type translation struct {
	translation     string
	customRegisFunc validator.RegisterTranslationsFunc
	customTransFunc validator.TranslationFunc
	override        bool
}

type translations map[string]translation

type TranslatorNotFoundError struct {
	Locale string
}

func (e *TranslatorNotFoundError) Error() string {
	return fmt.Sprintf("translator for %q locale not found", e.Locale)
}

func (r *registrar) registerDefaultTranslations(logger *zap.Logger) error {
	for locale, recorder := range map[string]func(validator *validator.Validate, translator ut.Translator) error{
		"en": en.RegisterDefaultTranslations,
		"ru": ru.RegisterDefaultTranslations,
	} {
		logger.Debug("registering default translations",
			zap.String("locale", locale),
		)

		translator, found := unitrans.GetTranslator(locale)
		if !found {
			return &TranslatorNotFoundError{Locale: locale}
		}

		if err := recorder(r.validator, translator); err != nil {
			return errors.Wrapf(err, "failed to register %q locale default translations", locale)
		}
	}

	return nil
}

func (r *registrar) overrideDefaultTranslations(logger *zap.Logger) error {
	for tag, translations := range map[string]translations{} {
		logger.Debug("overriding default translations",
			zap.String("tag", tag),
		)

		if err := r.registerTranslations(logger, tag, translations); err != nil {
			return errors.Wrapf(err, "failed to override %q tag translations", tag)
		}
	}

	return nil
}

func (r *registrar) registerCustomValidations(logger *zap.Logger) error {
	for tag, st := range map[string]struct {
		validationFn validator.Func
		translations translations
	}{
		"username": {
			validationFn: usernameValidation,
			translations: translations{
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
			translations: translations{
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
		logger.Debug("registering custom validation",
			zap.String("tag", tag),
		)

		if err := r.validator.RegisterValidation(tag, st.validationFn); err != nil {
			return errors.Wrapf(err, "failed to register custom %q tag validation", tag)
		}

		logger.Debug("registering translations",
			zap.String("tag", tag),
		)

		if err := r.registerTranslations(logger, tag, st.translations); err != nil {
			return errors.Wrapf(err, "failed to register custom %q tag translations", tag)
		}
	}

	return nil
}

func (r *registrar) registerAliases(logger *zap.Logger) error {
	for alias, st := range map[string]struct {
		tags         string
		translations translations
	}{
		"login": {
			tags: "username|email",
			translations: translations{
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
		logger.Debug("registering alias",
			zap.String("alias", alias),
			zap.String("tags", st.tags),
		)

		r.validator.RegisterAlias(alias, st.tags)

		logger.Debug("registering translations",
			zap.String("tag", alias),
		)

		if err := r.registerTranslations(logger, alias, st.translations); err != nil {
			return errors.Wrapf(err, "failed to register %q alias translations", alias)
		}
	}

	return nil
}

func (r *registrar) registerTranslations(logger *zap.Logger, tag string, translations translations) error {
	for locale, t := range translations {
		translator, found := unitrans.GetTranslator(locale)
		if !found {
			return &TranslatorNotFoundError{Locale: locale}
		}

		regisFunc := t.customRegisFunc
		if t.customRegisFunc == nil {
			regisFunc = registrationFunc(tag, t.translation, t.override)
		}

		transFunc := t.customTransFunc
		if t.customTransFunc == nil {
			transFunc = translateFunc(logger)
		}

		err := r.validator.RegisterTranslation(tag, translator, regisFunc, transFunc)
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

func translateFunc(logger *zap.Logger) func(ut.Translator, validator.FieldError) string {
	return func(translator ut.Translator, fieldError validator.FieldError) string {
		translation, err := translator.T(fieldError.Tag(), fieldError.Field())
		if err != nil {
			logger.Warn("failed to translate FieldError",
				zap.Error(fieldError),
			)

			return fieldError.(error).Error()
		}

		return translation
	}
}
