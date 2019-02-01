package r2

import "net/http"

// Close closes a response.
// It's meant to be nested in a `Do` call:
//    res, err := r2.Close(r2.New("https://foo.com").Do())
func Close(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		if err := res.Body.Close(); err != nil {
			return nil, err
		}
	}
	return res, nil
}
