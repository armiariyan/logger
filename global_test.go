package logger

import (
	"context"
	"testing"
)

func TestGlobalLogger_Nil(t *testing.T) {
	ctx := context.Background()
	message := "log message"
	fields := []Field{ToField("foo", "bar")}

	SetGlobalLogger(nil)
	Debug(ctx, message, fields...)
	Info(ctx, message, fields...)
	Warn(ctx, message, fields...)
	Error(ctx, message, fields...)
	Fatal(ctx, message, fields...)
	Panic(ctx, message, fields...)
	TDR(ctx, GenerateLogTDR(nil))
}

func TestGlobalLogger_Noop(t *testing.T) {
	ctx := context.Background()
	message := "log message"
	fields := []Field{ToField("foo", "bar")}

	SetGlobalLogger(NewNoopLogger())
	Debug(ctx, message, fields...)
	Info(ctx, message, fields...)
	Warn(ctx, message, fields...)
	Error(ctx, message, fields...)
	Fatal(ctx, message, fields...)
	Panic(ctx, message, fields...)
	TDR(ctx, GenerateLogTDR(nil))
}
