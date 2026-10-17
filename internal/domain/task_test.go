package domain

import (
	"testing"
)

func TestNewTask(t *testing.T) {

	t.Run("description is valid", func(t *testing.T) {
		task, err := NewTask("1", "some description", fixedClock{now: fixedNow})

		if err != nil {
			t.Fatalf("unexpected error creating task: %v", err)
		}

		got := task.Description
		want := "some description"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

		if !task.CreatedAt.Equal(fixedNow) {
			t.Errorf("got CreatedAt %v, want %v", task.CreatedAt, fixedNow)
		}
	})

	t.Run("description is invalid", func(t *testing.T) {
		_, err := NewTask("1", "", fixedClock{now: fixedNow})

		got := err
		want := ErrDescriptionInvalid

		if got == nil {
			t.Fatal("didn't get an error but wanted one")
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("description is too long", func(t *testing.T) {
		_, err := NewTask("1", "a"+string(make([]byte, 1025)), fixedClock{now: fixedNow})

		got := err
		want := ErrDescriptionTooLong

		if got == nil {
			t.Fatal("didn't get an error but wanted one")
		}

		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("userId is invalid", func(t *testing.T) {
		_, err := NewTask("", "some description", fixedClock{now: fixedNow})

		got := err
		want := ErrUserIdInvalid

		if got == nil {
			t.Fatal("didn't get an error but wanted one")
		}

		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

}
