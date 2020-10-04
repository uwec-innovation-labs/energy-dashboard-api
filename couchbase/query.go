package couchbase

import (
	"energy-dashboard-api/graph/model"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

/* DateRangeQuery queries the database on a date range, building name, energy type, and bucket (kw/kwh) */
func DateRangeQuery(bucketName string, dateLow int, dateHigh int, building string, energyType string) []model.EnergyDataPoint {
	cluster, err := gocb.Connect(
		"localhost",
		gocb.ClusterOptions{
			Username: "Administrator",
			Password: "password",
		})
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(5*time.Second, nil)

	query := fmt.Sprintf("SELECT doc.* FROM `%s` doc WHERE doc.buildingName = \"%s\" AND doc.energyType = \"%s\" AND doc.unixTimeValue >= %d AND doc.unixTimeValue <= %d",
		bucketName, building, energyType, dateLow, dateHigh)

	rows, err := cluster.Query(query, nil)

	if err != nil {
		panic(err)
	}

	var energyDataCouchbase []model.EnergyDataPoint

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

		energyDataCouchbase = append(energyDataCouchbase, model.EnergyDataPoint{
			Value:        energyPoint.EnergyValue,
			Building:     energyPoint.BuildingName,
			DateTimeUnix: energyPoint.UnixTimeValue,
			Unit:         energyPoint.EnergyUnit,
			Type:         energyPoint.EnergyType,
		})
	}

	return energyDataCouchbase
}
