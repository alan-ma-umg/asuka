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

type Queue struct {
	name                      string
	bloomFilterMutex          sync.Mutex
	bloomFilterInstance       *bloom.BloomFilter
	bloomFilterInstanceDoOnce *sync.Once
	bloomFilterSize           uint
}

func NewQueue(name string, size uint) (q *Queue) {
	return &Queue{name: name, bloomFilterInstanceDoOnce: new(sync.Once), bloomFilterSize: size}
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
		my.bloomFilterInstance = bloom.NewWithEstimates(my.bloomFilterSize, 0.004)
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

func (my *Queue) blTcp(db string, size uint, fun byte, s string) (res bool, err error) {
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
		//失败一定次数后停止项目
		return true, err
	}

	var result []byte
	json.Unmarshal(buf, &result)
	if len(result) == 0 || result[0] == 1 {
		return true, nil
	}

	return false, nil
}

//BlTestString if exists return true
func (my *Queue) BlTestString(s string) (bool, error) {
	if helper.Env().BloomFilterClient != "" {
		return my.blTcp(my.GetBlKey(), my.bloomFilterSize, 10, s)
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	return my.getBloomFilterInstance().TestString(s), nil
}

//BlTestAndAddString if exists return true
func (my *Queue) BlTestAndAddString(s string) (bool, error) {
	if helper.Env().BloomFilterClient != "" {
		return my.blTcp(my.GetBlKey(), my.bloomFilterSize, 20, s)
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	return my.getBloomFilterInstance().TestAndAddString(s), nil
}

func (my *Queue) Enqueue(rawUrl string) {
	database.Redis().RPush(my.GetKey(), rawUrl)
}

func (my *Queue) QueueLen() int64 {
	return database.Redis().LLen(my.GetKey()).Val()
}

func (my *Queue) Dequeue() (string, error) {
	return database.Redis().LPop(my.GetKey()).Result()
}

func (my *Queue) EnqueueForFailure(rawUrl string, retryTimes int) (success bool, tries int) {
	retryTimes *= 3 //todo control by project

	incrInt := int(database.Redis().HIncrBy(my.GetFailureKey(), rawUrl, 1).Val()) //incr failure
	if incrInt > retryTimes {
		//final failure
		return false, incrInt
	}
	my.Enqueue(rawUrl)
	return true, incrInt
}

func (my *Queue) CleanFailure() {
	database.Redis().Del(my.GetFailureKey())
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
