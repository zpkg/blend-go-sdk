package selector

// Error is a hard alias to string.
type Error string

// Error implements `error`
func (e Error) Error() string {
	return string(e)
}

// MarshalJSON implements json.Marshaler.
func (e Error) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e + "\""), nil
}
