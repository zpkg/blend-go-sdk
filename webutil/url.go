package webutil

import "net/url"

// MustParseURL parses a url and panics if there is an error.
func MustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

// URLWithScheme returns a copy url with a given scheme.
func URLWithScheme(u *url.URL, scheme string) *url.URL {
	copy := &(*u)
	copy.Scheme = scheme
	return copy
}

// URLWithHost returns a copy url with a given host.
func URLWithHost(u *url.URL, host string) *url.URL {
	copy := &(*u)
	copy.Host = host
	return copy
}

// URLWithPath returns a copy url with a given path.
func URLWithPath(u *url.URL, path string) *url.URL {
	copy := &(*u)
	copy.Path = path
	return copy
}

// URLWithRawQuery returns a copy url with a given raw query.
func URLWithRawQuery(u *url.URL, rawQuery string) *url.URL {
	copy := &(*u)
	copy.RawQuery = rawQuery
	return copy
}
