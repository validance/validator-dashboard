package services

import (
	goc "github.com/patrickmn/go-cache"
	"time"
)

var cache *goc.Cache

func AddToCache(key string, value interface{}) {
	if cache == nil {
		cache = newCache()
	}
	cache.Set(key, value, goc.DefaultExpiration)
}

func GetCache(key string) (interface{}, bool) {
	if cache == nil {
		return nil, false
	}
	return cache.Get(key)
}

func newCache() *goc.Cache {
	return goc.New(24*time.Hour, 24*time.Hour)
}
