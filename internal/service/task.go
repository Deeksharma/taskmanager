package service

import (
	"context"
	"github.com/Deeksharma/taskmanager/internal/enum"
	"github.com/Deeksharma/taskmanager/internal/log"
	"github.com/Deeksharma/taskmanager/internal/models"
	"github.com/Deeksharma/taskmanager/internal/repository"
	"github.com/Deeksharma/taskmanager/internal/repository/dbrepo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TaskManagementService struct {
	TaskRepo repository.TaskDatabaseRepo
}

func NewTaskManagementService(taskRepo repository.TaskDatabaseRepo) *TaskManagementService {
	return &TaskManagementService{TaskRepo: taskRepo}
}

func NewTestTaskManagementService() *TaskManagementService {
	testingRepo, _ := dbrepo.NewTestingRepo(context.TODO())
	return &TaskManagementService{TaskRepo: testingRepo}
}

// ById returns task by id
func (s *TaskManagementService) ById(ctx context.Context, taskId string) (*models.Task, error) {
	task, err := s.TaskRepo.ById(ctx, taskId)
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error fetching task")
	}
	return task, err
}

// All returns all tasks
func (s *TaskManagementService) All(ctx context.Context, filter map[string]interface{}) ([]*models.Task, error) {
	tasks, err := s.TaskRepo.All(ctx, filter)
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
		}, "error fetching tasks")
		return nil, err
	}
	return tasks, nil
}

// New creates a new task in 'created' state
func (s *TaskManagementService) New(ctx context.Context, taskRequest *models.Task) (*models.Task, error) {
	id := bson.NewObjectID()
	task := &models.Task{
		ID:          id,
		Title:       taskRequest.Title,
		Owner:       taskRequest.Owner,
		Description: taskRequest.Description,
		TaskID:      id.Hex(),
		Status:      enum.Created,
	}

	_, err := s.TaskRepo.New(ctx, task)
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
			"task":  task,
		}, "error creating new task")
		return nil, err
	}
	return task, nil
}

// Update updates an existing task
func (s *TaskManagementService) Update(ctx context.Context, taskId string, updatedFields map[string]interface{}) (*models.Task, error) {
	task, err := s.TaskRepo.ById(ctx, taskId)
	if err != nil && err == mongo.ErrNoDocuments {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "taskId not present in database")
		return nil, err
	} else if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error while fetching task")
		return nil, err
	}

	updatedTitle, ok := updatedFields["title"]
	if ok {
		task.Title = updatedTitle.(string)
	}
	updatedOwner, ok := updatedFields["owner"]
	if ok {
		task.Owner = updatedOwner.(string)
	}
	updatedDescription, ok := updatedFields["description"]
	if ok {
		task.Description = updatedDescription.(string)
	}
	updatedStatus, ok := updatedFields["status"]
	if ok {
		task.Status = updatedStatus.(enum.TaskStatus)
	}

	err = s.TaskRepo.Update(ctx, task.TaskID, updatedFields)
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
			"task":  task,
		}, "error updating task")
		return nil, err
	}
	return task, nil
}

// Delete deletes the task id
func (s *TaskManagementService) Delete(ctx context.Context, taskId string) error {
	err := s.TaskRepo.Delete(ctx, taskId)
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error deleting task")
	}
	return err
}
