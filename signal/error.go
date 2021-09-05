package signal

import (
	"errors"
	"os"

	"github.com/oklog/run"
)

// Error wraps run.SignalError{Signal: sig}.
func Error(sig os.Signal) error {
	return run.SignalError{Signal: sig}
}

// FromError returns os.Signal of run.SignalError.Signal if err is type of
// run.SignalError, otherwise nil.
func FromError(err error) (os.Signal, bool) {
	var se run.SignalError
	ok := errors.As(err, &se)
	return se.Signal, ok
}

// IgnoreError returns nil if err is type of run.SignalError.
func IgnoreError(err error) error {
	if _, ok := FromError(err); ok {
		return nil
	}
	return err
}
