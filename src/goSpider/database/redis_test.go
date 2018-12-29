package database

import (
	"fmt"
	"goSpider/helper"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	Redis().Set("go_text_key_123123", "i'am a string", time.Minute)
	res, _ := Redis().Get("go_text_key_123123").Result()
	fmt.Println(reflect.TypeOf(res), res)
	Redis().Del("go_text_key_123123")
}

func TestAddQueue(t *testing.T) {

	for i := 0; i < 100; i++ {
		AddUrlQueue(strconv.Itoa(i))
	}
	fmt.Println(Redis().LLen(helper.Env().Redis.URLQueueKey).Result())
	time.Sleep(10e9)
}
