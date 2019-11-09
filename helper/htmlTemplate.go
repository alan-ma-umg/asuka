package helper

import (
	"encoding/hex"
	"hash/fnv"
	"html/template"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var templates *template.Template

var templatesOnce sync.Once

var fileVersionCtlMap = make(map[string]string)
var fileVersionCtlMapMutex sync.Mutex
var startTime = time.Now()

func GetTemplates() *template.Template {
	//time.Since(startTime).Hours() > 12. waiting for jsdelivr cache
	if runtime.GOOS == "linux" && time.Since(startTime).Hours() > 12. {
		templatesOnce.Do(func() {
			templates = ParseTemplates()
		})
		return templates
	} else {
		return ParseTemplates()
	}
}

func ParseTemplates() *template.Template {
	fileVersionCtlMapMutex.Lock()
	fileVersionCtlMap = make(map[string]string)
	fileVersionCtlMapMutex.Unlock()
	templates = template.Must(template.Must(template.New("").Funcs(template.FuncMap{
		"FilePathBase":   filepath.Base,
		"FileVersionCtl": fileVersionCtl,
		"FileCdnCtl":     fileCdnCtl,
		"Incr": func(i int) int {
			return i + 1
		},
	}).ParseGlob("web/templates/*/*.html")).ParseGlob("web/templates/*.html"))
	return templates
}

func fileCdnCtl(src string) template.URL {
	versionSrc := fileVersionCtl(src)
	if runtime.GOOS == "linux" {
		if src == "/static/asuka.css" || src == "/static/asuka.js" {
			return "https://cdn.jsdelivr.net/gh/chenset/asuka@latest/web/templates" + versionSrc
		}
	}

	return versionSrc
}

func fileVersionCtl(src string) template.URL {
	fileVersionCtlMapMutex.Lock()
	defer fileVersionCtlMapMutex.Unlock()

	if m, ok := fileVersionCtlMap[src]; ok {
		return template.URL(m)
	}

	u, err := url.Parse(src)
	if err != nil {
		log.Println(err)
		return template.URL(src)
	}
	f, err := os.Open("web/templates/" + u.Path)
	if err != nil {
		log.Println(err)
		return template.URL(src)
	}
	defer f.Close()
	h := fnv.New32a()
	if _, err := io.Copy(h, f); err != nil {
		log.Println(err)
		return template.URL(src)
	}

	query := u.Query()
	query.Add("v", hex.EncodeToString(h.Sum(nil)))
	u.RawQuery = query.Encode()

	fileVersionCtlMap[src] = u.String()

	return template.URL(fileVersionCtlMap[src])
}
