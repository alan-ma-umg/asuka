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
		redis := flag.String("redis", "tcp://127.0.0.1:6379/9", "Redis connection url")
		mysql := flag.String("mysql", "root:11111111@(127.0.0.1:3306)/asuka?charset=utf8mb4", "Mysql DSN")
		webPassword := flag.String("webPassword", "", "WEB login password")
		bloomFilterPath := flag.String("bloomFilterPath", ".", "BloomFilter save path, don't using with bloomFilterClient at same time")
		bloomFilterClient := flag.String("bloomFilterClient", "tcp://127.0.0.1:7654", "BloomFilter tcp client, don't using with bloomFilterPath at same time")
		bloomFilterServer := flag.String("bloomFilterServer", "0.0.0.0:7654", "BloomFilter tcp server")
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
	return FileSizeH(b, 2)
}

func FileSizeH(b uint64, precision int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f%c", float64(b)/float64(div), "KMGTPE"[exp])
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

func CurrentFileAndLine() string {
	_, file, line, _ := runtime.Caller(1)
	return file + ":" + strconv.Itoa(line)
}

func stringAlign(field string, value interface{}) string {
	return field + strings.Repeat(" ", 20-len(field)) + fmt.Sprintln(value)
}

func PrintMemUsage(m runtime.MemStats) (str string) {
	s := time.Now()
	runtime.ReadMemStats(&m)
	// General statistics.

	str += stringAlign("RSS", ByteCountBinary(GetProgramRss()))

	// Alloc is bytes of allocated heap objects.
	//
	// This is the same as HeapAlloc (see below).
	str += stringAlign("Alloc", ByteCountBinary(m.Alloc))

	// TotalAlloc is cumulative bytes allocated for heap objects.
	//
	// TotalAlloc increases as heap objects are allocated, but
	// unlike Alloc and HeapAlloc, it does not decrease when
	// objects are freed.
	str += stringAlign("TotalAlloc", ByteCountBinary(m.TotalAlloc))

	// Sys is the total bytes of memory obtained from the OS.
	//
	// Sys is the sum of the XSys fields below. Sys measures the
	// virtual address space reserved by the Go runtime for the
	// heap, stacks, and other internal data structures. It's
	// likely that not all of the virtual address space is backed
	// by physical memory at any given moment, though in general
	// it all was at some point.
	str += stringAlign("Sys", ByteCountBinary(m.Sys))

	// Lookups is the number of pointer lookups performed by the
	// runtime.
	//
	// This is primarily useful for debugging runtime internals.
	str += stringAlign("Lookups", m.Lookups)

	// Mallocs is the cumulative count of heap objects allocated.
	// The number of live objects is Mallocs - Frees.
	str += stringAlign("Mallocs", m.Mallocs)

	// Frees is the cumulative count of heap objects freed.
	str += stringAlign("Frees", m.Frees)

	// Heap memory statistics.
	//
	// Interpreting the heap statistics requires some knowledge of
	// how Go organizes memory. Go divides the virtual address
	// space of the heap into "spans", which are contiguous
	// regions of memory 8K or larger. A span may be in one of
	// three states:
	//
	// An "idle" span contains no objects or other data. The
	// physical memory backing an idle span can be released back
	// to the OS (but the virtual address space never is), or it
	// can be converted into an "in use" or "stack" span.
	//
	// An "in use" span contains at least one heap object and may
	// have free space available to allocate more heap objects.
	//
	// A "stack" span is used for goroutine stacks. Stack spans
	// are not considered part of the heap. A span can change
	// between heap and stack memory; it is never used for both
	// simultaneously.

	// HeapAlloc is bytes of allocated heap objects.
	//
	// "Allocated" heap objects include all reachable objects, as
	// well as unreachable objects that the garbage collector has
	// not yet freed. Specifically, HeapAlloc increases as heap
	// objects are allocated and decreases as the heap is swept
	// and unreachable objects are freed. Sweeping occurs
	// incrementally between GC cycles, so these two processes
	// occur simultaneously, and as a result HeapAlloc tends to
	// change smoothly (in contrast with the sawtooth that is
	// typical of stop-the-world garbage collectors).
	str += stringAlign("HeapAlloc", ByteCountBinary(m.HeapAlloc))

	// HeapSys is bytes of heap memory obtained from the OS.
	//
	// HeapSys measures the amount of virtual address space
	// reserved for the heap. This includes virtual address space
	// that has been reserved but not yet used, which consumes no
	// physical memory, but tends to be small, as well as virtual
	// address space for which the physical memory has been
	// returned to the OS after it became unused (see HeapReleased
	// for a measure of the latter).
	//
	// HeapSys estimates the largest size the heap has had.
	str += stringAlign("HeapSys", ByteCountBinary(m.HeapSys))

	// HeapIdle is bytes in idle (unused) spans.
	//
	// Idle spans have no objects in them. These spans could be
	// (and may already have been) returned to the OS, or they can
	// be reused for heap allocations, or they can be reused as
	// stack memory.
	//
	// HeapIdle minus HeapReleased estimates the amount of memory
	// that could be returned to the OS, but is being retained by
	// the runtime so it can grow the heap without requesting more
	// memory from the OS. If this difference is significantly
	// larger than the heap size, it indicates there was a recent
	// transient spike in live heap size.
	str += stringAlign("HeapIdle", ByteCountBinary(m.HeapIdle))

	// HeapInuse is bytes in in-use spans.
	//
	// In-use spans have at least one object in them. These spans
	// can only be used for other objects of roughly the same
	// size.
	//
	// HeapInuse minus HeapAlloc estimates the amount of memory
	// that has been dedicated to particular size classes, but is
	// not currently being used. This is an upper bound on
	// fragmentation, but in general this memory can be reused
	// efficiently.
	str += stringAlign("HeapInuse", ByteCountBinary(m.HeapInuse))

	// HeapReleased is bytes of physical memory returned to the OS.
	//
	// This counts heap memory from idle spans that was returned
	// to the OS and has not yet been reacquired for the heap.
	str += stringAlign("HeapReleased", ByteCountBinary(m.HeapReleased))

	// HeapObjects is the number of allocated heap objects.
	//
	// Like HeapAlloc, this increases as objects are allocated and
	// decreases as the heap is swept and unreachable objects are
	// freed.
	str += stringAlign("HeapObjects", m.HeapObjects)

	// Stack memory statistics.
	//
	// Stacks are not considered part of the heap, but the runtime
	// can reuse a span of heap memory for stack memory, and
	// vice-versa.

	// StackInuse is bytes in stack spans.
	//
	// In-use stack spans have at least one stack in them. These
	// spans can only be used for other stacks of the same size.
	//
	// There is no StackIdle because unused stack spans are
	// returned to the heap (and hence counted toward HeapIdle).
	str += stringAlign("StackInuse", ByteCountBinary(m.StackInuse))

	// StackSys is bytes of stack memory obtained from the OS.
	//
	// StackSys is StackInuse, plus any memory obtained directly
	// from the OS for OS thread stacks (which should be minimal).
	str += stringAlign("StackSys", ByteCountBinary(m.StackSys))

	// Off-heap memory statistics.
	//
	// The following statistics measure runtime-internal
	// structures that are not allocated from heap memory (usually
	// because they are part of implementing the heap). Unlike
	// heap or stack memory, any memory allocated to these
	// structures is dedicated to these structures.
	//
	// These are primarily useful for debugging runtime memory
	// overheads.

	// MSpanInuse is bytes of allocated mspan structures.
	str += stringAlign("MSpanInuse", ByteCountBinary(m.MSpanInuse))

	// MSpanSys is bytes of memory obtained from the OS for mspan
	// structures.
	str += stringAlign("MSpanSys", ByteCountBinary(m.MSpanSys))

	// MCacheInuse is bytes of allocated mcache structures.
	str += stringAlign("MCacheInuse", ByteCountBinary(m.MCacheInuse))

	// MCacheSys is bytes of memory obtained from the OS for
	// mcache structures.
	str += stringAlign("MCacheSys", ByteCountBinary(m.MCacheSys))

	// BuckHashSys is bytes of memory in profiling bucket hash tables.
	str += stringAlign("BuckHashSys", ByteCountBinary(m.BuckHashSys))

	// GCSys is bytes of memory in garbage collection metadata.
	str += stringAlign("GCSys", ByteCountBinary(m.GCSys))

	// OtherSys is bytes of memory in miscellaneous off-heap
	// runtime allocations.
	str += stringAlign("OtherSys", ByteCountBinary(m.OtherSys))

	// Garbage collector statistics.

	// NextGC is the target heap size of the next GC cycle.
	//
	// The garbage collector's goal is to keep HeapAlloc ≤ NextGC.
	// At the end of each GC cycle, the target for the next cycle
	// is computed based on the amount of reachable data and the
	// value of GOGC.
	str += stringAlign("NextGC", ByteCountBinary(m.NextGC))

	// LastGC is the time the last garbage collection finished, as
	// nanoseconds since 1970 (the UNIX epoch).
	str += stringAlign("LastGC", time.Unix(0, int64(m.LastGC)).String())

	// PauseTotalNs is the cumulative nanoseconds in GC
	// stop-the-world pauses since the program started.
	//
	// During a stop-the-world pause, all goroutines are paused
	// and only the garbage collector can run.
	str += stringAlign("PauseTotalNs", time.Duration(m.PauseTotalNs).String())

	// NumGC is the number of completed GC cycles.
	str += stringAlign("NumGC", m.NumGC)
	str += stringAlign("NumForcedGC", m.NumForcedGC)
	str += stringAlign("GCCPUFraction", strconv.FormatFloat(m.GCCPUFraction, 'f', 8, 64))

	str += stringAlign("Time", time.Since(s).String())
	return
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

func FloatRound2(f float64) float64 {
	return math.Floor(f*100) / 100
}

func FloatRound4(f float64) float64 {
	return math.Floor(f*10000) / 10000
}

var onlyWhiteSpaceRex = regexp.MustCompile(`[ ]+`)

//GetNetTraffic
func GetNetTraffic(pid int) (rx, tx, rp, tp uint64) {
	if runtime.GOOS != "linux" {
		return
	}

	file := "/proc/net/dev"
	if pid > 0 {
		file = "/proc/" + strconv.Itoa(pid) + "/net/dev"
	}
	if dat, err := ioutil.ReadFile(file); err == nil {
		out := strings.ToLower(onlyWhiteSpaceRex.ReplaceAllString(string(dat), " "))
		lines := strings.Split(out, "\n")
		if len(lines) < 3 {
			return
		}
		lines = lines[2:]
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "lo:") {
				continue
			}
			nodes := strings.SplitN(line, " ", 12)
			if len(nodes) < 12 {
				return
			}

			//bytes
			r, _ := strconv.ParseUint(nodes[1], 10, 64)
			rx += r
			t, _ := strconv.ParseUint(nodes[9], 10, 64)
			tx += t

			//packets
			p1, _ := strconv.ParseUint(nodes[2], 10, 64)
			rp += p1
			p2, _ := strconv.ParseUint(nodes[10], 10, 64)
			tp += p2
		}
	}
	return
}

var RxSlice = make([]uint64, 600)
var TxSlice = make([]uint64, 600)
var RpSlice = make([]uint64, 600)
var TpSlice = make([]uint64, 600)
var getNetTrafficSliceDoOnce sync.Once

func GetNetTrafficSlice() ([]uint64, []uint64, []uint64, []uint64) {
	//net traffic counter
	getNetTrafficSliceDoOnce.Do(func() {
		//if runtime.GOOS == "linux" {
		go func() {
			time.Sleep(time.Second * 2)
			prevRx, prevTx, prevRp, prevTp := GetNetTraffic(0)
			for {
				time.Sleep(time.Second)
				rx, tx, rp, tp := GetNetTraffic(0)
				RxSlice = append(RxSlice[MaxInt(len(RxSlice)-599, 0):], rx-prevRx)
				TxSlice = append(TxSlice[MaxInt(len(TxSlice)-599, 0):], tx-prevTx)
				RpSlice = append(RpSlice[MaxInt(len(RpSlice)-599, 0):], rp-prevRp)
				TpSlice = append(TpSlice[MaxInt(len(TpSlice)-599, 0):], tp-prevTp)
				prevRx, prevTx, prevRp, prevTp = rx, tx, rp, tp
			}
		}()
		//}
	})

	return RxSlice, TxSlice, RpSlice, TpSlice
}
