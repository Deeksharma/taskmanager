package dbrepo

import (
	"context"
	"github.com/Deeksharma/taskmanager/internal/enum"
	"github.com/Deeksharma/taskmanager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

func (*testTaskDBRepo) New(ctx context.Context, task *models.Task) (int, error) {
	return 0, nil
}

func (*testTaskDBRepo) Update(ctx context.Context, id string, updateFields map[string]interface{}) error {
	return nil
}

func (*testTaskDBRepo) ById(ctx context.Context, id string) (*models.Task, error) {
	if id == "123456" {
		return &models.Task{
			ID:     bson.ObjectID{},
			Title:  "task 1: init",
			Status: enum.Created,
			Owner:  "deeksha",
		}, nil
	}
	if id == "1234567" {
		return nil, mongo.ErrNoDocuments
	}
	if id == "234567" {
		return &models.Task{
			ID:     bson.ObjectID{},
			Title:  "task 1: init",
			Status: enum.Created,
			Owner:  "deeksha",
		}, nil
	}
	return &models.Task{}, nil
}

func (*testTaskDBRepo) All(ctx context.Context, fields map[string]interface{}) ([]*models.Task, error) {
	return []*models.Task{}, nil
}

func (*testTaskDBRepo) Delete(ctx context.Context, id string) error {
	return nil
}
