package utils

import (
	"fmt"
	"time"
)

// CustomDate represents a date with time components, matching the original C++ Date class
type CustomDate struct {
	Year    int   `json:"year"`
	Month   int   `json:"month"`
	Day     int   `json:"day"`
	Hour    int   `json:"hour"`
	Minute  int   `json:"minute"`
	Second  int   `json:"second"`
	IntTime int64 `json:"intTime"` // Unix timestamp equivalent
	Week    int   `json:"week"`    // Day of week (0=Sunday)
	Updated bool  `json:"-"`       // Internal update flag
}

// NewCustomDate creates a new CustomDate with current time
func NewCustomDate() *CustomDate {
	now := time.Now()
	return NewCustomDateFromTime(now)
}

// NewCustomDateFromValues creates a new CustomDate from individual components
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

// NewCustomDateFromTime creates a new CustomDate from a time.Time
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

// Update calculates the internal time representation and week day
func (cd *CustomDate) Update() {
	if cd.Updated {
		return
	}

	// Calculate Unix timestamp equivalent (seconds since 1970-01-01)
	cd.IntTime = cd.ToTime().Unix()

	// Calculate day of week (0 = Sunday)
	refDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	cd.Week = int(cd.ToTime().Sub(refDate).Hours()/24) % 7

	cd.Updated = true
}

// ToTime converts CustomDate to time.Time
func (cd *CustomDate) ToTime() time.Time {
	return time.Date(cd.Year, time.Month(cd.Month), cd.Day, cd.Hour, cd.Minute, cd.Second, 0, time.UTC)
}

// GetTime returns the Unix timestamp
func (cd *CustomDate) GetTime() int64 {
	cd.Update()
	return cd.IntTime
}

// GetWeek returns the day of week (0 = Sunday)
func (cd *CustomDate) GetWeek() int {
	cd.Update()
	return cd.Week
}

// ToLastMonth moves the date to the previous month
func (cd *CustomDate) ToLastMonth() {
	cd.Month--
	if cd.Month <= 0 {
		cd.Month = 12
		cd.Year--
	}
	cd.Updated = false
	cd.Update()
}

// ToNextMonth moves the date to the next month
func (cd *CustomDate) ToNextMonth() {
	cd.Month++
	if cd.Month > 12 {
		cd.Month = 1
		cd.Year++
	}
	cd.Updated = false
	cd.Update()
}

// IsBefore returns true if this date is before the other date
func (cd *CustomDate) IsBefore(other *CustomDate) bool {
	return cd.ToTime().Before(other.ToTime())
}

// IsLeapYear checks if the year is a leap year
func IsLeapYear(year int) bool {
	return year%400 == 0 || (year%100 != 0 && year%4 == 0)
}

// DaysInYear returns the number of days in the given year
func DaysInYear(year int) int {
	if IsLeapYear(year) {
		return 366
	}
	return 365
}

// DaysInMonth returns the number of days in the given month
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

// GetCurrentDate returns the current date as CustomDate
func GetCurrentDate() *CustomDate {
	return NewCustomDate()
}

// GetCurrentDateWithOffset returns the current date with hour offset
func GetCurrentDateWithOffset(offsetHours int) *CustomDate {
	now := time.Now().Add(time.Duration(offsetHours) * time.Hour)
	return NewCustomDateFromTime(now)
}

// FormatDateKey creates a date key for file storage (YYYYMM format)
func FormatDateKey(year, month int) string {
	return fmt.Sprintf("%04d%02d", year, month)
}

// ParseDateKey parses a date key back to year and month
func ParseDateKey(dateKey string) (year, month int) {
	if len(dateKey) != 6 {
		return 0, 0
	}
	fmt.Sscanf(dateKey, "%04d%02d", &year, &month)
	return year, month
}
