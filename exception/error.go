package exception

// Error is a string wrapper that implements `error`.
// Use this to implement constant exception causes.
type Error string

// Error implements `error`.
func (e Error) Error() string {
	return string(e)
}
