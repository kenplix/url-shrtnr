package log

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once

func InitZap(options ...Option) error {
	var err error

	once.Do(func() {
		var l *zap.Logger
		l, err = newLogger(options...)
		if err != nil {
			return
		}

		zap.ReplaceGlobals(l)
	})

	return err
}

type logger struct {
	level    zapcore.Level
	mode     string
	encoding string
}

func newLogger(options ...Option) (*zap.Logger, error) {
	l := logger{
		level:    defaultLevel,
		mode:     defaultMode,
		encoding: defaultEncoding,
	}
	if err := Preset(options...).apply(&l); err != nil {
		return nil, err
	}

	var (
		cfg  zapcore.EncoderConfig
		opts []zap.Option
	)

	if l.mode == developmentMode {
		cfg = zap.NewDevelopmentEncoderConfig()
		cfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder

		opts = append(opts, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))
	} else {
		cfg = zap.NewProductionEncoderConfig()
		cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	cfg.TimeKey = "time"
	cfg.LevelKey = "level"
	cfg.NameKey = "service"
	cfg.CallerKey = "caller"
	cfg.FunctionKey = "function"
	cfg.MessageKey = "message"

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncodeName = zapcore.FullNameEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder

	var enc zapcore.Encoder
	if l.encoding == consoleEncoding {
		enc = zapcore.NewConsoleEncoder(cfg)
	} else {
		enc = zapcore.NewJSONEncoder(cfg)
	}

	core := zapcore.NewCore(enc, os.Stdout, zap.NewAtomicLevelAt(l.level))

	return zap.New(core, opts...), nil
}
