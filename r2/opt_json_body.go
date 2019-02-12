package r2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// OptJSONBody sets the post body on the request.
func OptJSONBody(obj interface{}) Option {
	return func(r *Request) error {
		contents, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationJSON)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(contents))
		return nil
	}
}
