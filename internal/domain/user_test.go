package domain

import (
	"testing"
)

func TestInvalidPassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError string
	}{
		{
			name:      "empty password",
			password:  "",
			wantError: "password must be at least 8 characters long",
		},
		{
			name:      "password less than 8 characters",
			password:  "1234567",
			wantError: "password must be at least 8 characters long",
		},
		{
			name:      "at least one special character",
			password:  "12345678",
			wantError: "password must contain at least one special character",
		},
		{
			name:      "password without a number",
			password:  "!passwordwithoutnumber",
			wantError: "password must contain at least one number",
		},
		{
			name:      "password without uppercase characters",
			password:  "!1abcdefghi",
			wantError: "password must contain at least one uppercase letter and one lowercase letter",
		},
		{
			name:      "password without lowercase characters",
			password:  "!1ABCDEFGHI",
			wantError: "password must contain at least one uppercase letter and one lowercase letter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUser("username", "username@email.com", tt.password)
			if err == nil {
				t.Fatalf("wanted an error but didn't get one")
			}

			if err.Error() != tt.wantError {
				t.Errorf("got %q, want %q", err, tt.wantError)
			}
		})
	}
}

func TestValidPassword(t *testing.T) {

	t.Run("check password is valid", func(t *testing.T) {
		user := User{
			Username: "username",
			Email:    "username@email.com",
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
