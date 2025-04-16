package repository

import (
	"context"
	"github.com/Deeksharma/taskmanager/internal/models"
)

type TaskDatabaseRepo interface {
	New(ctx context.Context, task *models.Task) (int, error)
	Update(ctx context.Context, id string, updateFields map[string]interface{}) error
	ById(ctx context.Context, id string) (*models.Task, error)
	All(ctx context.Context, filter map[string]interface{}) ([]*models.Task, error)
	Delete(ctx context.Context, id string) error
}
