// Package runner provides the way to use `func(context.Context) error` with github.com/oklog/run.
package runner

import "context"

// Run wraps RunContext(context.Background(), f).
func Run(f func(context.Context) error) (exec func() error, intr func(error)) {
	return RunContext(context.Background(), f)
}

// RunContext returns exec function and intr function, these can be used with run.Group#Add.
// exec function calls f with the parameter context.Context from context.WithCancel(ctx).
// intr function calls cancel function from context.WithCancel(ctx).
func RunContext(ctx context.Context, f func(context.Context) error) (exec func() error, intr func(error)) {
	ctx, cancel := context.WithCancel(ctx)
	exec = func() error {
		return f(ctx)
	}
	intr = func(error) {
		cancel()
	}
	return
}
