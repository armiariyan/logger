package logger

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, message string, fields ...Field)
	Info(ctx context.Context, message string, fields ...Field)
	Warn(ctx context.Context, message string, fields ...Field)
	Error(ctx context.Context, message string, fields ...Field)
	Fatal(ctx context.Context, message string, fields ...Field)
	Panic(ctx context.Context, message string, fields ...Field)
	TDR(ctx context.Context, tdr LogTdrModel)
	Close() error
}

type Field struct {
	Key string
	Val interface{}
}

type ctxKeyLogger struct{}

var ctxKey = ctxKeyLogger{}

type Context struct {
	ServiceName    string `json:"_app_name"`
	ServiceVersion string `json:"_app_version"`
	ServicePort    int    `json:"_app_port"`
	ThreadID       string `json:"_app_thread_id"`
	JourneyID      string `json:"_app_journey_id"`
	ChainID        string `json:"_app_chain_id"`
	Tag            string `json:"_app_tag"`

	ReqMethod string `json:"_app_method"`
	ReqURI    string `json:"_app_uri"`

	AdditionalData map[string]interface{} `json:"_app_data,omitempty"`
}

func InjectCtx(parent context.Context, ctx Context) context.Context {
	if parent == nil {
		return InjectCtx(context.Background(), ctx)
	}

	return context.WithValue(parent, ctxKey, ctx)
}

func ExtractCtx(ctx context.Context) Context {
	if ctx == nil {
		return Context{}
	}

	val, ok := ctx.Value(ctxKey).(Context)
	if !ok {
		return Context{}
	}

	return val
}
