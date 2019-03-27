package logger

import "context"

type subContextPathKey struct{}

// WithSubContextPath adds a sub context path to a context.
func WithSubContextPath(ctx context.Context, path []string) context.Context {
	if ctx != nil {
		return context.WithValue(ctx, subContextPathKey{}, path)
	}
	return context.WithValue(context.Background(), subContextPathKey{}, path)
}

// GetSubContextPath adds a sub context path to a context.
func GetSubContextPath(ctx context.Context) []string {
	if rawValue := ctx.Value(subContextPathKey{}); rawValue != nil {
		if typed, ok := rawValue.([]string); ok {
			return typed
		}
	}
	return nil
}
