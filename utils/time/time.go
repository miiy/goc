package time

import "time"

// NextDayBeginTime returns the next day's begin time.
func NextDayBeginTime(t time.Time) time.Time {
	nt := t.AddDate(0, 0, 1)
	return time.Date(nt.Year(), nt.Month(), nt.Day(), 0, 0, 0, 0, nt.Location())
}

// NextWeekBeginTime returns the next week's begin time.
func NextWeekBeginTime(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).
		AddDate(0, 0, offset).
		AddDate(0, 0, 7)
}

// NextMonthBeginTime returns the next month's begin time.
func NextMonthBeginTime(t time.Time) time.Time {
	nt := t.AddDate(0, 1, 0)
	return time.Date(nt.Year(), nt.Month(), 1, 0, 0, 0, 0, nt.Location())
}

// NextYearBeginTime returns the next year's begin time.
func NextYearBeginTime(t time.Time) time.Time {
	nt := t.AddDate(1, 0, 0)
	return time.Date(nt.Year(), 1, 1, 0, 0, 0, 0, nt.Location())
}
