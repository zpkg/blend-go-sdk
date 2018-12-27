package cron

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
func ParseString(cronString string) *StringSchedule {
	return &StringSchedule{}
}

// StringSchedule is a schedule generated from a cron string.
type StringSchedule struct {
	Original string
}
