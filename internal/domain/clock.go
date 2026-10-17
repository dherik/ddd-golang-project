package domain

import "time"

// Clock supplies the current time to domain constructors. Injecting it keeps
// entity creation deterministic in tests.
type Clock interface {
	Now() time.Time
}
