package web

import (
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/jwt"
)

const (
	// ErrSessionIDEmpty is thrown if a session id is empty.
	ErrSessionIDEmpty exception.Class = "auth session id is empty"
	// ErrSecureSessionIDEmpty is an error that is thrown if a given secure session id is invalid.
	ErrSecureSessionIDEmpty exception.Class = "auth secure session id is empty"
	// ErrUnsetViewTemplate is an error that is thrown if a given secure session id is invalid.
	ErrUnsetViewTemplate exception.Class = "view result template is unset"
	// ErrParameterMissing is an error on request validation.
	ErrParameterMissing exception.Class = "parameter is missing"
)

// NewParameterMissingError returns a new parameter missing error.
func NewParameterMissingError(paramName string) error {
	return exception.New(ErrParameterMissing).WithMessagef("`%s` parameter is missing", paramName)
}

// IsErrSessionInvalid returns if an error is a session invalid error.
func IsErrSessionInvalid(err error) bool {
	if err == nil {
		return false
	}
	if exception.Is(err, ErrSessionIDEmpty) ||
		exception.Is(err, ErrSecureSessionIDEmpty) ||
		exception.Is(err, jwt.ErrValidation) {
		return true
	}
	return false
}

// IsErrParameterMissing returns if an error is a session invalid error.
func IsErrParameterMissing(err error) bool {
	if err == nil {
		return false
	}
	return exception.Is(err, ErrParameterMissing)
}
