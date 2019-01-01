package web

import (
	"net/http"
	"goSpider/dispatcher"
	"strconv"
	"goSpider/database"
	"goSpider/helper"
	"runtime"
	"fmt"
	"os"
	"time"
	"log"
	"github.com/gorilla/websocket"
	"html/template"
	"math/rand"
	"io"
)

var upgrade = websocket.Upgrader{
	EnableCompression: true,
}
var startTime = time.Now()
var webSocketConnections = 0
var pwd, _ = os.Getwd()

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
		time.Sleep(1e9)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	MonitorHtml.Execute(w, "ws://"+r.Host+"/echo")
}

var dispatcherObj *dispatcher.Dispatcher

func Server(d *dispatcher.Dispatcher, address string) {
	dispatcherObj = d
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)

	http.HandleFunc("/forever", func(w http.ResponseWriter, r *http.Request) {
		str := ""
		for i := 0; i < rand.Intn(4); i++ {
			str += "<a href=\"/" + strconv.Itoa(rand.Int()) + "\">" + strconv.Itoa(i) + "</a>"
		}
		io.WriteString(w, str)
	})
	log.Fatal(http.ListenAndServe(address, nil))
}

var MonitorHtml = template.Must(template.ParseFiles(pwd + "/src/goSpider/web/templates/index.html"))

func html() string {
	html := "<style>th,td{border:1px solid #ccc}</style><table><tr><th>Server Address</th><th>Load Balance</th><th>Load Rate 5s</th><th>Load Rate 60s</th><th>Load Rate 5m</th><th>Load Rate 15m</th><th>Dispatcher Count</th><th>Access Count</th><th>Failure Count</th></tr>"
	start := time.Now()
	avgLoad := 0.0
	for _, s := range dispatcherObj.GetSpiders() {
		avgLoad += s.Transport.LoadRate(5)
		serAddr := s.Transport.S.ServerAddr
		if serAddr == "" {
			serAddr = "Localhost"
		}
		html += "<tr>"
		html += "<td>" + serAddr + " </td><td> " + strconv.FormatFloat(s.Transport.LoadBalanceRate(), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(5), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60*5), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60*15), 'f', 2, 64) + "</td><td>" + strconv.Itoa(s.Transport.LoopCount) + "</td><td>" + strconv.Itoa(s.Transport.GetAccessCount()) + "</td><td>" + strconv.Itoa(s.Transport.GetFailureCount()) + "</td>"
		html += "</tr>"
	}

	queueCount, _ := database.Redis().LLen(helper.Env().Redis.URLQueueKey).Result()

	//memory
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	//redis memory
	redisMem, err := database.Redis().MemoryUsage(helper.Env().Redis.URLQueueKey).Result()
	if err != nil {
		fmt.Println(err)
		redisMem = 0
	}

	//bloomFilter
	var fileSize int64 = 0
	fi, err := os.Stat(helper.Env().BloomFilterFile)
	if err == nil {
		fileSize = fi.Size()
	}

	html += "</table> queue: " + strconv.Itoa(int(queueCount)) + "<br> Redis mem: " + strconv.FormatFloat(helper.B2Mb(uint64(redisMem)), 'f', 2, 64) + " Mb<br>"
	html += "BloomFilter: " + strconv.FormatFloat(helper.B2Mb(uint64(fileSize)), 'f', 2, 64) + " Mb"
	html += "<br> Avg Load:" + strconv.FormatFloat(avgLoad/float64(len(dispatcherObj.GetSpiders())), 'f', 2, 64) + "</br>"
	html += "Alloc: " + strconv.FormatFloat(helper.B2Mb(mem.Alloc), 'f', 2, 64) + "Mb <br> TotalAlloc: " + strconv.FormatFloat(helper.B2Mb(mem.Alloc), 'f', 2, 64) + "Mb <br> Sys: " + strconv.FormatFloat(helper.B2Mb(mem.Sys), 'f', 2, 64) + "Mb <br>"
	html += "NumGC: " + strconv.Itoa(int(mem.NumGC)) + " <br> NumGoroutine: " + strconv.Itoa(runtime.NumGoroutine()) + "<br>"
	html += "webSocketConnections: " + strconv.Itoa(webSocketConnections) + "<br>"
	html += "time: " + time.Since(start).String() + "   " + time.Since(startTime).String()

	return html
}
