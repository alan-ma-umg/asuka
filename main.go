package main

import (
	"fmt"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/project"
	"github.com/chenset/asuka/queue"
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
	//todo 文件日记展示系统
	//todo chrome 插件 => pixiv.net js自动抓取报送
	//todo douban web page.
	//todo loading css效果
	//todo 开启swap
	//
	//list := []string{"https://book.douban.com/tag/%E8%87%AA%E5%8A%A9%E6%B8%B8", "https://book.douban.com/tag/%E8%87%AA%E5%8A%A9%E6%B8%B8"}
	//
	//for i := 0; i < 1000; i++ {
	//	list = append(list, "fsdjflskjdflkdf")
	//}

	//BloomFilterServer
	if helper.Env().BloomFilterServer != "" {
		go queue.GetTcpFilterInstance().ServerListen(helper.Env().BloomFilterServer)
	}

	asuka()

	//
	//time.Sleep(1e9)
	//var res []byte
	//
	//buf, err := queue.GetTcpFilterInstance().Cmd(10, &queue.Cmd10{
	//	Db:   "ccc",
	//	Fun:  10,
	//	Urls: []string{"https://book.douban.com/ta1g/%E8%87%AA%E5%8A%A9%E6%B8%B8", "https://book.douban.com/tag/%E8%87%AA%E5%8A%A9%E6%B8%B8"},
	//})
	//
	//json.Unmarshal(buf, &res)
	//log.Println(res)
	//
	//buf, err = queue.GetTcpFilterInstance().Cmd(10, &queue.Cmd10{
	//	Db:   "ccc",
	//	Fun:  20,
	//	Urls: []string{"https://book.douban.com/tag/%E8%87%AA%E5%8A%A9%E6%B8%B8", "https://book.douban.com/tag/%E8%87%AA%E5%8A%A9%E6%B8%B8"},
	//})
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//json.Unmarshal(buf, &res)
	//log.Println(res)

	//
	//s := time.Now()
	//for i := 0; i < 1000; i++ {
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//
	//	var res []byte
	//	json.Unmarshal(buf, &res)
	//}
	//
	//log.Println(time.Since(s))
	//log.Println(res)

	//time.Sleep(1e9)

	//s := time.Now()

	//length := 10000
	//for i := 0; i < length; i++ {

	//l := 0
	////
	//for ii := 0; ii < 10; ii++ {
	//	go func(i int) {
	//		time.Sleep(time.Duration(rand.Float64() * 1e9))
	//		for {
	//			queue.GetTcpFilterInstance().ClientBl(byte(i), 20, list)
	//			l++
	//		}
	//	}(ii)
	//}
	//}
	//log.Println(time.Since(s))

	//log.Println(queue.GetTcpFilterInstance().ClientOtherCmd(21, map[string]interface{}{"haha": "fsdfsdf"}))

	//time.Sleep(2e9)
	//log.Println(l)
	////
	//log.Println(queue.GetTcpFilterInstance().ServerHandleCount)
	////
	//time.Sleep(100e9)

	//
	//buf[lenOfData] = 1                  //cmd
	//buf[lenOfData+lenOfCmd] = 2         //Db
	//buf[lenOfData+lenOfCmd+lenOfDb] = 3 //fun
	//
	//jsonBytes, err := json.Marshal(list)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//copy(buf[lenOfData+lenOfCmd+lenOfDb+lenOfFun:], jsonBytes[:])
	//
	//dataLen := uint16(len(jsonBytes) + lenOfCmd + lenOfDb + lenOfFun)
	//binary.BigEndian.PutUint16(buf[:lenOfData], dataLen)
	//
	////34 93
	//
	////log.Println()
	//
	//go FilterServer()
	//time.Sleep(1e9)
	//FilterClient(buf[:dataLen+lenOfData])
	//time.Sleep(1e9)
	////bytes.NewBufferString(jsonStr)
	////copy(,)
	//
	////encoder := json.NewEncoder(bytes.NewBuffer(buf[:lenOfData+lenOfCmd+lenOfDb+lenOfFun]))
	////if err := encoder.Encode(list); err != nil {
	////	log.Println(err)
	////	return
	////}
	//
	////log.Println(buf)
	//
	////buf[lenOfData+lenOfCmd+lenOfDb+lenOfFun:], _ = json.Marshal(list)
	//
	////log.Println(len(buf))
	//
	////binary.BigEndian.PutUint16(buf[:lenOfData], uint16(len(buf)))
	//
	////var i uint16 = 1
	////log.Println(unsafe.Sizeof(i))
	////log.Println(buf)
	//
	////payload size
	////binary.BigEndian.PutUint16(writeBuf[headLength-PayloadSizeLength:headLength], uint16(len(writeBuf[currentHeadLength:writeBufLen])))
	////log.Println(binary.BigEndian.Uint16(buf[:lenOfData]))
	////log.Println(len(i))
	////log.Println(unsafe.Sizeof(i))
	//
	////log.Println(len([]byte("https://book.douban.com/tag/%E8%87%AA%E5%8A%A9%E6%B8%B8")))
	////go FilterServer()
	//
	////FilterClient()
}

func asuka() {
	fmt.Println("http://127.0.0.1:" + strings.Split(helper.Env().WEBListen, ":")[len(strings.Split(helper.Env().WEBListen, ":"))-1])
	log.Println(web.Server([]*project.Dispatcher{
		//project.New(&project.DouBan{}, time.Time{}).Run(),
		//project.New(&project.Pixiv{}, time.Time{}).Run(),
		project.New(&project.Test2{}, time.Time{}).Run(),
		//project.New(&project.ZhiHu{}, time.Now()).Run(),
		//project.New(&project.JianShu{}, time.Now()).Run(),
		//project.New(&project.Www{}, time.Now()).Run(),
		//
		//project.New(&project.Pixiv{}).CleanUp().Run(),
		//project.New(&project.DouBan{}).CleanUp().Run(),
		//project.New(&project.Test{}, time.Now()).CleanUp().Run(),
		//project.New(&project.Test2{}, time.Now().Add(time.Minute*10)).CleanUp().Run(),
		//project.New(&project.ZhiHu{}, time.Now().Add(time.Minute*15)).CleanUp().Run(),
		//project.New(&project.JianShu{}, time.Now().Add(time.Minute*5)).CleanUp().Run(),
		//project.New(&project.Www{}, time.Now().Add(time.Minute*20)).CleanUp().Run(),
	}, helper.Env().WEBListen))
	helper.ExitHandleFunc()
}
