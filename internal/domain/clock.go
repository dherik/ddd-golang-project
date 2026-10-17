package domain

import "time"

// Clock supplies the current time to domain constructors. Implementations may
// return any location; constructors normalize the result to UTC.
type Clock interface {
	Now() time.Time
}
