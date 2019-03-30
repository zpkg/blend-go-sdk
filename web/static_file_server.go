package web

import (
	"net/http"
	"os"
	"regexp"

	"github.com/blend/go-sdk/logger"
)

// NewStaticFileServer returns a new static file cache.
func NewStaticFileServer(searchPaths ...http.FileSystem) *StaticFileServer {
	return &StaticFileServer{
		SearchPaths: searchPaths,
	}
}

// StaticFileServer is a cache of static files.
type StaticFileServer struct {
	Log          logger.Log
	SearchPaths  []http.FileSystem
	RewriteRules []RewriteRule
	Middleware   Action
	Headers      http.Header
}

// GetHeaders implements part of the fileserver spec.
func (sc *StaticFileServer) GetHeaders() http.Header {
	return sc.Headers
}

// AddHeader adds a header to the static cache results.
func (sc *StaticFileServer) AddHeader(key, value string) {
	if sc.Headers == nil {
		sc.Headers = http.Header{}
	}
	sc.Headers[key] = append(sc.Headers[key], value)
}

// AddRewriteRule adds a static re-write rule.
func (sc *StaticFileServer) AddRewriteRule(match string, action RewriteAction) error {
	expr, err := regexp.Compile(match)
	if err != nil {
		return err
	}
	sc.RewriteRules = append(sc.RewriteRules, RewriteRule{
		MatchExpression: match,
		expr:            expr,
		Action:          action,
	})
	return nil
}

// Action is the entrypoint for the static server.
// It will run middleware if specified before serving the file.
func (sc *StaticFileServer) Action(r *Ctx) Result {
	if sc.Middleware != nil {
		return sc.Middleware(r)
	}
	return sc.ServeFile(r)
}

// ResolveFile resolves a file from rewrite rules and search paths.
func (sc *StaticFileServer) ResolveFile(filePath string) (f http.File, err error) {
	for _, rule := range sc.rewriteRules {
		if matched, newFilePath := rule.Apply(filePath); matched {
			filePath = newFilePath
		}
	}

	// for each searchpath, sniff if the file exists ...
	var openErr error
	for _, searchPath := range sc.searchPaths {
		f, openErr = searchPath.Open(filePath)
		if openErr == nil {
			break
		}
	}
	if openErr != nil && !os.IsNotExist(openErr) {
		err = openErr
		return
	}
	return
}

// ServeFile writes the file to the response without running middleware.
func (sc *StaticFileServer) ServeFile(r *Ctx) Result {
	filePath, err := r.RouteParam("filepath")
	if err != nil {
		return r.DefaultProvider.BadRequest(err)
	}

	for key, values := range sc.headers {
		for _, value := range values {
			r.Response.Header().Set(key, value)
		}
	}

	f, err := sc.ResolveFile(filePath)
	if f == nil || (err != nil && os.IsNotExist(err)) {
		return r.DefaultProvider.NotFound()
	}
	if err != nil {
		return r.DefaultProvider.InternalError(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return r.DefaultProvider.InternalError(err)
	}
	http.ServeContent(r.Response, r.Request, filePath, finfo.ModTime(), f)
	return nil
}
