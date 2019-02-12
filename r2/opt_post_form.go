package r2

import (
	"net/http"
	"net/url"
)

// OptPostForm sets the request post form and the content type.
func OptPostForm(postForm url.Values) Option {
	return func(r *Request) error {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationFormEncoded)
		r.PostForm = postForm
		return nil
	}
}

// OptPostFormValue sets a request post form value.
func OptPostFormValue(key, value string) Option {
	return func(r *Request) error {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationFormEncoded)
		if r.PostForm == nil {
			r.PostForm = url.Values{}
		}
		r.PostForm.Set(key, value)
		return nil
	}
}
