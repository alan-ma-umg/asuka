package project

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	return "<a href=\"/download/\"" + my.Name() + ">" + my.showingString + "</a>"
}

func (my *Pixiv) Init(d *Dispatcher) {
	my.showingString = "快, 我要营养快线 !!!"

	urlConvertRegex := regexp.MustCompile(`(?i)/img/[^.]+(_master1200|_square1200|_custom1200)(\.(jpg|png|jpeg|git))`)

	//http://127.0.0.1:666/project/pixiv/crawl/upload
	fmt.Println("pixiv upload server: http://127.0.0.1:666/project/pixiv/crawl/upload")

	http.HandleFunc("/project/pixiv/crawl/upload", func(w http.ResponseWriter, r *http.Request) {

		//fixme CORS policy 理解一下 !!!!!!!!!!!!!!!!!!!!!!!!!!!!
		w.Header().Set("access-control-allow-methods", "POST")
		w.Header().Set("access-control-allow-origin", "https://www.pixiv.net")

		if r.Method != "POST" {
			http.Error(w, "POST Required", 405)
			return
		}

		//https://px.flysay.com/https://i.pximg.net/img-original/img/2019/10/12/19/46/40/77248012_p0.jpg
		//https://px.flysay.com/https://i.pximg.net/img-original/img/2019/10/12/19/46/40/77248012_p0.png // todo png !!  retry
		//https://px.flysay.com/https://i.pximg.net/img-original/img/2019/10/12/19/46/40/77248012_p0.jpg

		var post []*Illusts
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&post)
		if err != nil {
			my.showingString = time.Now().Format("2006-01-02 15:04:05") + err.Error()
			http.Error(w, "Bad Request", 400)
			return
		}

		my.showingString = time.Now().Format("2006-01-02 15:04:05") + " upload succeed , len: " + strconv.Itoa(len(post))

		for _, item := range post {

			rawUrl := item.Url

			//todo if 404 need to try again with .png
			//"https://i.pximg.net/c/360x360_70/img-master/img/2019/06/14/09/00/01/75214268_p0_square1200.jpg"
			if re := urlConvertRegex.FindStringSubmatch(item.Url); len(re) >= 4 {
				rawUrl = "https://i.pximg.net/img-original" + strings.TrimSuffix(re[0], re[1]+re[2]) + re[2]
			} else {
				//send message to wx
				go func() {
					if helper.Env().WechatSendMessagePassword != "" {
						helper.DoOnceDurationHour(func() {
							http.Get("https://wx.flysay.com/send?password=" + helper.Env().WechatSendMessagePassword + "&touser=chen&content=" + url.QueryEscape("url convert failed: https://px.flysay.com/"+item.Url+" \nIllustId: "+item.IllustId))
						})
					}
				}()
			}

			if d.queue.BlTestAndAddString(rawUrl) {
				continue
			}
			d.queue.Enqueue(rawUrl)

			//img.onerror = function () {
			//	if (img.src.endsWith(".jpg")) {
			//		img.src = img.src.replace('.jpg', '.png');
			//	} else {
			//		console.log(img.src);
			//		nextImgGet();
			//	}
			//};
			//
			//try {

			//	img.src = "https://px.flysay.com/https://i.pximg.net/img-original" + (originImg.src.match(/\/img\/.*/)[0].replace('_master1200', '').replace('_square1200', ''))

			//todo change
			//https://i.pximg.net/c/360x360_70/img-master/img/2019/06/14/09/00/01/75214268_p0_square1200.jpg //todo if 404 enqueue again
		}

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

	spider.AddSleep(time.Duration(rand.Float64() * 50e9))

	if spider.FailureLevel > 40 {
		spider.Delete = true
	}
}

func (my *Pixiv) RequestBefore(spider *spider.Spider) {
	spider.CurrentRequest().Header.Set("Referer", "https://www.pixiv.net/artworks/77274818") //todo 动态 !!!!!!!!!!!!!!!!!
	spider.Client().Timeout = time.Minute
}

func (my *Pixiv) ResponseAfter(spider *spider.Spider) {
	//spider.ResetRequest() //todo !!!!!!!!!!!!!!! ???????????
	//spider.Transport.Close() //todo !!!!!!!!!!!!!!! ???????????

	spider.ResponseByte = nil //free memory
}

// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
// 一般用于根据Header过滤不想继续下载的response.content_type
func (my *Pixiv) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	if !strings.Contains(response.Header.Get("Content-type"), "image") {
		log.Println("not image: got " + response.Header.Get("Content-type") + " HTTP code: " + response.Status)
		return false, nil
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
        let url = "https://www.pixiv.net/ajax/illust/recommend/illusts?";
        Array.from(document.body.querySelectorAll('div[data-gtm-recommend-illust-id]:not([data-gtm-recommend-illust-id=""])')).forEach(function (element) {
            url += "illust_ids%5B%5D=" + element.getAttribute("data-gtm-recommend-illust-id") + "&"
        });
        ajaxDo_({
            method: "GET", url: url, success: function (res) {
                let jsonRes = JSON.parse(res.response);
                if (jsonRes.error) {
                    console.log(res.responseText);
                    return;
                }
                ajaxDo_({
                    method: "POST",
                    url: "http://127.0.0.1:666/project/pixiv/crawl/upload",
                    data: JSON.stringify(jsonRes.body.illusts),
                });
            }
        });
    })();
</pre>`
	w.Header().Set("Content-type", "text/html; charset=UTF-8")
	io.WriteString(w, htmlStr)
}
