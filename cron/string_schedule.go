package cron

import (
	"strings"
	"time"

	"github.com/blend/go-sdk/exception"
)

// ParseString parses a cron formatted string into a schedule.
/*
Field name     Mandatory?   Allowed values    Allowed special characters
----------     ----------   --------------    --------------------------
Seconds        No           0-59              * / , -
Minutes        Yes          0-59              * / , -
Hours          Yes          0-23              * / , -
Day of month   Yes          1-31              * / , - L W
Month          Yes          1-12 or JAN-DEC   * / , -
Day of week    Yes          0-6 or SUN-SAT    * / , - L #
Year           No           1970–2099         * / , -
*/
func ParseString(cronString string) (*StringSchedule, error) {
	parts := strings.Split(cronString, " ")
	if len(parts) != 7 {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(ErrStringScheduleComponents).WithMessagef("provided string; %s", cronString)
	}

	schedule := &StringSchedule{
		Original:    cronString,
		Seconds:     strings.Split(parts[0], string(cronSpecialComma)),
		Minutes:     strings.Split(parts[1], string(cronSpecialComma)),
		Hours:       strings.Split(parts[2], string(cronSpecialComma)),
		DaysOfMonth: strings.Split(parts[3], string(cronSpecialComma)),
		Months:      strings.Split(parts[4], string(cronSpecialComma)),
		DaysOfWeek:  strings.Split(parts[5], string(cronSpecialComma)),
		Years:       strings.Split(parts[6], string(cronSpecialComma)),
	}
	return schedule, nil
}

// Error Constants
const (
	ErrStringScheduleInvalid    exception.Class = "cron: schedule string invalid"
	ErrStringScheduleComponents exception.Class = "cron: must have (7) components space delimited"
)

// Interface assertions.
var (
	_ Schedule = (*StringSchedule)(nil)
)

// StringSchedule is a schedule generated from a cron string.
type StringSchedule struct {
	Original string

	Seconds     []string
	Minutes     []string
	Hours       []string
	DaysOfMonth []string
	Months      []string
	DaysOfWeek  []string
	Years       []string
}

// Next implements cron.Schedule.
func (ss *StringSchedule) Next(after *time.Time) *time.Time {
	return nil
}

// these are special characters
const (
	cronSpecialComma    = ',' //
	cronSpecialDash     = '-'
	cronSpecialStar     = '*'
	cronSpecialSlash    = '/'
	cronSpecialQuestion = '?' // sometimes used as the startup time, sometimes as a *

	cronSpecialLast       = 'L'
	cronSpecialWeekday    = 'W' // nearest weekday to the given day of the month
	cronSpecialDayOfMonth = '#' //
)
