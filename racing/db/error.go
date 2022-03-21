package db

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

var _ error = (*invalidArgumentError)(nil)

// invalidArgumentError indicates client specified an invalid argument.
type invalidArgumentError struct {
	field   string
	details string
}

// Code hints at the gRPC code that would be representative for this error.
func (e *invalidArgumentError) Code() codes.Code {
	return codes.InvalidArgument
}

// Field returns the field name that is invalid.
func (e *invalidArgumentError) Field() string {
	return e.field
}

// Field returns the details around the invalid argument.
func (e *invalidArgumentError) Details() string {
	return e.details
}

// Error satisfies the error interface.
func (e *invalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument \"%s\", %s", e.field, e.details)
}
