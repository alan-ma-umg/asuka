package spider

import (
	"goSpider/database"
	"goSpider/helper"
	"net/url"
	"strings"
	"testing"
)

func TestDispatcher_Run(t *testing.T) {

	database.Bl().ClearAll()                             //todo for test
	database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	c := &Dispatcher{}
	c.Run([]string{"https://www.douban.com/"}, func(s *Spider, l *url.URL) bool {
		pass := false
		for _, white := range []string{"movie.douban.com", "book.douban.com"} {
			if strings.Contains(strings.ToLower(l.Hostname()), white) {
				pass = true
			}
		}
		return pass
	})
}
