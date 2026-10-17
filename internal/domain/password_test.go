package domain

import (
	"errors"
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
			err := validatePassword(tt.password)

			if err == nil {
				t.Fatalf("wanted an error but didn't get one")
			}

			if err.Error() != tt.wantError {
				t.Errorf("got %q, want %q", err, tt.wantError)
			}
		})
	}
}

func TestPasswordHashingAndVerification(t *testing.T) {
	password, err := NewPassword("_1ValidPassword", fakeHasher{})
	if err != nil {
		t.Fatalf("unexpected error creating password: %v", err)
	}

	if got, want := password.Hash(), "hashed:_1ValidPassword"; got != want {
		t.Errorf("got hash %q, want %q", got, want)
	}

	if !password.Verify("_1ValidPassword", fakeHasher{}) {
		t.Error("the correct password did not verify")
	}

	if password.Verify("another_password", fakeHasher{}) {
		t.Error("a wrong password verified")
	}
}

func TestNewPasswordHasherError(t *testing.T) {
	_, err := NewPassword("_1ValidPassword", fakeHasher{hashErr: errors.New("no entropy")})
	if err == nil {
		t.Fatal("wanted an error but didn't get one")
	}
}

func TestRehydratedPasswordVerifies(t *testing.T) {
	password := RehydratePassword("hashed:_1ValidPassword")

	if !password.Verify("_1ValidPassword", fakeHasher{}) {
		t.Error("the rehydrated password did not verify")
	}

	if password.Verify("another_password", fakeHasher{}) {
		t.Error("a wrong password verified against a rehydrated password")
	}
}
