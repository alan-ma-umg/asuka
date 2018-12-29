package database

import (
	"github.com/go-redis/redis"
	"goSpider/helper"
	"sync"
	"time"
)

var redisOnce sync.Once
var redisInstance *redis.Client
var addQueueCollection []string

func init() {
	//save
	go func() {
		t := time.NewTicker(time.Second * 5)
		for {
			<-t.C
			UrlQueueSave()
		}
	}()
}

func UrlQueueSave() {
	count := len(addQueueCollection)
	if count > 0 {
		Redis().RPush(helper.Env().Redis.URLQueueKey, addQueueCollection)
		addQueueCollection = make([]string, 0, count)
	}
}

func Redis() *redis.Client {
	redisOnce.Do(func() {
		redisInstance = redis.NewClient(&redis.Options{Addr: helper.Env().Redis.Addr, Password: "", DB: helper.Env().Redis.DB})
	})
	return redisInstance
}

func AddUrlQueue(link string) {
	addQueueCollection = append(addQueueCollection, link)
}

func PopUrlQueue() (string, error) {
	return Redis().LPop(helper.Env().Redis.URLQueueKey).Result()
}
