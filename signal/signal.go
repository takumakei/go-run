// Package signal provides the way to handle signal with github.com/oklog/run.
package signal

import (
	"context"
	"os"
	"os/signal"

	"github.com/takumakei/go-run/runner"
)

// Signal has a channel type of `chan os.Signal` and
// a handler function type of `func(os.Signal) error`.
type Signal struct {
	C chan os.Signal
	f func(os.Signal) error
}

// Notify returns *Signal.
//   - Signal#C = make(chan os.Signal, 4)
//   - it calls signal.Notify(C, sig...)
//   - Signal#DefaultHandler is set to the signal handler of Signal returned.
func Notify(sig ...os.Signal) *Signal {
	c := make(chan os.Signal, 4)
	signal.Notify(c, sig...)
	si := &Signal{C: c}
	si.f = si.DefaultHandler
	return si
}

// Stop wraps signal.Stop(si.C).
func (si *Signal) Stop() {
	signal.Stop(si.C)
}

// Close wraps close(si.C).
func (si *Signal) Close() {
	close(si.C)
}

// DefaultHandler wraps si.Stop(), returns Error(sig).
func (si *Signal) DefaultHandler(sig os.Signal) error {
	si.Stop()
	return Error(sig)
}

// SetHandler set f to si's signal handler if f is not nil,
// otherwise set si.DefaultHandler to it.
func (si *Signal) SetHandler(f func(os.Signal) error) {
	if f == nil {
		f = si.DefaultHandler
	}
	si.f = f
}

// Receiver returns exec function and intr function, these can be used with run.Group#Add.
// exec waits a signal, then calls signal handler, returns its returned error if it is not nil.
// If the signal handler returns nil, exec does it again.
// exec returns immediately if intr is called.
//
// EXAMPLE
//
//   sigint := signal.Notify(syscall.SIGINT)
//   defer sigint.Stop()
//   var g run.Group
//   g.Add(sigint.Receiver())
func (si *Signal) Receiver() (exec func() error, intr func(error)) {
	return runner.Run(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case sig, ok := <-si.C:
				if !ok {
					return nil
				}
				if err := si.f(sig); err != nil {
					return err
				}
			}
		}
	})
}
