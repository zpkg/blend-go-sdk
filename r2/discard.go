package r2

import (
	"io"
	"io/ioutil"

	"github.com/blend/go-sdk/exception"
)

// Discard discards the response of a request.
func Discard(r *Request, err error) error {
	res, err := Do(r, err)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(ioutil.Discard, res.Body)
	return exception.New(err)
}
