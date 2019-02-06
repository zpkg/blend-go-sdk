package r2

import (
	"net/http"
	"net/url"
)

// PostForm sets the request post form and the content type.
func PostForm(postForm url.Values) Option {
	return func(r *Request) {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationFormEncoded)
		r.PostForm = postForm
	}
}

// PostFormValue sets a request post form value.
func PostFormValue(key, value string) Option {
	return func(r *Request) {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationFormEncoded)
		if r.PostForm == nil {
			r.PostForm = url.Values{}
		}
		r.PostForm.Set(key, value)
	}
}
