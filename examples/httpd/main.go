package main

import (
	"flag"
	"net/http"
	"syscall"

	"github.com/oklog/run"
	"github.com/takumakei/go-exit"
	"github.com/takumakei/go-run/httpd"
	"github.com/takumakei/go-run/signal"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "`address` to listen")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	sigint := signal.Notify(syscall.SIGINT, syscall.SIGTERM)
	defer sigint.Stop()

	var g run.Group
	g.Add(sigint.Receiver())
	g.Add(httpd.ListenAndServe(*addr, nil))
	err := g.Run()

	exit.ExitOnError(signal.IgnoreError(err))
}
