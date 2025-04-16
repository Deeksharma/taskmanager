package mongodb

import (
	"context"
	"github.com/Deeksharma/taskmanager/internal/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"time"
)

type ConnectionOptions struct {
	ConnectionURI     string
	ConnectionTimeout time.Duration
	MaxPoolSize       uint64
	MinPoolSize       uint64
}

type Collection struct {
	*mongo.Collection
}

func NewMongoConnection(ctx context.Context, conn *ConnectionOptions) (*mongo.Client, error) {
	// nolint
	ctx, cancel := context.WithTimeout(ctx,
		30*time.Second)
	client, err := mongo.Connect(options.Client().
		ApplyURI(conn.ConnectionURI).
		SetConnectTimeout(time.Second * conn.ConnectionTimeout).
		SetMaxPoolSize(conn.MaxPoolSize).
		SetMinPoolSize(conn.MinPoolSize))
	defer cancel()

	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	log.Info(ctx, "connected successfully to database")

	return client, err
}
