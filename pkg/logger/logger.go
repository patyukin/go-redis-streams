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

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Error(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Fatal(msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Panic(msg, fields...)
}
