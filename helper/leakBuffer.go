package helper

import (
	"sync"
)

//https://golang.org/doc/effective_go.html#leaky_buffer
type leakyBuf struct {
	freeList chan []byte
}

const bufSize = 4 + 12 + 4096

var newLeakBufOnce sync.Once
var leakyBufCache *leakyBuf

func LeakyBuf() *leakyBuf {
	newLeakBufOnce.Do(func() {
		leakyBufCache = &leakyBuf{make(chan []byte, 2048)}
	})
	return leakyBufCache
}

// Get Grab a buffer if available; allocate if not.
func (lb *leakyBuf) Get() (buf []byte) {
	// Grab a buffer if available; allocate if not.
	select {
	case buf = <-lb.freeList:
		// Got one; nothing more to do.
	default:
		// None free, so allocate a new one.
		buf = make([]byte, bufSize)
	}
	return buf
}

func (lb *leakyBuf) Put(buf []byte) {
	if len(buf) != bufSize {
		panic("invalid buffer size that's put into leaky buffer")
	}
	// Reuse buffer if there's room.
	select {
	case lb.freeList <- buf:
		// Buffer on free list; nothing more to do.
	default:
		// Free list full, just carry on.
	}
}
