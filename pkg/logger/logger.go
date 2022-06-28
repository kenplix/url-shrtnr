package logger

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"io"
	"os"
	"time"
)

const (
	DefaultLevel           = logrus.InfoLevel
	DefaultTimestampFormat = time.Stamp
)

type Config struct {
	Level           string `mapstructure:"level"`
	TimestampFormat string `mapstructure:"timestampFormat"`
}

func DefaultConfig() *Config {
	return &Config{
		Level:           DefaultLevel.String(),
		TimestampFormat: DefaultTimestampFormat,
	}
}

// Interface -.
type Interface interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Panic(format string, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *logrus.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(cfg *Config) (*Logger, error) {
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

// Debug -.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Info -.
func (l *Logger) Info(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warn -.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error -.
func (l *Logger) Error(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal -.
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Panic -.
func (l *Logger) Panic(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}
