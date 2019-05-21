package main

import (
	"flag"
	"fmt"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/project"
	"github.com/chenset/asuka/web"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	redis := flag.String("redis", "tcp://127.0.0.1:6379/10", "Redis connection url, default: tcp://127.0.0.1:6379/10")
	mysql := flag.String("mysql", "root:11111111@(127.0.0.1:3306)/asuka?charset=utf8mb4", "Mysql DSN, default: root:11111111@(127.0.0.1:3306)/asuka?charset=utf8mb4")
	webPassword := flag.String("webPassword", "", "WEB login password, default random")
	bloomFilterPath := flag.String("bloomFilterPath", ".", "BloomFilter save path, default:.")
	localTransport := flag.Bool("localTransport", false, "Enable http.DefaultTransport, default: false")
	flag.Parse()

	u, err := url.Parse(*redis)
	if err != nil {
		log.Fatal(err)
	}
	redisDB, _ := strconv.Atoi(strings.TrimLeft(u.Path, "/"))
	redisPassword, _ := u.User.Password()
	helper.GlobalEnvConfig = &helper.EnvConfig{
		BloomFilterPath: *bloomFilterPath,
		WEBPassword:     *webPassword,
		LocalTransport:  *localTransport,
		MysqlDSN:        *mysql,
		Redis: helper.Redis{
			Network:     u.Scheme,
			Addr:        u.Host,
			Password:    redisPassword,
			DB:          redisDB,
			URLQueueKey: "url_queue_key",
		},
	}
}

func main() {
	mainStart := time.Now()
	helper.ExitHandle()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	asuka()
}

func asuka() {
	fmt.Println("Monitor: http://127.0.0.1:666")

	log.Println(web.Server([]*project.Dispatcher{
		project.New(&project.DouBan{}).Run(),
		//project.New(&project.Pixiv{}).Run(),
		//project.New(&project.Test{}).Run(),
		//project.New(&project.Test2{}).Run(),
		//project.New(&project.ZhiHu{}).Run(),
		//project.New(&project.JianShu{}).Run(),
		//project.New(&project.Www{}).Run(),

		//project.New(&project.Pixiv{}).CleanUp().Run(),
		//project.New(&project.DouBan{}).CleanUp().Run(),
		//project.New(&project.Test2{}).CleanUp().Run(),
		//project.New(&project.ZhiHu{}).CleanUp().Run(),
		//project.New(&project.JianShu{}).CleanUp().Run(),
		//project.New(&project.Www{}).CleanUp().Run(),
	}, ":666")) // http://127.0.0.1:666
	helper.ExitHandleFunc()
}
