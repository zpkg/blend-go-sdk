package web

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestJSONResultProvider(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(http.StatusNotFound, JSON.NotFound().(*JSONResult).StatusCode)
	assert.Equal("Not Found", JSON.NotFound().(*JSONResult).Response)

	assert.Equal(http.StatusForbidden, JSON.NotAuthorized().(*JSONResult).StatusCode)
	assert.Equal("Not Authorized", JSON.NotAuthorized().(*JSONResult).Response)

	assert.Equal(http.StatusBadRequest, JSON.BadRequest(nil).(*JSONResult).StatusCode)
	assert.Equal("Bad Request", JSON.BadRequest(nil).(*JSONResult).Response)

	assert.Equal(http.StatusBadRequest, JSON.BadRequest(fmt.Errorf("bad-request")).(*JSONResult).StatusCode)
	assert.Equal("bad-request", JSON.BadRequest(fmt.Errorf("bad-request")).(*JSONResult).Response)

	assert.Equal(http.StatusOK, JSON.OK().(*JSONResult).StatusCode)
	assert.Equal("OK!", JSON.OK().(*JSONResult).Response)

	assert.Equal(http.StatusBadGateway, JSON.Status(http.StatusBadGateway, "test").(*JSONResult).StatusCode)
	assert.Equal("test", JSON.Status(http.StatusBadGateway, "test").(*JSONResult).Response)

	assert.Equal(http.StatusOK, JSON.Result("foo").(*JSONResult).StatusCode)
	assert.Equal("foo", JSON.Result("foo").(*JSONResult).Response)

	internalError := JSON.InternalError(fmt.Errorf("only a test"))

	typed, ok := internalError.(*loggedErrorResult)
	assert.True(ok)
	assert.Equal(fmt.Errorf("only a test"), typed.Error)
	inner := typed.Result.(*JSONResult)
	assert.Equal(http.StatusInternalServerError, inner.StatusCode)
	assert.Equal("only a test", inner.Response)
}
