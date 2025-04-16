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

func NewDeploymentService(taskRepo repository.TaskDatabaseRepo) *TaskManagementService {
	return &TaskManagementService{TaskRepo: taskRepo}
}

func NewTestDeploymentService() *TaskManagementService {
	testingRepo, _ := dbrepo.NewTestingRepo(context.TODO())
	return &TaskManagementService{TaskRepo: testingRepo}
}

// ById returns deployment by id
func (s *TaskManagementService) ById(ctx context.Context, taskId string) (*models.Task, error) {
	deployment, err := s.TaskRepo.ById(ctx, taskId)
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error fetching task")
	}
	return deployment, err
}

// All returns all deployment
func (s *TaskManagementService) All(ctx context.Context, filter map[string]interface{}) ([]*models.Task, error) {
	deployment, err := s.TaskRepo.All(ctx, filter)
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
		}, "error fetching deployments")
		return nil, err
	}
	return deployment, nil
}

// New creates a new task in 'created' state
func (s *TaskManagementService) New(ctx context.Context, taskId string, deployRequestBody *models.Task) (*models.Task, error) {
	portfolio, err := s.TaskRepo.ById(ctx, taskId)
	if err != nil && err == mongo.ErrNoDocuments {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "portfolio not present in database")
		return nil, err
	} else if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
		}, "error while fetching portfolio")
		return nil, err
	}

	task := &models.Task{
		ID:     bson.NewObjectID(),
		Owner:  portfolio.Owner,
		Status: enum.Created,
	}

	_, err = s.TaskRepo.New(ctx, task)
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
			"task":  task,
		}, "error inserting deployment")
		return nil, err
	}

	_ = s.TaskRepo.Update(ctx, taskId, map[string]interface{}{
		"latest": task,
	})

	return task, nil
}

// Update creates a new deployment in 'deployed' state
func (s *TaskManagementService) Update(ctx context.Context, taskId string, deployRequestBody *models.Task) (*models.Task, error) {
	task, err := s.TaskRepo.ById(ctx, taskId)
	if err != nil && err == mongo.ErrNoDocuments {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "taskId not present in database")
		return nil, err
	} else if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
		}, "error while fetching taskId")
		return nil, err
	}

	err = s.TaskRepo.Update(ctx, task.TaskID, make(map[string]interface{}))
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
			"task":  task,
		}, "error inserting deployment")
		return nil, err
	}

	_ = s.TaskRepo.Update(ctx, taskId, map[string]interface{}{
		"latest": task,
	})

	return task, nil
}

// Delete deletes the task id
func (s *TaskManagementService) Delete(ctx context.Context, taskId string) error {
	err := s.TaskRepo.Delete(ctx, taskId)
	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error delete task")
	}
	return err
}
