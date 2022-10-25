package grpc_prometheus

import (
	"context"
	"strconv"
)

type sourceContextKeyType struct{}

// sourceContextKey unique key for all packages
var sourceContextKey sourceContextKeyType

type Source int

const (
	// Internal for requests that have header referer
	Internal Source = 0

	// External for other requests
	External Source = 1
)

func (s Source) String() string {
	switch s {
	case Internal:
		return "internal"
	case External:
		return "external"
	default:
		return "source(" + strconv.FormatInt(int64(s), 10) + ")"
	}
}

// SetSourceToCtx returns derived context with source
func SetSourceToCtx(ctx context.Context, source Source) context.Context {
	return context.WithValue(ctx, sourceContextKey, source)
}

// GetSourceFromCtx returns source for given context
func GetSourceFromCtx(ctx context.Context) Source {
	v := ctx.Value(sourceContextKey)
	if v == nil {
		return Internal
	}
	return v.(Source)
}
