package configutil

var (
	_ ValueSource = (*ValueFunc)(nil)
)

// ValueFunc is a value source from a commandline flag.
type ValueFunc func() (string, error)

// Value returns the flag value.
// You *must* call `flag.Parse` before calling this function.
func (vf ValueFunc) Value() (string, error) {
	return vf()
}
