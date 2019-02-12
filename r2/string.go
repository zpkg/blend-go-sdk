package r2

// String reads the response and returns it as a string
func String(r *Request, err error) (string, error) {
	contents, err := Bytes(r, err)
	if err != nil {
		return "", err
	}
	return string(contents), err
}
