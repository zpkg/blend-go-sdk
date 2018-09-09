package web

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestViewFuncs(t *testing.T) {
	assert := assert.New(t)

	uuidValue := ViewFuncs()["uuid"].(func() string)()
	assert.NotEmpty(uuidValue)
	_, err := uuid.Parse(uuidValue)
	assert.Nil(err)

	jsonPrettyValue, err := ViewFuncs()["jsonPretty"].(func(interface{}) (string, error))("foo")
	assert.Nil(err)
	assert.Equal("\"foo\"\n", jsonPrettyValue)

	jsonValue, err := ViewFuncs()["jsonPretty"].(func(interface{}) (string, error))("foo")
	assert.Nil(err)
	assert.Equal("\"foo\"\n", jsonValue)

	assert.Equal("foo, bar, baz", ViewFuncs()["csv"].(func([]string) string)([]string{"foo", "bar", "baz"}))
	assert.Equal("55.44%", ViewFuncs()["pct"].(func(float64) string)(0.554433))
	assert.Equal("$44.33", ViewFuncs()["money"].(func(float64) string)(44.3322))

	formatDuration := ViewFuncs()["duration"].(func(time.Duration) string)
	for index, units := range []time.Duration{time.Hour, time.Minute, time.Second, time.Millisecond, time.Microsecond, time.Nanosecond} {
		assert.NotEmpty(formatDuration(time.Duration(index) * units))
	}

	assert.Equal("9/8", ViewFuncs()["monthDate"].(func(time.Time) string)(time.Date(2018, 9, 8, 13, 12, 11, 10, time.UTC)))
	assert.Equal("1:12PM", ViewFuncs()["kitchen"].(func(time.Time) string)(time.Date(2018, 9, 8, 13, 12, 11, 10, time.UTC)))
	assert.Equal("9/08/2018", ViewFuncs()["shortDate"].(func(time.Time) string)(time.Date(2018, 9, 8, 13, 12, 11, 10, time.UTC)))
	assert.Equal("9/08/2018 1:12:11 PM", ViewFuncs()["short"].(func(time.Time) string)(time.Date(2018, 9, 8, 13, 12, 11, 10, time.UTC)))
	assert.Equal("Sep 08, 2018 1:12:11 PM", ViewFuncs()["medium"].(func(time.Time) string)(time.Date(2018, 9, 8, 13, 12, 11, 10, time.UTC)))
}
