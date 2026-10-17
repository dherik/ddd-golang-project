package task

import (
	"fmt"
	"log/slog"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

// Service is the application seam for task operations. Both transports —
// HTTP and messaging — depend on it; it speaks domain types only, so neither
// transport drags the other's shapes across the seam.
type Service struct {
	taskRepository domain.TaskRepository
}

func NewService(taskRepository domain.TaskRepository) *Service {
	return &Service{taskRepository: taskRepository}
}

func (s *Service) AddTaskToUser(task domain.Task) error {
	if err := s.taskRepository.AddTaskToUser(task.UserId, task); err != nil {
		return fmt.Errorf("failed to add task to user: %w", err)
	}
	return nil
}

func (s *Service) FindTasks(query domain.TaskQuery) (domain.TaskPage, error) {
	slog.Info("find task page",
		"userId", query.UserId,
		"startDate", query.StartDate,
		"endDate", query.EndDate,
		"pageSize", query.PageSize,
		"offset", query.Offset,
	)
	page, err := s.taskRepository.FindTasks(query)
	if err != nil {
		return domain.TaskPage{}, fmt.Errorf("failed finding tasks: %w", err)
	}
	return page, nil
}

func (s *Service) GetTaskByID(userId string, id int) (domain.Task, error) {
	task, err := s.taskRepository.GetByIDAndUserID(id, userId)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed getting task by id: %w", err)
	}
	slog.Info("found task", "id", id, "userId", userId)
	return task, nil
}
