package r2

import (
	"net/http"

	"github.com/blend/go-sdk/exception"
)

// Close executes and closes the response.
// It returns the response for metadata purposes.
func Close(r *Request, err error) (*http.Response, error) {
	res, err := Do(r, err)
	if err != nil {
		return nil, err
	}
	return res, exception.New(res.Body.Close())
}
