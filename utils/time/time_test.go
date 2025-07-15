package time

import (
	"testing"
	"time"
)

func TestNextDayBeginTime(t *testing.T) {
	t.Log(NextDayBeginTime(time.Now()).Format(time.UnixDate))
}

func TestNextWeekBeginTime(t *testing.T) {
	t.Log(NextWeekBeginTime(time.Now()).Format(time.UnixDate))
}

func TestNextMonthBeginTime(t *testing.T) {
	t.Log(NextMonthBeginTime(time.Now()).Format(time.UnixDate))
}

func TestNextYearBeginTime(t *testing.T) {
	t.Log(NextYearBeginTime(time.Now()).Format(time.UnixDate))
}
