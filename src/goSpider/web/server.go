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

	//const (
	//	CloseNormalClosure           = 1000
	//	CloseGoingAway               = 1001
	//	CloseProtocolError           = 1002
	//	CloseUnsupportedData         = 1003
	//	CloseNoStatusReceived        = 1005
	//	CloseAbnormalClosure         = 1006
	//	CloseInvalidFramePayloadData = 1007
	//	ClosePolicyViolation         = 1008
	//	CloseMessageTooBig           = 1009
	//	CloseMandatoryExtension      = 1010
	//	CloseInternalServerErr       = 1011
	//	CloseServiceRestart          = 1012
	//	CloseTryAgainLater           = 1013
	//	CloseTLSHandshake            = 1015
	//)

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
				refreshRate, _ = strconv.ParseFloat(string(b), 64)
				if refreshRate < refreshRateMin {
					refreshRate = refreshRateMin
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
	template.Must(template.ParseFiles(helper.Env().TemplatePath+"index.html")).Execute(w, "ws://"+r.Host+"/echo")
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
	html := `<table><tr><th>Server</th><th>Avg Time</th><th>Traffic In</th><th>Traffic Out</th><th>Load 5s</th><th>60s</th><th>5min</th><th>15min</th><th>Dispatch</th><th>Access</th><th>Failure</th><th style="width:100px">Failure 5min</th></tr>`

	start := time.Now()
	avgLoad := 0.0
	for _, s := range dispatcherObj.GetSpiders() {
		avgLoad += s.Transport.LoadRate(5)
		if s.ConnectFail {
			html += "<tr style=\"background:yellow\">"
		} else {
			html += "<tr>"
		}
		html += "<td>" + s.Transport.S.Name + " </td><td>" + s.GetAvgTime().String() + "</td><td>" + helper.ByteCountBinary(s.Transport.TrafficIn) + "</td><td>" + helper.ByteCountBinary(s.Transport.TrafficOut) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(5), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60*5), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60*15), 'f', 2, 64) + "</td><td>" + strconv.Itoa(s.Transport.LoopCount) + "</td><td>" + strconv.Itoa(s.Transport.GetAccessCount()) + "</td><td>" + strconv.Itoa(s.Transport.GetFailureCount()) + "</td><td> " + strconv.FormatFloat(s.Transport.FailureRate(60*5)*100, 'f', 2, 64) + "%</td>"
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

	html += "</table><br>"
	overviewHtml := `
<table>
    <tr>
        <th>Queue</th>
        <td style="width:130px">` + strconv.Itoa(int(queueCount)) + `</td>
        <th>Redis Mem</th>
        <td style="width:130px">` + helper.ByteCountBinary(uint64(redisMem)) + `</td>
        <th>Bloom Filter</th>
        <td style="width:130px">` + helper.ByteCountBinary(uint64(fileSize)) + `</td>
        <th>Load</th>
        <td style="width:130px">` + strconv.FormatFloat(avgLoad/float64(len(dispatcherObj.GetSpiders())), 'f', 2, 64) + `</td>
        <th>Mem SYS</th>
        <td style="width:130px">` + helper.ByteCountBinary(mem.Sys) + `</td>
	</tr>
	<tr>
        <th>Goroutine</th>
        <td>` + strconv.Itoa(runtime.NumGoroutine()) + `</td>
        <th>Sockets</th>
        <td>` + strconv.Itoa(helper.GetSocketEstablishedCountLazy()) + `</td>
        <th>WebSockets</th>
        <td>` + strconv.Itoa(webSocketConnections) + `</td>
        <th>Time</th>
        <td>` + time.Since(start).String() + `</td>
        <th>Run</th>
        <td>` + time.Since(startTime).String() + `</td>
    </tr>
</table>
<br>
`

	html += "<table><tr><th style=\"width:100px\">Server</th><th style=\"width:100px\">Time</th><th>Current Url</th></tr>"
	for _, s := range dispatcherObj.GetSpiders() {
		if s.CurrentRequest != nil {
			html += "<tr><td>" + s.Transport.S.Name + "</td><td>" + time.Since(s.RequestStartTime).String() + "</td><td><a class=\"text-ellipsis\" target=\"_blank\" href=\"" + s.CurrentRequest.URL.String() + "\">" + helper.TruncateStr(s.CurrentRequest.URL.String(), 80, "...("+strconv.Itoa(len(s.CurrentRequest.URL.String()))+")") + "</a></td></tr>"
		}
	}
	html += "</table><br>"

	html += "<table><tr><th style=\"width:100px\">Server</th><th style=\"width:100px\">Status</th><th>Size</th><th style=\"width:120px\">Add At</th><th style=\"width:120px\">Time</th><th>Url</th></tr>"

	for i := len(spider.RecentFetchList); i > 0; i-- {
		l := spider.RecentFetchList[i-1]
		if l.StatusCode == 0 && l.ConsumeTime != 0 {
			html += "<tr style=\"background:#ff9d87\">"
		} else if l.ConsumeTime == 0 {
			html += "<tr style=\"background:#f2f2f2\">"
		} else if l.StatusCode != 200 {
			html += "<tr style=\"background:yellow\">"
		} else {
			html += "<tr>"
		}
		html += "<td>" + l.TransportName + "</td><td>" + strconv.Itoa(l.StatusCode) + " " + l.ErrType + "</td><td>" + helper.ByteCountBinary(l.ResponseSize) + "</td><td>" + l.AddTime.Format("01-02 15:04:05") + "</td><td>" + l.ConsumeTime.String() + "</td><td><a class=\"text-ellipsis\" target=\"_blank\" href=\"" + l.Url.String() + "\">" + helper.TruncateStr(l.Url.String(), 50, "...("+strconv.Itoa(len(l.Url.String()))+")") + "</a></td>"
		html += "</tr>"
	}
	html += "</table>"
	return overviewHtml + html
}
