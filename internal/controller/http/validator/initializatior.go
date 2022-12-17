package validator

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	once     sync.Once
	unitrans *ut.UniversalTranslator
)

func Init(logger *zap.Logger) (*ut.UniversalTranslator, error) {
	var err error

	once.Do(func() {
		unitrans, err = initialize(logger)
	})

	if err != nil {
		return nil, err
	}

	return unitrans, nil
}

func initialize(logger *zap.Logger) (*ut.UniversalTranslator, error) {
	english := en.New()
	russian := ru.New()

	unitrans = ut.New(english, english, russian)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		logger.Info("configuring gin validator instance")

		r := registrar{
			validator:  v,
			translator: unitrans,
		}
		if err := r.configure(logger); err != nil {
			return nil, err
		}

		return unitrans, nil
	}

	return nil, fmt.Errorf("unknown gin validator engine %T", binding.Validator.Engine())
}
