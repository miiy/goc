package time

import (
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFormatDateTime(t *testing.T) {
	loc, err := LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		in   any
		want string
	}{
		{
			name: "rfc3339 string",
			in:   "2025-12-31T16:00:00Z",
			want: "2026-01-01 00:00:00",
		},
		{
			name: "formatted string",
			in:   "2026-01-01 00:00:00",
			want: "2026-01-01 00:00:00",
		},
		{
			name: "time",
			in:   time.Date(2025, 12, 31, 16, 0, 0, 0, time.UTC),
			want: "2026-01-01 00:00:00",
		},
		{
			name: "protobuf timestamp",
			in:   timestamppb.New(time.Date(2025, 12, 31, 16, 0, 0, 0, time.UTC)),
			want: "2026-01-01 00:00:00",
		},
		{
			name: "unknown string",
			in:   "not a time",
			want: "not a time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDateTime(tt.in, loc); got != tt.want {
				t.Fatalf("FormatDateTime() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatTimeWithLayout(t *testing.T) {
	loc, err := LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}

	got := FormatTime("2025-12-31T16:00:00Z", loc, DateMinuteLayout)
	if got != "2026-01-01 00:00" {
		t.Fatalf("FormatTime() = %q, want %q", got, "2026-01-01 00:00")
	}
}
