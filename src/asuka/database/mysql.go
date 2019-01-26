package database

import (
	"asuka/helper"
	"fmt"
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

			MysqlDelayInsertTillSuccessQueueLock.Lock() //lock

			mSlice := make([]interface{}, len(MysqlDelayInsertTillSuccessQueue))
			copy(mSlice, MysqlDelayInsertTillSuccessQueue)
			MysqlDelayInsertTillSuccessQueue = []interface{}{}

			MysqlDelayInsertTillSuccessQueueLock.Unlock() //unlock

			//retry
			for _, m := range mSlice {
				_, err := Mysql().Insert(m)
				if err != nil {
					MysqlDelayInsertTillSuccess(m)
					log.Println(err)
				}
			}

			if len(MysqlDelayInsertTillSuccessQueue) > 0 {
				fmt.Println("MysqlDelayInsertTillSuccessQueue: ", len(MysqlDelayInsertTillSuccessQueue))
			}

		}
	}()

}

var MysqlDelayInsertTillSuccessQueue []interface{}
var MysqlDelayInsertTillSuccessQueueLock sync.Mutex

func MysqlDelayInsertTillSuccess(beans ...interface{}) {
	MysqlDelayInsertTillSuccessQueueLock.Lock()
	defer MysqlDelayInsertTillSuccessQueueLock.Unlock()
	for _, e := range beans {
		MysqlDelayInsertTillSuccessQueue = append(MysqlDelayInsertTillSuccessQueue, e)
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
