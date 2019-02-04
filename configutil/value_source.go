package configutil

// ValueSource is a type that can return a value.
type ValueSource interface {
	// Value should return a string if the source has a given value.
	// It should return empty string if the value is not present.
	// It should return an error if there was a problem fetching the value.
	Value() (string, error)
}
