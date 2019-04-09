package r2

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestRequestDo(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
	defer server.Close()

	res, err := New(server.URL).Do()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestDoHeaders(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value := r.Header.Get("foo"); value != "bar" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad value for foo: %#v\n", r.PostForm)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
	defer server.Close()

	res, err := New(server.URL, OptHeaderValue("foo", "bar")).Do()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestDoPostForm(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v!\n", err)
			return
		}
		if value := r.PostForm.Get("foo"); value != "bar" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad value for foo: %#v\n", r.PostForm)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
	defer server.Close()

	res, err := New(server.URL,
		OptPost(),
		OptPostFormValue("foo", "bar"),
	).Do()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode, readString(res.Body))
}

func readString(r io.Reader) string {
	contents, _ := ioutil.ReadAll(r)
	return string(contents)
}
