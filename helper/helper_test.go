package helper

import (
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestSSSubscriptionParse(t *testing.T) {
	//log.Println(net.LookupAddr("www.qq.com"))
	//ip, _ := net.ResolveIPAddr("ip4:icmp", "qq.com")
	//log.Println(Ping(ip, 2))
}

func TestHttpProxyParse(t *testing.T) {
	s := HttpProxyParse(`138.94.115.166:8080`)
	ioutil.WriteFile("C:/Users/41991/Desktop/http.json", []byte(s), 0644)
}
