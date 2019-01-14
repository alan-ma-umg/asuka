package web

import (
	"asuka/database"
	"asuka/dispatcher"
	"asuka/helper"
	"asuka/spider"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
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
	http.HandleFunc("/socket.io", IO)
	http.HandleFunc("/", commonHandleFunc(index))
	http.HandleFunc("/forever/", forever)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, helper.Env().TemplatePath+"/favicon.ico")
	})
	http.Handle("/js/", commonHandle(http.StripPrefix("/js", http.FileServer(http.Dir(helper.Env().TemplatePath+"js")))))

	log.Fatal(http.ListenAndServe(address, nil))
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

	responseJsonCommon(jsonMap, start)
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
		server["sleep"] = s.GetSleep().Truncate(time.Millisecond).String()
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
	responseJsonCommon(jsonMap, start)

	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b
}
func responseJsonCommon(jsonMap map[string]interface{}, start time.Time) {
	defer func() {
		jsonMap["basic"].(map[string]interface{})["time"] = time.Since(start).Truncate(time.Microsecond).String()
	}()

	jsonMap["basic"] = map[string]interface{}{
		"queue_bls": make(map[int]int),
	}

	sumLoad := .0
	pingFailureAvg := .0
	failureLevelZeroCount := 0
	var pingAvg time.Duration
	var sleepAvg time.Duration
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

		sleepAvg += s.GetSleep()
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
	jsonMap["basic"].(map[string]interface{})["sleep_avg"] = ""
	jsonMap["basic"].(map[string]interface{})["ping_avg"] = ""
	jsonMap["basic"].(map[string]interface{})["ping_failure_avg"] = ""
	jsonMap["basic"].(map[string]interface{})["load_avg"] = ""
	jsonMap["basic"].(map[string]interface{})["avg_time_avg"] = ""
	jsonMap["basic"].(map[string]interface{})["waiting_avg"] = ""

	spiderCount := len(dispatcherObj.GetSpiders())
	if spiderCount > 0 {
		for i, v := range dispatcherObj.GetSpiders()[0].Queue.BlsTestCount {
			jsonMap["basic"].(map[string]interface{})["queue_bls"].(map[int]int)[i] = v
		}
		jsonMap["basic"].(map[string]interface{})["sleep_avg"] = (sleepAvg / time.Duration(spiderCount)).Truncate(time.Millisecond).String()
		jsonMap["basic"].(map[string]interface{})["ping_avg"] = (pingAvg / time.Duration(spiderCount)).Truncate(time.Millisecond).String()
		jsonMap["basic"].(map[string]interface{})["ping_failure_avg"] = strconv.FormatFloat(pingFailureAvg/float64(spiderCount), 'f', 2, 64)
		jsonMap["basic"].(map[string]interface{})["load_avg"] = strconv.FormatFloat(sumLoad/float64(spiderCount), 'f', 2, 64)
	}
	if failureLevelZeroCount > 0 {
		jsonMap["basic"].(map[string]interface{})["avg_time_avg"] = (avgTimeAvg / time.Duration(failureLevelZeroCount)).Truncate(time.Millisecond).String()
		jsonMap["basic"].(map[string]interface{})["waiting_avg"] = (waitingAvg / time.Duration(failureLevelZeroCount)).Truncate(time.Millisecond).String()
	}
	jsonMap["basic"].(map[string]interface{})["queue"] = strconv.Itoa(int(queueCount))
	jsonMap["basic"].(map[string]interface{})["redis_mem"] = helper.ByteCountBinary(uint64(redisMem))
	jsonMap["basic"].(map[string]interface{})["load_sum"] = strconv.FormatFloat(sumLoad, 'f', 2, 64)
	jsonMap["basic"].(map[string]interface{})["traffic_in"] = helper.ByteCountBinary(TrafficIn)
	jsonMap["basic"].(map[string]interface{})["traffic_out"] = helper.ByteCountBinary(TrafficOut)
	jsonMap["basic"].(map[string]interface{})["mem_sys"] = helper.ByteCountBinary(mem.Sys)
	jsonMap["basic"].(map[string]interface{})["goroutine"] = runtime.NumGoroutine()
	jsonMap["basic"].(map[string]interface{})["connections"] = helper.GetSocketEstablishedCountLazy()
	jsonMap["basic"].(map[string]interface{})["ws_connections"] = webSocketConnections
	jsonMap["basic"].(map[string]interface{})["uptime"] = time.Since(startTime).Truncate(time.Second).String()
}
