package datacache

import (
	"energy-dashboard-api/graph/model"
	"time"

	"github.com/dgraph-io/ristretto"
)

func CreateCache() *ristretto.Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 60, // maximum cost of mainCache (2GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return cache
}

func CacheLookup(lookupCache *ristretto.Cache, lookupKey string) *model.EnergyDataPointsReturn {
	value, found := lookupCache.Get(lookupKey)

	if !found {
		return nil
	}
	dataPoints, err := value.([]*model.EnergyDataPoint)

	returnValue := model.EnergyDataPointsReturn{
		Data: dataPoints,
	}

	if !err {
		returnValue.Errors.Error = true
		returnValue.Errors.Errors = append(returnValue.Errors.Errors, "Error when grabbing from cache")
	}

	return &returnValue

}

func SetCache(setValue []*model.EnergyDataPoint, cache *ristretto.Cache, timeExpire string, cacheKey string) bool {
	expireTime, err := time.ParseDuration(timeExpire)
	if err != nil {
		panic(err)
	}
	return cache.SetWithTTL(cacheKey, setValue, 0, expireTime)
}
