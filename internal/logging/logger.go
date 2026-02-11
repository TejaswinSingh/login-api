package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/TejaswinSingh/login-api/internal/config"
	"github.com/TejaswinSingh/login-api/internal/constants"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(constants.SlogFieldsCtxKey).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

func NewLogger(config config.Config) *slog.Logger {
	var handler slog.Handler
	if config.Env.IsProduction() {
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	handler = &ContextHandler{Handler: handler}
	return slog.New(handler)
}

// attributes appended to ctx via this function will be extracted and logged
// when functions such as LogAttrs, InfoContext, ErrorContext etc are called
// with the same ctx
func AppendAttrToCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(constants.SlogFieldsCtxKey).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, constants.SlogFieldsCtxKey, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(parent, constants.SlogFieldsCtxKey, v)
}
