package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/stringutil"
)

type stringScheduleTestCase struct {
	Input       string
	ExpectedErr error
	Expected    time.Time
	After       time.Time
}

func TestParseString(t *testing.T) {
	assert := assert.New(t)

	testCases := []stringScheduleTestCase{
		{Input: "", ExpectedErr: ErrStringScheduleInvalid},
		{Input: stringutil.Random(stringutil.Letters, 10), ExpectedErr: ErrStringScheduleInvalid},
		{Input: "*/1 * * * * * *", After: time.Date(2018, 12, 29, 13, 12, 11, 10, time.UTC), Expected: time.Date(2018, 12, 29, 13, 12, 12, 0, time.UTC)},
		{Input: "*/5 * * * * * *", After: time.Date(2018, 12, 29, 13, 12, 11, 0, time.UTC), Expected: time.Date(2018, 12, 29, 13, 12, 15, 0, time.UTC)},
		{Input: "* 2 1 * * 1-6 *", After: time.Date(2019, 01, 01, 12, 0, 0, 0, time.UTC), Expected: time.Date(2019, 01, 02, 01, 02, 0, 0, time.UTC)},
		{Input: "* 2 1 * * MON-FRI *", After: time.Date(2019, 01, 01, 12, 0, 0, 0, time.UTC), Expected: time.Date(2019, 01, 02, 01, 02, 0, 0, time.UTC)},
		{Input: "* 9 10 * * SUN-TUE *", After: time.Date(2019, 01, 02, 12, 0, 0, 0, time.UTC), Expected: time.Date(2019, 01, 06, 10, 9, 0, 0, time.UTC)}, //
		{Input: "0 0 0 * * 0 *", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 06, 0, 0, 0, 0, time.UTC)},         // every week at midnight sat/sun
		{Input: "0 0 0 * * * *", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 03, 0, 0, 0, 0, time.UTC)},         // every day at midnight
		{Input: "0 0 * * * * *", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 02, 13, 0, 0, 0, time.UTC)},        // every hour on the hour
	}

	for _, tc := range testCases {
		parsed, err := ParseString(tc.Input)
		if tc.ExpectedErr != nil {
			assert.NotNil(err)
			assert.True(exception.Is(err, tc.ExpectedErr))
		} else {
			assert.Nil(err)
			next := parsed.Next(tc.After)
			assert.Equal(tc.Expected, next, fmt.Sprintf("%s vs. %s\n%v vs. %v", tc.Input, parsed.String(), tc.Expected.Format(time.RFC3339), next.Format(time.RFC3339)))
		}
	}
}

func TestStringScheduleEvery(t *testing.T) {
	assert := assert.New(t)

	schedule, err := ParseString("*/1 * * * * * *")
	assert.Nil(err)

	last := time.Date(2019, 01, 29, 0, 0, 0, 0, time.UTC)
	for x := 0; x < 60*60*6; x++ {
		last = schedule.Next(last)
	}
	assert.Equal(time.Date(2019, 01, 29, 6, 0, 0, 0, time.UTC), last)
}

func TestMapKeysToArray(t *testing.T) {
	assert := assert.New(t)

	assert.Equal([]int{1, 2, 3}, mapKeysToArray(map[int]bool{
		3: true,
		1: true,
		2: true,
	}))
	assert.Empty(mapKeysToArray(nil))
	assert.Empty(mapKeysToArray(map[int]bool{}))
}
