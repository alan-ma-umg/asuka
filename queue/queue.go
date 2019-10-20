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

	bloomFilterMutex    sync.Mutex
	bloomFilterInstance *bloom.BloomFilter
}

func NewQueue(name string) (q *Queue) {
	q = &Queue{name: name, BlsTestCount: make(map[int]int)}
	q.bloomFilterInstance = bloom.NewWithEstimates(30000000, 0.004)
	f, _ := os.Open(helper.Env().BloomFilterPath + q.GetBlKey())
	q.bloomFilterInstance.ReadFrom(f)
	f.Close()
	return
}

func (my *Queue) BlCleanUp() {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	os.Remove(helper.Env().BloomFilterPath + my.GetBlKey())

	for i := 0; i < helper.MaxInt(10, len(my.bls)); i++ {
		os.Remove(helper.Env().BloomFilterPath + my.name + "_enqueue_retry_" + strconv.Itoa(i) + ".db")
	}

	for _, e := range my.bls {
		e.ClearAll()
	}

	my.bloomFilterInstance.ClearAll()
}

//BlTestAndAddString if exists return true
func (my *Queue) BlTestAndAddString(s string) bool {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	return my.bloomFilterInstance.TestAndAddString(s)
}

//BlTestString if exists return true
func (my *Queue) BlTestString(s string) bool {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()
	return my.bloomFilterInstance.TestString(s)
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

func (my *Queue) BlSave() {
	f, _ := os.Create(helper.Env().BloomFilterPath + my.GetBlKey())
	my.bloomFilterInstance.WriteTo(f)
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
