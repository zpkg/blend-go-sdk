package configutil

import "time"

// CoalesceSources coalesces a given list of sources.
func CoalesceSources(sources ...ValueSource) (string, error) {
	var value string
	var err error
	for _, source := range sources {
		value, err = source.Value()
		if err != nil {
			return "", err
		}
		if value != "" {
			return value, nil
		}
	}
	return "", nil
}

// CoalesceSourcesVar coalesces a given list of sources.
func CoalesceSourcesVar(destination *string, sources ...ValueSource) error {
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

// CoalesceString returns a coalesced value.
func CoalesceString(value, defaultValue string, inheritedValues ...string) string {
	if len(value) > 0 {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceBool returns a coalesced value.
func CoalesceBool(value *bool, defaultValue bool, inheritedValues ...bool) bool {
	if value != nil {
		return *value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceInt returns a coalesced value.
func CoalesceInt(value, defaultValue int, inheritedValues ...int) int {
	if value > 0 {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceInt32 returns a coalesced value.
func CoalesceInt32(value, defaultValue int32, inheritedValues ...int32) int32 {
	if value > 0 {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceInt64 returns a coalesced value.
func CoalesceInt64(value, defaultValue int64, inheritedValues ...int64) int64 {
	if value > 0 {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceFloat32 returns a coalesced value.
func CoalesceFloat32(value, defaultValue float32, inheritedValues ...float32) float32 {
	if value > 0 {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceFloat64 returns a coalesced value.
func CoalesceFloat64(value, defaultValue float64, inheritedValues ...float64) float64 {
	if value > 0 {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceDuration returns a coalesced value.
func CoalesceDuration(value, defaultValue time.Duration, inheritedValues ...time.Duration) time.Duration {
	if value > 0 {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceTime returns a coalesced value.
func CoalesceTime(value, defaultValue time.Time, inheritedValues ...time.Time) time.Time {
	if !value.IsZero() {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceStrings returns a coalesced value.
func CoalesceStrings(value, defaultValue []string, inheritedValues ...[]string) []string {
	if len(value) > 0 {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}

// CoalesceBytes returns a coalesced value.
func CoalesceBytes(value, defaultValue []byte, inheritedValues ...[]byte) []byte {
	if len(value) > 0 {
		return value
	}
	if len(inheritedValues) > 0 {
		return inheritedValues[0]
	}
	return defaultValue
}
