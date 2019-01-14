package database

import (
	"asuka/helper"
	"github.com/willf/bloom"
	"os"
	"sync"
	"time"
)

var bloomFilterOnce sync.Once
var bloomFilterMutex = &sync.Mutex{}
var bloomFilterInstance *bloom.BloomFilter

func Bl() *bloom.BloomFilter {
	bloomFilterOnce.Do(func() {
		bloomFilterInstance = bloom.NewWithEstimates(10000000, 0.001)
		f, _ := os.Open(helper.Env().BloomFilterPath + "enqueue.db")
		bloomFilterInstance.ReadFrom(f)
		f.Close()

		//save
		go func() {
			t := time.NewTicker(time.Minute * 2)
			for {
				<-t.C
				save()
			}
		}()
	})
	return bloomFilterInstance
}

func save() {
	f, _ := os.Create(helper.Env().BloomFilterPath + "enqueue.db")
	Bl().WriteTo(f)
	f.Close()
}

func BlTestAndAddString(s string) bool {
	bloomFilterMutex.Lock()
	defer bloomFilterMutex.Unlock()
	return Bl().TestAndAddString(s)
}

func BlTestString(s string) bool {
	bloomFilterMutex.Lock()
	defer bloomFilterMutex.Unlock()
	return Bl().TestString(s)
}
