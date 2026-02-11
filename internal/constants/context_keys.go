package constants

type ctxKey string

const (
	SlogFieldsCtxKey ctxKey = "slog_fields"
	RequestIdCtxKey  ctxKey = "request_id"
)
