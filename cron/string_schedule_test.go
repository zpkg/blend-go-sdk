package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/stringutil"
)

type stringScheduleTestCase struct {
	Input       string
	ExpectedErr error
	Expected    *time.Time
	After       *time.Time
}

func TestParseString(t *testing.T) {
	assert := assert.New(t)

	testCases := []stringScheduleTestCase{
		{Input: "", ExpectedErr: ErrStringScheduleInvalid},
		{Input: stringutil.Random(stringutil.Letters, 10), ExpectedErr: ErrStringScheduleInvalid},
		{Input: "*/5 * * * * * *", After: Ref(time.Date(2018, 12, 29, 13, 12, 11, 0, time.UTC)), Expected: Ref(time.Date(2018, 12, 29, 13, 12, 16, 0, time.UTC))},
	}

	for _, tc := range testCases {
		parsed, err := ParseString(tc.Input)
		if tc.ExpectedErr != nil {
			assert.NotNil(err)
			assert.True(exception.Is(err, tc.ExpectedErr))
		} else {
			assert.Nil(err)
			assert.Equal(*tc.Expected, Deref(parsed.Next(tc.After)))
		}
	}
}

func TestFindNext(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(4, findNext(4, nil))
	assert.Equal(5, findNext(4, []int{1, 3, 5, 7, 9}))
	assert.Equal(10, findNext(10, []int{1, 2, 3, 4, 5}))
}
