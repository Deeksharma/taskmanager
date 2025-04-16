package dbrepo

import (
	"context"
	"github.com/Deeksharma/taskmanager/internal/log"
	"github.com/Deeksharma/taskmanager/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

func (md *taskDBRepo) New(ctx context.Context, submission *models.Task) (int, error) {
	submission.CreatedAt = time.Now()
	submission.UpdatedAt = time.Now()
	_, err := md.TaskCollection.InsertOne(ctx, submission)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (md *taskDBRepo) Update(ctx context.Context, id string, updateFields map[string]interface{}) error {
	updateFields["updated_at"] = time.Now()
	var update []bson.D
	for key, value := range updateFields {
		update = append(update, bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: key, Value: value}}}})
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateResult, err := md.TaskCollection.UpdateOne(ctx,
		bson.D{bson.E{Key: "_id", Value: objID}},
		update)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (md *taskDBRepo) ById(ctx context.Context, id string) (deployment *models.Task, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return deployment, err
	}
	err = md.TaskCollection.FindOne(ctx, bson.D{bson.E{Key: "_id", Value: objID}}).Decode(&deployment)
	if err != nil {
		return deployment, err
	}
	return deployment, nil
}

func (md *taskDBRepo) All(ctx context.Context, filter map[string]interface{}) (deployments []*models.Task, err error) {
	filterConditions := bson.M{}
	for key, value := range filter {
		filterConditions[key] = value
	}

	cursor, err := md.TaskCollection.Find(ctx, filterConditions, options.Find().SetSort(bson.D{bson.E{Key: "updated_at", Value: -1}}))
	if err != nil {
		return deployments, err
	}
	for cursor.Next(ctx) {
		var submission *models.Task
		err = cursor.Decode(&submission)
		if err != nil {
			log.InfoWithFields(ctx, map[string]interface{}{
				"error": err,
			}, "Error while decoding document")
		}
		deployments = append(deployments, submission)
	}

	if err := cursor.Err(); err != nil {
		log.InfoWithFields(ctx, map[string]interface{}{
			"error": err,
		}, "Error while decoding document")
	}

	//Close the cursor once finished
	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	return deployments, nil
}

func (*taskDBRepo) Delete(ctx context.Context, id string) error {
	return nil
}
