package spider

import (
	"goSpider/database"
	"goSpider/helper"
	"goSpider/project"
	"testing"
)

func TestDispatcher_Run(t *testing.T) {

	database.Bl().ClearAll()                             //todo for test
	database.Redis().Del(helper.Env().Redis.URLQueueKey) //todo for test

	c := &Dispatcher{}
	c.Run(&project.DouBan{})
}
