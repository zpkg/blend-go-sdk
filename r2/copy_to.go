package r2

import "io"

// CopyTo copies the response body to a given writer.
func CopyTo(r *Request, err error, dst io.Writer) (int64, error) {
	res, err := Do(r, err)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return io.Copy(dst, res.Body)
}
