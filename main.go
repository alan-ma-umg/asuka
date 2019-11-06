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

	//todo systemd
	//todo MUX
	//todo servers/crawling 切换使用ajax
	//todo tcp filter 支持压缩加密多个链接查询
	//todo TCP filter 链接池空闲释放
	//todo retires 重试成功后删除,  douban 的 url比较特殊
	//todo WEB输入频率限2制
	//todo chrome 插件 => pixiv.net js自动抓取报送
	//todo douban web page.

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
		project.New(&project.V2ex{}, time.Time{}).Run(),
		//project.New(&project.Test2{}, time.Time{}).Run(),
		//project.New(&project.ZhiHu{}, time.Now()).Run(),
		//project.New(&project.JianShu{}, time.Now()).Run(),
		//project.New(&project.Www{}, time.Now()).Run(),
		//project.New(&project.Forever{}, time.Now()).CleanUp().Run(),
		//project.New(&project.Test{}, time.Now()).CleanUp().Run(),
		//project.New(&project.Pixiv{}).CleanUp().Run(),
		//project.New(&project.DouBan{}).CleanUp().Run(),
		//project.New(&project.ZhiHu{}, time.Now()).CleanUp().Run(),
		//project.New(&project.JianShu{}, time.Now()).CleanUp().Run(),
		//project.New(&project.Www{}, time.Now().Add(time.Minute*20)).CleanUp().Run(),
		project.New(&project.Death{}, time.Now()).CleanUp().Run(),
		project.New(&project.JS{}, time.Now()).CleanUp().Run(),
	}, helper.Env().WEBListen))
	helper.ExitHandleFunc()
}
