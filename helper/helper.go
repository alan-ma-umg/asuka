package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"golang.org/x/net/publicsuffix"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var envParseOnce sync.Once
var envConfig *EnvConfig

func Env() *EnvConfig {
	envParseOnce.Do(func() {
		redis := flag.String("redis", "tcp://10.0.0.2:6379/9", "Redis connection url")
		mysql := flag.String("mysql", "root:11111111@(127.0.0.1:3306)/asuka?charset=utf8mb4", "Mysql DSN")
		webPassword := flag.String("webPassword", "", "WEB login password")
		bloomFilterPath := flag.String("bloomFilterPath", ".", "BloomFilter save path")
		localTransport := flag.Bool("localTransport", true, "Enable http.DefaultTransport")
		listen := flag.String("listen", "0.0.0.0:666", "WEB monitor listen address")
		wechatSendMessagePassword := flag.String("wechatSendMessagePassword", "", "Ignore it")
		singleProject := flag.String("singleProject", "", "Run only specific project")
		flag.Parse()

		u, err := url.Parse(*redis)
		if err != nil {
			log.Fatal(err)
		}
		redisDB, _ := strconv.Atoi(strings.TrimLeft(u.Path, "/"))
		redisPassword, _ := u.User.Password()
		envConfig = &EnvConfig{
			BloomFilterPath:           strings.TrimRight(*bloomFilterPath, "/") + "/",
			WEBPassword:               *webPassword,
			WEBListen:                 *listen,
			LocalTransport:            *localTransport,
			MysqlDSN:                  *mysql,
			WechatSendMessagePassword: *wechatSendMessagePassword,
			SingleProject:             *singleProject,
			Redis: Redis{
				Network:     u.Scheme,
				Addr:        u.Host,
				Password:    redisPassword,
				DB:          redisDB,
				URLQueueKey: "url_queue_key",
			},
		}
	})
	return envConfig
}

var ExitHandleFuncSlice []func()

// kill signal handing
func ExitHandle() {
	c := make(chan os.Signal, 3)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		for range c {
			ExitHandleFunc()
			os.Exit(0)
		}
	}()
}

func ExitHandleFunc() {
	for _, f := range ExitHandleFuncSlice {
		f()
	}
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ByteCountBinary(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f%c", float64(b)/float64(div), "KMGTPE"[exp])
}

// TldDomain return the Second-level domain and Top-level domain from url string
// https://www.domain.com => domain.com
// http://c.a.b.domain.com => domain.com
func TldDomain(u *url.URL) (tld string, err error) {
	hostname := u.Hostname()
	if IsIP(hostname) {
		return hostname, nil
	}

	tld, err = publicsuffix.EffectiveTLDPlusOne(hostname) // fixme,  failure: netlify.com | s3-ap-northeast-1.amazonaws.com
	if err != nil {
		return
	}

	s := strings.Split(tld, ".")

	if len(s) == 1 {
		return "", errors.New("tld长度不正确: " + tld)
	}

	last := strings.ToLower(s[len(s)-1])
	if !OnlyAlphabetCharacter(last) {
		return "", errors.New("顶级域名中包含非字母") //fixme 会导致不支持中文域名等
	}

	if last == "html" || last == "htm" || last == "php" || last == "jsp" || last == "json" || last == "xml" || last == "txt" || last == "shtml" || len(last) == 1 {
		return "", errors.New("无效tld: " + tld)
	}

	return
}

var OnlyDomainCharacter = regexp.MustCompile(`^[\-\.a-zA-Z0-9]+$`).MatchString

var OnlyAlphabetCharacter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

//IsIP  todo support IPv6, net.ResolveIPAddr("ip4:icmp", "")
func IsIP(host string) bool {
	parts := strings.Split(host, ".")

	if len(parts) < 4 {
		return false
	}

	for _, x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func TruncateStr(str []rune, length int, postfix string) string {
	cut := str
	if len(str) > length {
		cut = str[0:length]
	}
	return string(cut) + postfix
}

var GetSocketEstablishedCountLazyTicker = false
var GetSocketEstablishedCountLazyCacheCount = 0

func GetSocketEstablishedCountLazy() int {

	if GetSocketEstablishedCountLazyTicker {
		return GetSocketEstablishedCountLazyCacheCount
	}

	GetSocketEstablishedCountLazyTicker = true
	ticker := time.After(time.Second * 5)
	go func() {
		defer func() {
			<-ticker
			GetSocketEstablishedCountLazyTicker = false
		}()
		GetSocketEstablishedCountLazyCacheCount = 0

		if runtime.GOOS == "windows" {
			out, err := exec.Command("netstat", "-ano", "-p", "tcp").Output() //slower
			if err != nil {
				GetSocketEstablishedCountLazyCacheCount = 0
				return
			}
			pid := strconv.Itoa(os.Getpid())
			for _, s := range strings.Split(string(out), "\r\n") {
				if strings.Contains(s, "ESTABLISHED") && strings.Contains(s, pid) {
					GetSocketEstablishedCountLazyCacheCount++
				}
			}
		} else {
			pid := strconv.Itoa(os.Getpid())
			files, err := ioutil.ReadDir("/proc/" + pid + "/fd/") // faster than netstat
			if err != nil {
				GetSocketEstablishedCountLazyCacheCount = 0
				return
			}

			GetSocketEstablishedCountLazyCacheCount = len(files) - 5
			if GetSocketEstablishedCountLazyCacheCount < 0 {
				GetSocketEstablishedCountLazyCacheCount = 0
			}
		}
	}()

	return GetSocketEstablishedCountLazyCacheCount
}

func SpiderFailureRate(accessCount, failureCount int) float64 {
	if accessCount == 0 {
		if failureCount > 0 {
			return 100.0
		}
		return 0.0
	}
	return math.Min(float64(failureCount)/float64(accessCount)*100, 100.0)
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func MaxInt64(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func Ping(ip *net.IPAddr, times int) (avgRtt time.Duration, failureTimes int) {
	failureTimes = times
	p := fastping.NewPinger()
	p.AddIPAddr(ip)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		failureTimes--
		avgRtt += rtt
	}

	for i := 0; i < times; i++ {
		p.Run()
	}
	success := time.Duration(times - failureTimes)
	if success > 0 {
		avgRtt /= success
	}
	return
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", ByteCountBinary(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", ByteCountBinary(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", ByteCountBinary(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func KDF(password string, keyLen int) []byte {
	var b, prev []byte
	h := md5.New()
	for len(b) < keyLen {
		h.Write(prev)
		h.Write([]byte(password))
		b = h.Sum(b)
		prev = b[len(b)-h.Size():]
		h.Reset()
	}
	return b[:keyLen]
}

func Enc(plain []byte) (encData string, nonce []byte) {
	block, err := aes.NewCipher(KDF(Env().WEBPassword, 32))
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce = make([]byte, 12)
	io.ReadFull(rand.Reader, nonce)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	return hex.EncodeToString(gcm.Seal(nil, nonce, plain, nil)), nonce
}

func Dec(encData string, nonce []byte) (plain []byte, err error) {
	block, err := aes.NewCipher(KDF(Env().WEBPassword, 32))
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	encString, err := hex.DecodeString(encData)
	if err != nil {
		return
	}

	return gcm.Open(nil, nonce, encString, nil)
}

func TimeSince(t time.Duration) (str string) {
	num := int(t.Seconds())
	if num/3600/24 > 0 {
		str += strconv.Itoa(num/3600/24) + "d"
	}
	if num/3600%24 > 0 {
		str += strconv.Itoa(num/3600%24) + "h"
	}
	if num/60%60 > 0 {
		str += strconv.Itoa(num/60%60) + "m"
	}
	if num%60 > 0 {
		str += strconv.Itoa(num%60) + "s"
	}

	if num == 0 {
		str += t.String()
	}
	return
}

var DoOnceDurationHourInstance = &sync.Once{}

// DoOnceDurationHour global do once with reset in duration
func DoOnceDurationHour(fun func()) {
	DoOnceDurationHourInstance.Do(func() {
		fun()
		go func() {
			time.Sleep(time.Hour)
			DoOnceDurationHourInstance = &sync.Once{} //reset
		}()
	})
}

func SendTextToWXDoOnceDurationHour(content string) {
	if Env().WechatSendMessagePassword != "" {
		go func() {
			DoOnceDurationHour(func() {
				http.Get("https://wx.flysay.com/send?password=" + Env().WechatSendMessagePassword + "&touser=chen&content=" + url.QueryEscape(content))
			})
		}()
	}
}
