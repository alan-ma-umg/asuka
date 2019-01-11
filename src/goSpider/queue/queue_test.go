package queue

import (
	"log"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func TestNewQueue(t *testing.T) {
	q := NewQueue()
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 2))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 4))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 5))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 6))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 6))
	log.Println(q.EnqueueForFailure("sdfdsfds2", 3))
	time.Sleep(16e9)
}
