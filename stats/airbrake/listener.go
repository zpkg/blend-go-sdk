package airbrake

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/webutil"

	"github.com/airbrake/gobrake"
	"github.com/blend/go-sdk/logger"
)

func mustInt(value string) int64 {
	output, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return output
}

const (
	// ListenerAirbrake is the airbrake listener name.
	ListenerAirbrake = "airbrake"
)

// AddListeners adds airbrake listeners.
func AddListeners(log *logger.Logger, cfg *Config) {
	if log == nil || cfg == nil || cfg.IsZero() {
		return
	}

	// create a new reporter
	airbrake := gobrake.NewNotifierWithOptions(&gobrake.NotifierOptions{
		ProjectId:   mustInt(cfg.ProjectID),
		ProjectKey:  cfg.ProjectKey,
		Environment: cfg.Environment,
	})

	// filter airbrakes from `dev`, `ci`, and `test`.
	airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
		if noticeEnv := notice.Context["environment"]; noticeEnv == env.ServiceEnvDev ||
			noticeEnv == env.ServiceEnvCI ||
			noticeEnv == env.ServiceEnvTest {
			return nil
		}
		return notice
	})

	listener := logger.NewErrorEventListener(func(ee *logger.ErrorEvent) {
		// use our custom notice creator
		airbrake.SendNotice(NewNotice(airbrake, ee))
	})
	log.Listen(logger.Error, ListenerAirbrake, listener)
	log.Listen(logger.Fatal, ListenerAirbrake, listener)
}

var (
	defaultContextOnce sync.Once
	defaultContext     map[string]interface{}
)

func getDefaultContext() map[string]interface{} {
	defaultContextOnce.Do(func() {
		defaultContext = map[string]interface{}{
			"notifier": map[string]interface{}{
				"name":    "gobrake",
				"version": "3.4.0",
				"url":     "https://github.com/airbrake/gobrake",
			},
			"language":     runtime.Version(),
			"os":           runtime.GOOS,
			"architecture": runtime.GOARCH,
		}

		if s, err := os.Hostname(); err == nil {
			defaultContext["hostname"] = s
		}

		if wd, err := os.Getwd(); err == nil {
			defaultContext["rootDirectory"] = wd
		}

		if s := os.Getenv("GOPATH"); s != "" {
			defaultContext["gopath"] = s
		}
	})
	return defaultContext
}

// NewNotice returns a new gobrake notice.
func NewNotice(reporter *gobrake.Notifier, ee *logger.ErrorEvent) *gobrake.Notice {
	var req *http.Request
	if typed, ok := ee.State().(*http.Request); ok {
		req = typed
	}

	var notice *gobrake.Notice

	if ex := exception.As(ee.Err()); ex != nil {
		var errors []gobrake.Error
		errors = append(errors, gobrake.Error{
			Type:      exception.ErrClass(ex),
			Message:   ex.Message(),
			Backtrace: frames(ex.Stack()),
		})

		for inner := exception.As(ex.Inner()); inner != nil; inner = exception.As(inner.Inner()) {
			errors = append(errors, gobrake.Error{
				Type:      exception.ErrClass(inner),
				Message:   fmt.Sprintf("%+v", ex),
				Backtrace: frames(inner.Stack()),
			})
		}

		notice = &gobrake.Notice{
			Errors:  errors,
			Context: make(map[string]interface{}),
			Env:     make(map[string]interface{}),
			Session: make(map[string]interface{}),
			Params:  make(map[string]interface{}),
		}
	} else {
		notice = &gobrake.Notice{
			Errors: []gobrake.Error{{
				Type:    fmt.Sprint(ee.Err()),
				Message: fmt.Sprint(ee.Err()),
			}},
			Context: make(map[string]interface{}),
			Env:     make(map[string]interface{}),
			Session: make(map[string]interface{}),
			Params:  make(map[string]interface{}),
		}
	}
	for k, v := range getDefaultContext() {
		notice.Context[k] = v
	}
	notice.Context["severity"] = string(ee.Flag())

	// set requests minus headers
	if req != nil {
		notice.Context["url"] = req.URL.String()
		notice.Context["httpMethod"] = req.Method
		if ua := webutil.GetUserAgent(req); ua != "" {
			notice.Context["userAgent"] = ua
		}
		notice.Context["userAddr"] = webutil.GetRemoteAddr(req)
	}
	return notice
}

func frames(stack exception.StackTrace) (output []gobrake.StackFrame) {
	if typed, ok := stack.(exception.StackPointers); ok {
		var frame exception.Frame
		for _, ptr := range typed {
			frame = exception.Frame(ptr)
			output = append(output, gobrake.StackFrame{
				File: frame.File(),
				Func: frame.Func(),
				Line: frame.Line(),
			})
		}
	}
	return
}
