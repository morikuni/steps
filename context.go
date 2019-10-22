package steps

import (
	"context"
)

type key int

const (
	optionKey key = iota
)

func getOptions(ctx context.Context) []Option {
	return ctx.Value(optionKey).([]Option)
}

func appendOptions(ctx context.Context, opts []Option) (context.Context, []Option) {
	if len(opts) == 0 {
		return ctx, nil
	}

	// copy slice to prevent race condition when opts exists
	v := getOptions(ctx)
	l := len(v)
	opts = append(v[0:l:l], opts...)
	return context.WithValue(ctx, optionKey, opts), opts
}
