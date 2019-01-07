package helper

import (
	"encoding/json"
	"fmt"
	"github.com/jpillora/go-tld"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
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
