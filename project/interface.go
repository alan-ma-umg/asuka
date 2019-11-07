package project

import (
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/spider"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ThrottleInterface interface {
	SetThrottleSpeed(ThrottleSpeed float64)
	SetThrottleSleep(ThrottleSleepSecond int)
}

type IProject interface {
	// Init DoOnce func
	Init(my *Dispatcher)

	// EntryUrl 万恶的起源
	// Firstly
	EntryUrl() []string

	// Throttle 控制抓取速度的一个地方
	// 使用spider.AddSleep()方法, 而不是time.Sleep().
	// Secondly
	Throttle(spider *spider.Spider)

	// RequestBefore http请求发起之前
	// Thirdly
	RequestBefore(spider *spider.Spider)

	Fetch(spider *spider.Spider, u *url.URL) (summary *spider.Summary, err error)

	// EnqueueForFailure 请求或者响应失败时重新入失败队列, 可以修改这里修改加入失败队列的实现. 会在 Goroutine 中被异步调用
	// retryEnqueueUrl & spiderEnqueueUrl 两者一般一致即可
	// retryEnqueueUrl 用于检测失败次数,后加入retries计数. retryEnqueueUrl是为了缩短url长度减少retries的空间, 比如去掉HOST部分, 只保存与检测PATH部分
	// spiderEnqueueUrl 用于重新加入正常抓取队列.
	EnqueueForFailure(spider *spider.Spider, err error, retryEnqueueUrl, spiderEnqueueUrl string, retryTimes int) (success bool, tries int)

	// RequestAfter HTTP请求已经完成, Response Header已经获取到, 但是 Response.Body 未下载
	// 一般用于根据Header过滤不想继续下载的response.content_type
	// Fourth
	DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error)

	// ResponseSuccess HTTP请求成功(Response.Body下载完成)之后
	// 一般用于采集数据的地方
	// Fifth
	ResponseSuccess(spider *spider.Spider)

	// EnqueueFilter HTTP完成并成功后, 从HTML中解析的每条URL都会经过这个筛选和处理. 空字符串则不入队列. 异步执行
	// Sixth
	EnqueueFilter(spider *spider.Spider, l *url.URL) string

	// ResponseAfter HTTP请求失败/成功之后
	// At Last
	ResponseAfter(spider *spider.Spider)

	//Showing 在web监控上展示信息
	Showing() string

	//Name of project
	Name() string

	// BloomFilterSize 系统首次创建对应过滤器时的容量
	InitBloomFilterCapacity() uint

	//项目自定义WEB
	WEBSite(w http.ResponseWriter, r *http.Request)
	WEBSiteLoginRequired(w http.ResponseWriter, r *http.Request) bool //控制是否需要登录
}

type Implement struct{}

func (my *Implement) InitBloomFilterCapacity() uint       { return 5000000 }
func (my *Implement) Init(d *Dispatcher)                  {}
func (my *Implement) RequestBefore(spider *spider.Spider) {}
func (my *Implement) Fetch(spider *spider.Spider, u *url.URL) (summary *spider.Summary, err error) {
	return spider.HttpFetch(u)
}
func (my *Implement) DownloadFilter(spider *spider.Spider, response *http.Response) (bool, error) {
	return true, nil
}

// EnqueueForFailure 请求或者响应失败时重新入失败队列, 可以修改这里修改加入失败队列的实现. 会在 Goroutine 中被异步调用
// retryEnqueueUrl & spiderEnqueueUrl 两者一般一致即可
// retryEnqueueUrl 用于检测失败次数,后加入retries计数. retryEnqueueUrl是为了缩短url长度减少retries的空间, 比如去掉HOST部分, 只保存与检测PATH部分
// spiderEnqueueUrl 用于重新加入正常抓取队列.
func (my *Implement) EnqueueForFailure(spider *spider.Spider, err error, retryEnqueueUrl, spiderEnqueueUrl string, retryTimes int) (success bool, tries int) {
	return spider.GetQueue().EnqueueForFailure(retryEnqueueUrl, spiderEnqueueUrl, retryTimes)
}

// ResponseAfter HTTP请求失败/成功之后
// At Last
func (my *Implement) ResponseAfter(spider *spider.Spider) {
	//重置spider比较消耗性能
	spider.ResetSpider() //现在是每请求一次, 就重置一次. 请求代理也会重新连接
}
func (my *Implement) EnqueueFilter(spider *spider.Spider, l *url.URL) (enqueueUrl string) {
	return l.String()
}
func (my *Implement) Name() string { return "" }
func (my *Implement) WEBSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=UTF-8")
	io.WriteString(w, my.Name())
}
func (my *Implement) WEBSiteLoginRequired(w http.ResponseWriter, r *http.Request) bool { return true }

type SpeedShowing struct {
	DefaultShowingEnable bool
	DefaultShowing       string
	DefaultSpeedCount    uint
	DefaultSpeedTotal    time.Duration
	DefaultSpeedMin      time.Duration
	DefaultSpeedAvgDiv   time.Duration
	DefaultSpeedMax      time.Duration
}

func (my *SpeedShowing) Showing() string {
	my.DefaultShowingEnable = true
	if my.DefaultShowing == "" {
		return "Have a nice day !"
	} else {
		return my.DefaultShowing
	}
}

func (my *SpeedShowing) ResponseSuccess(spider *spider.Spider) {
	if my.DefaultShowingEnable {
		if my.DefaultShowing == "" {
			my.DefaultSpeedMin = time.Hour
		}
		duration := spider.RequestEndTime.Sub(spider.RequestStartTime)
		if duration < my.DefaultSpeedMin {
			my.DefaultSpeedMin = duration
		}
		if duration > my.DefaultSpeedMax {

			my.DefaultSpeedMax = duration
		}
		if my.DefaultSpeedAvgDiv == 0 {
			my.DefaultSpeedAvgDiv = duration
		}
		my.DefaultSpeedAvgDiv = (my.DefaultSpeedAvgDiv + duration) / 2

		my.DefaultSpeedTotal += duration
		my.DefaultSpeedCount++
		my.DefaultShowing = "MIN: " + my.DefaultSpeedMin.Truncate(time.Microsecond).String() + "  MAX: " + my.DefaultSpeedMax.Truncate(time.Microsecond).String() + "  AVG: " + (my.DefaultSpeedTotal / time.Duration(my.DefaultSpeedCount)).Truncate(time.Microsecond).String() + " / " + my.DefaultSpeedAvgDiv.Truncate(time.Microsecond).String()
	}
}

type SpiderThrottle struct {
	SpiderThrottleSpeed       float64
	SpiderThrottleSleepSecond int
}

func (my *SpiderThrottle) SetThrottleSpeed(ThrottleSpeed float64) {
	my.SpiderThrottleSpeed = ThrottleSpeed
}
func (my *SpiderThrottle) SetThrottleSleep(ThrottleSleepSecond int) {
	my.SpiderThrottleSleepSecond = ThrottleSleepSecond
}
func (my *SpiderThrottle) Throttle(spider *spider.Spider) {
	if spider.LoadRate(5) > my.SpiderThrottleSpeed {
		spider.AddSleep(time.Duration(helper.MaxInt(my.SpiderThrottleSleepSecond, 1)) * 1e9)
	}
}
