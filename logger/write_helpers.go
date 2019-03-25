package logger

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/timeutil"
	"github.com/blend/go-sdk/webutil"
)

// WriteHTTPRequest is a helper method to write request start events to a writer.
func WriteHTTPRequest(tf Colorizer, wr io.Writer, req *http.Request) {
	if ip := webutil.GetRemoteAddr(req); len(ip) > 0 {
		io.WriteString(wr, ip)
		io.WriteString(wr, Space)
	}
	io.WriteString(wr, tf.Colorize(req.Method, ansi.ColorBlue))
	io.WriteString(wr, Space)
	io.WriteString(wr, req.URL.Path)
}

// WriteHTTPResponse is a helper method to write request complete events to a writer.
func WriteHTTPResponse(tf Colorizer, wr io.Writer, req *http.Request, statusCode, contentLength int, contentType string, elapsed time.Duration) {
	io.WriteString(wr, webutil.GetRemoteAddr(req))
	io.WriteString(wr, Space)
	io.WriteString(wr, tf.Colorize(req.Method, ansi.ColorBlue))
	io.WriteString(wr, Space)
	io.WriteString(wr, req.URL.Path)
	io.WriteString(wr, Space)
	io.WriteString(wr, ColorizeStatusCode(statusCode))
	io.WriteString(wr, Space)
	io.WriteString(wr, elapsed.String())
	if len(contentType) > 0 {
		io.WriteString(wr, Space)
		io.WriteString(wr, contentType)
	}
	io.WriteString(wr, Space)
	io.WriteString(wr, stringutil.FileSize(contentLength))
}

// WriteFields writes fields.
func WriteFields(tf Colorizer, wr io.Writer, fields Fields) {
	for _, value := range fields {
		if typed, ok := value.(fmt.Stringer); ok {
			io.WriteString(wr, typed.String())
		}
	}
}

// HTTPRequestFields marshals a request start as json.
func HTTPRequestFields(req *http.Request) Fields {
	return Fields{
		"verb":      req.Method,
		"path":      req.URL.Path,
		"host":      req.Host,
		"ip":        webutil.GetRemoteAddr(req),
		"userAgent": webutil.GetUserAgent(req),
	}
}

// HTTPResponseFields marshals a request as json.
func HTTPResponseFields(req *http.Request, statusCode, contentLength int, contentType, contentEncoding string, elapsed time.Duration) Fields {
	return Fields{
		"ip":              webutil.GetRemoteAddr(req),
		"userAgent":       webutil.GetUserAgent(req),
		"verb":            req.Method,
		"path":            req.URL.Path,
		"host":            req.Host,
		"contentLength":   contentLength,
		"contentType":     contentType,
		"contentEncoding": contentEncoding,
		"statusCode":      statusCode,
		FieldElapsed:      timeutil.Milliseconds(elapsed),
	}
}
