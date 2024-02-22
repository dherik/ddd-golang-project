package domain

import (
	"testing"
)

func TestPassword(t *testing.T) {
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
}
