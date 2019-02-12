package r2

import (
	"io/ioutil"

	"github.com/blend/go-sdk/exception"
)

// Bytes reads the response and returns it as a byte array.
func Bytes(r *Request, err error) ([]byte, error) {
	res, err := Do(r, err)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, exception.New(err)
	}
	return contents, nil
}
