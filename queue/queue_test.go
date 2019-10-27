package queue

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func TestNewQueue(t *testing.T) {
	i := 2

	i *= 2

	log.Println(i)
}
