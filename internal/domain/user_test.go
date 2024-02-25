package domain

import (
	"testing"
)

func TestPassword(t *testing.T) {

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
