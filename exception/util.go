package exception

// Is is a helper function that returns if an error is an exception.
func Is(err, cause error) bool {
	if typed, isTyped := err.(Exception); isTyped {
		return typed.Class() == cause
	}
	return err == cause
}

// Inner returns an inner error if the error is an exception.
func Inner(err error) error {
	if typed := As(err); typed != nil {
		return typed.Inner()
	}
	return nil
}

// As is a helper method that returns an error as an exception.
func As(err error) Exception {
	if typed, typedOk := err.(Exception); typedOk {
		return typed
	}
	return nil
}
