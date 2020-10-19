package couchbase

import (
	"energy-dashboard-api/graph/model"

	"github.com/dgraph-io/ristretto"
)

func CreateMainCache() *ristretto.Cache {
	mainCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 60, // maximum cost of mainCache (2GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return mainCache
}

func CacheLookup(lookupCache *ristretto.Cache, lookupKey string) []*model.EnergyDataPoint {
	value, found := lookupCache.Get(lookupKey)

	if !found {
		return nil
	}
	returnVal, err := value.([]*model.EnergyDataPoint)

	if !err {
		return nil
	}

	return returnVal

}

func CacheUtil(returnValue chan []*model.EnergyDataPoint) {
	//TODO: Implement

}
