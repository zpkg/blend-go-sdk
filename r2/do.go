package r2

import (
	"net/http"
	"time"
)

// Do executes the request.
func Do(r *Request, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}
	started := time.Now().UTC()

	if r.OnRequest != nil {
		r.OnRequest(r.Request)
	}

	var res *http.Response
	if r.Client != nil {
		res, err = r.Client.Do(r.Request)
	} else {
		res, err = http.DefaultClient.Do(r.Request)
	}

	if r.OnResponse != nil {
		r.OnResponse(r.Request, res, started, err)
	}
	return res, err
}
