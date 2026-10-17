package security

import "time"

// SystemClock is the production clock adapter.
type SystemClock struct{}

func NewSystemClock() *SystemClock {
	return &SystemClock{}
}

func (c *SystemClock) Now() time.Time {
	return time.Now().UTC()
}
