package time

import (
	"testing"
	"time"
)

func TestBoundaryTimes(t *testing.T) {
	loc := time.FixedZone("UTC+8", 8*60*60)
	base := time.Date(2026, time.June, 11, 13, 20, 30, 123, loc)

	tests := []struct {
		name string
		got  time.Time
		want time.Time
	}{
		{
			name: "today begin",
			got:  TodayBeginTime(base),
			want: time.Date(2026, time.June, 11, 0, 0, 0, 0, loc),
		},
		{
			name: "today end",
			got:  TodayEndTime(base),
			want: time.Date(2026, time.June, 11, 23, 59, 59, int(time.Second-time.Nanosecond), loc),
		},
		{
			name: "next day begin",
			got:  NextDayBeginTime(base),
			want: time.Date(2026, time.June, 12, 0, 0, 0, 0, loc),
		},
		{
			name: "next week begin",
			got:  NextWeekBeginTime(base),
			want: time.Date(2026, time.June, 15, 0, 0, 0, 0, loc),
		},
		{
			name: "next month begin",
			got:  NextMonthBeginTime(base),
			want: time.Date(2026, time.July, 1, 0, 0, 0, 0, loc),
		},
		{
			name: "next year begin",
			got:  NextYearBeginTime(base),
			want: time.Date(2027, time.January, 1, 0, 0, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.got.Equal(tt.want) {
				t.Fatalf("got %s, want %s", tt.got, tt.want)
			}
		})
	}
}
