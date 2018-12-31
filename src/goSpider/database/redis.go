package database

import (
	"github.com/go-redis/redis"
	"goSpider/helper"
	"sync"
)

var redisOnce sync.Once
var redisInstance *redis.Client

func Redis() *redis.Client {
	redisOnce.Do(func() {
		redisInstance = redis.NewClient(&redis.Options{Addr: helper.Env().Redis.Addr, Password: "", DB: helper.Env().Redis.DB})
	})
	return redisInstance
}

func AddUrlQueue(link string) {
	count, _ := Redis().LLen(helper.Env().Redis.URLQueueKey).Result()
	if count > 100000 { //todo remove
		return
	}

	Redis().RPush(helper.Env().Redis.URLQueueKey, link)
}

func PopUrlQueue() (string, error) {
	return Redis().LPop(helper.Env().Redis.URLQueueKey).Result()
}
