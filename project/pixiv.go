package project

import (
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Pixiv struct {
	*Implement
	lastRequestUrl string
	queueUrlLen    int64
}

func (my *Pixiv) Name() string {
	return "Kumiko"
}

func (my *Pixiv) Init() {
	go func() {
		t := time.NewTicker(time.Second * 5)
		for {
			<-t.C
			my.queueUrlLen, _ = database.Redis().LLen(strings.Split(reflect.TypeOf(my).String(), ".")[1] + "_" + helper.Env().Redis.URLQueueKey).Result()
		}
	}()

}
func (my *Pixiv) EntryUrl() []string {
	return []string{
		"https://www.pixiv.net/tags.php?tag=%E7%BE%8E%E3%81%97%E3%81%84",
		"https://www.pixiv.net/member_illust.php?mode=medium&illust_id=63093148",
		"https://www.pixiv.net/tags.php?tag=%E7%BE%8E%E3%81%97%E3%81%84",
		"https://www.pixiv.net/member_illust.php?mode=medium&illust_id=63093148",
		"https://www.pixiv.net/tags.php?tag=%E7%BE%8E%E3%81%97%E3%81%84",
		"https://www.pixiv.net/member_illust.php?mode=medium&illust_id=63093148",
		"https://www.pixiv.net/tags.php?tag=%E7%BE%8E%E3%81%97%E3%81%84",
		"https://www.pixiv.net/member_illust.php?mode=medium&illust_id=63093148",
		"https://www.pixiv.net/tags.php?tag=%E7%BE%8E%E3%81%97%E3%81%84",
		"https://www.pixiv.net/member_illust.php?mode=medium&illust_id=63093148",
		"https://www.pixiv.net/tags.php?tag=%E7%BE%8E%E3%81%97%E3%81%84",
		"https://www.pixiv.net/member_illust.php?mode=medium&illust_id=63093148",
		"https://www.pixiv.net/tags.php?tag=%E7%BE%8E%E3%81%97%E3%81%84",
		"https://www.pixiv.net/member_illust.php?mode=medium&illust_id=63093148",
	}
}

// frequency
func (my *Pixiv) Throttle(spider *spider.Spider) {
	if spider.Transport.LoadRate(5) > 5.0 {
		spider.AddSleep(120e9)
	}

	spider.AddSleep(time.Duration(rand.Float64() * 40e9))

	if spider.FailureLevel > 1 {
		spider.ResetRequest()
		spider.Transport.Close()
	}
}

func (my *Pixiv) RequestBefore(spider *spider.Spider) {
	//accept
	if spider.CurrentRequest() != nil {
		spider.CurrentRequest().Header.Set("Accept", "text/html")
	}

	//Referer
	if spider.CurrentRequest() != nil && spider.CurrentRequest().Referer() == "" && my.lastRequestUrl != "" {
		spider.CurrentRequest().Header.Set("Referer", my.lastRequestUrl)
	}

	spider.Client().Timeout = 20 * time.Second
}

func (my *Pixiv) ResponseAfter(spider *spider.Spider) {
	spider.ResetRequest()
	spider.Transport.Close()

	spider.ResponseByte = []byte{} //free memory
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *Pixiv) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "text/html") {
		return false, nil
	}
	if strings.ToLower(response.Header.Get("Content-Encoding")) != "gzip" {
		return false, nil
	}
	return true, nil
}

// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
// 一般用于采集数据的地方
func (my *Pixiv) ResponseSuccess(spider *spider.Spider) {
}

// queue
func (my *Pixiv) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	//if !strings.HasPrefix(strings.ToLower(l.String()), "https://movie.douban.com/subject") && !strings.HasPrefix(strings.ToLower(l.String()), "https://book.douban.com/subject") && !strings.HasPrefix(strings.ToLower(l.String()), "https://book.douban.com/tag") && !strings.HasPrefix(strings.ToLower(l.String()), "https://movie.douban.com/tag") {
	//	return
	//}
	//
	//if strings.HasPrefix(strings.ToLower(l.String()), "https://book.douban.com/subject") && !isDouBanSubject(strings.ToLower(l.String())) {
	//	return
	//}
	//
	//if strings.HasPrefix(strings.ToLower(l.String()), "https://movie.douban.com/subject") && !isDouBanSubject(strings.ToLower(l.String())) {
	//	return
	//}

	if my.queueUrlLen > 20000 {
		return
	}

	if !strings.HasPrefix(strings.ToLower(l.String()), "https://www.pixiv.net/") {
		return
	}

	return l.String()
	//return l.Scheme + "://" + l.Host + l.Path
}
