package web

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLocalSessionCache(t *testing.T) {
	assert := assert.New(t)

	lsc := NewLocalSessionCache()

	session := &Session{UserID: "bailey", SessionID: NewSessionID()}
	assert.Nil(lsc.PersistHandler(nil, session, nil))

	fetched, err := lsc.FetchHandler(nil, session.SessionID, nil)
	assert.Nil(err)
	assert.Equal(session.UserID, fetched.UserID)

	assert.Nil(lsc.RemoveHandler(nil, session.SessionID, nil))

	removed, err := lsc.FetchHandler(nil, session.SessionID, nil)
	assert.Nil(err)
	assert.Nil(removed)
}
