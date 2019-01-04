package web

import (
	"compress/gzip"
	"github.com/gorilla/websocket"
	"goSpider/database"
	"goSpider/dispatcher"
	"goSpider/helper"
	"goSpider/spider"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var upgrade = websocket.Upgrader{
	EnableCompression: true,
}
var startTime = time.Now()
var webSocketConnections = 0
var pwd, _ = os.Getwd()

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func commonHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}

		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html")
		}
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Server", "spider")
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	webSocketConnections++

	defer func() {
		webSocketConnections--
		c.Close()
	}()
	for {
		err = c.WriteMessage(websocket.TextMessage, []byte(html()))
		if err != nil {
			//log.Println("write:", err)
			break
		}
		time.Sleep(0.3e9)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(pwd+"/src/goSpider/web/templates/index.html")).Execute(w, "ws://"+r.Host+"/echo")
}

var dispatcherObj *dispatcher.Dispatcher

func Server(d *dispatcher.Dispatcher, address string) {
	dispatcherObj = d
	http.HandleFunc("/queue", commonHandler(queue))
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", commonHandler(home))
	http.HandleFunc("/forever/", forever)
	log.Fatal(http.ListenAndServe(address, nil))
}

func queue(w http.ResponseWriter, r *http.Request) {
	list, _ := database.Redis().LRange(helper.Env().Redis.URLQueueKey, 0, 1000).Result()
	template.Must(template.ParseFiles(pwd+"/src/goSpider/web/templates/queue.html")).Execute(w, list)
}

func forever(w http.ResponseWriter, r *http.Request) {
	str := ""
	for i := 0; i < rand.Intn(4); i++ {
		str += "<a href=\"/forever/" + strconv.Itoa(rand.Int()) + "\">" + strconv.Itoa(i) + "</a>"
	}
	w.Header().Set("Content-type", "text/html")
	io.WriteString(w, str)
}

func html() string {
	html := "<table><tr><th>Server</th><th>Avg Time</th><th>Traffic In</th><th>Traffic Out</th><th>Load Rate 5s</th><th>Load Rate 60s</th><th>Load Rate 5m</th><th>Load Rate 15m</th><th>Dispatcher Count</th><th>Access Count</th><th>Failure Count</th></tr>"

	start := time.Now()
	avgLoad := 0.0
	for _, s := range dispatcherObj.GetSpiders() {
		avgLoad += s.Transport.LoadRate(5)
		if s.ConnectFail {
			html += "<tr style=\"background:yellow\">"
		} else {
			html += "<tr>"
		}
		html += "<td>" + s.Transport.S.Name + " </td><td>" + s.GetAvgTime().String() + "</td><td>" + helper.ByteCountBinary(s.Transport.TrafficIn) + "</td><td>" + helper.ByteCountBinary(s.Transport.TrafficOut) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(5), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60*5), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60*15), 'f', 2, 64) + "</td><td>" + strconv.Itoa(s.Transport.LoopCount) + "</td><td>" + strconv.Itoa(s.Transport.GetAccessCount()) + "</td><td>" + strconv.Itoa(s.Transport.GetFailureCount()) + "</td>"
		html += "</tr>"
	}

	queueCount, _ := database.Redis().LLen(helper.Env().Redis.URLQueueKey).Result()

	//memory
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	//redis memory
	redisMem, err := database.Redis().MemoryUsage(helper.Env().Redis.URLQueueKey).Result()
	if err != nil {
		//fmt.Println(err)
		redisMem = 0
	}

	//bloomFilter
	var fileSize int64 = 0
	fi, err := os.Stat(helper.Env().BloomFilterFile)
	if err == nil {
		fileSize = fi.Size()
	}

	html += "</table><a href=\"/queue\">queue: " + strconv.Itoa(int(queueCount)) + "</a><br> Redis mem: " + helper.ByteCountBinary(uint64(redisMem)) + " <br>"
	html += "BloomFilter: " + helper.ByteCountBinary(uint64(fileSize))
	html += "<br> Avg Load:" + strconv.FormatFloat(avgLoad/float64(len(dispatcherObj.GetSpiders())), 'f', 2, 64) + "</br>"
	html += "Alloc: " + helper.ByteCountBinary(mem.Alloc) + " <br> TotalAlloc: " + helper.ByteCountBinary(mem.Alloc) + " <br> Sys: " + helper.ByteCountBinary(mem.Sys) + " <br>"
	html += "NumGC: " + strconv.Itoa(int(mem.NumGC)) + " <br> NumGoroutine: " + strconv.Itoa(runtime.NumGoroutine()) + "<br>"
	html += "webSocketConnections: " + strconv.Itoa(webSocketConnections) + "<br>"
	html += "time: " + time.Since(start).String() + "   " + time.Since(startTime).String()

	html += "<table><tr><th style=\"width:100px\">Server</th><th style=\"width:100px\">Time</th><th>Current Url</th></tr>"
	for _, s := range dispatcherObj.GetSpiders() {
		if s.CurrentRequest != nil {
			html += "<tr><td>" + s.Transport.S.Name + "</td><td>" + time.Since(s.RequestStartTime).String() + "</td><td><a class=\"text-ellipsis\" href=\"" + s.CurrentRequest.URL.String() + "\">" + s.CurrentRequest.URL.String() + "</a></td></tr>"
		}
	}
	html += "</table><br>"

	html += "<table><tr><th style=\"width:100px\">Server</th><th style=\"width:100px\">Status Code</th><th style=\"width:120px\">Add At</th><th style=\"width:120px\">Time</th><th>Url</th></tr>"

	for i := len(spider.RecentFetchList); i > 0; i-- {
		l := spider.RecentFetchList[i-1]
		if l.StatusCode == 0 && l.ConsumeTime != 0 {
			html += "<tr style=\"background:red\">"
		} else if l.ConsumeTime == 0 {
			html += "<tr style=\"background:#f2f2f2\">"
		} else {
			html += "<tr>"
		}
		html += "<td>" + l.TransportName + "</td><td>" + strconv.Itoa(l.StatusCode) + "</td><td>" + l.AddTime.Format("01-02 15:04:05") + "</td><td>" + l.ConsumeTime.String() + "</td><td><a class=\"text-ellipsis\" href=\"" + l.Url.String() + "\">" + l.Url.String() + "</a></td>"
		html += "</tr>"
	}
	html += "</table>"
	return html
}
