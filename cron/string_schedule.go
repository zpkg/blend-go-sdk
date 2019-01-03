package cron

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/stringutil"
)

// ParseString parses a cron formatted string into a schedule.
// The string must be 7 components, whitespace separated.
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
	parts := stringutil.SplitSpace(cronString)
	if len(parts) != 7 {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(ErrStringScheduleComponents).WithMessagef("provided string; %s", cronString)
	}

	seconds, err := parsePart(parts[0], parseInt, below(60))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("seconds invalid")
	}

	minutes, err := parsePart(parts[1], parseInt, below(60))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("minutes invalid")
	}

	hours, err := parsePart(parts[2], parseInt, below(24))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("hours invalid")
	}

	days, err := parsePart(parts[3], parseInt, between(1, 32))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("days invalid")
	}

	months, err := parsePart(parts[4], parseMonth, between(1, 13))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("months invalid")
	}

	daysOfWeek, err := parsePart(parts[5], parseDayOfWeek, between(0, 7))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("days of week invalid")
	}

	years, err := parsePart(parts[6], parseInt, between(1970, 2100))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("years invalid")
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
	ErrStringScheduleInvalidRange    exception.Class = "cron: range (from-to) invalid"
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
	DaysOfWeek  []int
	DaysOfMonth []int
	Months      []int
	Years       []int
}

// Next implements cron.Schedule.
func (ss *StringSchedule) Next(after time.Time) time.Time {
	working := Now()
	if !after.IsZero() {
		working = after
	}

	// figure out the next year
	year := findNext(working.Year(), ss.Years)
	month := findNext(int(working.Month()), ss.Months)

	day := findNext(int(working.Day()), ss.DaysOfMonth)
	year, month, day = findNextDayOfWeek(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), ss.DaysOfWeek)
	hour := findNext(working.Hour(), ss.Hours)
	minute := findNext(working.Minute(), ss.Minutes)
	second := findNext(working.Second(), ss.Seconds)
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
}

func parsePart(values string, parser func(string) (int, error), validator func(int) bool) ([]int, error) {
	if values == string(cronSpecialStar) {
		return nil, nil
	}

	// check if we need to expand an "every" pattern
	if strings.HasPrefix(values, cronSpecialEvery) {
		return parseEvery(values, parseInt, validator)
	}

	components := strings.Split(values, string(cronSpecialComma))

	output := map[int]bool{}
	var component string
	for x := 0; x < len(components); x++ {
		component = components[x]
		if strings.Contains(component, string(cronSpecialDash)) {
			rangeValues, err := parseRange(values, parser, validator)
			if err != nil {
				return nil, err
			}

			for _, value := range rangeValues {
				output[value] = true
			}
		}

		part, err := parser(component)
		if err != nil {
			return nil, exception.New(err)
		}
		if validator != nil && !validator(part) {
			return nil, exception.New(err)
		}
		output[part] = true
	}
	return mapKeysToArray(output), nil
}

func parseEvery(values string, parser func(string) (int, error), validator func(int) bool) ([]int, error) {
	every, err := parser(strings.TrimPrefix(values, "*/"))
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

func parseRange(values string, parser func(string) (int, error), validator func(int) bool) ([]int, error) {
	parts := strings.Split(values, string(cronSpecialDash))

	if len(parts) != 2 {
		return nil, exception.New(ErrStringScheduleInvalidRange).WithMessagef("invalid range: %s")
	}

	from, err := parser(parts[0])
	if err != nil {
		return nil, exception.New(err)
	}
	to, err := parser(parts[1])
	if err != nil {
		return nil, exception.New(err)
	}

	if validator != nil && !validator(from) {
		return nil, exception.New(ErrStringScheduleValueOutOfRange).WithMessage("invalid range from")
	}
	if validator != nil && !validator(to) {
		return nil, exception.New(ErrStringScheduleValueOutOfRange).WithMessage("invalid range to")
	}

	if from >= to {
		return nil, exception.New(ErrStringScheduleInvalidRange).WithMessage("invalid range; from greater than to")
	}

	var output []int
	for x := from; x <= to; x++ {
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

func mapKeysToArray(values map[int]bool) []int {
	output := make([]int, len(values))
	var index int
	for key := range values {
		output[index] = key
		index++
	}
	sort.Ints(output)
	return output
}

func findNext(basis int, values []int) int {
	for _, value := range values {
		if value >= basis {
			return value
		}
	}
	return basis
}

func findNextDayOfWeek(t time.Time, daysOfWeek []int) (year, month, day int) {
	for index := 0; index < 7; index++ {
		trial := t.AddDate(0, 0, 1)
		for _, dow := range daysOfWeek {
			if dow == int(trial.Weekday()) {
				year = trial.Year()
				month = int(trial.Month())
				day = trial.Day()
				return
			}
		}
	}

	year = t.Year()
	month = int(t.Month())
	day = t.Day()
	return

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

	cronSpecialEvery = "*/"
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
