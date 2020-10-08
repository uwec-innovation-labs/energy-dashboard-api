package couchbase

import (
	"energy-dashboard-api/graph/model"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/joho/godotenv"
)

/* DateRangeQuery queries the database on a date range, building name, energy type, and bucket (kw/kwh) */
func DateRangeQuery(returnValue chan []*model.EnergyDataPoint, bucketName string, dateLow int, dateHigh int, building string, energyType string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file")
	}
	cluster, err := gocb.Connect(
		os.Getenv("COUCH_ADDR"),
		gocb.ClusterOptions{
			Username: os.Getenv("COUCH_USR"),
			Password: os.Getenv("COUCH_PASS"),
		})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(5*time.Second, nil)

	query := fmt.Sprintf("SELECT doc.* FROM `%s` doc WHERE doc.BuildingName = '%s' AND doc.EnergyType = '%s' AND doc.UnixTimeValue >= %d AND doc.UnixTimeValue <= %d",
		bucketName, building, energyType, dateLow, dateHigh)

	rows, err := cluster.Query(query, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("Query")

	var energyDataCouchbase []*model.EnergyDataPoint

	/*
	   Value        int    `json:"value"`
	       Building     string `json:"building"`
	       DateTimeUnix int    `json:"dateTimeUnix"`
	       Unit         string `json:"unit"`
	       Type         string `json:"type"`
	*/

	for rows.Next() {
		var energyPoint CouchData
		err := rows.Row(&energyPoint)

		if err != nil {
			panic(err)
		}

		energyDataCouchbase = append(energyDataCouchbase, &model.EnergyDataPoint{
			Value:        energyPoint.EnergyValue,
			Building:     energyPoint.BuildingName,
			DateTimeUnix: energyPoint.UnixTimeValue,
			Unit:         energyPoint.EnergyUnit,
			Type:         energyPoint.EnergyType,
		})
	}
	fmt.Println("returning")

	returnValue <- energyDataCouchbase
}
