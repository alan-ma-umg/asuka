package project

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"io"
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
	lastRequestUrl  string
	dbSpeed         int
	dbSpeedNum      int
	lastInsertId    int64
	lastInsertError string
	showingString   string
}

func (my *Pixiv) Name() string {
	return "Kumiko"
}

func (my *Pixiv) Showing() (str string) {
	return "<a href=\"/download/" + my.Name() + "\">" + my.showingString + "</a>"
}

func (my *Pixiv) Init(d *Dispatcher) {
	my.showingString = "快, 我要营养快线 !!!"

	urlConvertRegex := regexp.MustCompile(`(?i)/img/[^.]+(_master1200|_square1200|_custom1200)(\.(jpg|png|jpeg|git|webp))`)

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

			if d.queue.BlTestAndAddString(rawUrl) {
				continue
			}
			addCount++
			d.queue.Enqueue(rawUrl)
		}

		my.showingString = time.Now().Format("2006-01-02 15:04:05") + " upload succeed , len: " + strconv.Itoa(len(post)) + " added: " + strconv.Itoa(addCount)

		w.Header().Set("Content-type", "application/json")
		jsonMap := map[string]interface{}{}
		jsonMap["error"] = false
		b, _ := json.Marshal(jsonMap)
		w.Write(b)
	})
}

func (my *Pixiv) EntryUrl() []string {
	return nil
}

// frequency
func (my *Pixiv) Throttle(spider *spider.Spider) {
	if spider.Transport.LoadRate(5) > 5.0 {
		spider.AddSleep(120e9)
	}

	spider.AddSleep(time.Duration(rand.Float64() * 10e9))

	if spider.FailureLevel > 40 {
		spider.Delete = true
	}
}

func (my *Pixiv) RequestBefore(spider *spider.Spider) {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	spider.CurrentRequest().Header.Set("Referer", "https://www.pixiv.net/artworks/"+strconv.Itoa(rand.Intn(71245413)))
	spider.Client().Timeout = time.Minute * 10
}

func (my *Pixiv) ResponseAfter(spider *spider.Spider) {
	//spider.ResetRequest() //todo !!!!!!!!!!!!!!! ???????????
	//spider.Transport.Close() //todo !!!!!!!!!!!!!!! ???????????

	my.Implement.ResponseAfter(spider)
}

// EnqueueForFailure 请求或者响应失败时重新入失败队列, 可以修改这里修改加入失败队列的实现
func (my *Pixiv) EnqueueForFailure(spider *spider.Spider, err error, rawUrl string, retryTimes int) {

	//没有响应直接入正常的队列, fixme 如果一直没有响应意味着会无限下去
	if spider.CurrentResponse() == nil || spider.CurrentResponse().StatusCode == 0 {
		spider.Queue.Enqueue(rawUrl)
		return
	}

	//响应状态200,但是读取body失败. 这种情况一般时代理超时/错误之类的情况直接无限重试下去
	if spider.CurrentResponse().StatusCode == 200 && err != nil && strings.Contains(spider.CurrentResponse().Header.Get("Content-type"), "image") {
		spider.Queue.Enqueue(rawUrl)
		return
	}

	//404丢弃原链接,Retries.F不会增加.插入新格式的链接
	if spider.CurrentResponse().StatusCode == 404 {
		newUrl := regexp.MustCompile(`(?i)\.[^\.]{2,5}$`).ReplaceAllString(rawUrl, ".png")
		if spider.Queue.BlTestAndAddString(newUrl) {
			return
		}
		spider.Queue.Enqueue(newUrl)
		return
	}

	//常规加入失败队列
	my.Implement.EnqueueForFailure(spider, err, rawUrl, retryTimes)
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *Pixiv) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {

	if response.StatusCode == 403 {
		helper.SendTextToWXDoOnceDurationHour(response.Status + ": " + spider.Transport.S.Host + " => " + spider.CurrentRequest().URL.String())
	}

	if response.StatusCode != 403 && response.StatusCode != 404 && !strings.Contains(response.Header.Get("Content-type"), "image") {
		helper.SendTextToWXDoOnceDurationHour("not image: got " + response.Header.Get("Content-type") + " => " + response.Status + ": " + spider.Transport.S.Host + " => " + spider.CurrentRequest().URL.String())
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

	//create dir
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		my.showingString = time.Now().Format("2006-01-02 15:04:05") + err.Error()
		log.Println(err)
		return
	}

	//write file
	if err := ioutil.WriteFile(filePath+hex.EncodeToString(h.Sum(nil))+filepath.Ext(spider.CurrentRequest().URL.String()), spider.ResponseByte, 0); err != nil {
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

func (my *Pixiv) HttpExportResult(w http.ResponseWriter, r *http.Request) {
	htmlStr := `<pre>
    function ajaxDo_(option) {
        let url = option.url || '',
            data = option.data,
            method = option.method || 'get',
            headers = option.headers || {},
            success = option.success,
            timeout = option.timeout || 10000,
            error = option.error;

        let xhr = new XMLHttpRequest();
        xhr.timeout = timeout;
        xhr.onreadystatechange = function () {
            if (this.readyState === 4) {
                if (this.status === 200 || this.status === 304) {
                    success && success(this);
                } else {
                    error && error(this);
                }
            }
        };
        xhr.open(method, url, true);
        for (let k in headers) {
            xhr.setRequestHeader(k, headers[k]);
        }
        if (data) {
            xhr.send(data);
        } else {
            xhr.send();
        }
    }

    (function () {
        let bottomTimes = 0;
        let fsdfsdfgdfg = setInterval(function () {
            if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
                if (bottomTimes++ > 300) {
                    clearInterval(fsdfsdfgdfg);
                    window.scrollTo(0, 0); //to top
                    hentaiStart();
                }
            } else {
                window.scrollTo(0, document.body.scrollHeight); //to bottom
                bottomTimes = 0;
            }
        }, 10);

        function hentaiStart() {
            let postJson = [];
            let urlParams = new URLSearchParams(window.location.search);
            let IllustId = urlParams.get('id');
            Array.from(document.body.querySelectorAll('img[src*="1200.jpg"]')).forEach(function (element) {
                postJson.push({
                    url: element.src,
                    illustId: IllustId,
                });
            });

            postJson.length && ajaxDo_({
                method: "POST",
                url: "http://127.0.0.1:666/project/pixiv/crawl/upload",
                data: JSON.stringify(postJson),
            });
        }
    })()
</pre>`
	w.Header().Set("Content-type", "text/html; charset=UTF-8")
	io.WriteString(w, htmlStr)
}
