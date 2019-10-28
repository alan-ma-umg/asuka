package database

import (
	"github.com/chenset/asuka/helper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"sync"
	"time"
)

var mysqlOnce sync.Once
var mysqlInstance *xorm.Engine

func init() {
	go func() {
		s := time.NewTicker(time.Minute * 5)
		for {
			<-s.C

			MysqlDelayInsertQueueLock.Lock() //lock

			mSlice := make([]interface{}, len(MysqlDelayInsertQueue))
			copy(mSlice, MysqlDelayInsertQueue)
			MysqlDelayInsertQueue = []interface{}{}

			MysqlDelayInsertQueueLock.Unlock() //unlock

			//retry
			for _, m := range mSlice {
				_, err := Mysql().Insert(m)
				if err != nil {
					MysqlDelayInsertTillSuccess(m)
					log.Println(err)
				}
			}

			if len(MysqlDelayInsertQueue) > 0 {
				log.Println("MysqlDelayInsertQueue: ", len(MysqlDelayInsertQueue))
			}

		}
	}()

}

var MysqlDelayInsertQueue []interface{}
var MysqlDelayInsertQueueLock sync.Mutex

func MysqlDelayInsertTillSuccess(beans ...interface{}) {
	MysqlDelayInsertQueueLock.Lock()
	defer MysqlDelayInsertQueueLock.Unlock()
	for _, e := range beans {
		MysqlDelayInsertQueue = append(MysqlDelayInsertQueue, e)
	}
}

func Mysql() *xorm.Engine {
	mysqlOnce.Do(func() {
		engine, err := xorm.NewEngine("mysql", helper.Env().MysqlDSN)
		if err != nil {
			panic(err)
		}
		mysqlInstance = engine
	})
	return mysqlInstance
}
