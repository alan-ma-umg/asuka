package web

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/chenset/asuka/project"
	"github.com/chenset/asuka/proxy"
	"github.com/chenset/asuka/queue"
	"github.com/chenset/asuka/spider"
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
var mem runtime.MemStats
var tcpFilterDoOnceInDuration = helper.NewDoOnceInDuration(time.Second*6 + 234*time.Millisecond)
var tcpFilterDoOnceInDurationCache *queue.Cmd20Response
var AlwaysEmptyTcpFilterDoOnceInDuration = &queue.Cmd20Response{} //not nil

func Server(d []*project.Dispatcher, address string) error {
	dispatchers = d

	//init start time
	go func() {
		for _, d := range dispatchers {
			if StartTime.Unix() > d.StartTime.Unix() {
				StartTime = d.StartTime
			}
		}
	}()

	http.HandleFunc("/log", commonHandleFunc(fileLog))
	http.HandleFunc("/add/", commonHandleFunc(addServer))
	http.HandleFunc("/get/", commonHandleFunc(getServer))
	http.HandleFunc("/website/", commonHandleFunc(projectWebsite))
	http.HandleFunc("/queue/", commonHandleFunc(redisQueue))
	http.HandleFunc("/login", commonHandleFunc(login))
	http.HandleFunc("/logout", commonHandleFunc(logout))
	http.HandleFunc("/login/post", commonHandleFunc(loginPost))
	http.HandleFunc("/project.io", projectIO)
	http.HandleFunc("/index.io", indexIO)
	http.HandleFunc("/switchProject", commonHandleFunc(switchProject))
	http.HandleFunc("/switchServer", commonHandleFunc(switchServer))
	http.HandleFunc("/forever/", forever)
	http.HandleFunc("/", commonHandleFunc(home))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/favicon.ico")
	})
	http.Handle("/static/", commonHandle(http.StripPrefix("/static", http.FileServer(http.Dir("web/templates/static")))))

	return http.ListenAndServe(address, nil)
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func commonHandleFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//intent no httpOnly, js need it
		if r.Method == "GET" && strings.ToLower(r.URL.String()) != "/login" && r.Header.Get("X-Requested-With") == "" {
			http.SetCookie(w, &http.Cookie{Name: "intent", Value: r.URL.String(), Path: "/", Expires: time.Now().Add(time.Hour), HttpOnly: false})
		}

		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html")
		}
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Server", "Asuka")

		//gzip
		if !strings.Contains(strings.ToLower(r.Header.Get("Accept-Encoding")), "gzip") {
			fn(w, r)
			return
		}
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

func switchServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST Required", 405)
		return
	}

	post := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	//login check
	if !authCheckOrRedirect(w, r) {
		return
	}

	//login check
	//if cookie, err := r.Cookie("id"); err != nil || !authCheck(cookie.Value) {
	//	http.Error(w, "Login Required", 401)
	//	return
	//}

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

	b, _ := json.Marshal(jsonMap)
	w.Header().Set("Content-type", "application/json")
	w.Write(b)
}

func switchProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST Required", 405)
		return
	}

	post := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	//login check
	if !authCheckOrRedirect(w, r) {
		return
	}

	//login check
	//if cookie, err := r.Cookie("id"); err != nil || !authCheck(cookie.Value) {
	//	http.Error(w, "Login Required", 401)
	//	return
	//}

	if _, ok := post["name"]; !ok {
		http.NotFound(w, r)
		return
	}

	p := getProjectByName(post["name"])
	if p == nil {
		http.NotFound(w, r)
		return
	}

	//p.Stop = !p.IsStop()
	if p.StopTime.IsZero() {
		p.StopTime = time.Now().Add(-time.Second)
	} else {
		p.StopTime = time.Time{}
	}

	jsonMap := map[string]interface{}{
		"success": true,
		"data":    post["name"],
	}

	b, _ := json.Marshal(jsonMap)
	w.Header().Set("Content-type", "application/json")
	w.Write(b)
}

func index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		GOOS        string
		Check       bool
		PreloadJson template.JS
	}{
		GOOS: runtime.GOOS,
	}

	//login check
	if cookie, err := r.Cookie("id"); err == nil {
		data.Check = authCheck(cookie.Value)
	}

	data.PreloadJson = template.JS(indexJson(data.Check))

	template.Must(template.ParseFiles("web/templates/index.html")).Execute(w, data)
}
func fileLog(w http.ResponseWriter, r *http.Request) {
	//login check
	if !authCheckOrRedirect(w, r) {
		return
	}

	GetFileLogInstance().UpdateLogCheckTime()

	template.Must(template.ParseFiles("web/templates/log.html")).Execute(w, map[string]interface{}{
		"log": string(GetFileLogInstance().TailFile(102400)),
		"mod": GetFileLogInstance().GetLogModifyTime().Format(time.Stamp),
	})
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
		Check       bool
		PreloadJson template.JS
	}{
		GOOS:        runtime.GOOS,
		ProjectName: p.Name(),
	}

	//login check
	if cookie, err := r.Cookie("id"); err == nil {
		data.Check = authCheck(cookie.Value)
	}

	data.PreloadJson = template.JS(projectJson(data.Check, p, "home"))

	template.Must(template.ParseFiles("web/templates/project.html")).Execute(w, data)
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

	check := false
	//login check
	if cookie, err := r.Cookie("id"); err == nil {
		check = authCheck(cookie.Value)
	}

	var sleepSecondTimes int64 = 1
	for {
		messageType, b, err := c.ReadMessage()
		if err != nil {
			break
		}
		if messageType == 1 && check {
			input := strings.TrimSpace(string(b))
			switch input {
			case "free":
				debug.FreeOSMemory()
				fmt.Println("debug.FreeOsMemory")
			case "stop":
				for _, d := range dispatchers {
					for _, s := range d.GetSpiders() {
						if s != nil {
							s.Stop = true
						}
					}
				}
				fmt.Println("spider stop")
			case "start":
				for _, d := range dispatchers {
					for _, s := range d.GetSpiders() {
						if s != nil {
							s.Stop = false
						}
					}
				}
				fmt.Println("spider start")
			default:
				if speedInt, err := strconv.ParseInt(input, 10, 64); err == nil && speedInt > 0 {
					sleepSecondTimes = helper.MaxInt64(speedInt, 1)
				}
			}
		}

		err = c.WriteMessage(websocket.TextMessage, indexJson(check))
		if err != nil {
			//log.Println("write:", err)
			break
		}
		time.Sleep(time.Second * time.Duration(sleepSecondTimes))
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

	check := false
	//login check
	if cookie, err := r.Cookie("id"); err == nil {
		check = authCheck(cookie.Value)
	}

	webSocketConnections++

	defer func() {
		webSocketConnections--
		c.Close()
	}()

	responseContent := "home"
	var recentFetchIndex int64 = 0
	var sleepSecondTimes int64 = 1
	for {
		messageType, b, err := c.ReadMessage()
		if err != nil {
			//log.Println("read:", err)
			break
		}
		if messageType == 1 {
			input := strings.TrimSpace(string(b))
			switch input {
			case "free":
				debug.FreeOSMemory()
				fmt.Println("debug.FreeOsMemory")
			case "enqueue":
				if check {
					for _, l := range p.EntryUrl() {
						p.GetQueue().Enqueue(l)
					}
				}
				fmt.Println("enqueue")
			case "clear":
				if check {
					p.GetQueue().BlCleanUp()
					break
				}
				fmt.Println("bloomFilter clearAll")
			case "stop":
				if check {
					for _, s := range p.GetSpiders() {
						if s != nil {
							s.Stop = true
						}
					}
				}
				fmt.Println("spider stop")
			case "start":
				if check {
					for _, s := range p.GetSpiders() {
						if s != nil {
							s.Stop = false
						}
					}
				}
				fmt.Println("spider start")
			case "home":
				responseContent = strings.TrimSpace(string(b))
			case "recent":
				responseContent = strings.TrimSpace(string(b))
			default:
				if speedInt, err := strconv.ParseInt(input, 10, 64); err == nil && speedInt > 0 {
					sleepSecondTimes = helper.MaxInt64(speedInt, 1)
				}
			}
		}

		switch responseContent {
		case "home":
			err = c.WriteMessage(websocket.TextMessage, projectJson(check, p, responseContent))
		case "recent":
			jsonRes, n := recentJson(check, p, responseContent, recentFetchIndex)
			recentFetchIndex = n
			err = c.WriteMessage(websocket.TextMessage, jsonRes)
		}
		if err != nil {
			//log.Println("write:", err)
			break
		}
		time.Sleep(time.Second * time.Duration(sleepSecondTimes))
	}
}

func getDispatcher(name string) *project.Dispatcher {
	for _, d := range dispatchers {
		if name == d.Name() {
			return d
		}
	}

	return nil
}

func projectWebsite(w http.ResponseWriter, r *http.Request) {
	//login check
	if !authCheckOrRedirect(w, r) {
		return
	}

	ps := strings.Split(r.URL.Path, "/")
	if len(ps) != 3 {
		http.NotFound(w, r)
		return
	}

	p := getDispatcher(ps[2])
	if p == nil {
		http.NotFound(w, r)
		return
	}

	p.WEBSite(w, r)
}

func getServer(w http.ResponseWriter, r *http.Request) {
	//login check
	if !authCheckOrRedirect(w, r) {
		return
	}

	//if cookie, err := r.Cookie("id"); err != nil || !authCheck(cookie.Value) {
	//	http.Error(w, "Login Required", 401)
	//	return
	//}

	ps := strings.Split(r.URL.Path, "/")
	if len(ps) != 3 {
		http.NotFound(w, r)
		return
	}

	p := getDispatcher(ps[2])
	if p == nil {
		http.NotFound(w, r)
		return
	}

	//string.jo
	buf := &bytes.Buffer{}
	for _, e := range p.GetSpiders() {
		if e != nil {
			buf.WriteString(e.TransportUrl.String())
			buf.WriteString("<br>")
		}
	}
	w.Write(buf.Bytes())
}

func addServer(w http.ResponseWriter, r *http.Request) {
	//login check
	if !authCheckOrRedirect(w, r) {
		return
	}
	//if cookie, err := r.Cookie("id"); err != nil || !authCheck(cookie.Value) {
	//	http.Error(w, "Login Required", 401)
	//	return
	//}

	ps := strings.Split(r.URL.Path, "/")
	if len(ps) != 3 {
		http.NotFound(w, r)
		return
	}

	p := getDispatcher(ps[2])
	if p == nil {
		http.NotFound(w, r)
		return
	}

	oldSpiderCount := len(p.GetSpiders())

	if r.Method == "POST" {
		addServerPost(w, r, p)
	}

	template.Must(template.ParseFiles("web/templates/addServer.html")).Execute(w, struct {
		ProjectName     string
		FormValueServer string
		FormValueType   string
		AddNum          int
		Method          string
	}{
		ProjectName:     p.Name(),
		FormValueServer: strings.TrimSpace(r.FormValue("servers")),
		FormValueType:   strings.TrimSpace(r.FormValue("type")),
		AddNum:          len(p.GetSpiders()) - oldSpiderCount,
		Method:          r.Method,
	})
}

func addServerPost(_ http.ResponseWriter, r *http.Request, dispatcher *project.Dispatcher) {
	switch r.FormValue("type") {
	case "url":
		for _, line := range strings.Split(strings.TrimSpace(strings.Replace(strings.Replace(r.FormValue("servers"), "\r\n", "\n", len(r.FormValue("servers"))), "\r", "\n", len(r.FormValue("servers")))), "\n") {
			line = strings.ToLower(strings.TrimSpace(line))
			if strings.HasPrefix(line, "http") {
				for _, addr := range proxy.HttpProxyParse("http", line) {
					dispatcher.AddSpider(addr)
				}
			} else if strings.HasPrefix(line, "https") {
				for _, addr := range proxy.HttpProxyParse("https", line) {
					dispatcher.AddSpider(addr)
				}
			} else if strings.HasPrefix(line, "socks5") {
				for _, addr := range proxy.Socks5ProxyParse(line) {
					dispatcher.AddSpider(addr)
				}
			}
		}
	case "https":
		for _, addr := range proxy.HttpProxyParse("https", strings.TrimSpace(r.FormValue("servers"))) {
			dispatcher.AddSpider(addr)
		}
	case "http":
		for _, addr := range proxy.HttpProxyParse("http", strings.TrimSpace(r.FormValue("servers"))) {
			dispatcher.AddSpider(addr)
		}
	case "socks5":
		for _, addr := range proxy.Socks5ProxyParse(strings.TrimSpace(r.FormValue("servers"))) {
			dispatcher.AddSpider(addr)
		}
	}

	return
}

func login(w http.ResponseWriter, _ *http.Request) {
	template.Must(template.ParseFiles("web/templates/login.html")).Execute(w, nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST Required", 405)
		return
	}

	if cookie, err := r.Cookie("id"); err == nil {
		if authCheck(cookie.Value) {
			database.Redis().Del(cookie.Value)
		}
	}

	//clean login session
	http.SetCookie(w, &http.Cookie{Name: "id", Value: "", Path: "/", Expires: time.Unix(0, 0), HttpOnly: true})

	w.Header().Set("Content-type", "application/json")
	jsonMap := map[string]interface{}{}
	jsonMap["success"] = true
	jsonMap["message"] = "success"
	b, _ := json.Marshal(jsonMap)
	w.Write(b)
}

func loginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST Required", 405)
		return
	}

	post := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	if _, ok := post["password"]; !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-type", "application/json")
	jsonMap := map[string]interface{}{}
	if post["password"] == helper.Env().WEBPassword {
		jsonMap["success"] = true
		jsonMap["message"] = "success"
		jsonMap["url"] = "/"

		expireDuration := time.Hour * 24 * 7
		id, _ := helper.Enc([]byte(helper.Env().WEBPassword))
		database.Redis().Set(id, helper.Env().WEBPassword, expireDuration)
		//set login session
		http.SetCookie(w, &http.Cookie{Name: "id", Value: id, Path: "/", Expires: time.Now().Add(expireDuration), MaxAge: 0, HttpOnly: true})

		//intent
		if cookie, err := r.Cookie("intent"); err == nil {
			jsonMap["url"] = cookie.Value
		}

		//clean intent
		http.SetCookie(w, &http.Cookie{Name: "intent", Value: "", Path: "/", Expires: time.Unix(0, 0), HttpOnly: true})

		//send message to wx
		ip := ""
		if r.Header.Get("CF-Connecting-IP") != "" {
			ip += r.Header.Get("CF-Connecting-IP") + " <=> "
		}
		if r.Header.Get("X-Forwarded-For") != "" {
			ip += r.Header.Get("X-Forwarded-For") + " <=> "
		}
		ip += r.RemoteAddr

		helper.SendTextToWXDoOnceDurationHour("Asuka login: " + ip + "\n" + r.UserAgent() + "\n" + time.Now().Format("2006-01-02 15:04:05"))
	} else {
		jsonMap["success"] = false
		jsonMap["message"] = "Password incorrect"
	}

	b, _ := json.Marshal(jsonMap)
	w.Write(b)
}

func redisQueue(w http.ResponseWriter, r *http.Request) {
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

	//login check
	if !authCheckOrRedirect(w, r) {
		return
	}

	list, _ := database.Redis().LRange(p.GetQueueKey(), 0, 1000).Result()
	data := struct {
		List        []string
		ProjectName string
	}{
		List:        list,
		ProjectName: p.Name(),
	}

	template.Must(template.ParseFiles("web/templates/queue.html")).Execute(w, data)
}

func forever(w http.ResponseWriter, _ *http.Request) {
	str := ""
	for i := 0; i < rand.Intn(4); i++ {
		str += "<a href=\"/forever/" + strconv.Itoa(rand.Int()) + "\">" + strconv.Itoa(i) + "</a>"
	}
	w.Header().Set("Content-type", "text/html")
	io.WriteString(w, str)
}

func recentJson(check bool, p *project.Dispatcher, sType string, recentFetchIndex int64) ([]byte, int64) {
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
			ll := *l
			if !check {
				ll.TransportName = ""
				ll.RawUrl = ""
			}
			jsonMap["fetched"] = append(jsonMap["fetched"].([]*spider.Summary), &ll)
			lastIndex = helper.MaxInt64(lastIndex, l.Index)
		}
	}

	responseJsonCommon(check, []*project.Dispatcher{p}, jsonMap, start)
	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b, helper.MaxInt64(lastIndex, recentFetchIndex)
}

func indexJson(check bool) []byte {
	start := time.Now()
	var jsonMap = map[string]interface{}{
		"projects": []map[string]interface{}{},
	}

	periodOfFailureSecond := helper.MinInt(int(time.Since(StartTime).Seconds()), spider.PeriodOfFailureSecond)

	for _, p := range dispatchers {
		projectMap := map[string]interface{}{}

		loads := make(map[int]float64, 10)

		failureRatePeriodValue := 0.0
		var sleepDuration time.Duration
		var waiting time.Duration
		var TrafficIn uint64
		var TrafficOut uint64
		//var NetIn uint64
		//var NetOut uint64
		//var connections int
		var accessCount int
		var failureCount int
		var serverCount int
		var serverRun int
		var serverEnable int

		loads[5] += p.LoadRate(5)
		loads[60] += p.LoadRate(60)
		loads[60*15] += p.LoadRate(900)
		loads[60*30] += p.LoadRate(1800)
		loads[3600] += p.LoadRate(3600)
		loads[3600*6] += p.LoadRate(3600 * 6)
		loads[3600*12] += p.LoadRate(3600 * 12)
		loads[86400] += p.LoadRate(86400)
		loads[86400*2] += p.LoadRate(86400 * 2)
		loads[86400*3] += p.LoadRate(86400 * 3)

		failureRatePeriodValue += helper.SpiderFailureRate(p.AccessCount(periodOfFailureSecond))
		accessCount += p.GetAccessCount()
		failureCount += p.GetFailureCount()
		TrafficIn += p.TrafficIn
		TrafficOut += p.TrafficOut

		for _, s := range p.GetSpiders() {
			if s != nil {
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
				//NetIn += s.Transport.S.TrafficIn
				//NetOut += s.Transport.S.TrafficOut
				//connections += s.Transport.S.Connections
			}
		}

		projectMap["stop"] = p.IsStop()
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

		projectMap["failure_period"] = failureRatePeriodValue
		projectMap["failure_period_hsl"] = strconv.Itoa(int(100 - failureRatePeriodValue))
		projectMap["failure_all"] = strconv.FormatFloat(float64(p.GetFailureCount())/float64(p.GetAccessCount())*100, 'f', 2, 64)
		projectMap["failure_all_hsl"] = strconv.Itoa(int(100 - float64(p.GetFailureCount())/float64(p.GetAccessCount())*100))

		projectMap["traffic_in"] = helper.ByteCountBinary(TrafficIn)
		projectMap["traffic_out"] = helper.ByteCountBinary(TrafficOut)
		//projectMap["net_in"] = helper.ByteCountBinary(NetIn)
		//projectMap["net_out"] = helper.ByteCountBinary(NetOut)
		projectMap["loads"] = loads
		//projectMap["connections"] = connections
		projectMap["access_count"] = accessCount
		projectMap["failure_count"] = failureCount
		projectMap["name"] = p.Name()

		jsonMap["projects"] = append(jsonMap["projects"].([]map[string]interface{}), projectMap)
	}

	responseJsonCommon(check, dispatchers, jsonMap, start)
	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b
}

func projectJson(check bool, p *project.Dispatcher, sType string) []byte {
	start := time.Now()
	var jsonMap = map[string]interface{}{
		"type":    sType,
		"servers": []map[string]interface{}{},
	}

	periodOfFailureSecond := helper.MinInt(int(time.Since(StartTime).Seconds()), spider.PeriodOfFailureSecond)

	spiders := p.GetSpiders()
	if len(spiders) > 1 {
		if spiders[len(spiders)-1].TransportUrl.Scheme == "direct" {
			spiders = append([]*spider.Spider{spiders[len(spiders)-1]}, spiders[0:helper.MinInt(101, len(spiders))-1]...)
		} else {
			spiders = spiders[0:helper.MinInt(100, len(spiders))]
		}
	}
	for index, s := range spiders {
		//avgTime := s.GetAvgTime()

		failureRateAllValue := .0
		failureRatePeriodValue := .0
		failureRatePeriodValue = helper.SpiderFailureRate(s.AccessCount(periodOfFailureSecond))
		if s.GetAccessCount() > 0 {
			failureRateAllValue = float64(s.GetFailureCount()) / float64(s.GetAccessCount()) * 100
		}

		server := map[string]interface{}{}

		loads := make(map[int]float64, 9)
		loads[5] += s.LoadRate(5)
		loads[60] += s.LoadRate(60)
		loads[60*15] += s.LoadRate(900)
		//loads[60*30] += s.LoadRate(1800)
		//loads[3600] += s.LoadRate(3600)
		//loads[3600*6] += s.LoadRate(3600 * 6)
		//loads[3600*12] += s.LoadRate(3600 * 12)
		//loads[86400] += s.LoadRate(86400)
		//loads[86400*2] += s.LoadRate(86400 * 2)
		//loads[86400*3] += s.LoadRate(86400 * 3)
		server["loads"] = loads

		server["enable"] = !s.Stop
		server["stop"] = s.Stop
		server["idle"] = s.IsIdle()
		//server["proxy_status"] = s.Transport.S.Status
		server["failure_period"] = strconv.FormatFloat(failureRatePeriodValue, 'f', 2, 64)
		server["failure_period_hsl"] = strconv.Itoa(int(100 - failureRatePeriodValue))
		server["failure_all"] = strconv.FormatFloat(failureRateAllValue, 'f', 2, 64)
		server["failure_all_hsl"] = strconv.Itoa(int(100 - failureRateAllValue))
		server["failure_level"] = s.FailureLevel
		server["failure_level_hsl"] = 100 - s.FailureLevel
		server["index"] = index
		if check {
			server["name"] = s.TransportUrl.Host
		} else {
			server["name"] = ""
		}
		//if s.Transport.Ping == 0 {
		//	server["ping"] = "-"
		//} else {
		//	server["ping"] = s.Transport.Ping.Truncate(time.Millisecond).String()
		//}
		//server["ping_hsl"] = helper.MinInt(150, helper.MaxInt(150-int(s.Transport.Ping.Seconds()*1000/2), 0))
		//server["ping_failure"] = strconv.FormatFloat(s.Transport.PingFailureRate*100, 'f', 0, 64)
		//server["ping_failure_hsl"] = int(150 - s.Transport.PingFailureRate*150)
		//server["avg_time"] = avgTime.Truncate(time.Millisecond).String()
		server["sleep"] = s.GetSleep().Truncate(time.Millisecond).String()
		server["waiting"] = "0s"
		if !s.RequestStartTime.IsZero() {
			server["waiting"] = time.Since(s.RequestStartTime).Truncate(time.Millisecond).String()
		}
		//server["traffic_in"] = helper.ByteCountBinary(s.Transport.TrafficIn)
		//server["traffic_out"] = helper.ByteCountBinary(s.Transport.TrafficOut)
		//server["net_in"] = helper.ByteCountBinary(s.Transport.S.TrafficIn)
		//server["net_out"] = helper.ByteCountBinary(s.Transport.S.TrafficOut)
		//server["connections"] = s.Transport.S.Connections
		server["access_count"] = s.GetAccessCount()
		server["failure_count"] = s.GetFailureCount()
		jsonMap["servers"] = append(jsonMap["servers"].([]map[string]interface{}), server)
	}

	//basic
	responseJsonCommon(check, []*project.Dispatcher{p}, jsonMap, start)

	b, _ := json.Marshal(jsonMap)
	return b
}

func responseJsonCommon(check bool, ps []*project.Dispatcher, jsonMap map[string]interface{}, start time.Time) {
	defer func() {
		jsonMap["basic"].(map[string]interface{})["time"] = time.Since(start).Truncate(time.Microsecond).String()
	}()

	jsonMap["basic"] = map[string]interface{}{
		"queue_bls":  make(map[int]int),
		"tcp_filter": make(map[string]interface{}),
	}

	periodOfFailureSecond := helper.MinInt(int(time.Since(StartTime).Seconds()), spider.PeriodOfFailureSecond)
	//pingFailureAvg := .0
	failureLevelZeroCount := 0
	//var pingAvg time.Duration
	var sleepAvg time.Duration
	var waitingAvg time.Duration
	//var avgTimeAvg time.Duration
	var TrafficIn uint64
	var TrafficOut uint64
	//var NetIn uint64
	//var NetOut uint64
	var queueCount int64
	var redisMem int64
	var serverCount int
	var serverEnable int
	var accessCount int
	var failureCount int
	failureRatePeriodValue := 0.0

	loads := make(map[int]float64, 10)
	for _, p := range ps {
		loads[5] += p.LoadRate(5)
		loads[60] += p.LoadRate(60)
		loads[60*15] += p.LoadRate(900)
		loads[60*30] += p.LoadRate(1800)
		loads[3600] += p.LoadRate(3600)
		loads[3600*6] += p.LoadRate(3600 * 6)
		loads[3600*12] += p.LoadRate(3600 * 12)
		loads[86400] += p.LoadRate(86400)
		loads[86400*2] += p.LoadRate(86400 * 2)
		loads[86400*3] += p.LoadRate(86400 * 3)
		failureRatePeriodValue += helper.SpiderFailureRate(p.AccessCount(periodOfFailureSecond))
		accessCount += p.GetAccessCount()
		failureCount += p.GetFailureCount()
		TrafficIn += p.TrafficIn
		TrafficOut += p.TrafficOut
		for _, s := range p.GetSpiders() {
			if s != nil {
				if s.FailureLevel == 0 && !s.Stop {
					failureLevelZeroCount++
					if !s.RequestStartTime.IsZero() {
						waitingAvg += time.Since(s.RequestStartTime)
					}
					//avgTimeAvg += s.GetAvgTime()
				}

				serverCount++
				if !s.Stop {
					serverEnable++
				}
				sleepAvg += s.GetSleep()
				//pingFailureAvg += s.Transport.PingFailureRate
				//pingAvg += s.Transport.Ping

				//NetIn += s.Transport.S.TrafficIn
				//NetOut += s.Transport.S.TrafficOut
			}
		}

		if len(p.GetSpiders()) > 0 {
			indexSlice, valueSlice := p.GetQueue().GetBlsTestCount()
			for i, v := range indexSlice {
				jsonMap["basic"].(map[string]interface{})["queue_bls"].(map[int]int)[v] += valueSlice[i]
			}
		}

		//redis
		mem, _ := database.Redis().MemoryUsage(p.GetQueueKey()).Result() //about 1ms
		redisMem += mem
		num, _ := database.Redis().LLen(p.GetQueueKey()).Result() //about 1ms
		queueCount += num

		if check {
			jsonMap["basic"].(map[string]interface{})["showing"] = p.Showing()
		} else {
			jsonMap["basic"].(map[string]interface{})["showing"] = ""
		}
	}

	runtime.ReadMemStats(&mem)

	//_, memAvailable := helper.GetMemInfoFromProc()

	sysLoad := ""
	sysMemInfo := ""
	sysMemInfoPercent := .0
	if check && runtime.GOOS != "windows" {
		sysLoad = helper.GetSystemLoadFromProc()
		availableMemByte, totalMemByte := helper.GetMemInfoFromProc()
		sysMemInfoPercent = float64(totalMemByte-availableMemByte) / float64(totalMemByte) * 100
		sysMemInfo = helper.ByteCountBinary(totalMemByte-availableMemByte) + "/" + helper.ByteCountBinary(totalMemByte)
	}

	if check && helper.Env().BloomFilterClient != "" {
		tcpFilterDoOnceInDuration.Do(func() {
			go func() {
				reportBuf, err := queue.GetTcpFilterInstance().Cmd(20, nil)
				if err != nil {
					log.Println(err)
					return
				}
				json.Unmarshal(reportBuf, &tcpFilterDoOnceInDurationCache) //must be struct instead of map in this case
			}()
		})

		jsonMap["basic"].(map[string]interface{})["tcp_filter"] = tcpFilterDoOnceInDurationCache
	} else {
		jsonMap["basic"].(map[string]interface{})["tcp_filter"] = AlwaysEmptyTcpFilterDoOnceInDuration
	}

	jsonMap["basic"].(map[string]interface{})["log_mod"] = 0
	jsonMap["basic"].(map[string]interface{})["log_check"] = 0
	if check {
		jsonMap["basic"].(map[string]interface{})["log_mod"] = GetFileLogInstance().GetLogModifyTime().Unix()
		jsonMap["basic"].(map[string]interface{})["log_check"] = GetFileLogInstance().GetLogCheckTime().Unix()
	}

	jsonMap["basic"].(map[string]interface{})["filter_new_connections"] = queue.GetTcpFilterInstance().NewConnectionCount
	jsonMap["basic"].(map[string]interface{})["pool_size"] = queue.GetTcpFilterInstance().ConnPoolSize()
	//basic
	jsonMap["basic"].(map[string]interface{})["failure_period"] = strconv.FormatFloat(failureRatePeriodValue, 'f', 2, 64)
	jsonMap["basic"].(map[string]interface{})["sleep_avg"] = "0s"
	//jsonMap["basic"].(map[string]interface{})["ping_avg"] = "0s"
	//jsonMap["basic"].(map[string]interface{})["ping_failure_avg"] = ""
	//jsonMap["basic"].(map[string]interface{})["avg_time_avg"] = "0s"
	jsonMap["basic"].(map[string]interface{})["waiting_avg"] = "0s"

	if serverCount > 0 {
		jsonMap["basic"].(map[string]interface{})["sleep_avg"] = (sleepAvg / time.Duration(serverCount)).Truncate(time.Millisecond).String()
		//jsonMap["basic"].(map[string]interface{})["ping_avg"] = (pingAvg / time.Duration(serverCount)).Truncate(time.Millisecond).String()
		//jsonMap["basic"].(map[string]interface{})["ping_failure_avg"] = strconv.FormatFloat(pingFailureAvg/float64(serverCount), 'f', 2, 64)
		jsonMap["basic"].(map[string]interface{})["loads"] = loads
	}
	if failureLevelZeroCount > 0 {
		//jsonMap["basic"].(map[string]interface{})["avg_time_avg"] = (avgTimeAvg / time.Duration(failureLevelZeroCount)).Truncate(time.Millisecond).String()
		jsonMap["basic"].(map[string]interface{})["waiting_avg"] = (waitingAvg / time.Duration(failureLevelZeroCount)).Truncate(time.Millisecond).String()
	}
	jsonMap["basic"].(map[string]interface{})["servers"] = serverCount
	jsonMap["basic"].(map[string]interface{})["server_run"] = failureLevelZeroCount
	jsonMap["basic"].(map[string]interface{})["server_enable"] = serverEnable
	jsonMap["basic"].(map[string]interface{})["queue"] = queueCount
	jsonMap["basic"].(map[string]interface{})["redis_mem"] = helper.ByteCountBinary(uint64(redisMem))
	jsonMap["basic"].(map[string]interface{})["traffic_in"] = helper.ByteCountBinary(TrafficIn)
	jsonMap["basic"].(map[string]interface{})["traffic_out"] = helper.ByteCountBinary(TrafficOut)
	//jsonMap["basic"].(map[string]interface{})["net_in"] = helper.ByteCountBinary(NetIn)
	//jsonMap["basic"].(map[string]interface{})["net_out"] = helper.ByteCountBinary(NetOut)
	//jsonMap["basic"].(map[string]interface{})["net_in_int"] = NetIn
	//jsonMap["basic"].(map[string]interface{})["net_out_int"] = NetOut
	jsonMap["basic"].(map[string]interface{})["mem_sys"] = helper.ByteCountBinary(helper.GetProgramRss())
	jsonMap["basic"].(map[string]interface{})["sys_load"] = sysLoad
	jsonMap["basic"].(map[string]interface{})["sys_mem"] = sysMemInfo
	jsonMap["basic"].(map[string]interface{})["sys_mem_percent"] = sysMemInfoPercent
	jsonMap["basic"].(map[string]interface{})["runtime_mem"] = helper.ByteCountBinary(mem.Sys) + "/" + helper.ByteCountBinary(mem.Alloc)
	jsonMap["basic"].(map[string]interface{})["goroutine"] = runtime.NumGoroutine()
	jsonMap["basic"].(map[string]interface{})["connections"] = helper.GetSocketEstablishedCountLazy()
	jsonMap["basic"].(map[string]interface{})["ws_connections"] = webSocketConnections

	jsonMap["basic"].(map[string]interface{})["access_count"] = accessCount
	jsonMap["basic"].(map[string]interface{})["failure_count"] = failureCount

	jsonMap["basic"].(map[string]interface{})["uptime"] = helper.TimeSince(time.Since(StartTime))

}

func getProjectByName(name string) *project.Dispatcher {
	for _, e := range dispatchers {
		if e.Name() == name {
			return e
		}
	}
	return nil
}

func searchSpider(projectName string, serverName string) *spider.Spider {
	for _, e := range dispatchers {
		if e.Name() == projectName {
			for _, e := range e.GetSpiders() {
				if e != nil && e.TransportUrl.Host == serverName {
					return e
				}
			}
			return nil
		}
	}
	return nil
}

func authCheck(id string) bool {
	if res, err := database.Redis().Get(id).Result(); err == nil && res == helper.Env().WEBPassword {
		return true
	}
	return false
}

func authCheckOrRedirect(w http.ResponseWriter, r *http.Request) bool {
	if cookie, err := r.Cookie("id"); err != nil || !authCheck(cookie.Value) {
		if strings.ToLower(r.Header.Get("X-Requested-With")) == "xmlhttprequest" {
			http.Error(w, "Login Required", 401)
			return false
		}

		http.Redirect(w, r, "/login", 302)
		return false
	}

	return true
}
