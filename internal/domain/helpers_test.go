package domain

import "time"

var fixedNow = time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

type fakeHasher struct {
	hashErr error
}

func (h fakeHasher) Hash(plainPassword string) (string, error) {
	if h.hashErr != nil {
		return "", h.hashErr
	}
	return "hashed:" + plainPassword, nil
}

func (h fakeHasher) Verify(hash, plainPassword string) bool {
	return hash == "hashed:"+plainPassword
}
