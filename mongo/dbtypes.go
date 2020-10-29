package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnergyDataPointMongo struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	BuildingName  string             `bson:"buildingname,omitempty"`
	EnergyType    string             `bson:"energytype,omitempty"`
	EnergyUnit    string             `bson:"energyunit,omitempty"`
	EnergyValue   float64            `bson:"energyvalue,omitempty"`
	UnixTimeValue int                `bson:"unixtimevalue,omitempty"`
}
