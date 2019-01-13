package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"goSpider/helper"
	"sync"
)

var mysqlOnce sync.Once
var mysqlInstance *xorm.Engine

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
