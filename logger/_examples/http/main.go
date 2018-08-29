package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"

	"github.com/blend/go-sdk/logger"
)

var pool = logger.NewBufferPool(16)

func logged(log *logger.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		log.Trigger(logger.NewHTTPRequestEvent(req))
		rw := logger.NewResponseWriter(res)
		handler(rw, req)
		log.Trigger(logger.NewHTTPResponseEvent(req).WithStatusCode(rw.StatusCode()).WithContentLength(rw.ContentLength()).WithElapsed(time.Now().Sub(start)))
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
	res.Write([]byte(`{"status":"not ok."}`))
}

func errorHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte(`{"status":"not ok."}`))
}

func warningHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(`{"status":"not ok."}`))
}

func subContextHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"did sub-context things"}`))
}

func auditHandler(res http.ResponseWriter, req *http.Request) {
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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log := logger.NewFromEnv().WithEnabled(logger.Info, logger.Audit)

	http.HandleFunc("/", logged(log, indexHandler))

	http.HandleFunc("/sub-context", logged(log, subContextHandler))
	http.HandleFunc("/fatalerror", logged(log, fatalErrorHandler))
	http.HandleFunc("/error", logged(log, errorHandler))
	http.HandleFunc("/warning", logged(log, warningHandler))
	http.HandleFunc("/audit", logged(log, auditHandler))

	http.HandleFunc("/bench/logged", logged(log, indexHandler))
	http.HandleFunc("/bench/stdout", stdoutLogged(indexHandler))

	log.SyncInfof("Listening on :%s", port())
	log.SyncInfof("Events %s", log.Flags().String())
	log.Fatal(http.ListenAndServe(":"+port(), nil))
}
