run
======================================================================

[![Go Reference](https://pkg.go.dev/badge/github.com/takumakei/go-run.svg)](https://pkg.go.dev/github.com/takumakei/go-run)

Packages provide utilities for github.com/oklog/run.

httpd
----------------------------------------------------------------------

### example

```
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	sigint := signal.Notify(syscall.SIGINT, syscall.SIGTERM)
	defer sigint.Stop()

	var g run.Group
	g.Add(sigint.Receiver())
	g.Add(httpd.ListenAndServe(":80", nil))
	g.Run()
```

### go doc

```
package httpd // import "github.com/takumakei/go-run/httpd"

Package httpd provides the way to use net/http with github.com/oklog/run.

VARIABLES

var ShutdownTimeout = 7 * time.Second
    ShutdownTimeout is duration for waiting the server shutdown.


FUNCTIONS

func ListenAndServe(addr string, handler http.Handler) (exec func() error, intr func(error))
    ListenAndServe wraps http.ListenAndServe.

func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) (exec func() error, intr func(error))
    ListenAndServeTLS wraps http.ListenAndServeTLS.
```

runner
----------------------------------------------------------------------

```
package runner // import "github.com/takumakei/go-run/runner"

Package runner provides the way to use `func(context.Context) error` with
github.com/oklog/run.

FUNCTIONS

func Run(f func(context.Context) error) (exec func() error, intr func(error))
    Run wraps RunContext(context.Background(), f).

func RunContext(ctx context.Context, f func(context.Context) error) (exec func() error, intr func(error))
    RunContext returns exec function and intr function, these can be used with
    run.Group#Add. exec function calls f with the parameter context.Context from
    context.WithCancel(ctx). intr function calls cancel function from
    context.WithCancel(ctx).
```

signal
----------------------------------------------------------------------

```
package signal // import "github.com/takumakei/go-run/signal"

Package signal provides the way to handle signal with github.com/oklog/run.

FUNCTIONS

func Error(sig os.Signal) error
    Error wraps run.SignalError{Signal: sig}.

func FromError(err error) (os.Signal, bool)
    FromError returns os.Signal of run.SignalError.Signal if err is type of
    run.SignalError, otherwise nil.

func IgnoreError(err error) error
    IgnoreError returns nil if err is type of run.SignalError.


TYPES

type Signal struct {
	C chan os.Signal
	// Has unexported fields.
}
    Signal has a channel type of `chan os.Signal` and a handler function type of
    `func(os.Signal) error`.

func Notify(sig ...os.Signal) *Signal
    Notify returns *Signal.

        - Signal#C = make(chan os.Signal, 4)
        - it calls signal.Notify(C, sig...)
        - Signal#DefaultHandler is set to the signal handler of Signal returned.

func (si *Signal) Close()
    Close wraps close(si.C).

func (si *Signal) DefaultHandler(sig os.Signal) error
    DefaultHandler wraps si.Stop(), returns Error(sig).

func (si *Signal) Receiver() (exec func() error, intr func(error))
    Receiver returns exec function and intr function, these can be used with
    run.Group#Add. exec waits a signal, then calls signal handler, returns its
    returned error if it is not nil. If the signal handler returns nil, exec
    does it again. exec returns immediately if intr is called.

    EXAMPLE

        sigint := signal.Notify(syscall.SIGINT)
        defer sigint.Stop()
        var g run.Group
        g.Add(sigint.Receiver())

func (si *Signal) SetHandler(f func(os.Signal) error)
    SetHandler set f to si's signal handler if f is not nil, otherwise set
    si.DefaultHandler to it.

func (si *Signal) Stop()
    Stop wraps signal.Stop(si.C).
```

ticker
----------------------------------------------------------------------

```
package ticker // import "github.com/takumakei/go-run/ticker"

Package ticker provides the way to use time.Ticker with
github.com/oklog/run.

FUNCTIONS

func Run(d time.Duration, f func(context.Context) error) (exec func() error, intr func(error))
    Run wraps RunContext(context.Background(), f).

func RunContext(ctx context.Context, d time.Duration, f func(context.Context) error) (exec func() error, intr func(error))
    RunContext returns exec function and intr function, these can be used with
    run.Group#Add. exec function calls f after each tick duration d. f is called
    with the parameter context.Context from context.WithCancel(ctx). intr
    function calls cancel function from context.WithCancel(ctx).
```

timer
----------------------------------------------------------------------

```
package timer // import "github.com/takumakei/go-run/timer"

Package timer provides the way to use time.Timer with github.com/oklog/run.

FUNCTIONS

func Run(d time.Duration, f func(context.Context) error) (exec func() error, intr func(error))
    Run wraps RunContext(context.Background(), f).

func RunContext(ctx context.Context, d time.Duration, f func(context.Context) error) (exec func() error, intr func(error))
    RunContext returns exec function and intr function, these can be used with
    run.Group#Add. exec function calls f repeatedly with interval d. f is called
    with the parameter context.Context from context.WithCancel(ctx). intr
    function calls cancel function from context.WithCancel(ctx).
```
