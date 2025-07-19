package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	ms := float64(d.Microseconds()) / 1000
	switch {
	case ms < 1:
		return "< 1 ms"
	case ms < 1000:
		return formatMs(ms)
	default:
		return d.Round(time.Millisecond).String()
	}
}

func formatMs(ms float64) string {
	return fmt.Sprintf("%.2f ms", ms)
}
