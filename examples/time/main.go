package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-run/signal"
	"github.com/takumakei/go-run/ticker"
	"github.com/takumakei/go-run/timer"
)

const (
	C0 = "\x1b[97;44m"
	C1 = "\x1b[97;45m"
	CE = "\x1b[0m"
)

func main() {
	var start time.Time

	sigint := signal.Notify(syscall.SIGINT, syscall.SIGTERM)
	defer sigint.Stop()

	var g run.Group
	g.Add(sigint.Receiver())

	g.Add(timer.Run(3*time.Second, func(ctx context.Context) error {
		fmt.Println(C0 + "[" + since(start) + "] timer. sleep 1 sec." + CE)
		time.Sleep(time.Second)
		fmt.Println(C0 + "[" + since(start) + "] timer end. next is 3 sec later." + CE)
		return nil
	}))

	g.Add(ticker.Run(3*time.Second, func(ctx context.Context) error {
		fmt.Println(C1 + "[" + since(start) + "] ticker. next is 3 sec later. sleep 1 sec." + CE)
		time.Sleep(time.Second)
		fmt.Println(C1 + "[" + since(start) + "] ticker end." + CE)
		return nil
	}))

	fmt.Println("start")
	start = time.Now()
	err := g.Run()
	fmt.Println(err)
	exit.ExitOnError(signal.IgnoreError(err))
}

func since(start time.Time) string {
	return time.Since(start).Truncate(time.Second).String()
}
