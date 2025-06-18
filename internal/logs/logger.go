package logs

import (
	"log"

	"go.uber.org/zap"
)

type Log struct {
	Logger *zap.Logger
}

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

func NewLogger() Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Ошибка при билде zap logger: %v", err)
	}

	return Log{
		Logger: logger.WithOptions(zap.AddCallerSkip(1)),
	}
}

func (l Log) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l Log) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l Log) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l Log) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}
