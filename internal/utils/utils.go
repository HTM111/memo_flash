package utils

import (
	"crypto/md5"
	"fmt"
	"math"
	"strings"
	"time"
)

func UniqueId(input ...string) string {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(input, ":"))))
	return hash[:8]
}

func FormatDuration(date time.Time) string {
	since := time.Since(date)
	switch {
	case 60 > since.Seconds():
		return "Moments ago"
	case 60 > since.Minutes():
		return fmt.Sprintf("%.0f minutes ago", since.Round(time.Minute).Minutes())
	case 24 > since.Hours():
		return fmt.Sprintf("%.0f hours ago", since.Round(time.Hour).Hours())
	case 24*7 > since.Hours():
		return fmt.Sprintf("%.0f days ago", math.Round(since.Hours()/24))
	case 24*30 > since.Hours():
		return fmt.Sprintf("%.0f weeks ago", math.Round(since.Hours()/(24*7)))
	case 24*365 > since.Hours():
		return fmt.Sprintf("%.0f months ago", math.Round(since.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%.0f years ago", math.Round(since.Hours()/(24*365)))
	}
}
