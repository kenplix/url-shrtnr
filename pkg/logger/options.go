package logger

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Option configures a Logger.
type Option interface {
	apply(l *Logger) error
}

type optionFunc func(l *Logger) error

func (fn optionFunc) apply(l *Logger) error {
	return fn(l)
}

// Preset turns a list of Option instances into an Option
func Preset(options ...Option) Option {
	return optionFunc(func(l *Logger) error {
		for _, option := range options {
			if err := option.apply(l); err != nil {
				return err
			}
		}

		return nil
	})
}

func SetLevel(level string) Option {
	return optionFunc(func(l *Logger) error {
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			return errors.Wrap(err, "could not parse logger level")
		}

		l.logger.Level = lvl

		return nil
	})
}
