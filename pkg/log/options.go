package log

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

// Option configures a logger.
type Option interface {
	apply(l *logger) error
}

type optionFunc func(l *logger) error

func (fn optionFunc) apply(l *logger) error {
	return fn(l)
}

// Preset turns a list of Option instances into an Option
func Preset(options ...Option) Option {
	return optionFunc(func(l *logger) error {
		for _, option := range options {
			if err := option.apply(l); err != nil {
				return err
			}
		}

		return nil
	})
}

func SetLevel(level string) Option {
	return optionFunc(func(l *logger) error {
		if level == "" {
			l.level = defaultLevel
			return nil
		}

		lvl, err := zapcore.ParseLevel(level)
		if err != nil {
			return errors.Wrap(err, "failed to parse logger level")
		}

		l.level = lvl
		return nil
	})
}

func SetMode(mode string) Option {
	return optionFunc(func(l *logger) error {
		if mode == "" {
			l.mode = defaultMode
			return nil
		}

		if mode != developmentMode && mode != productionMode {
			return fmt.Errorf("unknown logger mode %q", mode)
		}

		l.mode = mode
		return nil
	})
}

func SetEncoding(encoding string) Option {
	return optionFunc(func(l *logger) error {
		if encoding == "" {
			l.encoding = defaultEncoding
			return nil
		}

		if encoding != consoleEncoding && encoding != jsonEncoding {
			return fmt.Errorf("unknown logger encoding %q", encoding)
		}

		l.encoding = encoding
		return nil
	})
}
