// Package ticker provides the way to use time.Ticker with github.com/oklog/run.
package ticker

import (
	"context"
	"time"

	"github.com/takumakei/go-run/runner"
)

// Run wraps RunContext(context.Background(), f).
func Run(d time.Duration, f func(context.Context) error) (exec func() error, intr func(error)) {
	return RunContext(context.Background(), d, f)
}

// RunContext returns exec function and intr function, these can be used with run.Group#Add.
// exec function calls f after each tick duration d.
// f is called with the parameter context.Context from context.WithCancel(ctx).
// intr function calls cancel function from context.WithCancel(ctx).
func RunContext(ctx context.Context, d time.Duration, f func(context.Context) error) (exec func() error, intr func(error)) {
	return runner.RunContext(ctx, func(ctx context.Context) error {
		t := time.NewTicker(d)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case <-t.C:
				if err := f(ctx); err != nil {
					return err
				}
			}
		}
	})
}
