package main

import (
	"fmt"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/project"
	"github.com/chenset/asuka/web"
	"log"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	mainStart := time.Now()
	helper.ExitHandle()
	defer func() {
		fmt.Println("Done: ", time.Since(mainStart))
	}()

	//todo url filter & queue按需初始化或释放
	//todo url filter支持远程
	//todo mysql/oracle 数据库写入支持
	//todo chrome 插件 => pixiv.net js自动抓取报送
	//todo douban web page.

	asuka()
}

func asuka() {
	fmt.Println("http://127.0.0.1:" + strings.Split(helper.Env().WEBListen, ":")[len(strings.Split(helper.Env().WEBListen, ":"))-1])
	log.Println(web.Server([]*project.Dispatcher{
		project.New(&project.DouBan{}, time.Time{}).Run(),
		project.New(&project.Pixiv{}, time.Time{}).Run(),
		//project.New(&project.Test2{}, time.Now()).Run(),
		//project.New(&project.ZhiHu{}, time.Now()).Run(),
		//project.New(&project.JianShu{}, time.Now()).Run(),
		//project.New(&project.Www{}, time.Now()).Run(),
		//
		//project.New(&project.Pixiv{}).CleanUp().Run(),
		//project.New(&project.DouBan{}).CleanUp().Run(),
		project.New(&project.Test{}, time.Now()).CleanUp().Run(),
		project.New(&project.Test2{}, time.Now().Add(time.Minute*10)).CleanUp().Run(),
		project.New(&project.ZhiHu{}, time.Now().Add(time.Minute*15)).CleanUp().Run(),
		project.New(&project.JianShu{}, time.Now().Add(time.Minute*5)).CleanUp().Run(),
		project.New(&project.Www{}, time.Now().Add(time.Minute*20)).CleanUp().Run(),
	}, helper.Env().WEBListen))
	helper.ExitHandleFunc()
}
