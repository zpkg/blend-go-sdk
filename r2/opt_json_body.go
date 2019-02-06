package r2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// JSONBody sets the post body on the request.
func JSONBody(obj interface{}) Option {
	return func(r *Request) {
		contents, err := json.Marshal(obj)
		if err != nil {
			r.Err = err
			return
		}
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationJSON)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(contents))
	}
}
