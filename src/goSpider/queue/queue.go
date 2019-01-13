package queue

import (
	"github.com/willf/bloom"
	"goSpider/database"
	"goSpider/helper"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Queue struct {
	bls                    []*bloom.BloomFilter
	BlsTestCount           map[int]int
	enqueueForFailureMutex *sync.Mutex
}

func NewQueue() (q *Queue) {
	q = &Queue{enqueueForFailureMutex: &sync.Mutex{}, BlsTestCount: make(map[int]int)}

	go func(q *Queue) {
		t := time.NewTicker(time.Minute * 3)
		for {
			<-t.C
			q.blSave()
		}
	}(q)
	return
}

func (my *Queue) Enqueue(rawUrl string) {
	database.Redis().RPush(helper.Env().Redis.URLQueueKey, rawUrl)
}

func (my *Queue) Dequeue() (string, error) {
	return database.Redis().LPop(helper.Env().Redis.URLQueueKey).Result()
}

func (my *Queue) EnqueueForFailure(rawUrl string, retryTimes int) bool {
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

	return false
}

func (my *Queue) getBl(index int) *bloom.BloomFilter {
	for i := len(my.bls); i <= index; i++ {
		bloomFilterInstance := bloom.NewWithEstimates(1000000, 0.01)
		f, _ := os.Open(helper.Env().BloomFilterPath + "enqueue_retry_" + strconv.Itoa(i) + ".db")
		bloomFilterInstance.ReadFrom(f)
		f.Close()

		my.bls = append(my.bls, bloomFilterInstance)
	}
	return my.bls[index]
}

func (my *Queue) blSave() {
	for i, bl := range my.bls {
		f, err := os.Create(helper.Env().BloomFilterPath + "enqueue_retry_" + strconv.Itoa(i) + ".db")
		if err != nil {
			log.Println(err)
			continue
		}
		bl.WriteTo(f)
		f.Close()
	}
}
