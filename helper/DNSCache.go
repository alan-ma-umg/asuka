package helper

import (
	"net"
	"sync"
	"time"
)

type dnsItem struct {
	ipAddr string
	hit    time.Time
	flush  time.Time
}

type DNSCache struct {
	length        int
	cacheMap      sync.Map
	cacheDuration time.Duration
	flushDuration time.Duration
}

var newDNSCacheOnce sync.Once
var newDNSCacheInstance *DNSCache

func GetDNSCache() *DNSCache {

	newDNSCacheOnce.Do(func() {
		newDNSCacheInstance = &DNSCache{cacheDuration: time.Hour * 5, flushDuration: time.Minute * 30}

		go func() {
			for {
				time.Sleep(time.Minute)
				newDNSCacheInstance.cacheMap.Range(func(key, value interface{}) bool {
					item := value.(*dnsItem)

					//delete
					if time.Since(item.hit) > newDNSCacheInstance.cacheDuration {
						newDNSCacheInstance.cacheMap.Delete(key)
						newDNSCacheInstance.length--
						return true
					}

					//flush dns
					if time.Since(item.flush) > newDNSCacheInstance.flushDuration {
						newDNSCacheInstance.flush(key.(string))
						return true
					}

					return true
				})
			}
		}()
	})

	return newDNSCacheInstance
}

func (my *DNSCache) Lookup(domain string) (ipAddr string) {
	if net.ParseIP(domain) != nil {
		return domain //It's IP address
	}

	//get from cache
	res, ok := my.cacheMap.Load(domain)
	if ok {
		res.(*dnsItem).hit = time.Now()
		return res.(*dnsItem).ipAddr
	}

	return my.flush(domain)
}

func (my *DNSCache) flush(domain string) (ipAddr string) {
	//lookup
	ip, err := net.LookupIP(domain)
	if err != nil || len(ip) == 0 {
		//err
		return domain
	}

	ipAddr = ip[0].String()

	//cache
	item, ok := my.cacheMap.LoadOrStore(domain, &dnsItem{
		hit: time.Now(),
	})
	if !ok {
		newDNSCacheInstance.length++
	}
	item.(*dnsItem).ipAddr = ipAddr
	item.(*dnsItem).flush = time.Now()
	return ipAddr
}

func (my *DNSCache) Len() int {
	return my.length
}
