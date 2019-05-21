package main

import (
	"flag"
	"fmt"
	"github.com/chenset/asuka/helper"
	"log"
	"net/url"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//helper.PathToEnvFile = os.Args[1]

	//C:\data\codes\asuka\env.json

	redis := flag.String("redis", "redis://:password@127.0.0.1:3306/10?timeout=0", "Redis connection url, default: redis://:password@127.0.0.1:6379/10?timeout=0")
	mysql := flag.String("mysql", "dns://root:11111111@(127.0.0.1:3306)/asuka?charset=utf8mb4", "Mysql DSN, default: dns://root:11111111@(127.0.0.1:3306)/asuka?charset=utf8mb4")
	webPassword := flag.String("webPassword", "", "WEB login password, default random")
	bloomFilterPath := flag.String("bloomFilterPath", ".", "BloomFilter save path, default:.")
	localTransport := flag.Bool("localTransport", false, "Enable http.DefaultTransport, default: false")

	flag.Parse()
	//if *webPassword == "" {
	//	rand.Seed(time.Now().Unix())
	//	flag.Set("webPassword", strconv.Itoa(rand.Int()))
	//}
	helper.GlobalEnvConfig = &helper.EnvConfig{
		BloomFilterPath: *bloomFilterPath,
		WEBPassword:     *webPassword,
		LocalTransport:  *localTransport,
		MysqlDSN:        *mysql,
		Redis: helper.Redis{
			Network:     "tcp",
			Addr:        "127.0.0.1",
			Password:    "",
			DB:          10,
			URLQueueKey: "url_queue_key",
		},
	}
	//helper.GlobalEnvConfig.BloomFilterPath = *bloomFilterPath
	//helper.GlobalEnvConfig.WEBPassword = *webPassword
	//helper.GlobalEnvConfig.LocalTransport = *localTransport
	//helper.GlobalEnvConfig.MysqlDSN = *mysql
	//
	u, err := url.Parse(*redis)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(u.Scheme)

	//Network     string
	//Addr        string
	//Password    string
	//DB          int
	//URLQueueKey string

	//log.Println(u.Host)
	//log.Println(u.Hostname())
	//log.Println(u.User.Password())
	//log.Println(u.User.Username())

	//helper.GlobalEnvConfig.Redis = *redis

	//log.Println(*redis)
	//log.Println(*mysql)
	//log.Println(*webPassword)
	//log.Println(*bloomFilterPath)
	//log.Println(*localTransport)
}

func main() {
	mainStart := time.Now()
	//helper.ExitHandle()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	//[: password@] host [: port] [/ database][? [timeout=timeout[d|h|m|s|ms|us|ns]]

	//l := "redis://:password@127.0.0.1:3306/10?timeout=0"

	//u, err := url.Parse(l)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//log.Println(u.Host)
	//log.Println(u.Hostname())
	//log.Println(u.User.Password())
	//log.Println(u.User.Username())

	//asuka()
}

//func asuka() {
//	fmt.Println("Monitor: http://127.0.0.1:666")
//
//	log.Println(web.Server([]*project.Dispatcher{
//		project.New(&project.DouBan{}).Run(),
//		//project.New(&project.Pixiv{}).Run(),
//		//project.New(&project.Test{}).Run(),
//		//project.New(&project.Test2{}).Run(),
//		//project.New(&project.ZhiHu{}).Run(),
//		//project.New(&project.JianShu{}).Run(),
//		//project.New(&project.Www{}).Run(),
//
//		//project.New(&project.Pixiv{}).CleanUp().Run(),
//		//project.New(&project.DouBan{}).CleanUp().Run(),
//		//project.New(&project.Test2{}).CleanUp().Run(),
//		//project.New(&project.ZhiHu{}).CleanUp().Run(),
//		//project.New(&project.JianShu{}).CleanUp().Run(),
//		//project.New(&project.Www{}).CleanUp().Run(),
//	}, ":666")) // http://127.0.0.1:666
//	helper.ExitHandleFunc()
//}

//2019/01/24 11:56:11 h2_bundle.go:8723: protocol error: received *http.http2GoAwayFrame before a SETTINGS frame
