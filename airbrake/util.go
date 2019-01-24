package airbrake

import "strconv"

func mustInt(value string) int64 {
	output, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return output
}
