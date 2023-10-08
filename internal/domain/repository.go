package domain

// UserRepository defines the contract for managing user data.
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	Delete(userID int) error
	FindByID(userID int) (*User, error)
	FindByEmail(email string) (*User, error)
	// Add more methods as needed for user-related data operations
}

// TaskRepository defines the contract for managing task data.
type TaskRepository interface {
	Create(task *Task) error
	Update(task *Task) error
	Delete(taskID int) error
	FindByID(taskID int) (*Task, error)
	FindAllByUserID(userID int) ([]*Task, error)
	// Add more methods as needed for task-related data operations
}
