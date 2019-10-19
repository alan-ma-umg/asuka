package database

//import (
//	"github.com/go-xorm/xorm"
//	_ "github.com/mattn/go-sqlite3"
//	"sync"
//)
//
//var sqliteOnce sync.Once
//var sqliteInstance *xorm.Engine
//
//func Sqlite() *xorm.Engine {
//	sqliteOnce.Do(func() {
//		engine, err := xorm.NewEngine("sqlite3", "./sqlite.db?_synchronous=OFF&_journal_mode=WAL")
//		if err != nil {
//			panic(err)
//		}
//		sqliteInstance = engine
//	})
//	return sqliteInstance
//}
