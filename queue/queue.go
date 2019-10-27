package queue

import (
	"encoding/json"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/willf/bloom"
	"log"
	"os"
	"sync"
)

const (
	mainBlSize = 30000000
)

type Queue struct {
	Retries                   []int
	name                      string
	bloomFilterMutex          sync.Mutex
	bloomFilterInstance       *bloom.BloomFilter
	bloomFilterInstanceDoOnce *sync.Once
}

func NewQueue(name string) (q *Queue) {
	return &Queue{name: name, Retries: make([]int, 1), bloomFilterInstanceDoOnce: new(sync.Once)}
}

//ResetBloomFilterInstance purpose for release memory usage
func (my *Queue) ResetBloomFilterInstance() {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	if my.bloomFilterInstance == nil {
		return
	}

	//save
	my.BlSave(false)

	//main
	my.bloomFilterInstance = nil

	//new once
	my.bloomFilterInstanceDoOnce = new(sync.Once)

	log.Println("Reset: " + my.GetKey())
}

func (my *Queue) getBloomFilterInstance() *bloom.BloomFilter {
	my.bloomFilterInstanceDoOnce.Do(func() {
		my.bloomFilterInstance = bloom.NewWithEstimates(mainBlSize, 0.004)
		f, _ := os.Open(my.mainBlFilename())
		my.bloomFilterInstance.ReadFrom(f)
		f.Close()
		log.Println("New: " + my.GetKey())
	})
	return my.bloomFilterInstance
}

func (my *Queue) GetBlKey() string {
	return my.GetKey() + "_bl"
}

func (my *Queue) GetKey() string {
	return my.name + "_" + helper.Env().Redis.URLQueueKey
}

func (my *Queue) GetFailureKey() string {
	return my.name + "_fail_" + helper.Env().Redis.URLQueueKey
}

func (my *Queue) mainBlFilename() string {
	return helper.Env().BloomFilterPath + my.GetBlKey() + ".db"
}

func (my *Queue) BlRemoveFile() {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	os.Remove(my.mainBlFilename())
}

func (my *Queue) BlCleanUp() {
	if helper.Env().BloomFilterClient != "" {
		go func() {
			my.bloomFilterMutex.Lock()
			defer my.bloomFilterMutex.Unlock()
			//main
			GetTcpFilterInstance().Cmd(11, &Cmd11{Db: my.GetBlKey()})
		}()
		return
	}

	my.BlRemoveFile()

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	my.getBloomFilterInstance().ClearAll()
}

func (my *Queue) blTcp(db string, size uint, fun byte, s string) (res bool) {
	buf, err := GetTcpFilterInstance().Cmd(10, &Cmd10{
		Db:   db,
		Size: size,
		Fun:  fun,
		Urls: []string{s},
	})
	if err != nil {
		TcpErrorPrintDoOnce.Do(func() {
			log.Println(err)
		})
		return true
	}

	var result []byte
	json.Unmarshal(buf, &result)
	if len(result) == 0 || result[0] == 1 {
		return true
	}

	return false
}

//BlTestString if exists return true
func (my *Queue) BlTestString(s string) bool {
	if helper.Env().BloomFilterClient != "" {
		return my.blTcp(my.GetBlKey(), mainBlSize, 10, s)
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	return my.getBloomFilterInstance().TestString(s)
}

//BlTestAndAddString if exists return true
func (my *Queue) BlTestAndAddString(s string) bool {
	if helper.Env().BloomFilterClient != "" {
		return my.blTcp(my.GetBlKey(), mainBlSize, 20, s)
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	return my.getBloomFilterInstance().TestAndAddString(s)
}

func (my *Queue) Enqueue(rawUrl string) {
	database.Redis().RPush(my.GetKey(), rawUrl)
}

func (my *Queue) Dequeue() (string, error) {
	return database.Redis().LPop(my.GetKey()).Result()
}

func (my *Queue) EnqueueForFailure(rawUrl string, retryTimes int) bool {
	retryTimes *= 3 //control by project

	incrInt := int(database.Redis().HIncrBy(my.GetFailureKey(), rawUrl, 1).Val()) //incr failure

	if incrInt > retryTimes {
		//final failure
		my.Retries[0]++
		return false
	}

	if len(my.Retries) <= incrInt {
		my.Retries = append(my.Retries, 0) //put 0 instead of 1
	}
	my.Retries[incrInt]++
	my.Enqueue(rawUrl)
	return true
}

func (my *Queue) BlSave(checkLock bool) {
	if checkLock {
		my.bloomFilterMutex.Lock()
		defer my.bloomFilterMutex.Unlock()
	}

	if my.bloomFilterInstance != nil {
		f, _ := os.Create(my.mainBlFilename())
		my.getBloomFilterInstance().WriteTo(f)
		f.Close()
	}
}
