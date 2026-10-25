package task

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

// Service is the application seam for task operations. Both transports —
// HTTP and messaging — depend on it; it speaks domain types only, so neither
// transport drags the other's shapes across the seam.
type Service struct {
	taskRepository domain.TaskRepository
	clock          domain.Clock
	publisher      Publisher
}

func NewService(taskRepository domain.TaskRepository, clock domain.Clock, publisher Publisher) *Service {
	return &Service{taskRepository: taskRepository, clock: clock, publisher: publisher}
}

func (s *Service) CreateTask(userId, description string) error {
	task, err := domain.NewTask(userId, description, s.clock)
	if err != nil {
		return err
	}
	saved, err := s.taskRepository.AddTaskToUser(userId, task)
	if err != nil {
		return fmt.Errorf("failed to add task to user: %w", err)
	}
	s.publishTaskCreated(saved)
	return nil
}

// publishTaskCreated announces the saved task as one task created fact.
// Publication is best-effort: a failure is logged and swallowed,
// so the operation's result reflects the database write alone — on the
// messaging path a failure here must not dead-letter the triggering event.
func (s *Service) publishTaskCreated(saved domain.Task) {
	userId, err := strconv.Atoi(saved.UserId)
	if err != nil {
		slog.Error("failed converting the user id for the task created event",
			"taskId", saved.Id, "userId", saved.UserId, "err", err)
		return
	}
	event := TaskCreated{
		TaskId:      saved.Id,
		UserId:      userId,
		Description: saved.Description,
		// The persisted timestamp arrives in the database session's timezone
		// (RETURNING created_at); the fact's contract is RFC 3339 UTC, the
		// same normalization the domain applies when constructing the task.
		CreatedAt: saved.CreatedAt.UTC(),
	}
	if err := s.publisher.PublishTaskCreated(event); err != nil {
		slog.Error("failed publishing the task created event", "taskId", saved.Id, "err", err)
	}
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
