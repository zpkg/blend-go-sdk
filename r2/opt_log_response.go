package r2

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/blend/go-sdk/logger"
)

// OptLogResponse adds an OnResponse listener to log the response of a call.
func OptLogResponse(log logger.Log) Option {
	return func(r *Request) error {
		r.OnResponse = func(req *http.Request, res *http.Response, started time.Time, err error) {
			if err != nil {
				return
			}
			defer res.Body.Close()
			buffer := new(bytes.Buffer)
			io.Copy(buffer, res.Body)
			res.Body = ioutil.NopCloser(bytes.NewBuffer(buffer.Bytes()))
			log.Trigger(NewEvent(logger.HTTPResponse, EventStarted(started), EventRequest(req), EventResponse(res), EventBody(bytes.NewBuffer(buffer.Bytes()))))
		}
		return nil
	}
}
