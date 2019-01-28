package jobkit

import "time"

// OptBool returns a reference to a value.
func OptBool(value bool) *bool {
	return &value
}

// OptDuration returns a reference to a value.
func OptDuration(value time.Duration) *time.Duration {
	return &value
}

// OptString returns a reference to a value.
func OptString(value string) *string {
	return &value
}

// OptStrings returns a reference to a value.
func OptStrings(values ...string) []*string {
	output := make([]*string, len(values))
	for index := range values {
		output[index] = &values[index]
	}
	return output
}
