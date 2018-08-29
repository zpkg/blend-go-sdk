package main

import (
	"encoding/json"
	"net/http"
	"os"
)

// Common constants.
const (
	HeaderContentLength        = "Content-Length"
	HeaderContentType          = "Content-Type"
	HeaderServer               = "Server"
	HeaderDate                 = "Date"
	ServerName                 = "golang"
	ContentTypeApplicationJSON = "application/json; charset=UTF-8"
	MessageText                = "Hello, World!"
	JSONContentLength          = "27"
)

// Message is a json message.
type Message struct {
	Message string `json:"message"`
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderContentType, ContentTypeApplicationJSON)
	w.Header().Set(HeaderServer, ServerName)
	json.NewEncoder(w).Encode(&Message{Message: MessageText})
}

func port() string {
	envPort := os.Getenv("PORT")
	if len(envPort) != 0 {
		return envPort
	}
	return "9090"
}

func main() {
	http.HandleFunc("/json", jsonHandler)
	http.ListenAndServe("127.0.0.1:"+port(), nil)
}
