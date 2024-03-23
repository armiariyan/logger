package logger

import (
	"context"
)

var instance Logger

func SetGlobalLogger(in Logger) {
	instance = in
}

func getInstance() Logger {
	if instance == nil {
		return NewNoopLogger()
	}

	return instance
}

func Debug(ctx context.Context, message string, fields ...Field) {
	getInstance().Debug(ctx, message, fields...)
}

func Info(ctx context.Context, message string, fields ...Field) {
	getInstance().Info(ctx, message, fields...)
}

func Warn(ctx context.Context, message string, fields ...Field) {
	getInstance().Warn(ctx, message, fields...)
}

func Error(ctx context.Context, message string, fields ...Field) {
	getInstance().Error(ctx, message, fields...)
}

func Fatal(ctx context.Context, message string, fields ...Field) {
	getInstance().Fatal(ctx, message, fields...)
}

func Panic(ctx context.Context, message string, fields ...Field) {
	getInstance().Panic(ctx, message, fields...)
}

func TDR(ctx context.Context, tdr LogTdrModel) {
	getInstance().TDR(ctx, tdr)
}
