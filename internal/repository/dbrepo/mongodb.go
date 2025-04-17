package dbrepo

import (
	"context"
	"github.com/Deeksharma/taskmanager/internal/log"
	"github.com/Deeksharma/taskmanager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

func (md *taskDBRepo) New(ctx context.Context, task *models.Task) (int, error) {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	_, err := md.TaskCollection.InsertOne(ctx, task)
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

	objID, err := bson.ObjectIDFromHex(id)
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

func (md *taskDBRepo) ById(ctx context.Context, id string) (task *models.Task, err error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return task, err
	}
	err = md.TaskCollection.FindOne(ctx, bson.D{bson.E{Key: "_id", Value: objID}}).Decode(&task)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (md *taskDBRepo) All(ctx context.Context, filter map[string]interface{}) (tasks []*models.Task, err error) {
	filterConditions := bson.M{}
	for key, value := range filter {
		filterConditions[key] = value
	}

	cursor, err := md.TaskCollection.Find(ctx, filterConditions, options.Find().SetSort(bson.D{bson.E{Key: "updated_at", Value: -1}}))
	if err != nil {
		return tasks, err
	}
	for cursor.Next(ctx) {
		var task *models.Task
		err = cursor.Decode(&task)
		if err != nil {
			log.InfoWithFields(ctx, map[string]interface{}{
				"error": err,
			}, "Error while decoding document")
		}
		tasks = append(tasks, task)
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

	return tasks, nil
}

func (md *taskDBRepo) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result := md.TaskCollection.FindOneAndDelete(ctx, bson.D{bson.E{Key: "_id", Value: objID}})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
