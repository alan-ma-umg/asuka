package proxy

import (
	"asuka/helper"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestTransport(t *testing.T) {
	t1, _ := NewTransport(&SsAddr{})

	t1.LoadRate(1360000)

	for i := 0; i < 1000000; i++ {
		t1.AddAccess("sfdsfsdf")
		t1.AddAccess("sfdsfsdf")
		t1.AddAccess("sfdsfsdf")
		t1.AddFailure("sfdsfsdf")
		//t1.AddFailure("sfdsfsdf")

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

		//time.Sleep(0.1e9)
		//fmt.Println(time.Since(s))
		//fmt.Println("Load: ", t1.LoadRate(5))
		//fmt.Println("Load: ", t1.LoadRate(5))
		//fmt.Println("Fail: ", t1.FailureRate(5))
		//fmt.Println("Load: ", t1.LoadRate(120))
		//fmt.Println("Fail: ", t1.FailureRate(120))
		//fmt.Println("Load: ", t1.LoadRate(300))
		//fmt.Println("Fail: ", t1.FailureRate(300))
		//fmt.Println("Load: ", t1.LoadRate(600))
		//fmt.Println("Fail: ", t1.FailureRate(600))

		//helper.PrintMemUsage()
		t1.recordAccessSecondCount()
		t1.recordFailureSecondCount()
	}

	helper.PrintMemUsage()
	//time.Sleep(2e9)
	//time.Sleep(2e9)
	s := time.Now()
	for i := 0; i < 10000; i++ {
		t1.LoadRate(5)
		t1.LoadRate(60)
		t1.LoadRate(900)
		t1.LoadRate(1800)
		t1.LoadRate(200000)
	}
	fmt.Println(time.Since(s))
	fmt.Println("Load: ", t1.LoadRate(30*60))
	fmt.Println("Load: ", t1.LoadRate(60*10))
	fmt.Println("Load: ", t1.LoadRate(60*10))
	fmt.Println("Load: ", t1.LoadRate(100000))

	log.Println(len(t1.accessCountSecondSlice))
}
