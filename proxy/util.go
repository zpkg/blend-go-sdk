package proxy

import "net/url"

// MustParseURL parses a url and panics if it's bad.
func MustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
