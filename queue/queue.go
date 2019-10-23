package queue

import (
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/willf/bloom"
	"log"
	"os"
	"strconv"
	"sync"
)

type Queue struct {
	name                   string
	bls                    []*bloom.BloomFilter
	BlsTestCount           map[int]int
	enqueueForFailureMutex sync.Mutex

	bloomFilterMutex          sync.Mutex
	bloomFilterInstance       *bloom.BloomFilter
	bloomFilterInstanceDoOnce *sync.Once
}

func NewQueue(name string) (q *Queue) {
	return &Queue{name: name, BlsTestCount: make(map[int]int), bloomFilterInstanceDoOnce: new(sync.Once)}
}

//ResetBloomFilterInstance purpose for release memory usage
func (my *Queue) ResetBloomFilterInstance() {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	if my.bloomFilterInstance == nil {
		return
	}

	my.BlSave(false)
	my.bloomFilterInstance = nil
	my.bloomFilterInstanceDoOnce = new(sync.Once)

	log.Println("Reset: " + my.GetKey())
}

func (my *Queue) getBloomFilterInstance() *bloom.BloomFilter {
	my.bloomFilterInstanceDoOnce.Do(func() {
		my.bloomFilterInstance = bloom.NewWithEstimates(30000000, 0.004)
		f, _ := os.Open(helper.Env().BloomFilterPath + my.GetBlKey())
		my.bloomFilterInstance.ReadFrom(f)
		f.Close()
		log.Println("New: " + my.GetKey())
	})
	return my.bloomFilterInstance
}

func (my *Queue) BlRemoveFile() {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	os.Remove(helper.Env().BloomFilterPath + my.GetBlKey())

	for i := 0; i < helper.MaxInt(10, len(my.bls)); i++ {
		os.Remove(helper.Env().BloomFilterPath + my.name + "_enqueue_retry_" + strconv.Itoa(i) + ".db")
	}
}

func (my *Queue) BlCleanUp() {
	my.BlRemoveFile()

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	for _, e := range my.bls {
		e.ClearAll()
	}

	my.getBloomFilterInstance().ClearAll()
}

//BlTestAndAddString if exists return true
func (my *Queue) BlTestAndAddString(s string) bool {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	return my.getBloomFilterInstance().TestAndAddString(s)
}

//BlTestString if exists return true
func (my *Queue) BlTestString(s string) bool {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	return my.getBloomFilterInstance().TestString(s)
}

func (my *Queue) GetBlKey() string {
	return my.GetKey() + "_bl.db"
}

func (my *Queue) GetKey() string {
	return my.name + "_" + helper.Env().Redis.URLQueueKey
}

func (my *Queue) Enqueue(rawUrl string) {
	database.Redis().RPush(my.GetKey(), rawUrl)
}

func (my *Queue) Dequeue() (string, error) {
	return database.Redis().LPop(my.GetKey()).Result()
}

func (my *Queue) GetBlsTestCount() (index, value []int) {
	//lock
	my.enqueueForFailureMutex.Lock()
	defer my.enqueueForFailureMutex.Unlock()

	for i, v := range my.BlsTestCount {
		index = append(index, i)
		value = append(value, v)
	}
	return
}

func (my *Queue) EnqueueForFailure(rawUrl string, retryTimes int) bool {
	//lock
	my.enqueueForFailureMutex.Lock()
	defer my.enqueueForFailureMutex.Unlock()

	if retryTimes < 1 {
		return false
	}

	for i := 0; i < retryTimes; i++ {
		if !my.getBl(i).TestAndAddString(rawUrl) {
			my.BlsTestCount[i]++
			my.Enqueue(rawUrl)
			return true
		}
	}

	my.BlsTestCount[-1]++
	return false
}

func (my *Queue) getBl(index int) *bloom.BloomFilter {
	for i := len(my.bls); i <= index; i++ {
		bloomFilterInstance := bloom.NewWithEstimates(10000000, 0.01)
		f, _ := os.Open(helper.Env().BloomFilterPath + my.name + "_enqueue_retry_" + strconv.Itoa(i) + ".db")
		bloomFilterInstance.ReadFrom(f)
		f.Close()

		my.bls = append(my.bls, bloomFilterInstance)
	}
	return my.bls[index]
}

func (my *Queue) BlSave(checkLock bool) {
	if checkLock {
		my.enqueueForFailureMutex.Lock()
	}
	defer func() {
		if checkLock {
			my.enqueueForFailureMutex.Unlock()
		}
	}()

	f, _ := os.Create(helper.Env().BloomFilterPath + my.GetBlKey())
	my.getBloomFilterInstance().WriteTo(f)
	f.Close()

	for i, bl := range my.bls {
		f, err := os.Create(helper.Env().BloomFilterPath + my.name + "_enqueue_retry_" + strconv.Itoa(i) + ".db")
		if err != nil {
			log.Println(err)
			continue
		}
		bl.WriteTo(f)
		f.Close()
	}
}
