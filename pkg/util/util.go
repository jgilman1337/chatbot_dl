package util

import (
	"time"
)

// Returns a timestamp according to the format: `yyyyMMdd_HHmmss`, eg: `20250723_133723`.
func Timestamp() string {
	layout := "20060102_150405"
	return time.Now().Format(layout)
}
