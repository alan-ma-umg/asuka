package helper

import (
	"encoding/json"
	"github.com/jpillora/go-tld"
	"log"
	"net/url"
	"os"
	"os/user"
	"runtime"
	"sync"
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

// b2Mb byte to mb
func B2Mb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

//workspace path
var downloadPathOnce sync.Once
var downloadSavePath string

func WorkspacePath() string {
	downloadPathOnce.Do(func() {
		if Env().WorkspacePath != "" {
			downloadSavePath = Env().WorkspacePath
			return
		}

		osUser, _ := user.Current()
		downloadSavePath = osUser.HomeDir + "/Downloads/"
		if runtime.GOOS == "windows" {
			downloadSavePath = osUser.HomeDir + "\\Downloads\\"
		}
	})

	return downloadSavePath
}

// TldDomain return the Second-level domain and Top-level domain from url string
// https://www.domain.com => domain.com
// http://c.a.b.domain.com => domain.com
func TldDomain(rawUrl string) (string, error) {
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
