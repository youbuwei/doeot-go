package domain

import "context"

// NotFoundError is a tiny domain error type to distinguish not-found cases.
type NotFoundError struct {
	msg string
}

func (e *NotFoundError) Error() string { return e.msg }

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{msg: msg}
}

// ErrUserNotFound is returned when user is missing in persistence layer.
var ErrUserNotFound = NewNotFoundError("user not found")

// Repo abstracts persistence operations for User.
type Repo interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, u *User) (*User, error)
}
