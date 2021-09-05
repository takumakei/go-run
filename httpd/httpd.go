// Package httpd provides the way to use net/http with github.com/oklog/run.
package httpd

import (
	"context"
	"net/http"
	"time"
)

// ListenAndServe wraps http.ListenAndServe.
func ListenAndServe(addr string, handler http.Handler) (exec func() error, intr func(error)) {
	server := &http.Server{Addr: addr, Handler: handler}
	exec = func() error {
		return server.ListenAndServe()
	}
	intr = func(error) {
		shutdown(server)
	}
	return
}

// ListenAndServeTLS wraps http.ListenAndServeTLS.
func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) (exec func() error, intr func(error)) {
	server := &http.Server{Addr: addr, Handler: handler}
	exec = func() error {
		return server.ListenAndServeTLS(certFile, keyFile)
	}
	intr = func(error) {
		shutdown(server)
	}
	return
}

// ShutdownTimeout is duration for waiting the server shutdown.
var ShutdownTimeout = 7 * time.Second

func shutdown(s *http.Server) {
	if ShutdownTimeout <= 0 {
		_ = s.Shutdown(context.Background())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	_ = s.Shutdown(ctx)
	cancel()
}
