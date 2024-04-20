package domain

import (
	"testing"
)

func TestPassword(t *testing.T) {

	t.Run("check password length is not empty", func(t *testing.T) {
		_, err := NewUser("username", "email@email.com", "")
		if err == nil {
			t.Fatalf("wanted an error but didn't get one")
		}

		want := "password must be at least 8 characters long"
		if err.Error() != want {
			t.Errorf("got %q, want %q", err, want)
		}
	})

	t.Run("check password length is bigger than 8 characters", func(t *testing.T) {
		_, err := NewUser("username", "email@email.com", "1234567")
		if err == nil {
			t.Fatalf("wanted an error but didn't get one")
		}

		want := "password must be at least 8 characters long"
		if err.Error() != want {
			t.Errorf("got %q, want %q", err, want)
		}
	})

	t.Run("check password is valid", func(t *testing.T) {
		user := User{
			Username: "username",
			Email:    "email@email.com",
		}
		user.SetPassword("some_password")

		got := user.CheckPasswordHash("some_password")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("check password is not valid", func(t *testing.T) {
		user := User{
			Username: "username",
			Email:    "email@email.com",
		}
		user.SetPassword("some_password")

		got := user.CheckPasswordHash("another_password")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

}
