package database

import (
	"asuka/helper"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

var redisOnce sync.Once
var redisInstance *redis.Client

func init() {
	//save
	go func() {
		t := time.NewTicker(time.Second * 10)
		for {
			<-t.C
			//todo remove
			count, _ := Redis().LLen(helper.Env().Redis.URLQueueKey).Result()
			if count > 100000 {
				Redis().LTrim(helper.Env().Redis.URLQueueKey, 0, 10000)
			}
		}
	}()
}

func Redis() *redis.Client {
	redisOnce.Do(func() {
		redisInstance = redis.NewClient(&redis.Options{Addr: helper.Env().Redis.Server, Password: "", DB: helper.Env().Redis.DB})
	})
	return redisInstance
}