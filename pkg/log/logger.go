package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerContext struct{}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerContext{}, logger)
}

// LoggerFromContext returns logger from context
func LoggerFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerContext{}).(*zap.Logger); ok {
		return logger
	}

	return zap.L()
}

type config struct {
	level    zapcore.Level
	mode     string
	encoding string
}

func NewLogger(options ...Option) (*zap.Logger, error) {
	c := config{
		level:    defaultLevel,
		mode:     defaultMode,
		encoding: defaultEncoding,
	}
	if err := Preset(options...).apply(&c); err != nil {
		return nil, err
	}

	var (
		cfg  zapcore.EncoderConfig
		opts []zap.Option
	)

	if c.mode == developmentMode {
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
	cfg.MessageKey = "message"

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncodeName = zapcore.FullNameEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder

	var enc zapcore.Encoder
	if c.encoding == consoleEncoding {
		enc = zapcore.NewConsoleEncoder(cfg)
	} else {
		enc = zapcore.NewJSONEncoder(cfg)
	}

	core := zapcore.NewCore(enc, os.Stdout, zap.NewAtomicLevelAt(c.level))

	return zap.New(core, opts...), nil
}
