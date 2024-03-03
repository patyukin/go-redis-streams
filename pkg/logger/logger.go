package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func InitLogger(path, typeEnv string) (*zap.Logger, error) {
	logWriter := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}

	level := zap.ErrorLevel
	if typeEnv == "dev" {
		level = zap.DebugLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(logWriter), level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level),
	)

	return zap.New(core), nil
}

func GetLogger(ctx context.Context) *zap.Logger {
	logger := ctx.Value("logger").(*zap.Logger)

	return logger
}

func WithLoggerFields(ctx context.Context, fields ...zap.Field) context.Context {
	l := GetLogger(ctx)
	l = l.With(fields...)
	wrap := context.WithValue(ctx, "logger", l)

	return wrap
}
