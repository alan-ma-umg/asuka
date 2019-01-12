package web

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"goSpider/database"
	"goSpider/dispatcher"
	"goSpider/helper"
	"goSpider/spider"
	"html/template"
	"io"
	"log"
	"math"
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

func commonHandleFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}

		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html")
		}
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Server", "Asuka")
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func commonHandle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/javascript")
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Server", "Asuka")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		h.ServeHTTP(gzr, r)
	})
}

func Server(d *dispatcher.Dispatcher, address string) {
	dispatcherObj = d //todo
	http.HandleFunc("/queue", commonHandleFunc(queue))
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/socket.io", IO)
	http.HandleFunc("/", commonHandleFunc(index))
	http.HandleFunc("/monitor/", commonHandleFunc(monitor))
	http.HandleFunc("/forever/", forever)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, helper.Env().TemplatePath+"/favicon.ico")
	})
	http.Handle("/js/", commonHandle(http.StripPrefix("/js", http.FileServer(http.Dir(helper.Env().TemplatePath+"js")))))

	log.Fatal(http.ListenAndServe(address, nil))
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
		err = c.WriteMessage(websocket.TextMessage, []byte(responseHtml()))
		if err != nil {
			log.Println("write:", err)
			break
		}
		time.Sleep(time.Duration(refreshRate * 1e9))
	}
}

func monitor(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(helper.Env().TemplatePath+"monitor.html")).Execute(w, nil)
}
func index(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(helper.Env().TemplatePath+"index.html")).Execute(w, runtime.GOOS)
}
func IO(w http.ResponseWriter, r *http.Request) {
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
	responseContent := "home"
	var recentFetchIndex int64 = 0

	for {
		messageType, b, err := c.ReadMessage()
		if err != nil {
			//log.Println("read:", err)
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
			case "home":
				responseContent = strings.TrimSpace(string(b))
			case "recent":
				responseContent = strings.TrimSpace(string(b))
			default:
				refreshRateTemp, err := strconv.ParseFloat(string(b), 64)
				if err == nil {
					refreshRate = math.Max(refreshRateTemp, refreshRateMin)
				}
			}
		}

		switch responseContent {
		case "home":
			err = c.WriteMessage(websocket.TextMessage, homeJson(responseContent))
		case "recent":
			jsonRes, n := recentJson(responseContent, recentFetchIndex)
			recentFetchIndex = n
			err = c.WriteMessage(websocket.TextMessage, jsonRes)
		}
		if err != nil {
			//log.Println("write:", err)
			break
		}
		time.Sleep(time.Duration(refreshRate * 1e9))
	}
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

func recentJson(sType string, recentFetchIndex int64) ([]byte, int64) {
	start := time.Now()
	var jsonMap = map[string]interface{}{
		"type":    sType,
		"basic":   map[string]interface{}{},
		"fetched": []*spider.RecentFetch{},
	}

	var lastIndex int64
	for _, l := range spider.RecentFetchList {
		if l == nil { //Change frequently, prevent nil pointer
			continue
		}
		if l.Index > recentFetchIndex {
			jsonMap["fetched"] = append(jsonMap["fetched"].([]*spider.RecentFetch), l)
			lastIndex = helper.MaxInt64(lastIndex, l.Index)
		}
	}

	responseJsonCommon(jsonMap)
	jsonMap["basic"].(map[string]interface{})["time"] = time.Since(start).Truncate(time.Microsecond).String()
	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b, helper.MaxInt64(lastIndex, recentFetchIndex)
}
func homeJson(sType string) []byte {
	start := time.Now()
	var jsonMap = map[string]interface{}{
		"type":    sType,
		"basic":   map[string]interface{}{},
		"servers": []map[string]interface{}{},
	}
	for index, s := range dispatcherObj.GetSpiders() {
		load5s := s.Transport.LoadRate(5)
		avgTime := s.GetAvgTime()
		failureRate60Value := helper.SpiderFailureRate(s.Transport.AccessCount(60))
		failureRateAllValue := .0
		if s.Transport.GetAccessCount() > 0 {
			failureRateAllValue = float64(s.Transport.GetFailureCount()) / float64(s.Transport.GetAccessCount()) * 100
		}

		server := map[string]interface{}{}
		server["failure_60"] = strconv.FormatFloat(failureRate60Value, 'f', 2, 64)
		server["failure_60_hsl"] = strconv.Itoa(int(100 - failureRate60Value))
		server["failure_all"] = strconv.FormatFloat(failureRateAllValue, 'f', 2, 64)
		server["failure_all_hsl"] = strconv.Itoa(int(100 - failureRateAllValue))
		server["failure_level"] = s.FailureLevel
		server["failure_level_hsl"] = 100 - s.FailureLevel
		server["index"] = index
		server["name"] = s.Transport.S.Name
		server["ping"] = s.Transport.Ping.Truncate(time.Millisecond).String()
		server["ping_hsl"] = helper.MinInt(150, helper.MaxInt(150-int(s.Transport.Ping.Seconds()*1000/2), 0))
		server["ping_failure"] = strconv.FormatFloat(s.Transport.PingFailureRate*100, 'f', 0, 64)
		server["ping_failure_hsl"] = int(150 - s.Transport.PingFailureRate*150)
		server["avg_time"] = avgTime.Truncate(time.Millisecond).String()
		server["waiting"] = "0s"
		if !s.RequestStartTime.IsZero() {
			server["waiting"] = time.Since(s.RequestStartTime).Truncate(time.Millisecond).String()
		}
		server["traffic_in"] = helper.ByteCountBinary(s.Transport.TrafficIn)
		server["traffic_out"] = helper.ByteCountBinary(s.Transport.TrafficOut)
		server["load_5s"] = strconv.FormatFloat(load5s, 'f', 2, 64)                    //todo slowly, make improvement
		server["load_60s"] = strconv.FormatFloat(s.Transport.LoadRate(60), 'f', 2, 64) //todo slowly, make improvement
		server["access_count"] = s.Transport.GetAccessCount()
		server["failure_count"] = s.Transport.GetFailureCount()
		jsonMap["servers"] = append(jsonMap["servers"].([]map[string]interface{}), server)
	}

	//basic
	responseJsonCommon(jsonMap)
	jsonMap["basic"].(map[string]interface{})["time"] = time.Since(start).Truncate(time.Microsecond).String()
	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b
}
func responseJsonCommon(jsonMap map[string]interface{}) {
	sumLoad := .0
	pingFailureAvg := .0
	failureLevelZeroCount := 0
	var pingAvg time.Duration
	var waitingAvg time.Duration
	var avgTimeAvg time.Duration
	var TrafficIn uint64
	var TrafficOut uint64
	for _, s := range dispatcherObj.GetSpiders() {
		if s.FailureLevel == 0 {
			failureLevelZeroCount++
			if !s.RequestStartTime.IsZero() {
				waitingAvg += time.Since(s.RequestStartTime)
			}
			avgTimeAvg += s.GetAvgTime()
		}

		load5s := s.Transport.LoadRate(5)
		pingFailureAvg += s.Transport.PingFailureRate
		pingAvg += s.Transport.Ping
		sumLoad += load5s
		TrafficIn += s.Transport.TrafficIn
		TrafficOut += s.Transport.TrafficOut
	}

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	//redis
	queueCount, _ := database.Redis().LLen(helper.Env().Redis.URLQueueKey).Result()      //about 1ms
	redisMem, _ := database.Redis().MemoryUsage(helper.Env().Redis.URLQueueKey).Result() //about 1ms
	//basic
	jsonMap["basic"].(map[string]interface{})["queue"] = strconv.Itoa(int(queueCount))
	jsonMap["basic"].(map[string]interface{})["redis_mem"] = helper.ByteCountBinary(uint64(redisMem))
	jsonMap["basic"].(map[string]interface{})["avg_time_avg"] = (avgTimeAvg / time.Duration(failureLevelZeroCount)).Truncate(time.Millisecond).String()      //fixme Divide by Zero
	jsonMap["basic"].(map[string]interface{})["waiting_avg"] = (waitingAvg / time.Duration(failureLevelZeroCount)).Truncate(time.Millisecond).String()       //fixme Divide by Zero
	jsonMap["basic"].(map[string]interface{})["ping_avg"] = (pingAvg / time.Duration(len(dispatcherObj.GetSpiders()))).Truncate(time.Millisecond).String()   //fixme Divide by Zero
	jsonMap["basic"].(map[string]interface{})["ping_failure_avg"] = strconv.FormatFloat(pingFailureAvg/float64(len(dispatcherObj.GetSpiders())), 'f', 2, 64) //fixme Divide by Zero
	jsonMap["basic"].(map[string]interface{})["load_sum"] = strconv.FormatFloat(sumLoad, 'f', 2, 64)
	jsonMap["basic"].(map[string]interface{})["load_avg"] = strconv.FormatFloat(sumLoad/float64(len(dispatcherObj.GetSpiders())), 'f', 2, 64) //fixme Divide by Zero
	jsonMap["basic"].(map[string]interface{})["traffic_in"] = helper.ByteCountBinary(TrafficIn)
	jsonMap["basic"].(map[string]interface{})["traffic_out"] = helper.ByteCountBinary(TrafficOut)
	jsonMap["basic"].(map[string]interface{})["mem_sys"] = helper.ByteCountBinary(mem.Sys)
	jsonMap["basic"].(map[string]interface{})["goroutine"] = runtime.NumGoroutine()
	jsonMap["basic"].(map[string]interface{})["connections"] = helper.GetSocketEstablishedCountLazy()
	jsonMap["basic"].(map[string]interface{})["ws_connections"] = webSocketConnections
	jsonMap["basic"].(map[string]interface{})["uptime"] = time.Since(startTime).Truncate(time.Second).String()
}

func responseHtml() string {
	html := `<table><tr><th style="width:1px">#</th><th style="width:1px">Server</th><th style="width:1px">Ping / Lost</th><th style="width:1px">Avg Time</th><th style="width:1px">Traffic I / O</th><th>Load 5s / 60s / 15min / 30min</th><th style="width:1px">Dispatch</th><th style="width:145px">Failure</th><th style="width:1px">Status</th></tr>`

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

		FailStr := "-"
		if s.Transport.GetAccessCount() > 0 {
			failureRate60Value := helper.SpiderFailureRate(s.Transport.AccessCount(60))
			failureRateAllValue := float64(s.Transport.GetFailureCount()) / float64(s.Transport.GetAccessCount()) * 100
			FailStr = `<div style="display:inline-block;width:50px" class="right"><span style="color: hsl(` + strconv.Itoa(int(100-failureRate60Value)) + `, 100%, 35%);">` + strconv.FormatFloat(failureRate60Value, 'f', 2, 64) + `%</span></div>` +
				" | " +
				`<div style="display:inline-block;width:50px" class="left"><span style="color: hsl(` + strconv.Itoa(int(100-failureRateAllValue)) + `, 100%, 35%);">` + strconv.FormatFloat(failureRateAllValue, 'f', 2, 64) + `%</span></div>`
		}

		html += `
<td>` + strconv.Itoa(index+1) + ` </td>
<td class="center">` + s.Transport.S.Name + `</td>
<td class="center">
	<div style="display:inline-block;width:50px" class="right">
	<span style="color: hsl(` + strconv.Itoa(helper.MinInt(150, helper.MaxInt(150-int(s.Transport.Ping.Seconds()*1000/2), 0))) + `, 100%, 35%);">` + s.Transport.Ping.Truncate(time.Millisecond).String() + `</span></div> | 
<div style="display:inline-block;width:50px" class="left"><span style="color: hsl(` + strconv.Itoa(int(150-s.Transport.PingFailureRate*150)) + `, 100%, 35%);">` + strconv.FormatFloat(s.Transport.PingFailureRate*100, 'f', 0, 64) + `%</span></div></td>
<td>` + s.GetAvgTime().Truncate(time.Millisecond).String() + `</td>
<td class="center"><div style="display:inline-block;width:60px" class="right">` + helper.ByteCountBinary(s.Transport.TrafficIn) + `</div> | <div style="display:inline-block;width:60px" class="left">` + helper.ByteCountBinary(s.Transport.TrafficOut) + `</div></td>
<td class="center"> ` + strconv.FormatFloat(s.Transport.LoadRate(5), 'f', 2, 64) + ` |
 ` + strconv.FormatFloat(s.Transport.LoadRate(60), 'f', 2, 64) + ` |
 ` + strconv.FormatFloat(s.Transport.LoadRate(60*15), 'f', 2, 64) + ` |
 ` + strconv.FormatFloat(s.Transport.LoadRate(60*30), 'f', 2, 64) + `</td>
<td class="center">	<div style="display:inline-block;width:50px" class="right">` + strconv.Itoa(s.Transport.GetAccessCount()) + `</div> | <div style="display:inline-block;width:50px" class="left">` + strconv.Itoa(s.Transport.GetFailureCount()) + `</div></td>
<td class="center"> ` + FailStr + `</td>
<td style="font-weight:800;color: hsl(` + strconv.Itoa(100-s.FailureLevel) + `, 100%, 35%);" class="center">` + strconv.Itoa(s.FailureLevel) + `</td>`

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
	//
	//html += "<table><tr><th style=\"width:100px\">Server</th><th style=\"width:100px\">Time</th><th>Current Url</th></tr>"
	//for _, s := range dispatcherObj.GetSpiders() {
	//	if s.CurrentRequest != nil && s.FailureLevel == 0 {
	//		html += "<tr><td>" + s.Transport.S.Name + "</td><td>" + time.Since(s.RequestStartTime).Truncate(time.Millisecond).String() + "</td><td><a class=\"text-ellipsis\" target=\"_blank\" href=\"" + s.CurrentRequest.URL.String() + "\">" + helper.TruncateStr([]rune(s.CurrentRequest.URL.String()), 60, "...("+strconv.Itoa(len([]rune(s.CurrentRequest.URL.String())))+")") + "</a></td></tr>"
	//	}
	//}
	//html += "</table><br>"

	//html += "<table><tr><th style=\"width:100px\">Server</th><th style=\"width:100px\">Status</th><th>Size</th><th style=\"width:120px\">Add At</th><th style=\"width:120px\">Time</th><th>Url</th></tr>"
	//
	//recentFetchList := make([]*spider.RecentFetch, helper.MinInt(len(spider.RecentFetchList), spider.RecentFetchCount))
	//copy(recentFetchList, spider.RecentFetchList)
	//for i := len(recentFetchList); i > 0; i-- {
	//	l := recentFetchList[i-1]
	//	if l.StatusCode == 0 && l.ConsumeTime != 0 {
	//		html += "<tr style=\"background:#ff9d87\">"
	//	} else if l.ConsumeTime == 0 {
	//		html += "<tr style=\"background:#f2f2f2\">"
	//	} else if l.StatusCode != 200 {
	//		html += "<tr style=\"background:yellow\">"
	//	} else {
	//		html += "<tr>"
	//	}
	//	html += "<td>" + l.TransportName + "</td><td>" + strconv.Itoa(l.StatusCode) + " " + l.ErrType + "</td><td>" + helper.ByteCountBinary(l.ResponseSize) + "</td><td>" + l.AddTime.Format("01-02 15:04:05") + "</td><td>" + l.ConsumeTime.Truncate(time.Millisecond).String() + "</td><td><a class=\"text-ellipsis\" target=\"_blank\" href=\"" + l.RawUrl + "\">" + helper.TruncateStr([]rune(l.RawUrl), 40, "...("+strconv.Itoa(len([]rune(l.RawUrl)))+")") + "</a></td>"
	//	html += "</tr>"
	//}
	//html += "</table>"

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
