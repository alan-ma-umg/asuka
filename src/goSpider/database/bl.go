package database

import (
	"github.com/willf/bloom"
	"goSpider/helper"
	"os"
	"sync"
	"time"
)

var bloomFilterOnce sync.Once
var bloomFilterInstance *bloom.BloomFilter

func init() {
	//save
	go func() {
		t := time.NewTicker(time.Second * 5)
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
