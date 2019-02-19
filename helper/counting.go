package helper

import (
	"math"
	"sync"
)

const SecondInterval = 1
const MinuteInterval = 40 * 60
const CountQueueSecondCap = MinuteInterval * 2

type Counting struct {
	countSliceMutex         sync.RWMutex
	CountSliceCursor        int
	AccessCountSecondSlice  []uint32
	FailureCountSecondSlice []uint32
	AccessCountMinuteSlice  []uint32
	FailureCountMinuteSlice []uint32
	AccessCountHistory      int
	FailureCountHistory     int
	CountQueueMinuteCap     int
}

// AddAccess 每次调用请求时增加一次记录, 无论是否成功
func (my *Counting) AddAccess() {
	my.AccessCountHistory++
}

// AddFailure 每次调用请求并失败时增加一次失败记录
func (my *Counting) AddFailure() {
	my.FailureCountHistory++
}

func (my *Counting) GetAccessCount() int {
	return my.AccessCountHistory
}

func (my *Counting) GetFailureCount() int {
	return my.FailureCountHistory
}

//updateCountQueueCap
func (my *Counting) updateCountQueueCap(second int) {
	if my.CountQueueMinuteCap <= second/MinuteInterval {
		my.CountQueueMinuteCap = second/MinuteInterval + 1
	}
}

// LoadRate 获取指定秒数内的负载值.参数最小值SecondInterval秒
func (my *Counting) LoadRate(second int) (rate float64) {
	//Read lock
	defer func() {
		my.countSliceMutex.RUnlock()
	}()
	my.countSliceMutex.RLock()

	my.updateCountQueueCap(second)

	sliceLen := len(my.AccessCountSecondSlice)
	if sliceLen == 0 || second == 0 {
		return
	}

	times := int(math.Ceil(float64(second) / SecondInterval))

	//SecondInterval
	if times <= CountQueueSecondCap {
		return float64(my.AccessCountSecondSlice[sliceLen-1]-my.AccessCountSecondSlice[MaxInt(sliceLen-times-1, 0)]) / float64(times)
	}

	minuteSliceLen := len(my.AccessCountMinuteSlice)
	minSecond := MinInt(second, minuteSliceLen*MinuteInterval+my.CountSliceCursor%(MinuteInterval))
	realTimeSecond := minSecond%(MinuteInterval) + MinuteInterval
	rate += float64(my.AccessCountSecondSlice[sliceLen-1] - my.AccessCountSecondSlice[MaxInt(sliceLen-realTimeSecond-1, 0)])
	minSecond -= realTimeSecond
	if minSecond > 0 {
		rate += float64(my.AccessCountMinuteSlice[minuteSliceLen-1] - my.AccessCountMinuteSlice[minuteSliceLen-minSecond/(MinuteInterval)-1])
	}
	return rate / float64(times)
}

//AccessCount  获取指定秒数内的访问数/失败j数量.参数最小值SecondInterval秒
func (my *Counting) AccessCount(second int) (accessTimes, failureTimes int) {
	//Read lock
	defer func() {
		my.countSliceMutex.RUnlock()
	}()
	my.countSliceMutex.RLock()

	my.updateCountQueueCap(second)

	accessSliceLen := len(my.AccessCountSecondSlice)
	failureSliceLen := len(my.FailureCountSecondSlice)

	if accessSliceLen+failureSliceLen == 0 {
		return
	}

	times := int(math.Ceil(float64(second) / SecondInterval))
	if times == 0 {
		return
	}

	if times <= CountQueueSecondCap {
		if accessSliceLen != 0 {
			accessTimes = int(my.AccessCountSecondSlice[accessSliceLen-1] - my.AccessCountSecondSlice[MaxInt(accessSliceLen-times-1, 0)])
		}

		if failureSliceLen != 0 {
			failureTimes = int(my.FailureCountSecondSlice[failureSliceLen-1] - my.FailureCountSecondSlice[MaxInt(failureSliceLen-times-1, 0)])
		}
		return
	}

	minuteSliceLen := len(my.AccessCountMinuteSlice) //len(my.FailureCountMinuteSlice)/
	minSecond := MinInt(second, minuteSliceLen*MinuteInterval+my.CountSliceCursor%(MinuteInterval))
	realTimeSecond := minSecond%(MinuteInterval) + MinuteInterval
	minSecond -= realTimeSecond

	accessTimes += int(my.AccessCountSecondSlice[accessSliceLen-1] - my.AccessCountSecondSlice[MaxInt(accessSliceLen-realTimeSecond-1, 0)])
	if minSecond > 0 {
		accessTimes += int(my.AccessCountMinuteSlice[minuteSliceLen-1] - my.AccessCountMinuteSlice[minuteSliceLen-minSecond/(MinuteInterval)-1])
	}

	failureTimes += int(my.FailureCountSecondSlice[accessSliceLen-1] - my.FailureCountSecondSlice[MaxInt(accessSliceLen-realTimeSecond-1, 0)])
	if minSecond > 0 {
		failureTimes += int(my.FailureCountMinuteSlice[minuteSliceLen-1] - my.FailureCountMinuteSlice[minuteSliceLen-minSecond/(MinuteInterval)-1])
	}

	return
}

func (my *Counting) RecordAccessSecondCount() {
	//Write lock
	defer func() {
		my.countSliceMutex.Unlock()
	}()
	my.countSliceMutex.Lock()
	//slice fifo
	my.AccessCountSecondSlice = append(my.AccessCountSecondSlice[MaxInt(len(my.AccessCountSecondSlice)-CountQueueSecondCap, 0):], uint32(my.GetAccessCount()))

	if my.CountSliceCursor%MinuteInterval == 0 {
		my.AccessCountMinuteSlice = append(my.AccessCountMinuteSlice[MaxInt(len(my.AccessCountMinuteSlice)-my.CountQueueMinuteCap, 0):], my.AccessCountSecondSlice[len(my.AccessCountSecondSlice)-MinuteInterval])
	}
}

func (my *Counting) RecordFailureSecondCount() {
	//Write lock
	defer func() {
		my.countSliceMutex.Unlock()
	}()
	my.countSliceMutex.Lock()

	//slice fifo
	my.FailureCountSecondSlice = append(my.FailureCountSecondSlice[MaxInt(len(my.FailureCountSecondSlice)-CountQueueSecondCap, 0):], uint32(my.GetFailureCount()))

	if my.CountSliceCursor%MinuteInterval == 0 {
		my.FailureCountMinuteSlice = append(my.FailureCountMinuteSlice[MaxInt(len(my.FailureCountMinuteSlice)-my.CountQueueMinuteCap, 0):], my.FailureCountSecondSlice[len(my.FailureCountSecondSlice)-MinuteInterval])
	}
}
