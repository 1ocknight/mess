package ctxkey

import "context"

type ctxKey string

const (
	RequestIDKey ctxKey = "request_id"
)

var publicKeys = map[ctxKey]bool{
    RequestIDKey: true,
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func Parse(ctx context.Context) map[string]any {
	m := make(map[string]any)
	for key, ok := range publicKeys {
		if !ok{
			continue
		}
		if v := ctx.Value(key); v != nil {
			m[string(key)] = v
		}
	}

	return m
}
