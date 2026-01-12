package errs

import (
	"errors"
	"fmt"
)

type Error struct {
	ErrCode Code
	Msg     string
	Cause   error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause == nil {
		if e.Msg != "" {
			return string(e.ErrCode) + ": " + e.Msg
		}
		return string(e.ErrCode)
	}
	if e.Msg != "" {
		return fmt.Sprintf("%s: %s: %v", e.ErrCode, e.Msg, e.Cause)
	}
	return fmt.Sprintf("%s: %v", e.ErrCode, e.Cause)
}

func (e *Error) Code() Code {
	if e == nil {
		return ""
	}
	return e.ErrCode
}

func (e *Error) Unwrap() error { return e.Cause }

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	if t.ErrCode != "" {
		return e.ErrCode == t.ErrCode
	}
	return false
}

func New(code Code, msg string) *Error {
	return &Error{ErrCode: code, Msg: msg}
}

func Wrap(code Code, msg string, cause error) *Error {
	if cause == nil {
		return New(code, msg)
	}
	return &Error{ErrCode: code, Msg: msg, Cause: cause}
}

func CodeOf(err error) (Code, bool) {
	var de *Error
	if errors.As(err, &de) && de != nil {
		return de.ErrCode, true
	}
	return "", false
}
