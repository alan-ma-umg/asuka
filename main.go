package main

import (
	"fmt"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/project"
	"github.com/chenset/asuka/queue"
	"github.com/chenset/asuka/web"
	"log"
	"math/rand"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	rand.Seed(time.Now().UnixNano()) //global effect
}

func main() {
	mainStart := time.Now()
	helper.ExitHandle()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	//todo websocket to Server-Sent Events

	//BloomFilterServer
	if helper.Env().BloomFilterServer != "" && helper.Env().BloomFilterClient != "" {
		go func() {
			queue.GetTcpFilterInstance().ServerListen(helper.Env().BloomFilterServer)
		}()
		asuka()
	} else if helper.Env().BloomFilterServer != "" {
		queue.GetTcpFilterInstance().ServerListen(helper.Env().BloomFilterServer)
	} else {
		asuka()
	}
}

func asuka() {
	fmt.Println("http://127.0.0.1:" + strings.Split(helper.Env().WEBListen, ":")[len(strings.Split(helper.Env().WEBListen, ":"))-1])
	log.Println(web.Server([]*project.Dispatcher{
		project.New(&project.DouBan{}, time.Now()).Run(),
		project.New(&project.Pixiv{}, time.Time{}).Run(),

		project.New(&project.V2ex{
			SpeedShowing: &project.SpeedShowing{},
		}, time.Time{}).Run(),
		//project.New(&project.Test2{}, time.Time{}).Run(),
		//project.New(&project.ZhiHu{}, time.Now()).Run(),
		//project.New(&project.JianShu{}, time.Now()).Run(),
		//project.New(&project.Www{}, time.Now()).Run(),
		//project.New(&project.Pixiv{}).CleanUp().Run(),
		project.New(&project.CDN{
			SpiderThrottle: &project.SpiderThrottle{.001},
			SpeedShowing:   &project.SpeedShowing{},
		}, time.Time{}).Run(),
		//project.New(&project.DouBan{}).CleanUp().Run(),
		//project.New(&project.ZhiHu{}, time.Now()).CleanUp().Run(),
		//project.New(&project.JianShu{}, time.Now()).CleanUp().Run(),
		//project.New(&project.Www{}, time.Now().Add(time.Minute*20)).CleanUp().Run(),
		project.New(&project.Test{
			SpiderThrottle: &project.SpiderThrottle{.1},
			SpeedShowing:   &project.SpeedShowing{},
		}, time.Now()).CleanUp().Run(),

		project.New(&project.Forever{
			SpiderThrottle: &project.SpiderThrottle{.1},
			SpeedShowing:   &project.SpeedShowing{},
		}, time.Now()).CleanUp().Run(),

		project.New(&project.JS{
			SpiderThrottle: &project.SpiderThrottle{.01},
			SpeedShowing:   &project.SpeedShowing{},
		}, time.Now()).CleanUp().Run(),

		//project.New(&project.Death{}, time.Now()).Run(),
		project.New(&project.Death{
			SpiderThrottle: &project.SpiderThrottle{1},
			SpeedShowing:   &project.SpeedShowing{},
			Setting:        &project.Setting{&project.DeathSettingOption{}},
		}, time.Now()).CleanUp().Run(),
	}, helper.Env().WEBListen))
	helper.ExitHandleFunc()
}
