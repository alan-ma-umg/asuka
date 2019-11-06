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

func (my *Queue) blTcp(db string, size uint, fun byte, s []string) (res []bool, err error) {
	buf, err := GetTcpFilterInstance().Cmd(10, &Cmd10{
		Db:   db,
		Size: size,
		Fun:  fun,
		Urls: s,
	})
	if err != nil {
		TcpErrorPrintDoOnce.Do(func() {
			log.Println(err)
		})
		//失败一定次数后停止项目
		return res, err
	}

	var result []byte
	err = json.Unmarshal(buf, &result)
	for _, b := range result {
		if b == 1 {
			res = append(res, true)
		} else {
			res = append(res, false)
		}
	}

	return
}

//BlTestString if exists return true
//func (my *Queue) BlTestString(s string) (bool, error) {
//	if helper.Env().BloomFilterClient != "" {
//		return my.blTcp(my.GetBlKey(), my.bloomFilterSize, 10, []string{s})
//	}
//
//	my.bloomFilterMutex.Lock()
//	defer my.bloomFilterMutex.Unlock()
//	return my.getBloomFilterInstance().TestString(s), nil
//}

//BlTestAndAddStrings if exists return true
func (my *Queue) BlTestAndAddStrings(s []string) (res []bool, err error) {
	if helper.Env().BloomFilterClient != "" {
		return my.blTcp(my.GetBlKey(), my.bloomFilterSize, 20, s)
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	for _, e := range s {
		res = append(res, my.getBloomFilterInstance().TestAndAddString(e))
	}

	return
}

//BlTestAndAddString if exists return true
func (my *Queue) BlTestAndAddString(s string) (bool, error) {
	if helper.Env().BloomFilterClient != "" {
		res, err := my.blTcp(my.GetBlKey(), my.bloomFilterSize, 20, []string{s})
		if len(res) == 0 {
			return true, err
		}

		return res[0], err
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	return my.getBloomFilterInstance().TestAndAddString(s), nil
}

func (my *Queue) Enqueue(values interface{}) {
	database.Redis().RPush(my.GetKey(), values)
}

func (my *Queue) QueueLen() int64 {
	return database.Redis().LLen(my.GetKey()).Val()
}

func (my *Queue) Dequeue() (string, error) {
	return database.Redis().LPop(my.GetKey()).Result()
}

func (my *Queue) EnqueueForFailure(retryEnqueueUrl, spiderEnqueueUrl string, retryTimes int) (success bool, tries int) {
	retryTimes *= 3 //todo control by project

	incrInt := int(database.Redis().HIncrBy(my.GetFailureKey(), retryEnqueueUrl, 1).Val()) //incr failure
	if incrInt > retryTimes {
		//final failure
		return false, incrInt
	}
	my.Enqueue(spiderEnqueueUrl)
	return true, incrInt
}

func (my *Queue) DequeueForFailure(rawUrl string) {
	database.Redis().HDel(my.GetFailureKey(), rawUrl)
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
