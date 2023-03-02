package main

import (
	"context"
	"fmt"
	"strings"
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
type ThemeEntry struct {
	AppId int    `bson:"AppID"`
	Pages []Page `bson:"Pages"`
}
type Page struct {
	Status string `bson:"Status"`
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
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Mongo Client Connect")
	defer client.Disconnect(ctx)
	database := client.Database(config.DatabaseName)
	collection := database.Collection(config.CollectionName)

	updateFilter := bson.M{
		"IsActive":     true,
		"Pages.Status": bson.M{"$type": 2},
	}
	deleteRecords(ctx, collection, logger, "", "", false, &updateFilter)
}

func deleteRecords(ctx context.Context, collection *mongo.Collection, logger _logger.Logger, key string, value string, isRegex bool, query *bson.M) {
	var result *mongo.Cursor
	var err error
	result, err = collection.Find(ctx, query)
	if err != nil {
		logger.Fatal(err)
	}
	count := 654
	i := 1
	for result.Next(ctx) {
		var elem ThemeEntry
		err = result.Decode(&elem)
		if err != nil {
			logger.Error(err)
		} else {
			logger.Info(elem)
		}
		var val string
		if strings.Contains(elem.Pages[0].Status, "1") {
			val = "0"
			if strings.Contains(elem.Pages[0].Status, "9") {
				val = "09"
			}
		}
		if strings.Contains(elem.Pages[0].Status, "0") {
			val = "1"
			if strings.Contains(elem.Pages[0].Status, "9") {
				val = "19"
			}
		}
		if len(val) != 0 && i > count {
			query := bson.M{
				"AppID": elem.AppId,
			}
			update := bson.M{
				"$set": bson.M{
					"Pages.$[].Status": val,
				},
			}
			fmt.Println("Updating the val: ", val, elem.AppId)
			if _, err := collection.UpdateOne(ctx, query, update); err != nil {
				logger.Error(err)
			}
		}
		i++
	}
}
