package txmanager

import (
	"fmt"
	"runtime"
)

// Error is an enhanced error type
type Error struct {
	code       int
	e          error
	msg        string
	calledFile string
	calledLine int
}

// MakeError makes a basic error message
func MakeError(msg string) *Error {
	_, f, l, _ := runtime.Caller(1)
	return &Error{msg: msg, calledFile: f, calledLine: l}
}

// MakeErrorf returns an error message with formatting
func MakeErrorf(msg string, args ...interface{}) *Error {
	_, f, l, _ := runtime.Caller(1)
	return &Error{
		msg:        fmt.Sprintf(msg, args...),
		calledFile: f,
		calledLine: l,
	}
}

// WrapError wraps an error
func WrapError(e error, msg string) *Error {
	_, f, l, _ := runtime.Caller(1)
	return &Error{msg: msg, calledFile: f, calledLine: l, e: e}
}

// Error returns the error message.
func (e *Error) Error() string {
	if e.e != nil {
		sub := e.e.Error()
		return fmt.Sprintf(
			"%s:%d %s\n%s",
			e.calledFile, e.calledLine, e.msg,
			sub,
		)
	}
	msg := e.msg
	if e.e != nil {
		msg += ": " + e.e.Error()
	}
	return fmt.Sprintf(
		"%s:%d %s", e.calledFile, e.calledLine, msg,
	)
}

// Unwrap returns the error this one wraps. Will return
// nil if this error doesn't wrap another
func (e *Error) Unwrap() error {
	return e.e
}

// Is returns true if the target is an instance of
// the supplied error
func (e *Error) Is(target error) bool {
	_, ok := target.(*Error)
	return ok
}

// Type returns the error code
func (e *Error) Type() int {
	return e.code
}
