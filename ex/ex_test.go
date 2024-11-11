/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestNewOfString(t *testing.T) {
	a := assert.New(t)
	ex := As(New("this is a test"))
	a.Equal("this is a test", fmt.Sprintf("%v", ex))
	a.NotNil(ex.StackTrace)
	a.Nil(ex.Inner)
}

func TestNewOfError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("This is an error")
	wrappedErr := New(err)
	a.NotNil(wrappedErr)
	typedWrapped := As(wrappedErr)
	a.NotNil(typedWrapped)
	a.Equal("This is an error", fmt.Sprintf("%v", typedWrapped))
}

func TestNewOfException(t *testing.T) {
	a := assert.New(t)
	ex := New(Class("This is an exception"))
	wrappedEx := New(ex)
	a.NotNil(wrappedEx)
	typedWrappedEx := As(wrappedEx)
	a.Equal("This is an exception", fmt.Sprintf("%v", typedWrappedEx))
	a.Equal(ex, typedWrappedEx)
}

func TestNewOfNil(t *testing.T) {
	a := assert.New(t)

	shouldBeNil := New(nil)
	a.Nil(shouldBeNil)
	a.Equal(nil, shouldBeNil)
	a.True(nil == shouldBeNil)
}

func TestNewOfTypedNil(t *testing.T) {
	a := assert.New(t)

	var nilError error
	a.Nil(nilError)
	a.Equal(nil, nilError)

	shouldBeNil := New(nilError)
	a.Nil(shouldBeNil)
	a.True(shouldBeNil == nil)
}

func TestNewOfReturnedNil(t *testing.T) {
	a := assert.New(t)

	returnsNil := func() error {
		return nil
	}

	shouldBeNil := New(returnsNil())
	a.Nil(shouldBeNil)
	a.True(shouldBeNil == nil)

	returnsTypedNil := func() error {
		return New(nil)
	}

	shouldAlsoBeNil := returnsTypedNil()
	a.Nil(shouldAlsoBeNil)
	a.True(shouldAlsoBeNil == nil)
}

func TestError(t *testing.T) {
	a := assert.New(t)

	ex := New(Class("this is a test"))
	message := ex.Error()
	a.NotEmpty(message)
}

func TestErrorOptions(t *testing.T) {
	a := assert.New(t)

	ex := New(Class("this is a test"), OptMessage("foo"))
	message := ex.Error()
	a.NotEmpty(message)

	typed := As(ex)
	a.NotNil(typed)
	a.Equal("foo", typed.Message)
}

func TestCallers(t *testing.T) {
	a := assert.New(t)

	callStack := func() StackTrace { return Callers(DefaultStartDepth) }()

	a.NotNil(callStack)
	callstackStr := callStack.String()
	a.True(strings.Contains(callstackStr, "TestCallers"), callstackStr)
}

func TestExceptionFormatters(t *testing.T) {
	assert := assert.New(t)

	// test the "%v" formatter with just the exception class.
	class := &Ex{Class: Class("this is a test")}
	assert.Equal("this is a test", fmt.Sprintf("%v", class))

	classAndMessage := &Ex{Class: Class("foo"), Message: "bar"}
	assert.Equal("foo; bar", fmt.Sprintf("%v", classAndMessage))
}

func TestMarshalJSON(t *testing.T) {

	type ReadableStackTrace struct {
		Class   string   `json:"Class"`
		Message string   `json:"Message"`
		Inner   error    `json:"Inner"`
		Stack   []string `json:"StackTrace"`
	}

	a := assert.New(t)
	message := "new test error"
	ex := As(New(message))
	a.NotNil(ex)
	stackTrace := ex.StackTrace
	typed, isTyped := stackTrace.(StackPointers)
	a.True(isTyped)
	a.NotNil(typed)
	stackDepth := len(typed)

	jsonErr, err := json.Marshal(ex)
	a.Nil(err)
	a.NotNil(jsonErr)

	ex2 := &ReadableStackTrace{}
	err = json.Unmarshal(jsonErr, ex2)
	a.Nil(err)
	a.Len(ex2.Stack, stackDepth)
	a.Equal(message, ex2.Class)

	ex = As(New(fmt.Errorf(message)))
	a.NotNil(ex)
	stackTrace = ex.StackTrace
	typed, isTyped = stackTrace.(StackPointers)
	a.True(isTyped)
	a.NotNil(typed)
	stackDepth = len(typed)

	jsonErr, err = json.Marshal(ex)
	a.Nil(err)
	a.NotNil(jsonErr)

	ex2 = &ReadableStackTrace{}
	err = json.Unmarshal(jsonErr, ex2)
	a.Nil(err)
	a.Len(ex2.Stack, stackDepth)
	a.Equal(message, ex2.Class)
}

func TestJSON(t *testing.T) {
	assert := assert.New(t)

	ex := New("this is a test",
		OptMessage("test message"),
		OptInner(New("inner exception", OptMessagef("inner test message"))),
	)

	contents, err := json.Marshal(ex)
	assert.Nil(err)

	var verify Ex
	err = json.Unmarshal(contents, &verify)
	assert.Nil(err)

	assert.Equal(ErrClass(ex), ErrClass(verify))
	assert.Equal(ErrMessage(ex), ErrMessage(verify))
	assert.NotNil(verify.Inner)
	assert.Equal(ErrClass(ErrInner(ex)), ErrClass(ErrInner(verify)))
	assert.Equal(ErrMessage(ErrInner(ex)), ErrMessage(ErrInner(verify)))
}

func TestNest(t *testing.T) {
	a := assert.New(t)

	ex1 := As(New("this is an error"))
	ex2 := As(New("this is another error"))
	err := As(Nest(ex1, ex2))

	a.NotNil(err)
	a.NotNil(err.Inner)
	a.NotEmpty(err.Error())

	a.True(Is(ex1, Class("this is an error")))
	a.True(Is(ex1.Inner, Class("this is another error")))
}

func TestNestNil(t *testing.T) {
	a := assert.New(t)

	var ex1 error
	var ex2 error
	var ex3 error

	err := Nest(ex1, ex2, ex3)
	a.Nil(err)
	a.Equal(nil, err)
	a.True(nil == err)
}

func TestExceptionFormat(t *testing.T) {
	assert := assert.New(t)

	e := &Ex{Class: fmt.Errorf("this is only a test")}
	output := fmt.Sprintf("%v", e)
	assert.Equal("this is only a test", output)

	output = fmt.Sprintf("%+v", e)
	assert.Equal("this is only a test", output)

	e = &Ex{
		Class: fmt.Errorf("this is only a test"),
		StackTrace: StackStrings([]string{
			"foo",
			"bar",
		}),
	}

	output = fmt.Sprintf("%+v", e)
	assert.Equal("this is only a test\nfoo\nbar", output)
}

func TestExceptionPrintsInner(t *testing.T) {
	assert := assert.New(t)

	ex := New("outer", OptInner(New("middle", OptInner(New("terminal")))))

	output := fmt.Sprintf("%v", ex)

	assert.Contains(output, "outer")
	assert.Contains(output, "middle")
	assert.Contains(output, "terminal")

	output = fmt.Sprintf("%+v", ex)

	assert.Contains(output, "outer")
	assert.Contains(output, "middle")
	assert.Contains(output, "terminal")
}

type structuredError struct {
	value string
}

func (err structuredError) Error() string {
	return err.value
}

func TestException_ErrorsIsCompatability(t *testing.T) {
	assert := assert.New(t)

	{ // Single nesting, Ex is outermost
		innerErr := errors.New("inner")
		outerErr := New("outer", OptInnerClass(innerErr))

		assert.True(errors.Is(outerErr, innerErr))
	}

	{ // Single nesting, Ex is innermost
		innerErr := New("inner")
		outerErr := fmt.Errorf("outer: %w", innerErr)

		assert.True(errors.Is(outerErr, Class("inner")))
	}

	{ // Triple nesting, including Ex and non-Ex
		firstErr := errors.New("inner most")
		secondErr := fmt.Errorf("standard err: %w", firstErr)
		thirdErr := New("ex err", OptInner(secondErr))
		fourthErr := New("outer most", OptInner(thirdErr))

		assert.True(errors.Is(fourthErr, firstErr))
		assert.True(errors.Is(fourthErr, secondErr))
		assert.True(errors.Is(fourthErr, Class("ex err")))
	}

	{ // Target is nested in an Ex class and not in Inner chain
		firstErr := errors.New("inner most")
		secondErr := fmt.Errorf("standard err: %w", firstErr)
		thirdErr := New(secondErr, OptInner(fmt.Errorf("another cause")))

		assert.True(errors.Is(thirdErr, firstErr))
		assert.True(errors.Is(thirdErr, secondErr))
	}
}

func TestException_ErrorsAsCompatability(t *testing.T) {
	assert := assert.New(t)

	{ // Single nesting, targeting non-Ex
		innerErr := structuredError{"inner most"}
		outerErr := New("outer", OptInner(innerErr))

		var matchedErr structuredError
		assert.True(errors.As(outerErr, &matchedErr))
		assert.Equal("inner most", matchedErr.value)
	}

	{ // Single nesting, targeting Ex
		innerErr := New("outer most")
		outerErr := fmt.Errorf("outer err: %w", innerErr)

		var matchedErr *Ex
		assert.True(errors.As(outerErr, &matchedErr))
		assert.Equal("outer most", matchedErr.Class.Error())
	}

	{ // Single nesting, targeting inner Ex class
		innerErr := New(structuredError{"inner most"})
		outerErr := New("outer most", OptInner(innerErr))

		var matchedErr structuredError
		assert.True(errors.As(outerErr, &matchedErr))
		assert.Equal("inner most", matchedErr.value)
	}

	{ // Triple Nesting, targeting non-Ex
		firstErr := structuredError{"inner most"}
		secondErr := fmt.Errorf("standard err: %w", firstErr)
		thirdErr := New("ex err", OptInner(secondErr))
		fourthErr := New("outer most", OptInner(thirdErr))

		var matchedErr structuredError
		assert.True(errors.As(fourthErr, &matchedErr))
		assert.Equal("inner most", matchedErr.value)
	}
}
