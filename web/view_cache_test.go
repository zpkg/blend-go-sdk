package web

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
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

func TestViewCacheCached(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	assert.True(vc.Cached())
	vc.AddLiterals(`{{ define "foo" }}bar{{ end }}`)
	assert.Nil(vc.Initialize())

	tmp, err := vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf := bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("bar", buf.String())

	vc.viewLiterals = []string{`{{ define "foo" }}baz{{ end }}`}
	tmp, err = vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf = bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("bar", buf.String())
}

func TestViewCacheCachingDisabled(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	assert.True(vc.Cached())
	vc.WithCached(false)
	vc.AddLiterals(`{{ define "foo" }}bar{{ end }}`)
	assert.Nil(vc.Initialize())

	tmp, err := vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf := bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("bar", buf.String())

	vc.viewLiterals = []string{`{{ define "foo" }}baz{{ end }}`}
	tmp, err = vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf = bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("baz", buf.String())
}

func TestViewCacheBadRequest(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	assert.Nil(vc.Initialize())

	vr, _ := vc.BadRequest(exception.Class("only a test")).(*ViewResult)
	assert.Equal(vc.BadRequestTemplateName(), vr.ViewName)
	assert.Equal(http.StatusBadRequest, vr.StatusCode)
	assert.NotNil(vr.Views)
	assert.NotNil(vr.ViewModel)
}

func TestViewCacheInternalError(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	assert.Nil(vc.Initialize())

	ler, _ := vc.InternalError(exception.Class("only a test")).(*loggedErrorResult)
	assert.NotNil(ler)
	vr := ler.Result.(*ViewResult)
	assert.Equal(vc.InternalErrorTemplateName(), vr.ViewName)
	assert.Equal(http.StatusInternalServerError, vr.StatusCode)
	assert.NotNil(vr.Views)
	assert.NotNil(vr.Template)
	assert.NotNil(vr.ViewModel)
}

func TestViewCacheNotFound(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	assert.Nil(vc.Initialize())

	vr, _ := vc.NotFound().(*ViewResult)
	assert.NotNil(vr)
	assert.Equal(vc.NotFoundTemplateName(), vr.ViewName)
	assert.Equal(http.StatusNotFound, vr.StatusCode)
	assert.NotNil(vr.Views)
	assert.NotNil(vr.Template)
	assert.Nil(vr.ViewModel)
}

func TestViewCacheNotAuthorized(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	assert.Nil(vc.Initialize())

	vr, _ := vc.NotAuthorized().(*ViewResult)
	assert.NotNil(vr)
	assert.Equal(vc.NotAuthorizedTemplateName(), vr.ViewName)
	assert.Equal(http.StatusForbidden, vr.StatusCode)
	assert.NotNil(vr.Views)
	assert.NotNil(vr.Template)
	assert.Nil(vr.ViewModel)
}

func TestViewCacheView(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	vc.AddLiterals(
		`{{ define "test" }}{{.ViewModel}}{{end}}`,
	)
	assert.Nil(vc.Initialize())

	vr, _ := vc.View("test", "foo").(*ViewResult)
	assert.NotNil(vr)
	assert.Equal("test", vr.ViewName)
	assert.Equal(http.StatusOK, vr.StatusCode)
	assert.Equal("foo", vr.ViewModel)
	assert.NotNil(vr.Template)
	assert.NotNil(vr.Views)

	// handle if the view is not found ...
	ler, _ := vc.View("not-test", "foo").(*loggedErrorResult)
	assert.NotNil(ler)
	vr, _ = ler.Result.(*ViewResult)
	assert.Equal(vc.InternalErrorTemplateName(), vr.ViewName)
	assert.Equal(http.StatusInternalServerError, vr.StatusCode)
	assert.NotNil(vr.Template)
	assert.NotNil(vr.Views)
}

func TestViewCacheViewError(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()

	vcr := vc.viewError(exception.Class("test error")).(*ViewResult)
	assert.NotNil(vcr)
	assert.Equal(DefaultTemplateNameInternalError, vcr.ViewName)
	assert.Equal(http.StatusInternalServerError, vcr.StatusCode)
	assert.Equal(exception.Class("test error"), vcr.ViewModel)
	assert.NotNil(vcr.Template)
	assert.NotNil(vcr.Views)
}
