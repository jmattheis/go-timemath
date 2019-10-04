package timemath

import (
	"time"
)

// Unit the time unit
type Unit rune

var (
	// Year the unit.
	Year Unit = 'y'
	// Month the unit.
	Month Unit = 'M'
	// Week the unit.
	Week Unit = 'w'
	// Day the unit.
	Day Unit = 'd'
	// Hour the unit.
	Hour Unit = 'h'
	// Minute the unit.
	Minute Unit = 'm'
	// Second the unit.
	Second Unit = 's'
)

// EndOf returns the end of the given unit.
func (u Unit) EndOf(date time.Time, endOfWeek time.Weekday) time.Time {
	switch u {
	case Second:
		return date
	case Minute:
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), 59, 999999999, date.Location())
	case Hour:
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), 59, 59, 999999999, date.Location())
	case Day:
		return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())
	case Week:
		temp := date
		for temp.Weekday() != endOfWeek {
			temp = Day.Add(temp, 1)
		}
		return Day.EndOf(temp, endOfWeek)
	case Month:
		return time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 999999999, date.Location())
	case Year:
		return time.Date(date.Year(), time.December, 31, 23, 59, 59, 999999999, date.Location())
	default:
		panic("unknown unit type")
	}
}

// StartOf returns the start of the given unit.
func (u Unit) StartOf(date time.Time, startOfWeek time.Weekday) time.Time {
	switch u {
	case Second:
		return date
	case Minute:
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, date.Location())
	case Hour:
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), 0, 0, 0, date.Location())
	case Day:
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	case Week:
		temp := date
		for temp.Weekday() != startOfWeek {
			temp = Day.Subtract(temp, 1)
		}
		return Day.StartOf(temp, startOfWeek)
	case Month:
		return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	case Year:
		return time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, date.Location())
	default:
		panic("unknown unit type")
	}
}

// Subtract subtracts the unit * amount to the given date.
func (u Unit) Subtract(date time.Time, amount int) time.Time {
	return u.Add(date, -amount)
}

// Add adds the unit * amount to the given date.
func (u Unit) Add(date time.Time, amount int) time.Time {
	switch u {
	case Second:
		return date.Add(time.Duration(amount) * time.Second)
	case Minute:
		return date.Add(time.Duration(amount) * time.Minute)
	case Hour:
		return date.Add(time.Duration(amount) * time.Hour)
	case Day:
		return date.AddDate(0, 0, amount)
	case Week:
		return date.AddDate(0, 0, amount*7)
	case Month:
		return date.AddDate(0, amount, 0)
	case Year:
		return date.AddDate(amount, 0, 0)
	default:
		panic("unknown unit type")
	}
}
