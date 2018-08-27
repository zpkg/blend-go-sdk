package db

import "github.com/blend/go-sdk/exception"

const (
	// ErrConfigUnset is an exception class.
	ErrConfigUnset exception.Class = "db: config is unset"
	// ErrUnsafeSSLMode is an error indicating unsafe ssl mode in production.
	ErrUnsafeSSLMode exception.Class = "db: unsafe ssl mode in prodlike environment"
	// ErrUsernameUnset is an error indicating there is no username set in a prodlike environment.
	ErrUsernameUnset exception.Class = "db: username is unset in prodlike environment"
	// ErrPasswordUnset is an error indicating there is no password set in a prodlike environment.
	ErrPasswordUnset exception.Class = "db: password is unset in prodlike environment"
	// ErrConnectionAlreadyOpen is an error indicating the db connection was already opened.
	ErrConnectionAlreadyOpen exception.Class = "db: the connection is already opened"
	// ErrConnectionClosed is an error indicating the db connection hasn't been opened.
	ErrConnectionClosed exception.Class = "db: the connection is closed"
	// ErrStatementCacheUnset is an error indicating the statement cache is unset.
	ErrStatementCacheUnset exception.Class = "db: the statement cache is unset"
	// ErrCollectionNotSlice is an error returned by OutMany if the destination is not a slice.
	ErrCollectionNotSlice exception.Class = "db: outmany destination collection is not a slice"
)

// IsConfigUnset returns if the error is an `ErrConfigUnset`.
func IsConfigUnset(err error) bool {
	return exception.Is(err, ErrConfigUnset)
}

// IsUnsafeSSLMode returns if an error is an `ErrUnsafeSSLMode`.
func IsUnsafeSSLMode(err error) bool {
	return exception.Is(err, ErrUnsafeSSLMode)
}

// IsUsernameUnset returns if an error is an `ErrUsernameUnset`.
func IsUsernameUnset(err error) bool {
	return exception.Is(err, ErrUsernameUnset)
}

// IsPasswordUnset returns if an error is an `ErrPasswordUnset`.
func IsPasswordUnset(err error) bool {
	return exception.Is(err, ErrPasswordUnset)
}

// IsConnectionClosed returns if the error is an `ErrConnectionClosed`.
func IsConnectionClosed(err error) bool {
	return exception.Is(err, ErrConnectionClosed)
}

// IsStatementCacheUnset returns if the error is an `ErrConnectionClosed`.
func IsStatementCacheUnset(err error) bool {
	return exception.Is(err, ErrStatementCacheUnset)
}
