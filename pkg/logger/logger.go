package logger

import (
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

const (
	defaultLevel           = logrus.InfoLevel
	defaultTimestampFormat = time.Stamp
)

type Config struct {
	Level           string `mapstructure:"level"`
	TimestampFormat string `mapstructure:"timestampFormat"`
}

func DefaultConfig() Config {
	return Config{
		Level:           defaultLevel.String(),
		TimestampFormat: defaultTimestampFormat,
	}
}

// Interface -.
type Interface interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *logrus.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(cfg Config) (*Logger, error) {
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse logger level")
	}

	logger := &logrus.Logger{
		Level: level,
		Out:   io.Discard,
		Formatter: &logrus.TextFormatter{
			ForceColors:     true,
			DisableQuote:    true,
			FullTimestamp:   true,
			TimestampFormat: cfg.TimestampFormat,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg: "message",
			},
		},
		Hooks:    make(logrus.LevelHooks),
		ExitFunc: os.Exit,
	}

	for _, hook := range []logrus.Hook{
		// Send logs with level higher than warning to stderr
		&writer.Hook{
			Writer: os.Stderr,
			LogLevels: []logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
			},
		},
		// Send info and debug logs to stdout
		&writer.Hook{
			Writer: os.Stdout,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
				logrus.DebugLevel,
			},
		},
	} {
		logger.AddHook(hook)
	}

	return &Logger{logger: logger}, nil
}

// Debugf -.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Infof -.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warnf -.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Errorf -.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatalf -.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Panicf -.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}
