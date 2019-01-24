package database

import (
	"asuka/helper"
	"fmt"
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
		bloomFilterInstance = bloom.NewWithEstimates(20000000, 0.001)
		f, _ := os.Open(helper.Env().BloomFilterPath + "enqueue.db")
		bloomFilterInstance.ReadFrom(f)
		f.Close()

		// kill signal handing
		helper.ExitHandleFuncSlice = append(helper.ExitHandleFuncSlice, func() {
			blSave()
			fmt.Println("bl saved")
		})

		//save
		go func() {
			t := time.NewTicker(time.Minute * 5)
			for {
				<-t.C
				blSave()
			}
		}()
	})
	return bloomFilterInstance
}

func blSave() {
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
