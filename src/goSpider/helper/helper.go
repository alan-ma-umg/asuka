package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jpillora/go-tld"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"net/http"
)

// env config

var envParseOnce sync.Once
var envConfig *EnvConfig

var pwd, _ = os.Getwd()

func Env() *EnvConfig {
	envParseOnce.Do(func() {
		filename := pwd + "/env.json"

		file, err := os.Open(filename)
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
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %c", float64(b)/float64(div), "KMGTPE"[exp])
}

// TldDomain return the Second-level domain and Top-level domain from url string
// https://www.domain.com => domain.com
// http://c.a.b.domain.com => domain.com
func TldDomain(rawUrl string) (str string, err error) {
	defer func() {
		if r := recover(); r != nil {
			str = rawUrl
		}
	}()

	tldU, err := tld.Parse(rawUrl)
	if err != nil {
		u, err := url.Parse(rawUrl)
		if err != nil {
			return "", err
		}

		return u.Hostname(), nil
	}

	return tldU.Domain + "." + tldU.TLD, nil
}

func TruncateStr(str string, length int, postfix string) string {
	//todo support utf8
	cut := str
	if len(str) > length {
		cut = str[0:length] + postfix
	}
	return cut
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

			u, err := url.Parse(strings.Replace(string(dec), "/>", "/?", len(string(dec)))) //fix typo
			if err != nil {
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
