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
)

var templates *template.Template

var templatesOnce sync.Once

var fileCacheCtlMap = make(map[string]string)
var fileCacheCtlMapMutex sync.Mutex

func GetTemplates() *template.Template {
	if runtime.GOOS == "linux" {
		templatesOnce.Do(func() {
			templates = ParseTemplates()
		})
		return templates
	} else {
		return ParseTemplates()
	}
}

func ParseTemplates() *template.Template {
	//templatesOnce.Do(func() {
	templates = template.Must(template.Must(template.New("").Funcs(template.FuncMap{
		"FilePathBase": filepath.Base,
		"FileCacheCtl": fileCacheCtl,
		"Incr": func(i int) int {
			return i + 1
		},
	}).ParseGlob("web/templates/*/*.html")).ParseGlob("web/templates/*.html"))
	//})
	return templates
}

func fileCacheCtl(src string) template.URL {
	fileCacheCtlMapMutex.Lock()
	defer fileCacheCtlMapMutex.Unlock()

	if m, ok := fileCacheCtlMap[src]; ok {
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

	fileCacheCtlMap[src] = u.String()

	return template.URL(fileCacheCtlMap[src])
}
