package configutil

import (
	"flag"
	"fmt"
)

// FlagSource returns a new flag source.
func FlagSource(name, usage string) ValueSource {
	return &FlagSourceValue{
		name:  name,
		usage: usage,
	}
}

// FlagSourceValue is a value source from a commandline flag.
type FlagSourceValue struct {
	set   *flag.FlagSet
	name  string
	usage string
	value string
}

// Register registers the value provider.
func (fsv *FlagSourceValue) Register() error {
	if fsv.set == nil {
		return fmt.Errorf("FlagSourceValue: flag set not provided")
	}
	fsv.set.StringVar(&fsv.value, fsv.name, "", fsv.usage)
	return nil
}

// Value returns the flag value.
// You *must* call `flag.Parse` before calling this function.
func (fsv *FlagSourceValue) Value() (string, error) {
	return fsv.value, nil
}
