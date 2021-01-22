package mongo

import (
	"context"
	"energy-dashboard-api/datacache"
	"energy-dashboard-api/graph/model"
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CampusHomeKWQuery(returnValue chan *model.EnergyDataPointsReturn, cache *ristretto.Cache) {
	errors := model.Errors{}
	data := []*model.EnergyDataPoint{}

	err := godotenv.Load(".env")
	if err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Unable to load environment variable")
		returnData := model.EnergyDataPointsReturn{
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
	cachedData := datacache.CacheLookup(cache, "campus-kw")
	if cachedData != nil {
		returnValue <- cachedData
	} else {
		// Replace the uri string with your MongoDB deployment's connection string.
		uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s", os.Getenv("MONGO_USR"), os.Getenv("MONGO_PASS"), os.Getenv("MONGO_URL"), os.Getenv("MONGO_AUTHDB"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			errors.Error = true
			errors.Errors = append(errors.Errors, "Errors connecting to Mongo instance")
			returnData := model.EnergyDataPointsReturn{
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
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				errors.Error = true
				errors.Errors = append(errors.Errors, "Error on client disconnect")
				returnData := model.EnergyDataPointsReturn{
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
			returnData := model.EnergyDataPointsReturn{
				Data:   data,
				Errors: &errors,
			}
			returnValue <- &returnData
		}

		collection := client.Database("energy-dashboard").Collection("kw")

		dateLow := time.Now().Local().AddDate(0, 0, -1).Unix()
		dateHigh := time.Now().Local().Unix()

		filter := bson.M{
			"$and": []interface{}{
				bson.M{"BuildingName": "campus"},
				bson.M{"$and": []interface{}{
					bson.M{"EnergyType": "electric"},
					bson.M{"$and": []interface{}{
						bson.M{"UnixTimeValue": bson.M{"$gte": dateLow}},
						bson.M{"$and": []interface{}{
							bson.M{"UnixTimeValue": bson.M{"$lte": dateHigh}},
						}},
					}},
				}},
			},
		}

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
				Data:   data,
				Errors: &errors,
			}
			returnValue <- &returnData
		}

		if err = energyDocs.All(ctx, &energyDataPointsJSON); err != nil {
			errors.Error = true
			errors.Errors = append(errors.Errors, "Error when parsing Mongo data")
			returnData := model.EnergyDataPointsReturn{
				Data:   data,
				Errors: &errors,
			}
			returnValue <- &returnData
		}

		addCache := datacache.SetCache(energyDataPointsJSON, cache, "0h20m", "campus-kw")
		if !addCache {
			errors.Error = true
			errors.Errors = append(errors.Errors, "Error creating cache. NOTE: This is not a fatal error")
		}
		returnData := model.EnergyDataPointsReturn{
			Data:   energyDataPointsJSON,
			Errors: &errors,
		}
		returnValue <- &returnData
	}
}

func CampusHomeKWHQuery(returnValue chan *model.EnergyDataPointsReturn, cache *ristretto.Cache) {
	errors := model.Errors{}
	data := []*model.EnergyDataPoint{}

	err := godotenv.Load(".env")
	if err != nil {
		errors.Error = true
		errors.Errors = append(errors.Errors, "Unable to load environment variable")
		returnData := model.EnergyDataPointsReturn{
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
	cachedData := datacache.CacheLookup(cache, "campus-kwh")
	if cachedData != nil {
		returnValue <- cachedData
	} else {
		// Replace the uri string with your MongoDB deployment's connection string.
		uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s", os.Getenv("MONGO_USR"), os.Getenv("MONGO_PASS"), os.Getenv("MONGO_URL"), os.Getenv("MONGO_AUTHDB"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			errors.Error = true
			errors.Errors = append(errors.Errors, "Errors connecting to Mongo instance")
			returnData := model.EnergyDataPointsReturn{
				Data:   data,
				Errors: &errors,
			}
			returnValue <- &returnData
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				errors.Error = true
				errors.Errors = append(errors.Errors, "Error on client disconnect")
				returnData := model.EnergyDataPointsReturn{
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
			returnData := model.EnergyDataPointsReturn{
				Data:   data,
				Errors: &errors,
			}
			returnValue <- &returnData
		}

		collection := client.Database("energy-dashboard").Collection("kwh")

		dateLow := time.Now().Local().AddDate(0, 0, -1).Unix()
		dateHigh := time.Now().Local().Unix()

		filter := bson.M{
			"$and": []interface{}{
				bson.M{"BuildingName": "campus"},
				bson.M{"$and": []interface{}{
					bson.M{"EnergyType": "electric"},
					bson.M{"$and": []interface{}{
						bson.M{"UnixTimeValue": bson.M{"$gte": dateLow}},
						bson.M{"$and": []interface{}{
							bson.M{"UnixTimeValue": bson.M{"$lte": dateHigh}},
						}},
					}},
				}},
			},
		}

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
				Data:   data,
				Errors: &errors,
			}
			returnValue <- &returnData
		}

		if err = energyDocs.All(ctx, &energyDataPointsJSON); err != nil {
			errors.Error = true
			errors.Errors = append(errors.Errors, "Error when parsing Mongo data")
			returnData := model.EnergyDataPointsReturn{
				Data:   data,
				Errors: &errors,
			}
			returnValue <- &returnData
		}

		addCache := datacache.SetCache(energyDataPointsJSON, cache, "24h20m", "campus-kwh")
		if !addCache {
			errors.Error = true
			errors.Errors = append(errors.Errors, "Error creating cache. NOTE: This is not a fatal error")
		}
		returnData := model.EnergyDataPointsReturn{
			Data:   energyDataPointsJSON,
			Errors: &errors,
		}

		returnValue <- &returnData
	}
}

func BuildingQuery(returnValue chan *[]model.BuildingInfo, cache *ristretto.Cache) {
}
