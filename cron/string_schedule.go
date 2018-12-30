package cron

import (
	"strconv"
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

	seconds, err := parsePart(parts[0], parseInt, below(60))
	if err != nil {
		return nil, err
	}

	minutes, err := parsePart(parts[1], parseInt, below(60))
	if err != nil {
		return nil, err
	}

	hours, err := parsePart(parts[2], parseInt, below(24))
	if err != nil {
		return nil, err
	}

	days, err := parsePart(parts[3], parseInt, between(1, 32))
	if err != nil {
		return nil, err
	}

	months, err := parsePart(parts[4], parseMonth, between(1, 13))
	if err != nil {
		return nil, err
	}

	daysOfWeek, err := parsePart(parts[5], parseDayOfWeek, between(0, 7))
	if err != nil {
		return nil, err
	}

	years, err := parsePart(parts[6], parseInt, nil)
	if err != nil {
		return nil, err
	}

	schedule := &StringSchedule{
		Original:    cronString,
		Seconds:     seconds,
		Minutes:     minutes,
		Hours:       hours,
		DaysOfMonth: days,
		Months:      months,
		DaysOfWeek:  daysOfWeek,
		Years:       years,
	}
	return schedule, nil
}

// Error Constants
const (
	ErrStringScheduleInvalid         exception.Class = "cron: schedule string invalid"
	ErrStringScheduleComponents      exception.Class = "cron: must have (7) components space delimited"
	ErrStringScheduleValueOutOfRange exception.Class = "cron: string schedule part out of range"
)

// Interface assertions.
var (
	_ Schedule = (*StringSchedule)(nil)
)

// StringSchedule is a schedule generated from a cron string.
type StringSchedule struct {
	Original string

	Seconds     []int
	Minutes     []int
	Hours       []int
	DaysOfMonth []int
	Months      []int
	DaysOfWeek  []int
	Years       []int
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

var (
	validMonths = map[string]int{
		"JAN": 1,
		"FEB": 2,
		"MAR": 3,
		"APR": 4,
		"MAY": 5,
		"JUN": 6,
		"JUL": 7,
		"AUG": 8,
		"SEP": 9,
		"OCT": 10,
		"NOV": 11,
		"DEC": 12,
	}

	validDaysOfWeek = map[string]int{
		"SUN": 0,
		"MON": 1,
		"TUE": 2,
		"WED": 3,
		"THU": 4,
		"FRI": 5,
		"SAT": 6,
	}
)

func parsePart(values string, parser func(string) (int, error), validator func(int) bool) ([]int, error) {
	if values == "*" {
		return nil, nil
	}

	// check if we need to expand an "every" pattern
	if strings.HasPrefix(values, "*/") {
		return parseEvery(values, validator)
	}

	components := strings.Split(values, string(cronSpecialComma))

	output := make([]int, len(components))
	for x := 0; x < len(components); x++ {
		part, err := strconv.Atoi(components[x])
		if err != nil {
			return nil, exception.New(err)
		}
		if validator != nil && !validator(part) {
			return nil, exception.New(err)
		}
		output[x] = part
	}
	return output, nil
}

func parseEvery(values string, validator func(int) bool) ([]int, error) {
	every, err := strconv.Atoi(strings.TrimPrefix(values, "*/"))
	if err != nil {
		return nil, exception.New(err)
	}
	if validator != nil && !validator(every) {
		return nil, exception.New(ErrStringScheduleValueOutOfRange)
	}

	var output []int
	for x := 0; x < 60; x = x + every {
		output = append(output, x)
	}
	return output, nil
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func parseMonth(s string) (int, error) {
	if value, ok := validMonths[s]; ok {
		return value, nil
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, exception.New(err).WithMessage("month not a valid integer")
	}
	if value < 1 || value > 12 {
		return 0, exception.New(ErrStringScheduleValueOutOfRange).WithMessagef("month out of range (1-12): %s", s)
	}
	return value, nil
}

func parseDayOfWeek(s string) (int, error) {
	if value, ok := validDaysOfWeek[s]; ok {
		return value, nil
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, exception.New(err).WithMessage("day of week not a valid integer")
	}
	if value < 0 || value > 6 {
		return 0, exception.New(ErrStringScheduleValueOutOfRange).WithMessagef("day of week out of range (0-6): %s", s)
	}
	return value, nil
}

func below(max int) func(int) bool {
	return between(0, max)
}

func between(min, max int) func(int) bool {
	return func(value int) bool {
		return value >= min && value < max
	}
}
