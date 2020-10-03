package couchbase

type CouchData struct {
	BuildingName  string `json:"buildingName"`
	EnergyType    string `json:"energyType"`
	EnergyUnit    string `json:"energyUnit"`
	EnergyValue   int    `json:"energyValue"`
	UnixTimeValue int    `json:"unixTimeValue"`
}
