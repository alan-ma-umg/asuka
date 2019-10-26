package helper

import (
	"os"
	"sync"
	"time"
)

type FileLog struct {
	filename  string
	checkTime time.Time
}

var fileLogInstanceDoOnce sync.Once
var fileLogInstance *FileLog

func GetFileLogInstance() *FileLog {
	fileLogInstanceDoOnce.Do(func() {
		fileLogInstance = &FileLog{
			filename: "./nohup.out",
		}
	})
	return fileLogInstance
}

func (my *FileLog) UpdateLogCheckTime() {
	my.checkTime = time.Now()
}

func (my *FileLog) GetLogCheckTime() time.Time {
	return my.checkTime
}

func (my *FileLog) GetLogModifyTime() (t time.Time) {
	if fi, err := os.Stat(my.filename); err == nil {
		return fi.ModTime()
	}
	return
}

func (my *FileLog) FileSize() uint64 {
	stat, err := os.Stat(my.filename)
	if err != nil {
		return 0
	}
	return uint64(stat.Size())
}

func (my *FileLog) TailFile(tailSize int64) (buf []byte) {
	stat, err := os.Stat(my.filename)
	if err != nil || stat.Size() == 0 {
		return
	}

	file, err := os.Open(my.filename)
	if err != nil {
		return
	}
	defer file.Close()

	buf = make([]byte, MinInt64(stat.Size(), tailSize))

	file.ReadAt(buf, stat.Size()-int64(len(buf)))
	return
}
