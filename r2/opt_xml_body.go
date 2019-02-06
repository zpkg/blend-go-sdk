package r2

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// XMLBody sets the post body on the request.
func XMLBody(obj interface{}) Option {
	return func(r *Request) {
		contents, err := xml.Marshal(obj)
		if err != nil {
			r.Err = err
			return
		}
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationXML)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(contents))
	}
}
