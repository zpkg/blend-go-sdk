package r2

import (
	"encoding/xml"

	"github.com/blend/go-sdk/exception"
)

// XML reads the response as json into a given object.
func XML(r *Request, err error, dst interface{}) error {
	res, err := Do(r, err)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return exception.New(xml.NewDecoder(res.Body).Decode(dst))
}
