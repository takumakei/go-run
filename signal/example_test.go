package signal_test

import (
	"fmt"
	"io"
	"syscall"

	"github.com/oklog/run"
	"github.com/takumakei/go-run/signal"
)

func Example() {
	sigint := signal.Notify(syscall.SIGINT)
	defer sigint.Stop()
	var g run.Group
	g.Add(sigint.Receiver())
	g.Add(func() error { return io.EOF }, func(error) {})
	if err := signal.IgnoreError(g.Run()); err != nil {
		fmt.Println(err)
	}
	// Output:
	// EOF
}
