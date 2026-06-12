package time

import (
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// DateTimeLayout is the common "YYYY-MM-DD HH:mm:ss" display layout.
	DateTimeLayout = "2006-01-02 15:04:05"
	// DateMinuteLayout is the common "YYYY-MM-DD HH:mm" display layout.
	DateMinuteLayout = "2006-01-02 15:04"
)

// LoadLocation loads a timezone by name and returns time.Local for an empty name.
func LoadLocation(name string) (*time.Location, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return time.Local, nil
	}
	return time.LoadLocation(name)
}

// FormatDateTime formats supported time-like values with DateTimeLayout.
func FormatDateTime(v any, loc *time.Location) string {
	return FormatTime(v, loc, DateTimeLayout)
}

// FormatTime formats supported time-like values with the given layout.
func FormatTime(v any, loc *time.Location, layout string) string {
	if strings.TrimSpace(layout) == "" {
		layout = DateTimeLayout
	}
	t, ok := ParseTimeInLocation(v, loc)
	if !ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
		return ""
	}
	if t.IsZero() {
		return ""
	}
	if loc != nil {
		t = t.In(loc)
	}
	return t.Format(layout)
}

// ParseTime parses supported time-like values using time.Local for zoneless strings.
func ParseTime(v any) (time.Time, bool) {
	return ParseTimeInLocation(v, nil)
}

// ParseTimeInLocation parses supported time-like values using loc for zoneless strings.
func ParseTimeInLocation(v any, loc *time.Location) (time.Time, bool) {
	if loc == nil {
		loc = time.Local
	}
	switch t := v.(type) {
	case time.Time:
		return t, true
	case *time.Time:
		if t == nil {
			return time.Time{}, false
		}
		return *t, true
	case *timestamppb.Timestamp:
		if t == nil || !t.IsValid() {
			return time.Time{}, false
		}
		return t.AsTime(), true
	case string:
		return parseTimeString(t, loc)
	default:
		return time.Time{}, false
	}
}

func parseTimeString(s string, loc *time.Location) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, true
	}
	for _, layout := range []string{DateTimeLayout, DateMinuteLayout} {
		if t, err := time.ParseInLocation(layout, s, loc); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
