package configutil

// ConstantSource returns a new constant source.
func ConstantSource(value string) ValueSource {
	return ConstantSourceValue(value)
}

// ConstantSourceValue returns
type ConstantSourceValue string

// Value returns the value for a constant.
func (csv ConstantSourceValue) Value() (string, error) {
	return string(csv), nil
}
