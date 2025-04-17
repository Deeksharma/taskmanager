package dbrepo

import (
	"context"
	"github.com/Deeksharma/taskmanager/internal/config"
	"github.com/Deeksharma/taskmanager/internal/log"
	"github.com/Deeksharma/taskmanager/internal/models"
	"github.com/Deeksharma/taskmanager/internal/repository"
	"github.com/Deeksharma/taskmanager/pkg/mongodb"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type taskDBRepo struct {
	TaskCollection *mongodb.Collection
}

type testTaskDBRepo struct {
	TaskCollection *mongodb.Collection
}

func NewMongoDB(ctx context.Context) (*mongo.Database, func() error) {
	db, err := mongodb.NewMongoConnection(ctx, &mongodb.ConnectionOptions{
		ConnectionURI:     config.GetString("database.connection_uri"),
		ConnectionTimeout: time.Duration(config.GetInt32("database.connection_timeout")),
		MaxPoolSize:       uint64(config.GetInt32("database.max_pool_size")),
		MinPoolSize:       uint64(config.GetInt32("database.min_pool_size")),
	})

	if err != nil {
		log.Fatal(ctx, err.Error())
	}
	return db.Database(config.GetString("database.database")), func() error {
		return db.Disconnect(ctx)
	}
}

func NewDBRepo(ctx context.Context) (repository.TaskDatabaseRepo, func() error) {
	dB, disconnect := NewMongoDB(ctx)
	return &taskDBRepo{
		TaskCollection: &mongodb.Collection{Collection: dB.Collection(models.TaskCollection)},
	}, disconnect
}

func NewTestingRepo(ctx context.Context) (repository.TaskDatabaseRepo, func() error) {
	return &testTaskDBRepo{
		nil,
	}, nil
}
