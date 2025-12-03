package errs

import "fmt"

// Code represents a business error code.
type Code string

const (
    CodeOK         Code = "OK"
    CodeBadRequest Code = "BAD_REQUEST"
    CodeNotFound   Code = "NOT_FOUND"
    CodeInternal   Code = "INTERNAL"
)

// Error is a structured error used across HTTP/RPC boundaries.
type Error struct {
    Code  Code   `json:"code"`
    Msg   string `json:"message"`
    cause error  `json:"-"`
}

func (e *Error) Error() string {
    if e.cause != nil {
        return fmt.Sprintf("%s: %s (cause=%v)", e.Code, e.Msg, e.cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

// WithCause attaches an underlying error for logging/debugging.
func (e *Error) WithCause(err error) *Error {
    e.cause = err
    return e
}

// Helpers to construct typed errors.

func BadRequest(msg string) *Error {
    return &Error{Code: CodeBadRequest, Msg: msg}
}

func NotFound(msg string) *Error {
    return &Error{Code: CodeNotFound, Msg: msg}
}

func Internal(msg string) *Error {
    return &Error{Code: CodeInternal, Msg: msg}
}
