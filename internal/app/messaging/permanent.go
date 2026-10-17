package messaging

import "errors"

// PermanentError marks a failure no retry can fix: the event is parked in
// the dead-letter queue without being retried.
type PermanentError struct {
	Err error
}

func (e PermanentError) Error() string { return e.Err.Error() }

func (e PermanentError) Unwrap() error { return e.Err }

// Permanent wraps a failure that no retry can fix, such as an event with
// invalid fields.
func Permanent(err error) error {
	return PermanentError{Err: err}
}

// IsPermanent reports whether the failure was marked permanent.
func IsPermanent(err error) bool {
	var p PermanentError
	return errors.As(err, &p)
}
