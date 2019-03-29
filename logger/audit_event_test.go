package logger

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestAuditEventMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	ae := NewAuditEvent("bailey", "pooped", OptAuditEventMetaOptions(OptEventMetaTimestamp(time.Date(2016, 01, 02, 03, 04, 05, 06, time.UTC))))

	contents, err := json.Marshal(ae)
	assert.Nil(err)

	assert.Contains(string(contents), "bailey")
	assert.Contains(string(contents), "pooped")

	assert.True(strings.HasPrefix(string(contents), `{"_timestamp":"2016-01-02T03:04:05`), string(contents))
}
