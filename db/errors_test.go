package db

import (
	"encoding/json"
	"testing"

	"github.com/lib/pq"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
)

const (
	errorData = `{"Severity":"ERROR", "Code":"42P01", "Message":"relation \"foo\" does not exist", "Detail":"", "Hint":"", "Position":"15", "InternalPosition":"", "InternalQuery":"", "Where":"", "Schema":"", "Table":"", "Column":"", "DataTypeName":"", "Constraint":"", "File":"parse_relation.c", "Line":"1180", "Routine":"parserOpenTable"}`
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(Error(nil))

	var err error
	assert.Nil(Error(err))

	err = exception.New("this is only a test")
	assert.True(exception.Is(Error(err), exception.Class("this is only a test")))

	var pqError pq.Error
	json.Unmarshal([]byte(errorData), &pqError)
	err = &pqError

	parsed := Error(err)
	assert.NotNil(parsed)
	assert.Equal("undefined_table", exception.As(parsed).Class())
	assert.Equal("relation \"foo\" does not exist", exception.As(parsed).Message())
}
