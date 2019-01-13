package queue

import (
	"fmt"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func TestNewQueue(t *testing.T) {
	q := NewQueue()
	log.Println(q.EnqueueForFailure("sdfdsfds12", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds11", 3))

	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))

	log.Println(q.EnqueueForFailure("sdfdsfds12", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds12", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds12", 3))

	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))
	log.Println(q.EnqueueForFailure("sdfdsfds1", 3))

	log.Println(q.EnqueueForFailure("sdfdsfds11", 3))

	fmt.Println(q.BlsTestCount)
	//time.Sleep(16e9)
}
