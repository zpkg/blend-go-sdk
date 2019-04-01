package secrets

import "context"

// Client is the general interface for a Secrets client
type Client interface {
	Put(ctx context.Context, key string, data Values, options ...RequestOption) error
	Get(ctx context.Context, key string, options ...RequestOption) (Values, error)
	Delete(ctx context.Context, key string, options ...RequestOption) error
	List(ctx context.Context, path string, options ...RequestOption) ([]string, error)
}
