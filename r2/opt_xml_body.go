package r2

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// OptXMLBody sets the post body on the request.
func OptXMLBody(obj interface{}) Option {
	return func(r *Request) error {
		contents, err := xml.Marshal(obj)
		if err != nil {
			return err
		}
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationXML)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(contents))
		return nil
	}
}
