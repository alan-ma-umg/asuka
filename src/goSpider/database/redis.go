package database

import (
	"github.com/go-redis/redis"
	"goSpider/helper"
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
		redisInstance = redis.NewClient(&redis.Options{Addr: helper.Env().Redis.Addr, Password: "", DB: helper.Env().Redis.DB})
	})
	return redisInstance
}

func AddUrlQueue(link string) {
	Redis().RPush(helper.Env().Redis.URLQueueKey, link)
}

func PopUrlQueue() (string, error) {
	return Redis().LPop(helper.Env().Redis.URLQueueKey).Result()
}
