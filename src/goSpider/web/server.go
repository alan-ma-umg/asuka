package web

import (
	"compress/gzip"
	"fmt"
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
	"runtime/debug"
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
var dispatcherObj *dispatcher.Dispatcher

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

	refreshRateMin := 0.2
	refreshRate := refreshRateMin
	go func() {
		for {
			messageType, b, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			if messageType == 1 {

				switch strings.TrimSpace(string(b)) {
				case "free":
					debug.FreeOSMemory()
					fmt.Println("debug.FreeOsMemory")
				case "stop":
					for _, s := range dispatcherObj.GetSpiders() {
						s.Stop = true
					}
					fmt.Println("spider stop")
				case "start":
					for _, s := range dispatcherObj.GetSpiders() {
						s.Stop = false
					}
					fmt.Println("spider start")
				default:
					refreshRate, _ = strconv.ParseFloat(string(b), 64)
					if refreshRate < refreshRateMin {
						refreshRate = refreshRateMin
					}
				}
			}
		}
	}()

	for {
		err = c.WriteMessage(websocket.TextMessage, []byte(html()))
		if err != nil {
			log.Println("write:", err)
			break
		}
		time.Sleep(time.Duration(refreshRate * 1e9))
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(helper.Env().TemplatePath+"index.html")).Execute(w, nil)
}

func Server(d *dispatcher.Dispatcher, address string) {
	dispatcherObj = d //todo
	http.HandleFunc("/queue", commonHandler(queue))
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", commonHandler(home))
	http.HandleFunc("/forever/", forever)
	log.Fatal(http.ListenAndServe(address, nil))
}

func queue(w http.ResponseWriter, r *http.Request) {
	list, _ := database.Redis().LRange(helper.Env().Redis.URLQueueKey, 0, 1000).Result()
	template.Must(template.ParseFiles(helper.Env().TemplatePath+"queue.html")).Execute(w, list)
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
	html := `<table><tr><th style="width:1px">#</th><th style="width:1px">Server</th><th style="width:1px">Avg Time</th><th>Traffic In</th><th>Traffic Out</th><th>Load 5s</th><th>60s</th><th>5min</th><th>15min</th><th>Access</th><th>Failure</th><th style="width:145px">Failure 60s</th></tr>`

	start := time.Now()
	sumLoad := 0.0
	var TrafficIn uint64 = 0
	var TrafficOut uint64 = 0
	for index, s := range dispatcherObj.GetSpiders() {
		sumLoad += s.Transport.LoadRate(5)
		TrafficIn += s.Transport.TrafficIn
		TrafficOut += s.Transport.TrafficOut
		if s.FailureLevel > 0 {
			html += `<tr style="background:#ffffd2">`
		} else {
			html += "<tr>"
		}

		FailStr := ""
		if s.Transport.GetAccessCount() > 0 {
			failureRate60Value := helper.SpiderFailureRate(s.Transport.AccessCount(60))
			failureRate60 := strconv.FormatFloat(failureRate60Value, 'f', 2, 64)
			failureRate60Html := ""
			if failureRate60Value > 30.0 {
				failureRate60Html = `<span style="color: rgb(` + failureRate60 + `%, 0%, 0%)">` + failureRate60 + `%</span>`
			} else {
				failureRate60Html = `<span style="color: rgb(0%, ` + strconv.FormatFloat(100.0-failureRate60Value, 'f', 2, 64) + `%, 0%)">` + failureRate60 + `%</span>`
			}

			failureRateAllValue := float64(s.Transport.GetFailureCount()) / float64(s.Transport.GetAccessCount()) * 100
			failureRateAll := strconv.FormatFloat(failureRateAllValue, 'f', 2, 64)
			failureRateAllHtml := ""
			if failureRateAllValue > 30.0 {
				failureRateAllHtml = `<span style="color: rgb(` + failureRateAll + `%, 0%, 0%)">` + failureRateAll + "%</span>"
			} else {
				failureRateAllHtml = `<span style="color: rgb(0%, ` + strconv.FormatFloat(100.0-failureRateAllValue, 'f', 2, 64) + `%, 0%)">` + failureRateAll + "%</span>"
			}

			FailStr = failureRate60Html + " | " + failureRateAllHtml
		} else {
			FailStr = strconv.FormatFloat(helper.SpiderFailureRate(s.Transport.AccessCount(60)), 'f', 2, 64)
		}

		html += `
<td>` + strconv.Itoa(index+1) + ` </td>
<td class="center">` + helper.TruncateStr([]rune(s.Transport.S.Name), 10, "") + `(F. ` + strconv.Itoa(s.FailureLevel) + `) </td>
<td>` + s.GetAvgTime().Truncate(time.Millisecond).String() + `</td>
<td>` + helper.ByteCountBinary(s.Transport.TrafficIn) + `</td>
<td>` + helper.ByteCountBinary(s.Transport.TrafficOut) + `</td>
<td> ` + strconv.FormatFloat(s.Transport.LoadRate(5), 'f', 2, 64) + `</td>
<td> ` + strconv.FormatFloat(s.Transport.LoadRate(60), 'f', 2, 64) + `</td>
<td> ` + strconv.FormatFloat(s.Transport.LoadRate(60*5), 'f', 2, 64) + `</td>
<td> ` + strconv.FormatFloat(s.Transport.LoadRate(60*15), 'f', 2, 64) + `</td>
<td>` + strconv.Itoa(s.Transport.GetAccessCount()) + `</td>
<td>` + strconv.Itoa(s.Transport.GetFailureCount()) + `</td>
<td class="center"> ` + FailStr + `</td>`
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
	html += "</table><br>"

	html += "<table><tr><th style=\"width:100px\">Server</th><th style=\"width:100px\">Time</th><th>Current Url</th></tr>"
	for _, s := range dispatcherObj.GetSpiders() {
		if s.CurrentRequest != nil && s.FailureLevel == 0 {
			html += "<tr><td>" + s.Transport.S.Name + "</td><td>" + time.Since(s.RequestStartTime).Truncate(time.Millisecond).String() + "</td><td><a class=\"text-ellipsis\" target=\"_blank\" href=\"" + s.CurrentRequest.URL.String() + "\">" + helper.TruncateStr([]rune(s.CurrentRequest.URL.String()), 60, "...("+strconv.Itoa(len([]rune(s.CurrentRequest.URL.String())))+")") + "</a></td></tr>"
		}
	}
	html += "</table><br>"

	html += "<table><tr><th style=\"width:100px\">Server</th><th style=\"width:100px\">Status</th><th>Size</th><th style=\"width:120px\">Add At</th><th style=\"width:120px\">Time</th><th>Url</th></tr>"

	recentFetchList := make([]*spider.RecentFetch, helper.MinInt(len(spider.RecentFetchList), spider.RecentFetchCount))
	copy(recentFetchList, spider.RecentFetchList)
	for i := len(recentFetchList); i > 0; i-- {
		l := recentFetchList[i-1]
		if l.StatusCode == 0 && l.ConsumeTime != 0 {
			html += "<tr style=\"background:#ff9d87\">"
		} else if l.ConsumeTime == 0 {
			html += "<tr style=\"background:#f2f2f2\">"
		} else if l.StatusCode != 200 {
			html += "<tr style=\"background:yellow\">"
		} else {
			html += "<tr>"
		}
		html += "<td>" + l.TransportName + "</td><td>" + strconv.Itoa(l.StatusCode) + " " + l.ErrType + "</td><td>" + helper.ByteCountBinary(l.ResponseSize) + "</td><td>" + l.AddTime.Format("01-02 15:04:05") + "</td><td>" + l.ConsumeTime.Truncate(time.Millisecond).String() + "</td><td><a class=\"text-ellipsis\" target=\"_blank\" href=\"" + l.Url.String() + "\">" + helper.TruncateStr([]rune(l.Url.String()), 40, "...("+strconv.Itoa(len([]rune(l.Url.String())))+")") + "</a></td>"
		html += "</tr>"
	}
	html += "</table>"

	overviewHtml := `
<table>
    <tr>
        <th>Queue</th>
        <td style="width:140px"><a href="/queue">` + strconv.Itoa(int(queueCount)) + `</a></td>
        <th>Redis Mem</th>
        <td style="width:140px">` + helper.ByteCountBinary(uint64(redisMem)) + `</td>
        <th>Load</th>
        <td style="width:140px">` + strconv.FormatFloat(sumLoad/float64(len(dispatcherObj.GetSpiders())), 'f', 2, 64) + ` | ` + strconv.FormatFloat(sumLoad, 'f', 2, 64) + `</td>
        <th>Traffic</th>
        <td style="width:140px">` + helper.ByteCountBinary(TrafficIn) + ` | ` + helper.ByteCountBinary(TrafficOut) + `</td>
        <th>Mem.SYS</th>
        <td style="width:140px">` + helper.ByteCountBinary(mem.Sys) + `</td>
	</tr>
	<tr>
        <th>Goroutine</th>
        <td>` + strconv.Itoa(runtime.NumGoroutine()) + `</td>
        <th>Connection</th>
        <td>` + strconv.Itoa(helper.GetSocketEstablishedCountLazy()) + `</td>
        <th>WebSockets</th>
        <td>` + strconv.Itoa(webSocketConnections) + `</td>
        <th>Time</th>
        <td>` + time.Since(start).Truncate(time.Microsecond).String() + `</td>
        <th>Uptime</th>
        <td>` + time.Since(startTime).Truncate(time.Second).String() + `</td>
    </tr>
</table>
<br>
`
	return overviewHtml + html
}
