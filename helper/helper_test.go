package helper

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestSSSubscriptionParse(t *testing.T) {
	base := 3000000
	log.Println(base * (10. - 0) / 10.)
	log.Println(base * (10. - 3) / 10.)
	log.Println(base * (10. - 6) / 10.)
	log.Println(base * (10. - 9) / 10.)
	log.Println(base * (10. - 9) / 10.)
}
