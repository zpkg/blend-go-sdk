package proxy

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/blend/go-sdk/logger"
)

// NewUpstream returns a new upstram.
func NewUpstream(target *url.URL) *Upstream {
	return &Upstream{
		URL: target,
		// Hop-by-hop headers. These are removed when sent to the backend.
		// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
		HopHeaders: []string{
			"Connection",
			"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
			"Keep-Alive",
			"Proxy-Authenticate",
			"Proxy-Authorization",
			"Te",      // canonicalized version of "TE"
			"Trailer", // not Trailers per URL above; http://www.rfc-editor.org/errata_search.php?eid=4522
			"Transfer-Encoding",
			"Upgrade",
		},
	}
}

// Upstream represents a proxyable server.
type Upstream struct {
	// Name is the name of the upstream.
	Name string
	// Log is a logger agent.
	Log *logger.Logger
	// URL represents the target of the upstream.
	URL *url.URL
	// Transport represents the underlying connection to the upstream.
	Transport *http.Transport

	// TLSInsecureSkipVerify skips tls cert ferification on upstreams.
	// It is typically used in debugging.
	TLSInsecureSkipVerify bool
	// Close closes the connections on completion.
	Close bool
	// FlushInterval defines the buffer flush interval for the default transports.
	FlushInterval time.Duration

	// BufferPool allows re-use of data buffers between requests.
	BufferPool BufferPool

	// HopHeaders indicate headers we should strip from upstream requests.
	HopHeaders []string

	// StaticDestination prevents forwarding of url or querystring params to the upstream.
	StaticDestination bool
}

// WithName sets the name field of the upstream.
func (u *Upstream) WithName(name string) *Upstream {
	u.Name = name
	return u
}

// WithLogger sets the logger agent for the upstream.
func (u *Upstream) WithLogger(log *logger.Logger) *Upstream {
	u.Log = log
	return u
}

// WithTLSInsecureSkipVerify sets a property and returns the upstream reference.
func (u *Upstream) WithTLSInsecureSkipVerify(insecureSkipVerify bool) *Upstream {
	u.TLSInsecureSkipVerify = insecureSkipVerify
	return u
}

// WithClose sets a property and returns the upstream reference.
func (u *Upstream) WithClose(close bool) *Upstream {
	u.Close = close
	return u
}

// WithFlushInterval sets a property and returns the upstream reference.
func (u *Upstream) WithFlushInterval(interval time.Duration) *Upstream {
	u.FlushInterval = interval
	return u
}

// WithStaticDestination sets a property and returns the upstream reference.
func (u *Upstream) WithStaticDestination(value bool) *Upstream {
	u.StaticDestination = value
	return u
}

// WithBufferPool sets a property and returns the upstream reference.
func (u *Upstream) WithBufferPool(bufferPool BufferPool) *Upstream {
	u.BufferPool = bufferPool
	return u
}

// WithHopHeaders sets a property and returns the upstream reference.
func (u *Upstream) WithHopHeaders(headers ...string) *Upstream {
	u.HopHeaders = append(u.HopHeaders, headers...)
	return u
}

// WithoutHopHeaders sets a property and returns the upstream reference.
func (u *Upstream) WithoutHopHeaders(headers ...string) *Upstream {
	var included []string
	for _, h := range u.HopHeaders {
		include := true
		for _, strike := range headers {
			include = include && strike != h
		}
		if include {
			included = append(included, h)
		}
	}
	u.HopHeaders = included
	return u
}

func (u *Upstream) setDestination(req *http.Request) error {
	req.URL.Scheme = u.URL.Scheme
	req.URL.Host = u.URL.Host
	if u.StaticDestination {
		req.Host = u.URL.Host
		req.URL.Path = u.URL.Path
		req.URL.RawQuery = u.URL.RawQuery
	} else {
		req.URL.Path = singleJoiningSlash(u.URL.Path, req.URL.Path)
		if u.URL.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = u.URL.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = u.URL.RawQuery + "&" + req.URL.RawQuery
		}
	}
	return nil
}

// ServeHTTP
func (u *Upstream) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var statusCode int
	var contentLength int64

	if u.Log != nil {
		// log the start event
		u.Log.Trigger(logger.NewWebRequestStartEvent(req))
	}

	start := time.Now()

	// create the transport if it doesn't exist.
	if u.Transport == nil {
		u.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: u.TLSInsecureSkipVerify,
			},
		}
	}

	ctx := req.Context()
	if cn, ok := rw.(http.CloseNotifier); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
		notifyChan := cn.CloseNotify()
		go func() {
			select {
			case <-notifyChan:
				statusCode = http.StatusRequestTimeout
				cancel()
			case <-ctx.Done():
			}
		}()
	}

	if u.Log != nil {
		// log events.
		defer func() {
			wre := logger.NewWebRequestEvent(req).WithStatusCode(statusCode).WithContentLength(contentLength).WithElapsed(time.Since(start))
			if value := rw.Header().Get("Content-Type"); len(value) > 0 {
				wre = wre.WithContentType(value)
			}
			if value := rw.Header().Get("Content-Encoding"); len(value) > 0 {
				wre = wre.WithContentEncoding(value)
			}
			u.Log.Trigger(wre)
		}()
	}

	outreq := new(http.Request)
	*outreq = *req // includes shallow copies of maps, but okay
	if req.ContentLength == 0 {
		outreq.Body = nil // Issue 16036: nil Body for http.Transport retries
	}

	u.setDestination(outreq)
	outreq = outreq.WithContext(ctx)
	outreq.Close = u.Close

	// We are modifying the same underlying map from req (shallow
	// copied above) so we only copy it if necessary.
	copiedHeaders := false

	// Remove hop-by-hop headers listed in the "Connection" header.
	// See RFC 2616, section 14.10.
	if c := outreq.Header.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				if !copiedHeaders {
					outreq.Header = make(http.Header)
					copyHeader(outreq.Header, req.Header)
					copiedHeaders = true
				}
				outreq.Header.Del(f)
			}
		}
	}

	// Remove hop-by-hop headers to the backend. Especially
	// important is "Connection" because we want a persistent
	// connection, regardless of what the client sent to us.
	for _, h := range u.HopHeaders {
		if outreq.Header.Get(h) != "" {
			if !copiedHeaders {
				outreq.Header = make(http.Header)
				copyHeader(outreq.Header, req.Header)
				copiedHeaders = true
			}
			outreq.Header.Del(h)
		}
	}

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := outreq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outreq.Header.Set("X-Forwarded-For", clientIP)
	}

	// Add the forwarded proto if it doesn't exist.
	if _, ok := outreq.Header["X-Forwarded-Proto"]; !ok {
		outreq.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
	}

	res, err := u.Transport.RoundTrip(outreq)
	if err != nil {
		statusCode = http.StatusBadGateway
		if u.Log != nil {
			u.Log.Error(err)
		}
		rw.WriteHeader(statusCode)
		return
	}

	// Remove hop-by-hop headers listed in the
	// "Connection" header of the response.
	if c := res.Header.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				res.Header.Del(f)
			}
		}
	}

	for _, h := range u.HopHeaders {
		res.Header.Del(h)
	}

	copyHeader(rw.Header(), res.Header)

	// The "Trailer" header isn't included in the Transport's response,
	// at least for *http.Transport. Build it up from Trailer.
	if len(res.Trailer) > 0 {
		trailerKeys := make([]string, 0, len(res.Trailer))
		for k := range res.Trailer {
			trailerKeys = append(trailerKeys, k)
		}
		rw.Header().Add("Trailer", strings.Join(trailerKeys, ", "))
	}

	statusCode = res.StatusCode
	rw.WriteHeader(res.StatusCode)
	if len(res.Trailer) > 0 {
		// Force chunking if we saw a response trailer.
		// This prevents net/http from calculating the length for short
		// bodies and adding a Content-Length.
		if fl, ok := rw.(http.Flusher); ok {
			fl.Flush()
		}
	}
	contentLength, err = u.copyResponse(rw, res.Body)
	if err != nil && err != io.EOF && u.Log != nil {
		u.Log.Error(err)
	}
	res.Body.Close() // close now, instead of defer, to populate res.Trailer
	copyHeader(rw.Header(), res.Trailer)
}

func (u *Upstream) copyResponse(dst io.Writer, src io.Reader) (contentLength int64, err error) {
	if u.FlushInterval != 0 {
		if wf, ok := dst.(writeFlusher); ok {
			mlw := &maxLatencyWriter{
				dst:     wf,
				latency: u.FlushInterval,
				done:    make(chan bool),
			}
			go mlw.flushLoop()
			defer mlw.stop()
			dst = mlw
		}
	}

	var buf []byte
	if u.BufferPool != nil {
		buf = u.BufferPool.Get()
	}
	contentLength, err = u.copyBuffer(dst, src, buf)
	if u.BufferPool != nil {
		u.BufferPool.Put(buf)
	}
	return
}

func (u *Upstream) copyBuffer(dst io.Writer, src io.Reader, buf []byte) (int64, error) {
	if len(buf) == 0 {
		buf = make([]byte, 32*1024)
	}
	var written int64
	for {
		nr, rerr := src.Read(buf)
		if rerr != nil && rerr != io.EOF {
			return written, rerr
		}
		if nr > 0 {
			nw, werr := dst.Write(buf[:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if werr != nil {
				return written, werr
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
		}
		if rerr != nil {
			return written, rerr
		}
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

type writeFlusher interface {
	io.Writer
	http.Flusher
}

type maxLatencyWriter struct {
	dst     writeFlusher
	latency time.Duration

	mu   sync.Mutex // protects Write + Flush
	done chan bool
}

func (m *maxLatencyWriter) Write(p []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.dst.Write(p)
}

func (m *maxLatencyWriter) flushLoop() {
	t := time.NewTicker(m.latency)
	defer t.Stop()
	for {
		select {
		case <-m.done:
			return
		case <-t.C:
			m.mu.Lock()
			m.dst.Flush()
			m.mu.Unlock()
		}
	}
}

func (m *maxLatencyWriter) stop() { m.done <- true }
