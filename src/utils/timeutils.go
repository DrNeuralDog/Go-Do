package utils

import (
	"fmt"
	"strconv"
	"time"
)

// CustomDate keeps old C++-style date fields just in case
type CustomDate struct {
	Year    int   `json:"year"`
	Month   int   `json:"month"`
	Day     int   `json:"day"`
	Hour    int   `json:"hour"`
	Minute  int   `json:"minute"`
	Second  int   `json:"second"`
	IntTime int64 `json:"intTime"` // Unix timestamp
	Week    int   `json:"week"`    // Sunday == 0
	Updated bool  `json:"-"`       // cached state
}

// NewCustomDate returns the current date
func NewCustomDate() *CustomDate {
	now := time.Now()

	return NewCustomDateFromTime(now)
}

// NewCustomDateFromValues builds a date from raw fields
func NewCustomDateFromValues(year, month, day, hour, minute, second int) *CustomDate {
	cd := &CustomDate{
		Year:   year,
		Month:  month,
		Day:    day,
		Hour:   hour,
		Minute: minute,
		Second: second,
	}

	cd.Update()

	return cd
}

// NewCustomDateFromTime converts time.Time to CustomDate
func NewCustomDateFromTime(t time.Time) *CustomDate {
	cd := &CustomDate{
		Year:   t.Year(),
		Month:  int(t.Month()),
		Day:    t.Day(),
		Hour:   t.Hour(),
		Minute: t.Minute(),
		Second: t.Second(),
	}

	cd.Update()

	return cd
}

// Update refreshes cached timestamp and weekday fields
func (cd *CustomDate) Update() {
	if cd.Updated {
		return
	}

	date := cd.ToTime()
	cd.IntTime = date.Unix()
	cd.Week = int(date.Weekday())

	cd.Updated = true
}

// ToTime converts CustomDate to time.Time
func (cd *CustomDate) ToTime() time.Time {
	return time.Date(cd.Year, time.Month(cd.Month), cd.Day, cd.Hour, cd.Minute, cd.Second, 0, time.UTC)
}

// GetTime returns cached Unix timestamp
func (cd *CustomDate) GetTime() int64 {
	cd.Update()

	return cd.IntTime
}

// GetWeek returns weekday with Sunday as zero
func (cd *CustomDate) GetWeek() int {
	cd.Update()

	return cd.Week
}

// ToLastMonth moves date to previous month
func (cd *CustomDate) ToLastMonth() {
	cd.Month--

	if cd.Month <= 0 {
		cd.Month = 12
		cd.Year--
	}

	cd.Updated = false
	cd.Update()
}

// ToNextMonth moves date to next month
func (cd *CustomDate) ToNextMonth() {
	cd.Month++

	if cd.Month > 12 {
		cd.Month = 1
		cd.Year++
	}

	cd.Updated = false
	cd.Update()
}

// IsBefore reports whether this date is earlier than another
func (cd *CustomDate) IsBefore(other *CustomDate) bool {
	return cd.ToTime().Before(other.ToTime())
}

// IsLeapYear reports whether year is leap
func IsLeapYear(year int) bool {
	return year%400 == 0 || (year%100 != 0 && year%4 == 0)
}

// DaysInYear returns days count for year
func DaysInYear(year int) int {
	if IsLeapYear(year) {
		return 366
	}

	return 365
}

// DaysInMonth returns days count for month
func DaysInMonth(month int, year int) int {
	switch month {
	case 4, 6, 9, 11:
		return 30
	case 2:
		if IsLeapYear(year) {
			return 29
		}

		return 28
	default:
		return 31
	}
}

// GetCurrentDate returns current date as CustomDate
func GetCurrentDate() *CustomDate {
	return NewCustomDate()
}

// GetCurrentDateWithOffset returns current date shifted by hours
func GetCurrentDateWithOffset(offsetHours int) *CustomDate {
	now := time.Now().Add(time.Duration(offsetHours) * time.Hour)

	return NewCustomDateFromTime(now)
}

// FormatDateKey builds YYYYMM storage key
func FormatDateKey(year, month int) string {
	return fmt.Sprintf("%04d%02d", year, month)
}

// ParseDateKey parses YYYYMM storage key
func ParseDateKey(dateKey string) (year, month int) {
	if len(dateKey) != 6 {
		return 0, 0
	}

	year, err := strconv.Atoi(dateKey[:4])
	if err != nil {
		return 0, 0
	}

	month, err = strconv.Atoi(dateKey[4:])
	if err != nil || month < 1 || month > 12 {
		return 0, 0
	}

	return year, month
}
