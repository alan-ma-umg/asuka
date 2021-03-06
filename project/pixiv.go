package project

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Illusts struct {
	IllustId string
	Url      string
	UserId   string
	Width    int
	Height   int
}

type Pixiv struct {
	*Implement
	lastRequestUrl    string
	dbSpeed           int
	dbSpeedNum        int
	lastInsertId      int64
	lastInsertError   string
	showingString     string
	lastHttpCodeIs404 bool
}

func (my *Pixiv) Name() string {
	return "Kumiko"
}

func (my *Pixiv) Showing() (str string) {
	return "<a href=\"/website/" + my.Name() + "\">" + my.showingString + "</a>"
}

func (my *Pixiv) Init(d *Dispatcher) {
	my.showingString = "快, 我要营养快线 !!!"

	urlConvertRegex := regexp.MustCompile(`(?i)/img/[^.]+(_master1200|_square1200|_custom1200)(\.(jpg|png|jpeg|git|webp))`)

	//images http index , todo 以后要删除掉 !!!!!!!!!!!!!!!
	http.Handle("/project/pixiv/images/", http.StripPrefix("/project/pixiv/images/", http.FileServer(http.Dir("project/pixiv"))))

	http.HandleFunc("/project/pixiv/crawl/upload", func(w http.ResponseWriter, r *http.Request) {

		//fixme CORS policy 理解一下 !!!!!!!!!!!!!!!!!!!!!!!!!!!!
		w.Header().Set("access-control-allow-methods", "POST")
		w.Header().Set("access-control-allow-origin", "https://www.pixiv.net")

		if r.Method != "POST" {
			http.Error(w, "POST Required", 405)
			return
		}

		var post []*Illusts
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&post)
		if err != nil {
			my.showingString = time.Now().Format("2006-01-02 15:04:05") + err.Error()
			http.Error(w, "Bad Request", 400)
			return
		}

		addCount := 0
		for _, item := range post {

			rawUrl := item.Url

			if re := urlConvertRegex.FindStringSubmatch(item.Url); len(re) >= 4 {
				rawUrl = "https://i.pximg.net/img-original" + strings.TrimSuffix(re[0], re[1]+re[2]) + re[2]
			} else {
				//send message to wx
				helper.SendTextToWXDoOnceDurationHour("url convert failed: https://px.flysay.com/" + item.Url + " \nIllustId: " + item.IllustId)
			}

			if exists, _ := d.GetQueue().BlTestAndAddString(rawUrl); exists {
				continue
			}
			addCount++
			d.GetQueue().Enqueue(rawUrl)
		}

		my.showingString = time.Now().Format("2006-01-02 15:04:05") + " upload succeed , len: " + strconv.Itoa(len(post)) + " added: " + strconv.Itoa(addCount)

		w.Header().Set("Content-type", "application/json")
		jsonMap := map[string]interface{}{}
		jsonMap["success"] = true
		b, _ := json.Marshal(jsonMap)
		w.Write(b)
	})
}

func (my *Pixiv) EntryUrl() []string {
	return nil
}

// frequency
func (my *Pixiv) Throttle(spider *spider.Spider) {
	if spider.LoadRate(5) > 5.0 {
		spider.AddSleep(120e9)
	}

	//if my.lastHttpCodeIs404 { // 404 is normal
	spider.ResetSleep() //无视spider的限制
	//}

	//spider.AddSleep(time.Duration(rand.Float64() * 10e9))

	if spider.FailureLevel > 40 {
		spider.Delete = true
	}
}

func (my *Pixiv) RequestBefore(spider *spider.Spider) {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	spider.CurrentRequest().Header.Set("Referer", "https://www.pixiv.net/artworks/"+strconv.Itoa(rand.Intn(71245413)))
	spider.SetRequestTimeout(time.Minute * 10)
}

func (my *Pixiv) ResponseAfter(spider *spider.Spider) {
	//下载或者耗时过长的删除掉, 无论成功失败都删除
	if !spider.RequestStartTime.IsZero() && time.Since(spider.RequestStartTime).Seconds() > 180 {
		spider.Delete = true
		spider.FailureLevel = 100
	}

	if spider.CurrentResponse() == nil {
		spider.Delete = true
	}

	if spider.CurrentResponse() != nil && spider.CurrentResponse().StatusCode != 404 && spider.CurrentResponse().StatusCode != 200 {
		spider.Delete = true
	}

	//} else if spider.CurrentResponse().StatusCode == 404 && spider.FailureLevel <= 10 { //404 is normal
	//	spider.FailureLevel = 0

	if spider.CurrentResponse() != nil && spider.CurrentResponse().StatusCode == 404 {
		my.lastHttpCodeIs404 = true
	} else {
		my.lastHttpCodeIs404 = false
	}

	my.Implement.ResponseAfter(spider)
}

// EnqueueForFailure 请求或者响应失败时重新入失败队列, 可以修改这里修改加入失败队列的实现
func (my *Pixiv) EnqueueForFailure(spider *spider.Spider, err error, retryEnqueueUrl, spiderEnqueueUr string, retryTimes int) (success bool, tries int) {

	//没有响应直接入正常的队列, fixme 如果一直没有响应意味着会无限下去
	if spider.CurrentResponse() == nil || spider.CurrentResponse().StatusCode == 0 {
		spider.GetQueue().Enqueue(spiderEnqueueUr)
		return
	}

	//响应状态200,但是读取body失败. 这种情况一般时代理超时/错误之类的情况直接无限重试下去
	if spider.CurrentResponse().StatusCode == 200 && err != nil && strings.Contains(spider.CurrentResponse().Header.Get("Content-type"), "image") {
		spider.GetQueue().Enqueue(spiderEnqueueUr)
		return
	}

	//404丢弃原链接,Retries.F不会增加.插入新格式的链接
	if spider.CurrentResponse().StatusCode == 404 {
		newUrl := regexp.MustCompile(`(?i)\.[^\.]{2,5}$`).ReplaceAllString(spiderEnqueueUr, ".png")
		if exists, _ := spider.GetQueue().BlTestAndAddString(newUrl); exists {
			return
		}
		spider.GetQueue().Enqueue(newUrl)
		return
	}

	//常规加入失败队列
	return my.Implement.EnqueueForFailure(spider, err, retryEnqueueUrl, spiderEnqueueUr, retryTimes)
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *Pixiv) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {

	if response.StatusCode == 403 {
		helper.SendTextToWXDoOnceDurationHour(response.Status + ": " + spider.TransportUrl.Host + " => " + spider.CurrentRequest().URL.String())
	}

	if response.StatusCode != 403 && response.StatusCode != 404 && !strings.Contains(response.Header.Get("Content-type"), "image") {
		helper.SendTextToWXDoOnceDurationHour("not image: got " + response.Header.Get("Content-type") + " => " + response.Status + ": " + spider.TransportUrl.Host + " => " + spider.CurrentRequest().URL.String())
		return false, errors.New("not image")
	}
	return true, nil
}

// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
// 一般用于采集数据的地方
func (my *Pixiv) ResponseSuccess(spider *spider.Spider) {
	//sha1
	h := sha1.New()
	h.Write(spider.ResponseByte)

	filePath := "project/pixiv/"
	filename := filePath + hex.EncodeToString(h.Sum(nil)) + filepath.Ext(spider.CurrentRequest().URL.String())

	if _, err := os.Stat(filename); os.IsExist(err) {
		return
	}

	//create dir
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		my.showingString = time.Now().Format("2006-01-02 15:04:05") + err.Error()
		log.Println(err)
		return
	}

	//write file
	if err := ioutil.WriteFile(filename, spider.ResponseByte, 0); err != nil {
		my.showingString = time.Now().Format("2006-01-02 15:04:05") + err.Error()
		log.Println(err)
		return
	}

	//f, err := os.Create("project/pixiv/" + hex.EncodeToString(h.Sum(nil)) + filepath.Ext(spider.CurrentRequest().URL.String()))
	//if err != nil {
	//	my.showingString = time.Now().Format("2006-01-02 15:04:05") + err.Error()
	//	log.Println(err)
	//	return
	//}
	//defer f.Close()
	//_, err = f.Write(spider.ResponseByte)
	//if err != nil {
	//	my.showingString = time.Now().Format("2006-01-02 15:04:05") + err.Error()
	//	log.Println(err)
	//	return
	//}

}

// queue
func (my *Pixiv) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	return l.String()
}

func (my *Pixiv) WEBSite(w http.ResponseWriter, r *http.Request) {
	files, _ := filepath.Glob("project/pixiv/*")

	helper.GetTemplates().ExecuteTemplate(w, "pixiv.html", struct {
		ProjectName string
		Files       []string
	}{
		ProjectName: my.Name(),
		Files:       files,
	})
}
