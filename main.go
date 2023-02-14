package main

import (
	"context"
	"time"
	_config "userprofile-delete-script/config"
	"userprofile-delete-script/logger"
	_logger "userprofile-delete-script/logger"

	"go.mongodb.org/mongo-driver/bson"
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
	filter := bson.M{
		"CreatedDate": bson.M{"$lt": time.Now().AddDate(-1, 0, 0)},
	}
	deleteRecords(ctx, collection, filter, logger)
}

func deleteRecords(ctx context.Context, collection *mongo.Collection, filter bson.M, logger _logger.Logger) {
	result, err := collection.Find(ctx, filter)
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
