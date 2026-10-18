package api

import (
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

type fakePasswordHasher struct{}

func (fakePasswordHasher) Hash(plainPassword string) (string, error) {
	return "hashed:" + plainPassword, nil
}

func (fakePasswordHasher) Verify(hash, plainPassword string) bool {
	return hash == "hashed:"+plainPassword
}

var testClock = fixedClock{now: time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)}

type fakeUserRepository struct {
	findByUsername func(username string) (domain.User, error)
	add            func(user domain.User) error
}

func (f *fakeUserRepository) FindUserByUsername(username string) (domain.User, error) {
	if f.findByUsername != nil {
		return f.findByUsername(username)
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (f *fakeUserRepository) Add(user domain.User) error {
	if f.add != nil {
		return f.add(user)
	}
	return nil
}
