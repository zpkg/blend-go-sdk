package logger

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/webutil"
)

// FormatFileSize returns a string representation of a file size in bytes.
func FormatFileSize(sizeBytes int) string {
	if sizeBytes >= 1<<30 {
		return strconv.Itoa(sizeBytes/Gigabyte) + "gb"
	} else if sizeBytes >= 1<<20 {
		return strconv.Itoa(sizeBytes/Megabyte) + "mb"
	} else if sizeBytes >= 1<<10 {
		return strconv.Itoa(sizeBytes/Kilobyte) + "kb"
	}
	return strconv.Itoa(sizeBytes)
}

// WriteHTTPRequest is a helper method to write request start events to a writer.
func WriteHTTPRequest(wr io.Writer, req *http.Request) {
	if ip := webutil.GetRemoteAddr(req); len(ip) > 0 {
		buf.WriteString(ip)
		buf.WriteRune(RuneSpace)
	}
	buf.WriteString(tf.Colorize(req.Method, ansi.ColorBlue))
	buf.WriteRune(RuneSpace)
	buf.WriteString(req.URL.Path)
}

// WriteHTTPResponse is a helper method to write request complete events to a writer.
func WriteHTTPResponse(wr io.Writer, req *http.Request, statusCode, contentLength int, contentType string, elapsed time.Duration) {
	buf.WriteString(webutil.GetRemoteAddr(req))
	buf.WriteRune(RuneSpace)
	buf.WriteString(tf.Colorize(req.Method, ansi.ColorBlue))
	buf.WriteRune(RuneSpace)
	buf.WriteString(req.URL.Path)
	buf.WriteRune(RuneSpace)
	buf.WriteString(tf.ColorizeByStatusCode(statusCode, strconv.Itoa(statusCode)))
	buf.WriteRune(RuneSpace)
	buf.WriteString(elapsed.String())
	if len(contentType) > 0 {
		buf.WriteRune(RuneSpace)
		buf.WriteString(contentType)
	}
	buf.WriteRune(RuneSpace)
	buf.WriteString(FormatFileSize(contentLength))
}

// JSONWriteHTTPRequest marshals a request start as json.
func HTTPRequestAsJSON(req *http.Request) JSONObj {
	return JSONObj{
		"ip":   webutil.GetRemoteAddr(req),
		"verb": req.Method,
		"path": req.URL.Path,
		"host": req.Host,
	}
}

// JSONWriteHTTPResponse marshals a request as json.
func HTTPResponseAsJSON(req *http.Request, statusCode, contentLength int, contentType, contentEncoding string, elapsed time.Duration) JSONObj {
	return JSONObj{
		"ip":              webutil.GetRemoteAddr(req),
		"verb":            req.Method,
		"path":            req.URL.Path,
		"host":            req.Host,
		"contentLength":   contentLength,
		"contentType":     contentType,
		"contentEncoding": contentEncoding,
		"statusCode":      statusCode,
		JSONFieldElapsed:  Milliseconds(elapsed),
	}
}
