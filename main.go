package main

import (
	"context"
	"time"
	_config "userprofile-delete-script/config"
	"userprofile-delete-script/logger"
	_logger "userprofile-delete-script/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogEntry struct {
	ID          string    `bson:"_id"`
	CreatedDate time.Time `bson:"CreatedDate"`
	Uid         string    `bson:"Uid"`
	Email       []Email   `bson:"Email" json:"Email"`
}

type Email struct {
	Type  string `bson:"Type" json:"Type"`
	Value string `bson:"Value" json:"Value"`
}

func main() {
	config := _config.GetInstance()
	logger := _logger.NewLogger(logger.Config{Name: "delete-user-profile-script"})

	client, err := mongo.NewClient(options.Client().ApplyURI(config.ConnectionString))
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Mongo Client created")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Mongo Client Connect")
	defer client.Disconnect(ctx)
	database := client.Database(config.DatabaseName)
	collection := database.Collection(config.CollectionName)

	// ** This is Regex Query
	deleteRecords(ctx, collection, logger, "Email.Value", "test\\w*", true, nil)

	// ** This is simple Query
	filter := bson.M{
		"CreatedDate": bson.M{"$lt": time.Now().AddDate(-1, 0, 0)},
	}
	deleteRecords(ctx, collection, logger, "", "", false, &filter)
}

func deleteRecords(ctx context.Context, collection *mongo.Collection, logger _logger.Logger, key string, value string, isRegex bool, query *bson.M) {
	var filter bson.D
	if isRegex {
		filter = bson.D{
			{Key: key, Value: bson.D{
				{"$regex", primitive.Regex{Pattern: value}},
			}},
		}
	}
	var result *mongo.Cursor
	var err error
	if query != nil {
		result, err = collection.Find(ctx, query)
	} else {
		result, err = collection.Find(ctx, filter)
	}
	if err != nil {
		logger.Fatal(err)
	}
	for result.Next(ctx) {
		var elem LogEntry
		err = result.Decode(&elem)
		if err != nil {
			logger.Error(err)
		} else {
			logger.Info(elem)
		}
	}
	// Delete records
	/***
	delRes, dErr := collection.DeleteMany(ctx, filter)
	if dErr != nil {
		logger.Fatal(err)
	}
	logger.Info(fmt.Sprintf("%v records deleted", delRes.DeletedCount))
	*/
}
