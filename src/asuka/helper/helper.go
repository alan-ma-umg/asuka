package helper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// env config

var envParseOnce sync.Once
var envConfig *EnvConfig
var PathToEnvFile string

func Env() *EnvConfig {
	envParseOnce.Do(func() {
		file, err := os.Open(PathToEnvFile)
		if err != nil {
			log.Fatal(err)
		}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&envConfig)
		if err != nil {
			log.Fatal(err)
		}
	})

	return envConfig
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

func SSSubscriptionParse(rawUrl string) {
	resp, err := http.Get(rawUrl)
	if err != nil {

		return
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	bStr, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(string(res))
	if err != nil {
		fmt.Println(err)
		return
	}

	var SSServerArr []SsServer
	for _, line := range strings.Split(string(bStr), "\n") {
		if strings.HasPrefix(line, "ssr://") {
			//line=strings.Replace(line,"_","+",0)
			line = strings.Replace(line, "_", "+", len(line))

			if strings.Contains(line[6:], "/") || strings.Contains(line, "_") || strings.Contains(line, "-") {
				log.Println(line)
			}
			dec, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(line[6:])
			if err != nil {
				log.Println(line[6:])
				log.Println(err)
				return
			}
			s := strings.Split(string(dec), ":")

			part5 := strings.Split(s[5], "/")

			password, _ := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(part5[0])
			//fmt.Println(password)

			u, err := url.Parse(strings.Replace(strings.TrimSpace(s[5]), "/>", "/?", len(s[5]))) //fix typo
			if err != nil {
				log.Println(string(dec))
				log.Println(err)
				return
			}
			queryMap := u.Query()

			ssServer := SsServer{
				Enable:     true,
				Name:       s[0],
				Server:     s[0],
				ServerPort: s[1],
				Password:   string(password),
				Method:     s[3],
				//	ssr only
				Obfs: s[4],
				//ObfsParam     string
				//ProtocolParam string
				Protocol: s[2],
			}

			if v, ok := queryMap["obfsparam"]; ok {
				obfsparam, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(string(v[0]))
				if err != nil {
					log.Println(err)
					return
				}
				ssServer.ObfsParam = string(obfsparam)
			}

			if v, ok := queryMap["protoparam"]; ok {
				protoparam, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(string(v[0]))
				if err != nil {
					log.Println(err)
					return
				}
				ssServer.ProtocolParam = string(protoparam)
			}
			if ssServer.Protocol != "" && ssServer.ProtocolParam == "" {
				log.Println(queryMap)
				log.Println(string(dec))
			}

			if strings.ToLower(ssServer.Method) == "chacha20" {
				log.Println("chacha20: Not support yet.")
				continue
			}
			SSServerArr = append(SSServerArr, ssServer)
		}

		if strings.HasPrefix(line, "ss://") {
			//todo
		}
	}

	b, err := json.Marshal(SSServerArr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
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
