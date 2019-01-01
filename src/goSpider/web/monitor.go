package web

import (
	"net/http"
	"io"
	"goSpider/dispatcher"
	"strconv"
	"goSpider/database"
	"goSpider/helper"
	"runtime"
	"fmt"
	"os"
	"time"
)

func Monitor(dispatcher *dispatcher.Dispatcher, startTime time.Time) {
	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		html := "<meta http-equiv=\"refresh\" content=\"1\"><style>th,td{border:1px solid #ccc}</style><table><tr><th>Server Address</th><th>Load Balance</th><th>Load Rate 5s</th><th>Load Rate 60s</th><th>Load Rate 5m</th><th>Load Rate 15m</th><th>Dispatcher Count</th><th>Access Count</th><th>Failure Count</th></tr>"
		start := time.Now()
		avgLoad := 0.0
		for _, s := range dispatcher.GetSpiders() {
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
		html += "<br> Avg Load:" + strconv.FormatFloat(avgLoad/float64(len(dispatcher.GetSpiders())), 'f', 2, 64) + "</br>"
		html += "Alloc: " + strconv.FormatFloat(helper.B2Mb(mem.Alloc), 'f', 2, 64) + "Mb <br> TotalAlloc: " + strconv.FormatFloat(helper.B2Mb(mem.Alloc), 'f', 2, 64) + "Mb <br> Sys: " + strconv.FormatFloat(helper.B2Mb(mem.Sys), 'f', 2, 64) + "Mb <br>"
		html += "NumGC: " + strconv.Itoa(int(mem.NumGC)) + " <br> NumGoroutine: " + strconv.Itoa(runtime.NumGoroutine()) + "<br>"
		html += "time: " + time.Since(start).String() + "   " + time.Since(startTime).String()

		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		io.WriteString(w, html)
	})
	http.ListenAndServe(":88", nil)
}
