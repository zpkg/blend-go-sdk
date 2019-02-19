package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestLogWriterWrite(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.showHeadings = false
	writer.useColor = false
	writer.Write(Messagef(Info, "test string"))
	assert.Equal("[info] test string\n", string(buffer.Bytes()))
}

func TestLogWriterWriteWithLabel(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.showHeadings = true
	writer.useColor = false
	writer.Write(Messagef(Info, "test string").WithHeadings("unit-test"))
	assert.Equal("[unit-test] [info] test string\n", string(buffer.Bytes()))
}

func TestLogWriterWriteWithLabelColorized(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.showHeadings = true
	writer.useColor = true
	writer.Write(Messagef(Info, "test string").WithHeadings("unit-test"))
	assert.Equal("["+ColorBlue.Apply("unit-test")+"] ["+ColorLightWhite.Apply("info")+"] test string\n", string(buffer.Bytes()))
}

func TestWriterErrorOutputCoalesced(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.useColor = false
	writer.WriteError(Messagef(Error, "test %s", "string"))
	assert.Equal("[error] test string\n", string(buffer.Bytes()))
}

func TestWriterErrorOutput(t *testing.T) {
	assert := assert.New(t)

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	writer := NewTextWriter(stdout).WithErrorOutput(stderr)
	writer.showTimestamp = false
	writer.useColor = false

	writer.WriteError(Messagef(Error, "test %s", "string"))
	assert.Equal(0, stdout.Len())
	assert.Equal("[error] test string\n", string(stderr.Bytes()))
}

func TestWriterLabels(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.useColor = false
	writer.WriteError(Messagef(Error, "test %s", "string").WithLabel("foo", "bar").WithLabel("moo", "boo"))
	assert.True(strings.HasPrefix(buffer.String(), "[error] test string"))
}

func TestFormatTimestamp(t *testing.T) {
	assert := assert.New(t)

	tsValues := [7]time.Time{
		time.Date(2019, time.February, 19, 15, 12, 47, 123000000, time.UTC),
		time.Date(2019, time.February, 19, 15, 12, 47, 123400000, time.UTC),
		time.Date(2019, time.February, 19, 15, 12, 47, 123450000, time.UTC),
		time.Date(2019, time.February, 19, 15, 12, 47, 123456000, time.UTC),
		time.Date(2019, time.February, 19, 15, 12, 47, 123456700, time.UTC),
		time.Date(2019, time.February, 19, 15, 12, 47, 123456780, time.UTC),
		time.Date(2019, time.February, 19, 15, 12, 47, 123456789, time.UTC),
	}
	expectedLogs := [7]string{
		"2019-02-19T15:12:47.123Z       [error] test string\n",
		"2019-02-19T15:12:47.1234Z      [error] test string\n",
		"2019-02-19T15:12:47.12345Z     [error] test string\n",
		"2019-02-19T15:12:47.123456Z    [error] test string\n",
		"2019-02-19T15:12:47.1234567Z   [error] test string\n",
		"2019-02-19T15:12:47.12345678Z  [error] test string\n",
		"2019-02-19T15:12:47.123456789Z [error] test string\n",
	}

	for i, ts := range tsValues {
		expected := expectedLogs[i]
		buffer := bytes.NewBuffer(nil)
		writer := NewTextWriter(buffer)
		writer.showTimestamp = true
		writer.useColor = false
		writer.WriteError(Messagef(Error, "test %s", "string").WithTimestamp(ts))
		assert.Equal(string(buffer.Bytes()), expected)
	}
}
