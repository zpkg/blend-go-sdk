package configutil

var (
	_ ValueSource = (*ConstantValue)(nil)
)

// Const returns a new constant source.
func Const(value string) ValueSource {
	return ConstantValue(value)
}

// ConstantValue returns
type ConstantValue string

// Value returns the value for a constant.
func (cv ConstantValue) Value() (string, error) {
	return string(cv), nil
}
