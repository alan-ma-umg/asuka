package database

import (
	"github.com/chenset/asuka/helper"
	"github.com/go-redis/redis"
	"sync"
)

var redisOnce sync.Once
var redisInstance *redis.Client

func Redis() *redis.Client {
	redisOnce.Do(func() {
		redisInstance = redis.NewClient(&redis.Options{Network: helper.Env().Redis.Network, Addr: helper.Env().Redis.Addr, Password: helper.Env().Redis.Password, DB: helper.Env().Redis.DB})
	})
	return redisInstance
}
