package web

import (
	"asuka/database"
	"asuka/helper"
	"asuka/project"
	"asuka/spider"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"log"
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
var StartTime = time.Now()
var webSocketConnections = 0
var dispatchers []*project.Dispatcher

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
		//if w.Header().Get("Content-Type") == "" {
		//w.Header().Set("Content-Type", "application/javascript")
		//}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Server", "Asuka")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		h.ServeHTTP(gzr, r)
	})
}

func Server(d []*project.Dispatcher, address string) {
	dispatchers = d

	//init start time
	go func() {
		for _, d := range dispatchers {
			for _, sp := range d.GetSpiders() {
				if StartTime.Unix() > sp.StartTime.Unix() {
					StartTime = sp.StartTime
				}
			}
		}
	}()

	http.HandleFunc("/queue/", commonHandleFunc(queue))
	http.HandleFunc("/project.io", projectIO)
	http.HandleFunc("/index.io", indexIO)
	http.HandleFunc("/switchProject", commonHandleFunc(switchProject))
	http.HandleFunc("/switchServer", commonHandleFunc(switchServer))
	http.HandleFunc("/forever/", forever)
	http.HandleFunc("/", commonHandleFunc(home))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, helper.Env().TemplatePath+"/favicon.ico")
	})
	http.Handle("/static/", commonHandle(http.StripPrefix("/static", http.FileServer(http.Dir(helper.Env().TemplatePath+"static")))))

	log.Fatal(http.ListenAndServe(address, nil))
}

func switchServer(w http.ResponseWriter, r *http.Request) {
	post := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	if _, ok := post["name"]; !ok {
		http.NotFound(w, r)
		return
	}

	if _, ok := post["project"]; !ok {
		http.NotFound(w, r)
		return
	}

	s := searchSpider(post["project"], post["name"])
	if s == nil {
		http.NotFound(w, r)
		return
	}
	s.Stop = !s.Stop

	jsonMap := map[string]interface{}{
		"success": true,
		"data":    post["name"],
	}

	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Header().Set("Content-type", "application/json")
	io.WriteString(w, string(b))
}

func switchProject(w http.ResponseWriter, r *http.Request) {
	post := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	if _, ok := post["name"]; !ok {
		http.NotFound(w, r)
		return
	}

	p := getProjectByName(post["name"])
	if p == nil {
		http.NotFound(w, r)
		return
	}

	p.Stop = !p.Stop

	jsonMap := map[string]interface{}{
		"success": true,
		"data":    post["name"],
	}

	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Header().Set("Content-type", "application/json")
	io.WriteString(w, string(b))
}

func index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		GOOS string
	}{
		GOOS: runtime.GOOS,
	}

	template.Must(template.ParseFiles(helper.Env().TemplatePath + "index.html")).Execute(w, data)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		index(w, r)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 2 {
		http.NotFound(w, r)
		return
	}

	p := getDispatcher(path[1])
	if p == nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		GOOS        string
		ProjectName string
	}{
		GOOS:        runtime.GOOS,
		ProjectName: p.GetProjectName(),
	}

	template.Must(template.ParseFiles(helper.Env().TemplatePath + "project.html")).Execute(w, data)
}

func indexIO(w http.ResponseWriter, r *http.Request) {
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
		messageType, b, err := c.ReadMessage()
		if err != nil {
			break
		}
		if messageType == 1 {
			switch strings.TrimSpace(string(b)) {
			case "free":
				debug.FreeOSMemory()
				fmt.Println("debug.FreeOsMemory")
			case "stop":
				for _, d := range dispatchers {
					for _, s := range d.GetSpiders() {
						s.Stop = true
					}
				}
				fmt.Println("spider stop")
			case "start":
				for _, d := range dispatchers {
					for _, s := range d.GetSpiders() {
						s.Stop = false
					}
				}
				fmt.Println("spider start")
			}
		}

		err = c.WriteMessage(websocket.TextMessage, indexJson())
		if err != nil {
			//log.Println("write:", err)
			break
		}
		time.Sleep(time.Second)
	}

}
func projectIO(w http.ResponseWriter, r *http.Request) {
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	ps, ok := r.URL.Query()["project"]
	if !ok || len(ps) != 1 {
		c.Close()
		return
	}

	p := getDispatcher(ps[0])
	if p == nil {
		c.Close()
		return
	}

	webSocketConnections++

	defer func() {
		webSocketConnections--
		c.Close()
	}()

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
			//case "reconnect":
			//	for _, s := range dispatcherObj.GetSpiders() {
			//		s.Transport.Reconnect()
			//	}
			//	fmt.Println("reconnect")
			case "stop":
				for _, d := range dispatchers {
					for _, s := range d.GetSpiders() {
						s.Stop = true
					}
				}
				fmt.Println("spider stop")
			case "start":
				for _, d := range dispatchers {
					for _, s := range d.GetSpiders() {
						s.Stop = false
					}
				}
				fmt.Println("spider start")
			case "home":
				responseContent = strings.TrimSpace(string(b))
			case "recent":
				responseContent = strings.TrimSpace(string(b))
			}
		}

		switch responseContent {
		case "home":
			err = c.WriteMessage(websocket.TextMessage, projectJson(p, responseContent))
		case "recent":
			jsonRes, n := recentJson(p, responseContent, recentFetchIndex)
			recentFetchIndex = n
			err = c.WriteMessage(websocket.TextMessage, jsonRes)
		}
		if err != nil {
			//log.Println("write:", err)
			break
		}
		time.Sleep(time.Second)
	}
}

func getDispatcher(name string) *project.Dispatcher {
	for _, d := range dispatchers {
		if name == d.GetProjectName() {
			return d
		}
	}

	return nil
}

func queue(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 3 {
		http.NotFound(w, r)
		return
	}

	p := getDispatcher(path[2])
	if p == nil {
		http.NotFound(w, r)
		return
	}

	list, _ := database.Redis().LRange(p.GetQueueKey(), 0, 1000).Result()
	data := struct {
		List        []string
		ProjectName string
	}{
		List:        list,
		ProjectName: p.GetProjectName(),
	}

	template.Must(template.ParseFiles(helper.Env().TemplatePath + "queue.html")).Execute(w, data)
}

func forever(w http.ResponseWriter, r *http.Request) {
	str := ""
	for i := 0; i < rand.Intn(4); i++ {
		str += "<a href=\"/forever/" + strconv.Itoa(rand.Int()) + "\">" + strconv.Itoa(i) + "</a>"
	}
	w.Header().Set("Content-type", "text/html")
	io.WriteString(w, str)
}

func recentJson(p *project.Dispatcher, sType string, recentFetchIndex int64) ([]byte, int64) {
	start := time.Now()
	var jsonMap = map[string]interface{}{
		"type":    sType,
		"fetched": []*spider.Summary{},
	}

	var lastIndex int64
	for _, l := range p.RecentFetchList {
		if l == nil { //Change frequently, prevent nil pointer
			continue
		}
		if l.Index > recentFetchIndex {
			jsonMap["fetched"] = append(jsonMap["fetched"].([]*spider.Summary), l)
			lastIndex = helper.MaxInt64(lastIndex, l.Index)
		}
	}

	responseJsonCommon([]*project.Dispatcher{p}, jsonMap, start)
	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b, helper.MaxInt64(lastIndex, recentFetchIndex)
}

func indexJson() []byte {
	start := time.Now()
	var jsonMap = map[string]interface{}{
		"projects": []map[string]interface{}{},
	}

	periodOfFailureSecond := helper.MinInt(int(time.Since(StartTime).Seconds()), spider.PeriodOfFailureSecond)

	for _, p := range dispatchers {
		projectMap := map[string]interface{}{}

		loads := make(map[int]float64, 9)

		failureRatePeriodValue := 0.0
		failureRateAllValue := .0
		var sleepDuration time.Duration
		var waiting time.Duration
		var TrafficIn uint64
		var TrafficOut uint64
		var NetIn uint64
		var NetOut uint64
		var connections int
		var accessCount int
		var failureCount int
		var serverCount int
		var serverRun int
		var serverEnable int
		for _, s := range p.GetSpiders() {

			failureRatePeriodValue += helper.SpiderFailureRate(s.Transport.AccessCount(periodOfFailureSecond))
			if s.Transport.GetAccessCount() > 0 {
				failureRateAllValue += float64(s.Transport.GetFailureCount()) / float64(s.Transport.GetAccessCount()) * 100
			}

			loads[5] += s.Transport.LoadRate(5)
			loads[60] += s.Transport.LoadRate(60)
			loads[60*15] += s.Transport.LoadRate(900)
			loads[60*30] += s.Transport.LoadRate(1800)
			loads[3600] += s.Transport.LoadRate(3600)
			loads[3600*5] += s.Transport.LoadRate(3600 * 5)
			loads[36000] += s.Transport.LoadRate(36000)
			loads[86400] += s.Transport.LoadRate(86400)
			loads[86400*2] += s.Transport.LoadRate(172800)

			sleepDuration += s.GetSleep()

			if !s.RequestStartTime.IsZero() {
				waiting += time.Since(s.RequestStartTime)
			}
			serverCount++
			if !s.Stop {
				serverEnable++

				if s.FailureLevel == 0 {
					serverRun++
				}
			}

			TrafficIn += s.Transport.TrafficIn
			TrafficOut += s.Transport.TrafficOut
			NetIn += s.Transport.S.TrafficIn
			NetOut += s.Transport.S.TrafficOut
			connections += s.Transport.S.Connections
			accessCount += s.Transport.GetAccessCount()
			failureCount += s.Transport.GetFailureCount()
		}

		projectMap["stop"] = p.Stop
		projectMap["servers"] = serverCount
		projectMap["server_run"] = serverRun
		projectMap["server_enable"] = serverEnable

		projectMap["waiting"] = "0s"
		projectMap["sleep"] = "0s"
		if serverEnable > 0 {
			projectMap["sleep"] = (sleepDuration / time.Duration(serverEnable)).Truncate(time.Millisecond).String()
			if waiting != 0 {
				projectMap["waiting"] = (waiting / time.Duration(serverEnable)).Truncate(time.Millisecond).String()
			}
		}

		if serverCount > 0 {
			projectMap["failure_period"] = failureRatePeriodValue / float64(serverCount)
			projectMap["failure_period_hsl"] = strconv.Itoa(int(100 - failureRatePeriodValue/float64(serverCount)))
			projectMap["failure_all"] = strconv.FormatFloat(failureRateAllValue/float64(serverCount), 'f', 2, 64)
			projectMap["failure_all_hsl"] = strconv.Itoa(int(100 - failureRateAllValue/float64(serverCount)))
		}

		projectMap["traffic_in"] = helper.ByteCountBinary(TrafficIn)
		projectMap["traffic_out"] = helper.ByteCountBinary(TrafficOut)
		projectMap["net_in"] = helper.ByteCountBinary(NetIn)
		projectMap["net_out"] = helper.ByteCountBinary(NetOut)
		projectMap["loads"] = loads
		projectMap["connections"] = connections
		projectMap["access_count"] = accessCount
		projectMap["failure_count"] = failureCount
		projectMap["name"] = p.GetProjectName()

		jsonMap["projects"] = append(jsonMap["projects"].([]map[string]interface{}), projectMap)
	}

	responseJsonCommon(dispatchers, jsonMap, start)
	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b
}

func projectJson(p *project.Dispatcher, sType string) []byte {
	start := time.Now()
	var jsonMap = map[string]interface{}{
		"type":    sType,
		"servers": []map[string]interface{}{},
	}

	periodOfFailureSecond := helper.MinInt(int(time.Since(StartTime).Seconds()), spider.PeriodOfFailureSecond)

	for index, s := range p.GetSpiders() {
		avgTime := s.GetAvgTime()

		failureRatePeriodValue := helper.SpiderFailureRate(s.Transport.AccessCount(periodOfFailureSecond))
		failureRateAllValue := .0
		if s.Transport.GetAccessCount() > 0 {
			failureRateAllValue = float64(s.Transport.GetFailureCount()) / float64(s.Transport.GetAccessCount()) * 100
		}

		server := map[string]interface{}{}

		loads := make(map[int]float64, 9)
		loads[5] += s.Transport.LoadRate(5)
		loads[60] += s.Transport.LoadRate(60)
		loads[60*15] += s.Transport.LoadRate(900)
		loads[60*30] += s.Transport.LoadRate(1800)
		loads[3600] += s.Transport.LoadRate(3600)
		loads[3600*5] += s.Transport.LoadRate(3600 * 5)
		loads[36000] += s.Transport.LoadRate(36000)
		loads[86400] += s.Transport.LoadRate(86400)
		loads[86400*2] += s.Transport.LoadRate(172800)

		server["enable"] = !s.Stop
		server["stop"] = s.Stop
		server["loads"] = loads
		server["failure_period"] = strconv.FormatFloat(failureRatePeriodValue, 'f', 2, 64)
		server["failure_period_hsl"] = strconv.Itoa(int(100 - failureRatePeriodValue))
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
		server["net_in"] = helper.ByteCountBinary(s.Transport.S.TrafficIn)
		server["net_out"] = helper.ByteCountBinary(s.Transport.S.TrafficOut)
		server["connections"] = s.Transport.S.Connections
		server["access_count"] = s.Transport.GetAccessCount()
		server["failure_count"] = s.Transport.GetFailureCount()
		jsonMap["servers"] = append(jsonMap["servers"].([]map[string]interface{}), server)
	}

	//basic
	responseJsonCommon([]*project.Dispatcher{p}, jsonMap, start)

	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b
}
func responseJsonCommon(ps []*project.Dispatcher, jsonMap map[string]interface{}, start time.Time) {
	defer func() {
		jsonMap["basic"].(map[string]interface{})["time"] = time.Since(start).Truncate(time.Microsecond).String()
	}()

	jsonMap["basic"] = map[string]interface{}{
		"queue_bls": make(map[int]int),
	}

	pingFailureAvg := .0
	failureLevelZeroCount := 0
	var pingAvg time.Duration
	var sleepAvg time.Duration
	var waitingAvg time.Duration
	var avgTimeAvg time.Duration
	var TrafficIn uint64
	var TrafficOut uint64
	var NetIn uint64
	var NetOut uint64
	var queueCount int64
	var redisMem int64
	var serverCount int
	var serverEnable int

	loads := make(map[int]float64, 9)
	for _, p := range ps {
		for _, s := range p.GetSpiders() {
			if s.FailureLevel == 0 && !s.Stop {
				failureLevelZeroCount++
				if !s.RequestStartTime.IsZero() {
					waitingAvg += time.Since(s.RequestStartTime)
				}
				avgTimeAvg += s.GetAvgTime()
			}

			loads[5] += s.Transport.LoadRate(5)
			loads[60] += s.Transport.LoadRate(60)
			loads[60*15] += s.Transport.LoadRate(900)
			loads[60*30] += s.Transport.LoadRate(1800)
			loads[3600] += s.Transport.LoadRate(3600)
			loads[3600*5] += s.Transport.LoadRate(3600 * 5)
			loads[36000] += s.Transport.LoadRate(36000)
			loads[86400] += s.Transport.LoadRate(86400)
			loads[86400*2] += s.Transport.LoadRate(172800)

			serverCount++
			if !s.Stop {
				serverEnable++
			}
			sleepAvg += s.GetSleep()
			pingFailureAvg += s.Transport.PingFailureRate
			pingAvg += s.Transport.Ping
			TrafficIn += s.Transport.TrafficIn
			TrafficOut += s.Transport.TrafficOut
			NetIn += s.Transport.S.TrafficIn
			NetOut += s.Transport.S.TrafficOut
		}

		if len(p.GetSpiders()) > 0 {
			for i, v := range p.GetSpiders()[0].Queue.BlsTestCount {
				jsonMap["basic"].(map[string]interface{})["queue_bls"].(map[int]int)[i] += v
			}
		}

		//redis
		mem, _ := database.Redis().MemoryUsage(p.GetQueueKey()).Result() //about 1ms
		redisMem += mem
		num, _ := database.Redis().LLen(p.GetQueueKey()).Result() //about 1ms
		queueCount += num

	}
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	//basic
	jsonMap["basic"].(map[string]interface{})["sleep_avg"] = "0s"
	jsonMap["basic"].(map[string]interface{})["ping_avg"] = "0s"
	jsonMap["basic"].(map[string]interface{})["ping_failure_avg"] = ""
	jsonMap["basic"].(map[string]interface{})["avg_time_avg"] = "0s"
	jsonMap["basic"].(map[string]interface{})["waiting_avg"] = "0s"

	if serverCount > 0 {
		jsonMap["basic"].(map[string]interface{})["sleep_avg"] = (sleepAvg / time.Duration(serverCount)).Truncate(time.Millisecond).String()
		jsonMap["basic"].(map[string]interface{})["ping_avg"] = (pingAvg / time.Duration(serverCount)).Truncate(time.Millisecond).String()
		jsonMap["basic"].(map[string]interface{})["ping_failure_avg"] = strconv.FormatFloat(pingFailureAvg/float64(serverCount), 'f', 2, 64)
		jsonMap["basic"].(map[string]interface{})["loads"] = loads
	}
	if failureLevelZeroCount > 0 {
		jsonMap["basic"].(map[string]interface{})["avg_time_avg"] = (avgTimeAvg / time.Duration(failureLevelZeroCount)).Truncate(time.Millisecond).String()
		jsonMap["basic"].(map[string]interface{})["waiting_avg"] = (waitingAvg / time.Duration(failureLevelZeroCount)).Truncate(time.Millisecond).String()
	}
	jsonMap["basic"].(map[string]interface{})["servers"] = serverCount
	jsonMap["basic"].(map[string]interface{})["server_run"] = failureLevelZeroCount
	jsonMap["basic"].(map[string]interface{})["server_enable"] = serverEnable
	jsonMap["basic"].(map[string]interface{})["queue"] = queueCount
	jsonMap["basic"].(map[string]interface{})["redis_mem"] = helper.ByteCountBinary(uint64(redisMem))
	jsonMap["basic"].(map[string]interface{})["traffic_in"] = helper.ByteCountBinary(TrafficIn)
	jsonMap["basic"].(map[string]interface{})["traffic_out"] = helper.ByteCountBinary(TrafficOut)
	jsonMap["basic"].(map[string]interface{})["net_in"] = helper.ByteCountBinary(NetIn)
	jsonMap["basic"].(map[string]interface{})["net_out"] = helper.ByteCountBinary(NetOut)
	jsonMap["basic"].(map[string]interface{})["net_in_int"] = NetIn
	jsonMap["basic"].(map[string]interface{})["net_out_int"] = NetOut
	jsonMap["basic"].(map[string]interface{})["mem_sys"] = helper.ByteCountBinary(mem.Sys)
	jsonMap["basic"].(map[string]interface{})["goroutine"] = runtime.NumGoroutine()
	jsonMap["basic"].(map[string]interface{})["connections"] = helper.GetSocketEstablishedCountLazy()
	jsonMap["basic"].(map[string]interface{})["ws_connections"] = webSocketConnections

	jsonMap["basic"].(map[string]interface{})["uptime"] = time.Since(StartTime).Truncate(time.Second).String()
}

func getProjectByName(name string) *project.Dispatcher {
	for _, e := range dispatchers {
		if e.GetProjectName() == name {
			return e
		}
	}
	return nil
}

func searchSpider(projectName string, serverName string) *spider.Spider {
	for _, e := range dispatchers {
		if e.GetProjectName() == projectName {
			for _, e := range e.GetSpiders() {
				if e.Transport.S.Name == serverName {
					return e
				}
			}
			return nil
		}
	}
	return nil
}
