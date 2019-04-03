package web

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestViewModelWrap(t *testing.T) {
	assert := assert.New(t)

	indexTemplate := `{{ define "index" }}{{ range $index, $obj := .ViewModel }}<div>{{ template "control" ( $.Wrap $obj ) }}</div>{{ end }}{{ end }}`
	controlTemplate := `{{ define "control" }}{{ if .Ctx }}{{ .ViewModel }}{{ end }}{{ end }}`

	app := New()
	app.Views.AddLiterals(indexTemplate, controlTemplate)

	app.GET("/", func(r *Ctx) Result {
		return r.Views.View("index", []string{"foo", "bar", "baz"})
	})

	meta, err := MockGet(app, "/")
	assert.Nil(err)
	contents, err := MockReadBytes(meta, err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(contents))
	assert.Equal("<div>foo</div><div>bar</div><div>baz</div>", string(contents))
}
