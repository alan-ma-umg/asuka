package proxy

import (
	"fmt"
	"testing"
	"time"
)

func TestSsLocalHandler(t *testing.T) {
	fmt.Println(SsLocalHandler())
	time.Sleep(1e9)
}
