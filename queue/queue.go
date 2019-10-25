package queue

import (
	"encoding/json"
	"github.com/chenset/asuka/database"
	"github.com/chenset/asuka/helper"
	"github.com/willf/bloom"
	"log"
	"os"
	"strconv"
	"sync"
)

const (
	mainBlSize         = 30000000
	retriesBlsSizeBase = 3000000
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

	if my.bloomFilterInstance == nil && len(my.bls) == 0 {
		return
	}

	//save
	my.BlSave(false)

	//main
	my.bloomFilterInstance = nil

	//retries
	my.bls = nil
	my.BlsTestCount = make(map[int]int)

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

func (my *Queue) GetBlsKey(i int) string {
	return my.name + "_enqueue_retry_" + strconv.Itoa(i)
}

func (my *Queue) mainBlFilename() string {
	return helper.Env().BloomFilterPath + my.GetBlKey() + ".db"
}

func (my *Queue) blsFilename(i int) string {
	return helper.Env().BloomFilterPath + my.GetBlsKey(i) + ".db"
}

func (my *Queue) BlRemoveFile() {
	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	os.Remove(my.mainBlFilename())

	for i := 0; i < helper.MaxInt(10, len(my.bls)); i++ {
		os.Remove(my.blsFilename(i))
	}
}

func (my *Queue) BlCleanUp() {
	if helper.Env().BloomFilterClient != "" {
		return
	}

	my.BlRemoveFile()

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	for _, e := range my.bls {
		e.ClearAll()
	}

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
		log.Println(err)
		return true
	}

	var result []byte
	json.Unmarshal(buf, &result)
	if len(result) == 0 || result[0] == 1 {
		return true
	}

	return false
}

func (my *Queue) getBlsRetriesBlSize(i int) uint {
	i *= 3

	if i >= 9 {
		i = 9
	}

	return uint(retriesBlsSizeBase * (10. - i) / 10.)
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
	if retryTimes < 1 {
		return false
	}

	for i := 0; i < retryTimes; i++ {
		res := false
		if helper.Env().BloomFilterClient == "" {
			my.enqueueForFailureMutex.Lock()
			res = my.getBl(i).TestAndAddString(rawUrl)
			my.enqueueForFailureMutex.Unlock()
		} else {
			res = my.blTcp(my.GetBlsKey(i), my.getBlsRetriesBlSize(i), 20, rawUrl)
		}

		if !res {
			my.enqueueForFailureMutex.Lock()
			my.BlsTestCount[i]++
			my.enqueueForFailureMutex.Unlock()
			my.Enqueue(rawUrl)
			return true
		}
	}
	my.enqueueForFailureMutex.Lock()
	my.BlsTestCount[-1]++
	my.enqueueForFailureMutex.Unlock()
	return false
}

func (my *Queue) getBl(index int) *bloom.BloomFilter {
	for i := len(my.bls); i <= index; i++ {
		bloomFilterInstance := bloom.NewWithEstimates(my.getBlsRetriesBlSize(i), 0.01)
		f, _ := os.Open(my.blsFilename(i))
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

	if my.bloomFilterInstance != nil {
		f, _ := os.Create(my.mainBlFilename())
		my.getBloomFilterInstance().WriteTo(f)
		f.Close()
	}

	for i, bl := range my.bls {
		f, err := os.Create(my.blsFilename(i))
		if err != nil {
			log.Println(err)
			continue
		}
		bl.WriteTo(f)
		f.Close()
	}
}
