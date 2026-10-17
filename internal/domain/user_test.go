package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		user, err := NewUser("john", "john@example.com", "_1ValidPassword", fakeHasher{}, fixedClock{now: fixedNow})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Username != "john" {
			t.Errorf("got %q, want %q", user.Username, "john")
		}
		if user.Email != "john@example.com" {
			t.Errorf("got %q, want %q", user.Email, "john@example.com")
		}
		if !user.CreatedAt.Equal(fixedNow) {
			t.Errorf("got CreatedAt %v, want %v", user.CreatedAt, fixedNow)
		}
		if got, want := user.Password.Hash(), "hashed:_1ValidPassword"; got != want {
			t.Errorf("got hash %q, want %q", got, want)
		}
	})

	t.Run("hasher failure", func(t *testing.T) {
		_, err := NewUser("john", "john@example.com", "_1ValidPassword", fakeHasher{hashErr: errors.New("no entropy")}, fixedClock{now: fixedNow})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty username", func(t *testing.T) {
		_, err := NewUser("", "john@example.com", "_1ValidPassword", fakeHasher{}, fixedClock{now: fixedNow})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("username too long", func(t *testing.T) {
		_, err := NewUser(strings.Repeat("a", 257), "john@example.com", "_1ValidPassword", fakeHasher{}, fixedClock{now: fixedNow})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty email", func(t *testing.T) {
		_, err := NewUser("john", "", "_1ValidPassword", fakeHasher{}, fixedClock{now: fixedNow})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("email too long", func(t *testing.T) {
		_, err := NewUser("john", strings.Repeat("a", 257)+"@example.com", "_1ValidPassword", fakeHasher{}, fixedClock{now: fixedNow})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("weak password", func(t *testing.T) {
		_, err := NewUser("john", "john@example.com", "weak", fakeHasher{}, fixedClock{now: fixedNow})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUserValidate(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		u := User{Username: "john", Email: "john@example.com"}
		if err := u.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty username", func(t *testing.T) {
		u := User{Username: "", Email: "john@example.com"}
		if err := u.Validate(); err == nil {
			t.Error("expected error for empty username")
		}
	})

	t.Run("username too long", func(t *testing.T) {
		u := User{Username: strings.Repeat("a", 257), Email: "john@example.com"}
		if err := u.Validate(); err == nil {
			t.Error("expected error for long username")
		}
	})

	t.Run("empty email", func(t *testing.T) {
		u := User{Username: "john", Email: ""}
		if err := u.Validate(); err == nil {
			t.Error("expected error for empty email")
		}
	})

	t.Run("email too long", func(t *testing.T) {
		u := User{Username: "john", Email: strings.Repeat("a", 257)}
		if err := u.Validate(); err == nil {
			t.Error("expected error for long email")
		}
	})
}

func TestNewUserCreatedAt(t *testing.T) {
	user, err := NewUser("john", "john@example.com", "_1ValidPassword", fakeHasher{}, fixedClock{now: fixedNow})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !user.CreatedAt.Equal(fixedNow) {
		t.Errorf("got CreatedAt %v, want %v", user.CreatedAt, fixedNow)
	}
}
