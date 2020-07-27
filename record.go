package typhoon

import (
	"time"
)

// Record respresent the information of HTTP connection
type Record struct {
	Err        error
	Duration   time.Duration
	StatusCode int
	Length     int64
}

// NewRecord returns a record with no-error
func NewRecord(d time.Duration, code int, length int64) *Record {
	return &Record{
		Err:        nil,
		Duration:   d,
		StatusCode: code,
		Length:     length,
	}
}

// NewErrorRecord returns a record with error
func NewErrorRecord(err error) *Record {
	return &Record{
		Err: err,
	}
}
