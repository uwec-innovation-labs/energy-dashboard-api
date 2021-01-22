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

func BuildingsQuery(returnValue chan *model.BuildingsQueryReturn) {
	maxDate := 1611333000
	minDate := 1554138000
	errors := model.Errors{}
	data := []*model.BuildingInfo{}

	err := godotenv.Load(".env")
	if err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Unable to load environment variable")
		returnData := model.BuildingsQueryReturn{
			Data:   data,
			Errors: &errors,
		}
		returnValue <- &returnData
	}
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s", os.Getenv("MONGO_USR"), os.Getenv("MONGO_PASS"), os.Getenv("MONGO_URL"), os.Getenv("MONGO_AUTHDB"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Errors connecting to Mongo instance")
		returnData := model.BuildingsQueryReturn{
			Data:   data,
			Errors: &errors,
		}
		returnValue <- &returnData
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			errors.Error = true
			errors.Errors = append(errors.Errors, "Error on client disconnect")
			returnData := model.BuildingsQueryReturn{
				Data:   data,
				Errors: &errors,
			}
			returnValue <- &returnData
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Error pinging Mongo instance")
		returnData := model.BuildingsQueryReturn{
			Data:   data,
			Errors: &errors,
		}
		returnValue <- &returnData
	}

	collection := client.Database("energy-dashboard").Collection("buildings")

	filter := bson.M{}

	var buildingDataReturn []*model.BuildingInfo


	ctx, cancelFilter := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelFilter()
	buildingDocs, err := collection.Find(ctx, filter)

	if err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Error when querying Mongo instance")
		returnData := model.BuildingsQueryReturn{
			Data:   data,
			Errors: &errors,
		}
		returnValue <- &returnData
	}

	if err = buildingDocs.All(ctx, &buildingDataReturn); err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Error when parsing Mongo data")
		fmt.Print(err)
		returnData := model.BuildingsQueryReturn{
			Data:   data,
			Errors: &errors,
		}
		returnValue <- &returnData
	}

	for _, i := range buildingDataReturn {
		for _, j := range i.EnergyInfo {
			j.MinDate = minDate
			j.MaxDate = maxDate
		}
	}

	returnData := model.BuildingsQueryReturn{
		Data:   buildingDataReturn,
		Errors: &errors,
	}
	returnValue <- &returnData
}
