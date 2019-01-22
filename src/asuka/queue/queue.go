package queue

import (
	"asuka/database"
	"asuka/helper"
	"fmt"
	"github.com/willf/bloom"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Queue struct {
	name                   string
	bls                    []*bloom.BloomFilter
	BlsTestCount           map[int]int
	enqueueForFailureMutex *sync.Mutex
}

func NewQueue(name string) (q *Queue) {
	q = &Queue{name: name, enqueueForFailureMutex: &sync.Mutex{}, BlsTestCount: make(map[int]int)}

	// kill signal handing
	helper.ExitHandleFuncSlice = append(helper.ExitHandleFuncSlice, func() {
		q.BlSave()
		fmt.Println(q.name + " bls saved")
	})

	go func(q *Queue) {
		t := time.NewTicker(time.Minute * 6)
		for {
			<-t.C
			q.BlSave()
		}
	}(q)
	return
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
		bloomFilterInstance := bloom.NewWithEstimates(1000000, 0.01)
		f, _ := os.Open(helper.Env().BloomFilterPath + my.name + "_enqueue_retry_" + strconv.Itoa(i) + ".db")
		bloomFilterInstance.ReadFrom(f)
		f.Close()

		my.bls = append(my.bls, bloomFilterInstance)
	}
	return my.bls[index]
}

func (my *Queue) BlSave() {
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
