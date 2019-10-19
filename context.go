package steps

import (
	"context"
)

type key int

const (
	optionKey key = iota
)

func getOptions(ctx context.Context) []RunOption {
	return ctx.Value(optionKey).([]RunOption)
}

func appendOptions(ctx context.Context, opts []RunOption) (context.Context, []RunOption) {
	if len(opts) == 0 {
		return ctx, nil
	}

	// copy slice to prevent race condition when opts exists
	v := getOptions(ctx)
	l := len(v)
	opts = append(v[0:l:l], opts...)
	return context.WithValue(ctx, optionKey, opts), opts
}
