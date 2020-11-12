package mongo

import (
	"context"
	"energy-dashboard-api/graph/model"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DateRangeQuery(returnValue chan *model.EnergyDataPointsReturn, bucketName string, dateLow int64, dateHigh int64, building string, energyType string) {
	errors := model.Errors{}
	err := godotenv.Load(".env")
	if err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Unable to load environment variable")
		returnData := model.EnergyDataPointsReturn{
			Data:   nil,
			Errors: &errors,
		}
		returnValue <- &returnData
	}
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s", os.Getenv("MONGO_USR"), os.Getenv("MONGO_PASS"), os.Getenv("MONGO_URL"), os.Getenv("MONGO_AUTHDB"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Errors connecting to Mongo instance")
		returnData := model.EnergyDataPointsReturn{
			Data:   nil,
			Errors: &errors,
		}
		returnValue <- &returnData
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			errors.Error = true
			errors.Errors = append(errors.Errors, "Error on client disconnect")
			returnData := model.EnergyDataPointsReturn{
				Data:   nil,
				Errors: &errors,
			}
			returnValue <- &returnData
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Error pinging Mongo instance")
	}
	fmt.Println("Successfully connected and pinged.")

	collection := client.Database("energy-dashboard").Collection(bucketName)

	filter := bson.M{
		"$and": []interface{}{
			bson.M{"BuildingName": building},
			bson.M{"$and": []interface{}{
				bson.M{"EnergyType": energyType},
				bson.M{"$and": []interface{}{
					bson.M{"UnixTimeValue": bson.M{"$gte": dateLow}},
					bson.M{"$and": []interface{}{
						bson.M{"UnixTimeValue": bson.M{"$lte": dateHigh}},
					}},
				}},
			}},
		},
	}

	var energyDataPointsBSON []EnergyDataPointMongo
	var energyDataPointsJSON []*model.EnergyDataPoint

	opts := options.Find()
	opts.SetSort(bson.D{{"UnixTimeValue", -1}})
	ctx, cancelFilter := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelFilter()
	energyDocs, err := collection.Find(ctx, filter, opts)

	if err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Error when querying Mongo instance")
		returnData := model.EnergyDataPointsReturn{
			Data:   nil,
			Errors: &errors,
		}
		returnValue <- &returnData
	}

	if err = energyDocs.All(ctx, &energyDataPointsBSON); err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Error when parsing Mongo data")
		returnData := model.EnergyDataPointsReturn{
			Data:   nil,
			Errors: &errors,
		}
		returnValue <- &returnData
	}

	fmt.Println("Query successful")

	fmt.Println("Format loop")
	for _, doc := range energyDataPointsBSON {
		dataPoint := &model.EnergyDataPoint{
			Building:     doc.BuildingName,
			Value:        doc.EnergyValue,
			DateTimeUnix: doc.UnixTimeValue,
			Unit:         doc.EnergyUnit,
			Type:         doc.EnergyType,
		}
		energyDataPointsJSON = append(energyDataPointsJSON, dataPoint)
	}

	returnData := model.EnergyDataPointsReturn{
		Data:   energyDataPointsJSON,
		Errors: &errors,
	}
	returnValue <- &returnData

}
