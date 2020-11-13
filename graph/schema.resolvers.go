package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"energy-dashboard-api/datacache"
	"energy-dashboard-api/graph/generated"
	"energy-dashboard-api/graph/model"
	"energy-dashboard-api/mongo"
	"fmt"
	"time"
)

func (r *queryResolver) EnergyDataPoints(ctx context.Context, input model.EnergyDataPointQueryInput) (*model.EnergyDataPointsReturn, error) {
	var returnValue chan *model.EnergyDataPointsReturn
	if returnValue == nil {
		returnValue = make(chan *model.EnergyDataPointsReturn)
	}
	go mongo.DateRangeQuery(returnValue, input.EnergyUnit, int64(input.DateLow), int64(input.DateHigh), input.Building, input.EnergyType)

	return <-returnValue, nil
}

func (r *queryResolver) Past24Hours(ctx context.Context, input model.Past24HoursInput) (*model.EnergyDataPointsReturn, error) {
	var returnValue chan *model.EnergyDataPointsReturn
	if returnValue == nil {
		returnValue = make(chan *model.EnergyDataPointsReturn)
	}
	go mongo.DateRangeQuery(returnValue, input.EnergyUnit, (time.Now().Unix() - 86400), time.Now().Unix(), input.Building, input.EnergyType)

	return <-returnValue, nil
}

func (r *queryResolver) DashboardHomePage(ctx context.Context) (*model.DashboardHomePage, error) {
	var kwEnergyValues chan *model.EnergyDataPointsReturn
	if kwEnergyValues == nil {
		kwEnergyValues = make(chan *model.EnergyDataPointsReturn)
	}
	var kwhEnergyValues chan *model.EnergyDataPointsReturn
	if kwhEnergyValues == nil {
		kwhEnergyValues = make(chan *model.EnergyDataPointsReturn)
	}
	go mongo.CampusHomeKWQuery(kwEnergyValues, datacache.GetDataPointCache())
	go mongo.CampusHomeKWQuery(kwhEnergyValues, datacache.GetDataPointCache())

	returnData := model.DashboardHomePage{
		CampusKw:  <-kwEnergyValues,
		CampusKwh: <-kwhEnergyValues,
	}
	return &returnData, nil
}

func (r *queryResolver) BuildingInfo(ctx context.Context, input model.BuildingInfoInput) (*model.BuildingInfo, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.

type mutationResolver struct{ *Resolver }
