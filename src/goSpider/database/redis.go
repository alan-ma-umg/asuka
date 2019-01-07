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
		redisInstance = redis.NewClient(&redis.Options{Addr: helper.Env().Redis.Server, Password: "", DB: helper.Env().Redis.DB})
	})
	return redisInstance
}

//var list = make([]string, 0, 100000)

func AddUrlQueue(link string) {
	//list = append(list, link)
	Redis().RPush(helper.Env().Redis.URLQueueKey, link)
}

func PopUrlQueue() (string, error) {
	//var list = make([]string, 100000)
	//x, list := list[len(list)-1], list[:len(list)-1]

	//if x == "" {
	//	x, list = list[len(list)-1], list[:len(list)-1]
	//}
	//
	//if x == "" {
	//	return "", errors.New("hahahahaa ")
	//}
	//
	//x, list := list[len(list)-1], list[:len(list)-1]
	//fmt.Println(x)
	//x, l := list[0], list[1:]
	//fmt.Println(l)
	//list = list
	//fmt.Println(x)
	//return x, nil
	return Redis().LPop(helper.Env().Redis.URLQueueKey).Result()
}
