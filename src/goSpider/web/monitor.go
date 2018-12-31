package web

import (
	"net/http"
	"io"
	"goSpider/dispatcher"
	"strconv"
	"goSpider/database"
	"goSpider/helper"
)

func Monitor(dispatcher *dispatcher.Dispatcher) {
	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		html := "<meta http-equiv=\"refresh\" content=\"1\"><style>th,td{border:1px solid #ccc}</style><table><tr><th>Server Address</th><th>Load Balance</th><th>Load Rate 5s</th><th>Load Rate 60s</th><th>Load Rate 5m</th><th>Load Rate 15m</th><th>Dispatcher Count</th><th>Access Count</th><th>Failure Count</th></tr>"

		avgLoad := 0.0
		for _, s := range dispatcher.GetSpiders() {
			avgLoad += s.Transport.LoadRate(5)
			serAddr := s.Transport.S.ServerAddr
			if serAddr == "" {
				serAddr = "Localhost"
			}
			html += "<tr>"
			html += "<td>" + serAddr + " </td><td> " + strconv.FormatFloat(s.Transport.LoadBalanceRate(), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(5), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60*5), 'f', 2, 64) + "</td><td> " + strconv.FormatFloat(s.Transport.LoadRate(60*15), 'f', 2, 64) + "</td><td>" + strconv.Itoa(s.Transport.LoopCount) + "</td><td>" + strconv.Itoa(len(s.Transport.AccessList)) + "</td><td>" + strconv.Itoa(len(s.Transport.FailureList)) + "</td>"
			html += "</tr>"
		}

		queueCount, _ := database.Redis().LLen(helper.Env().Redis.URLQueueKey).Result()
		html += "</table> queue: " + strconv.Itoa(int(queueCount)) + "<br> Avg Load:" + strconv.FormatFloat(avgLoad/float64(len(dispatcher.GetSpiders())), 'f', 2, 64)

		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		io.WriteString(w, html)
	})
	http.ListenAndServe(":88", nil)
}
