package domain

import "time"

// Clock supplies the current UTC time to domain constructors. Constructors
// normalize defensively, so an adapter returning a non-UTC location cannot
// leak a local timestamp into an entity.
type Clock interface {
	Now() time.Time
}
