package web

import (
	"net/http"
	"strconv"
	"io"
	"math/rand"
)

func Forever() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		str := ""
		for i := 0; i < 10; i++ {
			str += "<a href=\"/" + strconv.Itoa(rand.Int()) + "\">" + strconv.Itoa(i) + "</a>"
			str += "<a href=\"/" + strconv.Itoa(i) + "\">" + strconv.Itoa(i) + "</a>"
		}

		io.WriteString(w, str)
	})
	http.ListenAndServe(":888", nil)
}
