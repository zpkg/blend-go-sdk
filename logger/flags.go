package logger

import (
	"strings"
)

// ParseFlags returns a new flag set from an array of flag values.
func ParseFlags(flags ...string) *Flags {
	flagSet := &Flags{
		flags: make(map[string]bool),
	}

	for _, flag := range flags {
		parsedFlag := strings.Trim(strings.ToLower(flag), " \t\n")
		if parsedFlag == All {
			flagSet.all = true
		}

		if parsedFlag == None {
			flagSet.none = true
			return flagSet
		}

		if strings.HasPrefix(string(parsedFlag), "-") {
			flag := strings.TrimPrefix(string(parsedFlag), "-")
			flagSet.flags[flag] = false
		} else {
			flagSet.flags[parsedFlag] = true
		}
	}

	return flagSet
}

// NewFlags returns a new FlagSet with the given flags enabled.
func NewFlags(enabled ...string) *Flags {
	efs := &Flags{
		flags: make(map[string]bool),
	}
	for _, flag := range enabled {
		efs.flags[flag] = true
	}
	return efs
}

// Flags is a set of event flags.
type Flags struct {
	flags map[string]bool
	all   bool
	none  bool
}

// Enable enables an event flag.
func (efs *Flags) Enable(flag string) {
	efs.none = false
	efs.flags[flag] = true
}

// Disable disables a flag.
func (efs *Flags) Disable(flag string) {
	efs.flags[flag] = false
}

// SetAll flips the `all` bit on the flag set to true.
func (efs *Flags) SetAll() {
	efs.all = true
	efs.flags = make(map[string]bool)
	efs.none = false
}

// All returns if the all bit is flipped to true.
func (efs *Flags) All() bool {
	return efs.all
}

// SetNone flips the `none` bit on the flag set to true.
// It also disables the `all` bit.
func (efs *Flags) SetNone() {
	efs.all = false
	efs.flags = make(map[string]bool)
	efs.none = true
}

// None returns if the none bit is flipped to true.
func (efs *Flags) None() bool {
	return efs.none
}

// IsEnabled checks to see if an event is enabled.
func (efs Flags) IsEnabled(flag string) bool {
	if efs.all {
		// figure out if we explicitly disabled the flag.
		if enabled, hasEvent := efs.flags[flag]; hasEvent && !enabled {
			return false
		}
		return true
	}
	if efs.none {
		return false
	}
	if efs.flags != nil {
		if enabled, hasFlag := efs.flags[flag]; hasFlag {
			return enabled
		}
	}
	return false
}

func (efs Flags) String() string {
	if efs.none {
		return None
	}

	var flags []string
	if efs.all {
		flags = []string{All}
	}
	for key, enabled := range efs.flags {
		if key != All {
			if enabled {
				if !efs.all {
					flags = append(flags, string(key))
				}
			} else {
				flags = append(flags, "-"+string(key))
			}
		}
	}
	return strings.Join(flags, ", ")
}

// MergeWith sets the set from another, with the other taking precedence.
func (efs Flags) MergeWith(other *Flags) {
	if other.all {
		efs.all = true
	}
	if other.none {
		efs.none = true
	}
	for key, value := range other.flags {
		efs.flags[key] = value
	}
}