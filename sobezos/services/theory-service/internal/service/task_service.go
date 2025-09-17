package service

import (
	"sobezos/services/theory-service/internal/repository"
	"sobezos/services/theory-service/pkg/models"
	"time"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetRandomTask(tags []string) (*models.Task, error) {
	task, err := s.repo.GetRandomTaskRow(tags)
	if err != nil {
		return nil, err
	}
	tagObjs, err := s.repo.GetTagsByTaskID(task.ID)
	if err != nil {
		return nil, err
	}
	task.Tags = make([]string, len(tagObjs))
	for i, tg := range tagObjs {
		task.Tags[i] = tg.Name
	}
	return task, nil
}

// CreateTask creates a new task (бизнес-логика)
func (s *TaskService) CreateTask(tags []string, question, answer string) error {
	task := &models.Task{
		Tags:      tags,
		Question:  question,
		Answer:    answer,
		CreatedAt: time.Now(),
	}
	return s.repo.CreateTask(task)
}

// UpdateTask updates an existing task (бизнес-логика)
func (s *TaskService) UpdateTask(id int, tags []string, question, answer string) error {
	task := &models.Task{
		ID:       id,
		Tags:     tags,
		Question: question,
		Answer:   answer,
	}
	return s.repo.UpdateTask(task)
}

// GetTaskByID retrieves a task by its ID and attaches tags (бизнес-логика)
func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	tagObjs, err := s.repo.GetTagsByTaskID(task.ID)
	if err != nil {
		return nil, err
	}
	task.Tags = make([]string, len(tagObjs))
	for i, tg := range tagObjs {
		task.Tags[i] = tg.Name
	}
	return task, nil
}

// GetAllTags retrieves all tags (бизнес-логика)
func (s *TaskService) GetAllTags() ([]*models.Tag, error) {
	return s.repo.GetAllTags()
}
