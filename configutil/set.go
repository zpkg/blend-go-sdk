package configutil

// Set coalesces a given list of sources into a variable.
func Set(destination *string, sources ...ValueSource) error {
	var value string
	var err error
	for _, source := range sources {
		value, err = source.Value()
		if err != nil {
			return err
		}
		if value != "" {
			*destination = value
			return nil
		}
	}
	return nil
}
