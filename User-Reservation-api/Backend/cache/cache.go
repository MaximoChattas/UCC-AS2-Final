package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
)

var cache *memcache.Client

func InitCache() {
	cache = memcache.New("memcached:11211")
}

func Set(key string, value []byte) {
	err := cache.Set(&memcache.Item{Key: key, Value: value, Expiration: 10})

	if err != nil {
		log.Error("Error setting cache: ", err)
	}
}

func Get(key string) ([]byte, error) {

	resp, err := cache.Get(key)

	if err != nil {
		log.Error("Error getting cache data: ", err)
		return []byte{}, err
	}

	log.Info("Retrieved data from cache")
	return resp.Value, nil
}
