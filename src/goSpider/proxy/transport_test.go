package proxy

import (
	"fmt"
	"testing"
	"time"
)

func TestTransport(t *testing.T) {
	t1, _ := NewTransport(nil)
	s := time.Now()
	for i := 0; i < 100000; i++ {
		t1.AddAccess("sfdsfsdf")
		t1.AddAccess("sfdsfsdf")
		t1.AddFailure("sfdsfsdf")
		//t1.AddFailure("sfdsfsdf")
		//
		//t1.AddAccess("sfdsfsdf")
		//if rand.Intn(2) == 2 {
		//	t1.AddFailure("sfdsfsdf")
		//	t1.AddAccess("sfdsfsdf")
		//}
		//if rand.Intn(3) == 2 {
		//	t1.AddAccess("sfdsfsdf")
		//	t1.AddAccess("sfdsfsdf")
		//	t1.AddFailure("sfdsfsdf")
		//}

		time.Sleep(1e9)
		fmt.Println(time.Since(s))
		//fmt.Println("Load: ", t1.LoadRate(5))
		fmt.Println("Load: ", t1.LoadRate(5))
		fmt.Println("Fail: ", t1.FailureRate(5))
		fmt.Println("Load: ", t1.LoadRate(120))
		fmt.Println("Fail: ", t1.FailureRate(120))
		fmt.Println("Load: ", t1.LoadRate(300))
		fmt.Println("Fail: ", t1.FailureRate(300))
		fmt.Println("Load: ", t1.LoadRate(600))
		fmt.Println("Fail: ", t1.FailureRate(600))
	}
}
