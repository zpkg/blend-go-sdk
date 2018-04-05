package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

var pool = logger.NewBufferPool(16)

func logged(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		logger.Default().Trigger(logger.NewWebRequestStartEvent(req))
		rw := logger.NewResponseWriter(res)
		handler(rw, req)
		logger.Default().Trigger(logger.NewWebRequestEvent(req).WithStatusCode(rw.StatusCode()).WithContentLength(int64(rw.ContentLength())).WithElapsed(time.Now().Sub(start)))
	}
}

func stdoutLogged(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		handler(res, req)
		fmt.Printf("%s %s %s %s %s %s %s\n",
			time.Now().UTC().Format(time.RFC3339),
			"web.request",
			req.Method,
			req.URL.Path,
			"200",
			time.Since(start).String(),
			"??",
		)
	}
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"ok!"}`))
}

func fatalErrorHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	logger.Default().Fatal(exception.New("this is a fatal exception"))
	res.Write([]byte(`{"status":"not ok."}`))
}

func errorHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	logger.Default().Error(exception.New("this is an exception"))
	res.Write([]byte(`{"status":"not ok."}`))
}

func warningHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	logger.Default().Warning(exception.New("this is warning"))
	res.Write([]byte(`{"status":"not ok."}`))
}

func auditHandler(res http.ResponseWriter, req *http.Request) {
	logger.Default().Trigger(logger.NewAuditEvent(logger.GetIP(req), "viewed", "audit route").WithExtra(map[string]string{
		"remoteAddr": req.RemoteAddr,
		"host":       req.Host,
	}))
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"audit logged!"}`))
}

func port() string {
	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		return envPort
	}
	return "8888"
}

func main() {
	logger.SetDefault(logger.NewFromEnv().WithEnabled(logger.Audit))

	http.HandleFunc("/", logged(indexHandler))

	http.HandleFunc("/fatalerror", logged(fatalErrorHandler))
	http.HandleFunc("/error", logged(errorHandler))
	http.HandleFunc("/warning", logged(warningHandler))
	http.HandleFunc("/audit", logged(auditHandler))

	http.HandleFunc("/bench/logged", logged(indexHandler))
	http.HandleFunc("/bench/stdout", stdoutLogged(indexHandler))

	logger.Default().Infof("Listening on :%s", port())
	logger.Default().Infof("Events %s", logger.Default().Flags().String())
	log.Fatal(http.ListenAndServe(":"+port(), nil))
}
