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
	// deleteRecords(ctx, collection, logger, "Email.Value", "test\\w*", true, nil)

	// ** This query check the cyrilic charcter
	// fromTimeStr := "2019-12-23 10:26:34" // Date is YYYY-MM-DD HH:MM:SS fomat
	// toTimeStr := "2021-12-31 07:17:11"   // Date is YYYY-MM-DD HH:MM:SS fomat
	// toTime, err := time.Parse("2006-01-02 03:04:05", toTimeStr)
	// if err != nil {
	// 	logger.Error(err)
	// }
	// fromTime, err := time.Parse("2006-01-02 03:04:05", fromTimeStr)
	// if err != nil {
	// 	logger.Error(err)
	// }
	// filter := bson.M{
	// 	// "CreatedDate": bson.M{"$lt": time.Now().AddDate(-1, 0, 0)},
	// 	"CreatedDate": bson.M{"$gt": fromTime, "$lt": toTime},
	// 	"$or": []bson.M{{"FirstName": bson.D{
	// 		{"$regex", primitive.Regex{Pattern: "[\\p{Cyrillic}\\d]+"}}, // This regex is used to check the Cyrillic letter
	// 	}},
	// 		{"Fullname": bson.D{
	// 			{"$regex", primitive.Regex{Pattern: "[\\p{Cyrillic}\\d]+"}},
	// 		}},
	// 	},
	// }
	// deleteRecords(ctx, collection, logger, "", "", false, &filter)

	updateFilter := bson.M{
		// "CreatedDate": bson.M{"$lt": time.Now().AddDate(-1, 0, 0)},
		"IsActive":     true,
		"Pages.Status": bson.M{"$type": 2},
		// "AppID":        4498,
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
	count := 10
	i := 0
	for result.Next(ctx) {
		if i == count {
			var w1 string
			fmt.Println("Wish to see more records? (press y else n)")
			_, err := fmt.Scanln(&w1)
			if err == nil && w1 == "y" {
				i = 0
			} else {
				break
			}
		}
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
		if len(val) != 0 {
			logger.Info("Updating the val", val)
			query := bson.M{
				"AppID": elem.AppId,
			}
			update := bson.M{
				"$set": bson.M{
					"Pages.$[].Status": val,
				},
			}
			var w1 string
			fmt.Println("Wish to see more records? (press y else n)")
			_, err := fmt.Scanln(&w1)
			if err == nil && w1 == "y" {
				if _, err := collection.UpdateOne(ctx, query, update); err != nil {
					logger.Error(err)
				}
			}

			break
		}
		i++
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
