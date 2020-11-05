package mongo

import (
	"context"
	"energy-dashboard-api/datacache"
	"energy-dashboard-api/graph/model"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CampusHomeKWQuery(returnValue chan []*model.EnergyDataPoint, cache *ristretto.Cache) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file")
	}
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
			panic(err)
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()
		// Ping the primary
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			panic(err)
		}
		fmt.Println("Successfully connected and pinged.")

		collection := client.Database("energy-dashboard").Collection("kw")

		dateLow := time.Now().AddDate(0, 0, -1).Unix()

		dateHigh := time.Now().Unix()

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
		var energyDataPointsBSON []EnergyDataPointMongo
		var energyDataPointsJSON []*model.EnergyDataPoint

		opts := options.Find()
		opts.SetSort(bson.D{{"UnixTimeValue", -1}})
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
		energyDocs, err := collection.Find(ctx, filter, opts)

		if err != nil {
			panic(err)
		}

		if err = energyDocs.All(ctx, &energyDataPointsBSON); err != nil {
			panic(err)
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
		addCache := datacache.SetCache(energyDataPointsJSON, cache, "0h20m", "campus-kw")
		if !addCache {
			log.Fatal("Cache Creation Failed")
		}
		returnValue <- energyDataPointsJSON
	}

}

func CampusHomeKWHQuery(returnValue chan []*model.EnergyDataPoint, cache *ristretto.Cache) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file")
	}
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
			panic(err)
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()
		// Ping the primary
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			panic(err)
		}
		fmt.Println("Successfully connected and pinged.")

		collection := client.Database("energy-dashboard").Collection("kwh")

		dateLow := time.Now().AddDate(0, 0, -1).Unix()

		dateHigh := time.Now().Unix()

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
		var energyDataPointsBSON []EnergyDataPointMongo
		var energyDataPointsJSON []*model.EnergyDataPoint

		opts := options.Find()
		opts.SetSort(bson.D{{"UnixTimeValue", -1}})
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
		energyDocs, err := collection.Find(ctx, filter, opts)

		if err != nil {
			panic(err)
		}

		if err = energyDocs.All(ctx, &energyDataPointsBSON); err != nil {
			panic(err)
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
		addCache := datacache.SetCache(energyDataPointsJSON, cache, "24h30m", "campus-kwh")
		if !addCache {
			log.Fatal("Cache Creation Failed")
		}
		returnValue <- energyDataPointsJSON
	}
}
