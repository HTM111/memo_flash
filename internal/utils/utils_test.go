package utils_test

import (
	"memoflash/internal/utils"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	now := time.Now()

	tests := []struct {
		input    time.Time
		expected string
	}{
		{now.Add(-30 * time.Second), "Moments ago"},
		{now.Add(-59 * time.Second), "Moments ago"},

		{now.Add(-1 * time.Minute), "1 minutes ago"},
		{now.Add(-2 * time.Minute), "2 minutes ago"},
		{now.Add(-30 * time.Minute), "30 minutes ago"},
		{now.Add(-59 * time.Minute), "59 minutes ago"},

		{now.Add(-1 * time.Hour), "1 hours ago"},
		{now.Add(-2 * time.Hour), "2 hours ago"},
		{now.Add(-12 * time.Hour), "12 hours ago"},
		{now.Add(-23 * time.Hour), "23 hours ago"},

		{now.Add(-24 * time.Hour), "1 days ago"},
		{now.Add(-48 * time.Hour), "2 days ago"},
		{now.Add(-144 * time.Hour), "6 days ago"},

		{now.Add(-7 * 24 * time.Hour), "1 weeks ago"},
		{now.Add(-14 * 24 * time.Hour), "2 weeks ago"},
		{now.Add(-21 * 24 * time.Hour), "3 weeks ago"},
		{now.Add(-28 * 24 * time.Hour), "4 weeks ago"},

		{now.Add(-30 * 24 * time.Hour), "1 months ago"},
		{now.Add(-60 * 24 * time.Hour), "2 months ago"},
		{now.Add(-180 * 24 * time.Hour), "6 months ago"},
		{now.Add(-330 * 24 * time.Hour), "11 months ago"},

		{now.Add(-365 * 24 * time.Hour), "1 years ago"},
		{now.Add(-730 * 24 * time.Hour), "2 years ago"},
		{now.Add(-1095 * 24 * time.Hour), "3 years ago"},
	}

	for _, tt := range tests {
		if got := utils.FormatDuration(tt.input); got != tt.expected {
			t.Errorf("TimeFormat(%v) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}
