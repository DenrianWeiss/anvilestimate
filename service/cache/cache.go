package cache

import "sync"

var respCache = sync.Map{}

func SetRespCache(key string, value interface{}) {
	respCache.Store(key, value)
}

func ReadRespCache(key string) (interface{}, bool) {
	return respCache.Load(key)
}

func DelRespCache(key string) {
	respCache.Delete(key)
}
