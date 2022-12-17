package log

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

// Option configures a logger.
type Option interface {
	apply(c *config) error
}

type optionFunc func(c *config) error

func (fn optionFunc) apply(c *config) error {
	return fn(c)
}

// Preset turns a list of Option instances into an Option
func Preset(options ...Option) Option {
	return optionFunc(func(c *config) error {
		for _, option := range options {
			if err := option.apply(c); err != nil {
				return errors.Wrap(err, "failed to apply option")
			}
		}

		return nil
	})
}

func SetLevel(level string) Option {
	return optionFunc(func(c *config) error {
		if level == "" {
			c.level = defaultLevel
			return nil
		}

		lvl, err := zapcore.ParseLevel(level)
		if err != nil {
			return errors.Wrap(err, "failed to parse logger level")
		}

		c.level = lvl
		return nil
	})
}

func SetMode(mode string) Option {
	return optionFunc(func(c *config) error {
		if mode == "" {
			c.mode = defaultMode
			return nil
		}

		if mode != developmentMode && mode != productionMode {
			return fmt.Errorf("unknown logger mode %q", mode)
		}

		c.mode = mode
		return nil
	})
}

func SetEncoding(encoding string) Option {
	return optionFunc(func(c *config) error {
		if encoding == "" {
			c.encoding = defaultEncoding
			return nil
		}

		if encoding != consoleEncoding && encoding != jsonEncoding {
			return fmt.Errorf("unknown logger encoding %q", encoding)
		}

		c.encoding = encoding
		return nil
	})
}
