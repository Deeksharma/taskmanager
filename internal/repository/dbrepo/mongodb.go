package dbrepo

import (
	"context"
	"github.com/Deeksharma/taskmanager/internal/log"
	"github.com/Deeksharma/taskmanager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (md *taskDBRepo) All(ctx context.Context, filter map[string]interface{}, pagination map[string]int32, sort map[string]interface{}) ([]*models.Task, error) {
	filterConditions := bson.D{}
	for key, value := range filter {
		filterConditions = append(filterConditions, bson.E{Key: key, Value: value})
	}
	matchStage := bson.D{
		{"$match", filterConditions},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}},
		},
		}}
	projectStage := bson.D{{"$project", bson.D{
		{"_id", 0},
		{"total_count", 1},
		{"tasks", bson.D{{"$slice", []interface{}{"$data", pagination["startIndex"], pagination["recordPerPage"]}}}},
	}}}
	sortStage := bson.D{
		{"$sort", bson.D{
			{sort["sortBy"].(string), sort["sortOrder"].(int)},
		}},
	}

	result, err := md.TaskCollection.Aggregate(ctx,
		mongo.Pipeline{matchStage, sortStage, groupStage, projectStage})

	if err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
		}, "Error while aggregating documents")
		return nil, err
	}
	type TaskGroup struct {
		Tasks      []*models.Task `bson:"tasks"`
		TotalCount int            `bson:"total_count"`
	}
	var allTasks []TaskGroup
	if err := result.All(ctx, &allTasks); err != nil {
		log.ErrorWithFields(ctx, map[string]interface{}{
			"error": err,
		}, "Error while fetching documents")
		return nil, err
	}
	log.InfoWithFields(ctx, map[string]interface{}{
		"allTasks": allTasks,
	}, "Successfully fetched documents")
	return allTasks[0].Tasks, nil
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
