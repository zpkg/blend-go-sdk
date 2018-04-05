package web

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestViewCacheAddRawViews(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	vc.AddLiterals(`{{ define "test" }}<h1> This is a test. </h1>{{ end }}`)

	view, err := vc.Parse()
	assert.Nil(err)
	assert.NotNil(view)

	buf := bytes.NewBuffer(nil)
	assert.Nil(view.ExecuteTemplate(buf, "test", nil))
	assert.NotEmpty(buf.String())
}
