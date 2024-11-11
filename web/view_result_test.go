/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/env"
	"github.com/zpkg/blend-go-sdk/ex"
	"github.com/zpkg/blend-go-sdk/webutil"
)

type testViewModel struct {
	Text string
}

type errorWriter struct {
	ResponseWriter
}

func (ew *errorWriter) Write(_ []byte) (int, error) {
	return -1, fmt.Errorf("error")
}

func TestViewResultRender(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	rc := NewCtx(webutil.NewMockResponse(buffer), nil)

	assert.True(ex.Is((&ViewResult{}).Render(nil), ErrUnsetViewTemplate))

	testView := template.New("testView")
	_, err := testView.Parse("{{ .ViewModel.Text }}")
	assert.Nil(err)

	vr := &ViewResult{
		StatusCode: http.StatusOK,
		ViewModel:  testViewModel{Text: "bar"},
		Template:   testView,
	}

	err = vr.Render(rc)
	assert.Nil(err)

	assert.NotZero(buffer.Len())
	assert.True(strings.Contains(buffer.String(), "bar"))

	testView = template.New("testView")
	_, err = testView.Parse("{{ .Env.String \"HELLO\" }}")
	assert.Nil(err)

	expected := "world"

	curr := env.Env()
	defer func() {
		env.SetEnv(curr)
	}()

	env.SetEnv(env.New())
	env.Env().Set("HELLO", expected)

	vr = &ViewResult{
		StatusCode: http.StatusOK,
		ViewModel:  nil,
		Template:   testView,
	}

	err = vr.Render(rc)
	assert.Nil(err)

	assert.NotZero(buffer.Len())
	assert.True(strings.Contains(buffer.String(), expected))
}

func TestViewResultRenderErrorResponse(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	rc := NewCtx(webutil.NewMockResponse(buffer), nil)

	assert.True(ex.Is((&ViewResult{}).Render(nil), ErrUnsetViewTemplate))

	testView := template.New("testView")
	_, err := testView.Parse("{{ .ViewModel.Text }}")
	assert.Nil(err)

	vr := &ViewResult{
		StatusCode: http.StatusOK,
		ViewModel:  testViewModel{Text: "bar"},
		Template:   testView,
	}

	rc.Response = &errorWriter{ResponseWriter: rc.Response}

	err = vr.Render(rc)
	assert.NotNil(err)
}

func TestViewResultRenderError(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	rc := NewCtx(webutil.NewMockResponse(buffer), nil)

	testView := template.New("testView")
	_, err := testView.Parse("{{.ViewModel.Foo}}")
	assert.Nil(err)

	vr := &ViewResult{
		StatusCode: http.StatusOK,
		ViewModel:  testViewModel{Text: "bar"},
		Template:   testView,
	}

	err = vr.Render(rc)
	assert.NotNil(err)
	assert.NotZero(buffer.Len())
}

func TestViewResultRenderErrorTemplate(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	rc := NewCtx(webutil.NewMockResponse(buffer), nil)

	views := template.New("main")

	errorTemplate := views.New(DefaultTemplateNameInternalError)
	_, err := errorTemplate.Parse("{{.}}")
	assert.Nil(err)

	testView := views.New("testView")
	_, err = testView.Parse("Foo: {{.ViewModel.Foo}}")
	assert.Nil(err)

	vr := &ViewResult{
		StatusCode: http.StatusOK,
		ViewModel:  testViewModel{Text: "bar"},
		Template:   testView,
	}

	err = vr.Render(rc)
	assert.NotNil(err)
}

func TestViewResultErrorNestedViews(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	rc := NewCtx(webutil.NewMockResponse(buffer), nil)

	views := template.New("main")

	errorTemplate := views.New(DefaultTemplateNameInternalError)
	_, err := errorTemplate.Parse("{{.}}")
	assert.Nil(err)

	outerView := views.New("outerView")
	_, err = outerView.Parse("Inner: {{ template \"innerView\" . }}")
	assert.Nil(err)

	testView := views.New("innerView")
	_, err = testView.Parse("Foo: {{.ViewModel.Foo}}")
	assert.Nil(err)

	vr := &ViewResult{
		StatusCode: http.StatusOK,
		ViewModel:  testViewModel{Text: "bar"},
		Template:   outerView,
	}

	err = vr.Render(rc)
	assert.NotNil(err)
}
