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
		bloomFilterPath := flag.String("bloomFilterPath", ".", "BloomFilter save path, don't using with bloomFilterClient at same time")
		bloomFilterClient := flag.String("bloomFilterClient", "tcp://127.0.0.1:7654", "BloomFilter tcp client, don't using with bloomFilterPath at same time")
		bloomFilterServer := flag.String("bloomFilterServer", "0.0.0.0:7654", "BloomFilter tcp server")
		localTransport := flag.Bool("localTransport", true, "Enable http.DefaultTransport")
		listen := flag.String("listen", "0.0.0.0:666", "WEB monitor listen address")
		wechatSendMessagePassword := flag.String("wechatSendMessagePassword", "", "Ignore it")
		flag.Parse()

		u, err := url.Parse(*redis)
		if err != nil {
			log.Fatal(err)
		}
		redisDB, _ := strconv.Atoi(strings.TrimLeft(u.Path, "/"))
		redisPassword, _ := u.User.Password()
		envConfig = &EnvConfig{
			BloomFilterPath:           strings.TrimRight(*bloomFilterPath, "/") + "/",
			BloomFilterClient:         *bloomFilterClient,
			BloomFilterServer:         *bloomFilterServer,
			WEBPassword:               *webPassword,
			WEBListen:                 *listen,
			LocalTransport:            *localTransport,
			MysqlDSN:                  *mysql,
			WechatSendMessagePassword: *wechatSendMessagePassword,
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

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
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

type DoOnceInDuration struct {
	duration time.Duration
	once     *sync.Once
}

func NewDoOnceInDuration(duration time.Duration) *DoOnceInDuration {
	return &DoOnceInDuration{duration: duration, once: new(sync.Once)}
}
func (my *DoOnceInDuration) Do(f func()) (isRun bool) {
	my.once.Do(func() {
		f()
		isRun = true
		go func() {
			time.Sleep(my.duration)
			my.once = new(sync.Once)
		}()
	})

	return
}

var DoOnceDurationHourInstance = NewDoOnceInDuration(time.Hour + (123 * time.Millisecond))

// DoOnceDurationHour global do once with reset in duration
func DoOnceDurationHour(fun func()) {
	DoOnceDurationHourInstance.Do(func() {
		fun()
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

var intRex = regexp.MustCompile("[0-9]+")

func GetMemInfoFromProc() (available, total uint64) {
	if runtime.GOOS == "windows" {
		return
	}

	// system mem info
	if dat, err := ioutil.ReadFile("/proc/meminfo"); err == nil {
		for index, line := range strings.SplitN(string(dat), "\n", 6) {

			if strings.Contains(strings.ToLower(line), "memtotal") {
				if res := intRex.FindAllString(line, 1); len(res) == 1 {
					if kb, err := strconv.ParseUint(res[0], 10, 64); err == nil {
						total = kb * 1024
					}
				}
			}
			if strings.Contains(strings.ToLower(line), "memavailable") {
				if res := intRex.FindAllString(line, 1); len(res) == 1 {
					if kb, err := strconv.ParseUint(res[0], 10, 64); err == nil {
						available = kb * 1024
					}
				}
			}
			if index > 4 {
				break
			}
		}
	}

	return
}

var taskListDoOnceInDuration = NewDoOnceInDuration(time.Second*9 + (323 * time.Millisecond))
var getProgramRssWindowsCache uint64
var windowsFindRssRex = regexp.MustCompile("(?i)([\\d,]+)\\s?K$")

func GetProgramRss() (rss uint64) {
	if runtime.GOOS == "windows" {
		taskListDoOnceInDuration.Do(func() {
			go func() {
				if out, err := exec.Command("tasklist", "/fi", "pid  eq "+strconv.Itoa(os.Getpid()), "/FO", "LIST").Output(); err == nil { //slow
					for _, line := range strings.Split(string(out), "\n") {
						if res := windowsFindRssRex.FindStringSubmatch(strings.TrimSpace(line)); len(res) == 2 {
							if kb, err := strconv.ParseFloat(strings.ReplaceAll(res[1], ",", ""), 64); err == nil {
								getProgramRssWindowsCache = uint64(kb) * 1024
								return
							}
						}
					}
				}
			}()
		})

		return getProgramRssWindowsCache
	}

	// program mem info
	if dat, err := ioutil.ReadFile("/proc/" + strconv.Itoa(os.Getpid()) + "/status"); err == nil {
		for _, line := range strings.Split(string(dat), "\n") {
			if strings.Contains(strings.ToLower(line), "vmrss") {
				if res := intRex.FindAllString(line, 1); len(res) == 1 {
					if kb, err := strconv.ParseUint(res[0], 10, 64); err == nil {
						rss = kb * 1024
					}
				}
				break
			}
		}
	}
	return
}

func GetSystemLoadFromProc() (loadStr string) {
	if runtime.GOOS == "windows" {
		return
	}

	if dat, err := ioutil.ReadFile("/proc/loadavg"); err == nil {
		for index, str := range strings.SplitN(string(dat), " ", 4) {
			loadStr += str + " "
			if index == 2 {
				break
			}
		}
	}
	return
}

var getSocketEstablishedCountLazyCacheCount = 0
var getSocketEstablishedCountLazyCacheCountDoOnceInDuration = NewDoOnceInDuration(time.Second*12 + time.Millisecond*200) //时间错开

func GetSocketEstablishedCountLazy() int {
	getSocketEstablishedCountLazyCacheCountDoOnceInDuration.Do(func() {
		go func() {
			if runtime.GOOS == "windows" {
				out, err := exec.Command("netstat", "-ano", "-p", "tcp").Output() //slower
				if err != nil {
					getSocketEstablishedCountLazyCacheCount = 0
					return
				}
				pid := strconv.Itoa(os.Getpid())
				for _, s := range strings.Split(string(out), "\r\n") {
					if strings.Contains(s, "ESTABLISHED") && strings.Contains(s, pid) {
						getSocketEstablishedCountLazyCacheCount++
					}
				}
			} else {
				pid := strconv.Itoa(os.Getpid())
				files, err := ioutil.ReadDir("/proc/" + pid + "/fd/") // faster than netstat
				if err != nil {
					getSocketEstablishedCountLazyCacheCount = 0
					return
				}

				getSocketEstablishedCountLazyCacheCount = len(files) - 5
				if getSocketEstablishedCountLazyCacheCount < 0 {
					getSocketEstablishedCountLazyCacheCount = 0
				}
			}
		}()
	})
	return getSocketEstablishedCountLazyCacheCount
}
