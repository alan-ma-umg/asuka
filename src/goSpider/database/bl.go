package database

import (
	"github.com/willf/bloom"
	"goSpider/helper"
	"os"
	"sync"
	"time"
)

var bloomFilterOnce sync.Once
var bloomFilterMutex = &sync.Mutex{}
var bloomFilterInstance *bloom.BloomFilter

func init() {
	//save
	go func() {
		t := time.NewTicker(time.Minute * 2)
		for {
			<-t.C
			save()
		}
	}()
}

func Bl() *bloom.BloomFilter {
	bloomFilterOnce.Do(func() {
		bloomFilterInstance = bloom.NewWithEstimates(10000000, 0.001)
		f, _ := os.Open(helper.Env().BloomFilterFile)
		bloomFilterInstance.ReadFrom(f)
		f.Close()
	})
	return bloomFilterInstance
}

func save() {
	f, _ := os.Create(helper.Env().BloomFilterFile)
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
