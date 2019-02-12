package r2

import (
	"encoding/json"

	"github.com/blend/go-sdk/exception"
)

// JSON reads the response as json into a given object.
func JSON(r *Request, err error, dst interface{}) error {
	res, err := Do(r, err)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return exception.New(json.NewDecoder(res.Body).Decode(dst))
}
