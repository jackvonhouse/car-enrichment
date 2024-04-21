package errors

import "github.com/jackvonhouse/car-enrichment/pkg/errors"

var (
	ErrInvalid       = errors.NewType("invalid")
	ErrInternal      = errors.NewType("internal")
	ErrNotFound      = errors.NewType("not found")
	ErrAlreadyExists = errors.NewType("already exists")
	ErrFailed        = errors.NewType("failed")
)
