package queue

import (
	"encoding/binary"
	"encoding/json"
	"github.com/chenset/asuka/helper"
	"github.com/willf/bloom"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

//format
//request: length[4],cmd[1],json[*]
//response:length[4],json[*]
const (
	lenOfDataLen = 2
	lenOfCmd     = 1
)

//Cmd10 bloomFilter TestString & TestAndAddString
type Cmd10 struct {
	Urls []string
	Fun  byte //fun 10=TestString 20=TestAndAddString
	Db   string
	Size uint
}

//Cmd11 bloomFilter clear & remove file
type Cmd11 struct {
	Db string
}

//Cmd12 BlsItem status
type Cmd12 struct {
	Db string
}

//Cmd21 fileLog.TailFile
type Cmd21 struct {
	TailSize int64
}

//Cmd21 fileLog.TailFile
type Cmd21Response struct {
	TailContent []byte
	LogMod      int64
	LogSize     uint64
}

//Cmd20Response use struct instead of map. map may cause "fatal error: concurrent map iteration and map write" error when using json.Marshal with another Goroutine in some case
type Cmd20Response struct {
	BlAliveSize  int
	BlTestCount  int
	Goroutine    int
	Sockets      int
	Load         string
	MemRss       uint64
	MemAvailable uint64
	MemTotal     uint64
	MemSys       uint64
	MemAlloc     uint64
	StartTime    int64
	LogMod       int64
	LogCheck     int64
	LogSize      uint64
}

type BlsItem struct {
	bl        *bloom.BloomFilter //must be private to Cmd12.json.Marshal
	LastUse   time.Time
	TestCount int // TestString & TestAddString
	AddCount  int // TestAddString
}

type TcpFilter struct {
	blsItems           map[string]*BlsItem
	bloomFilterMutex   sync.Mutex
	serverAddress      string
	connPool           chan net.Conn
	mem                runtime.MemStats
	startTime          time.Time
	NewConnectionCount int //for client
	blTestCount        int //for server
}

var tcpFilterInstanceOnce sync.Once
var tcpFilterInstance *TcpFilter
var TcpErrorPrintDoOnce = helper.NewDoOnceInDuration(time.Minute)

func GetTcpFilterInstance() *TcpFilter {
	tcpFilterInstanceOnce.Do(func() {

		tcpFilterInstance = &TcpFilter{blsItems: make(map[string]*BlsItem), connPool: make(chan net.Conn, 1024), startTime: time.Now()}

		//for client mode
		if helper.Env().BloomFilterClient != "" {
			u, err := url.Parse(helper.Env().BloomFilterClient)
			if err != nil {
				log.Println(err)
			} else {
				tcpFilterInstance.serverAddress = u.Host
			}
		}

		//for server mode
		if helper.Env().BloomFilterServer != "" {
			//release idle bl
			go func() {
				for {
					time.Sleep(time.Minute * 33)
					tcpFilterInstance.bloomFilterMutex.Lock()

					for name, blItem := range tcpFilterInstance.blsItems {
						s := time.Now()
						tcpFilterInstance.blSave(name, blItem)
						if time.Since(blItem.LastUse).Seconds() > 3600 {
							delete(tcpFilterInstance.blsItems, name)
							log.Println("Release: " + name + " lastUse:" + blItem.LastUse.Format(time.Stamp) + " useCount:" + strconv.Itoa(blItem.TestCount) + " addCount:" + strconv.Itoa(blItem.AddCount) + " time:" + time.Since(s).String())
						}
					}

					tcpFilterInstance.bloomFilterMutex.Unlock()
				}
			}()

			// kill signal handing
			helper.ExitHandleFuncSlice = append(helper.ExitHandleFuncSlice, func() {

				tcpFilterInstance.bloomFilterMutex.Lock()
				defer tcpFilterInstance.bloomFilterMutex.Unlock()

				//save to file
				for name, blItem := range tcpFilterInstance.blsItems {
					tcpFilterInstance.blSave(name, blItem)
					log.Println("SAVE: " + name)
				}
			})
		}
	})
	return tcpFilterInstance
}
func (my *TcpFilter) blSave(name string, blItem *BlsItem) {
	f, _ := os.Create(tcpFilterInstance.getBlFileName(name))
	blItem.bl.WriteTo(f)
	f.Close()
}

func (my *TcpFilter) getBlFileName(name string) string {
	return helper.Env().BloomFilterPath + name + ".db"
}

//ClientOtherCmd db 0~255, fun 10=TestString 20=TestAndAddString
//请求不超过len(buf)的长度~4000
//响应不超过uint16的长度~65535
func (my *TcpFilter) Cmd(cmd byte, cmdData interface{}) (res []byte, err error) {
	buf := helper.LeakyBuf().Get()
	defer helper.LeakyBuf().Put(buf)

	jsonBytes, err := json.Marshal(cmdData)
	if err != nil {
		return res, err
	}
	dataLen := uint16(len(jsonBytes) + lenOfCmd)

	newBuf := buf
	if int(dataLen+lenOfDataLen) > len(buf) {
		newBuf = make([]byte, dataLen+lenOfDataLen)
	}

	newBuf[lenOfDataLen] = cmd

	copy(newBuf[lenOfDataLen+lenOfCmd:], jsonBytes[:])

	binary.BigEndian.PutUint16(newBuf[:lenOfDataLen], dataLen)
	result, err := my.client(newBuf, dataLen+lenOfDataLen)
	if err != nil {
		return res, err
	}

	res = make([]byte, len(result))
	copy(res, result) //must make a copy of buf or "panic: JSON decoder out of sync - data changing underfoot?"
	return
}

func (my *TcpFilter) getConn() (conn net.Conn, err error) {
	// Grab a buffer if available; allocate if not.
	select {
	case conn = <-my.connPool:
		// Got one; nothing more to do.
	default:
		// None free, so allocate a new one.
		conn, err = (&net.Dialer{Timeout: time.Second * 5}).Dial("tcp", my.serverAddress)
		my.NewConnectionCount++
		if err != nil {
			break
		}
		//defer conn.Close()
		conn.(*net.TCPConn).SetKeepAlive(true)
		conn.(*net.TCPConn).SetKeepAlivePeriod(time.Second * 58)
	}
	return
}

//ConnPoolSize for client
func (my *TcpFilter) ConnPoolSize() int {
	return len(my.connPool)
}

func (my *TcpFilter) putConn(conn net.Conn) {
	// Reuse buffer if there's room.
	select {
	case my.connPool <- conn:
		// Buffer on free list; nothing more to do.
	default:
		// Free list full, just carry on.
	}
}

func (my *TcpFilter) client(buf []byte, writeLen uint16) (response []byte, err error) {
	conn, err := my.getConn()
	defer func() {
		if err != nil {
			name, _ := os.Hostname()
			helper.SendTextToWXDoOnceDurationHour(name + " TcpFilter connection failed: " + err.Error())
			TcpErrorPrintDoOnce.Do(func() {
				log.Println(err)
			})
			return
		}
		my.putConn(conn)
	}()

	if err != nil {
		return nil, err
	}

	//write
	_, err = conn.Write(buf[:writeLen])
	if err != nil {
		return nil, err
	}

	//read
	n, err := io.ReadAtLeast(conn, buf, lenOfDataLen)
	if err != nil {
		return nil, err
	}

	dataLen := binary.BigEndian.Uint16(buf[:lenOfDataLen])

	newBuf := buf
	if int(dataLen+lenOfDataLen) > len(buf) {
		newBuf = make([]byte, dataLen+lenOfDataLen)
		copy(newBuf, buf)
	}

	// read continue
	if uint16(n) < lenOfDataLen+dataLen {
		nn, err := io.ReadAtLeast(conn, newBuf[n:], int(lenOfDataLen+dataLen)-n)
		n += nn
		if err != nil {
			return nil, err
		}
	}

	return newBuf[lenOfDataLen:n], err
}

func (my *TcpFilter) ServerListen(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		//go handleServerConnection(socket.NewCatConn(conn, 0))
		go my.handleServerConnection(conn)
	}
}

func (my *TcpFilter) serverReply(conn net.Conn, buf, data []byte) (err error) {
	dataLen := len(data)
	newBuf := buf
	if dataLen+lenOfDataLen > len(buf) {
		newBuf = make([]byte, dataLen+lenOfDataLen)
	}

	copy(newBuf[lenOfDataLen:], data[:])

	binary.BigEndian.PutUint16(newBuf[:lenOfDataLen], uint16(dataLen))
	_, err = conn.Write(newBuf[:lenOfDataLen+dataLen])
	if err != nil {
		log.Println(err)
	}

	return
}

func (my *TcpFilter) handleServerConnection(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	buf := helper.LeakyBuf().Get()
	defer helper.LeakyBuf().Put(buf)

	for {
		n, err := io.ReadAtLeast(conn, buf, lenOfDataLen+lenOfCmd)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), syscall.ECONNRESET.Error()) {
				return
			}

			log.Println(err)
			return
		}

		dataLen := binary.BigEndian.Uint16(buf[:lenOfDataLen])

		newBuf := buf
		if int(dataLen+lenOfDataLen) > len(buf) {
			newBuf = make([]byte, dataLen+lenOfDataLen)
			copy(newBuf, buf)
		}

		// read continue
		if uint16(n) < lenOfDataLen+dataLen {
			_, err := io.ReadAtLeast(conn, newBuf[n:], int(lenOfDataLen+dataLen)-n)
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				return
			}
		}

		//cmd
		var replyData []byte
		switch newBuf[lenOfDataLen] { //cmd
		case 10: //bloomFilter test
			replyData, err = my.serverBl(newBuf[lenOfDataLen+lenOfCmd : lenOfDataLen+dataLen])
		case 11: //bloomFilter clear & remove file.db
			replyData, err = my.serverBlClear(newBuf[lenOfDataLen+lenOfCmd : lenOfDataLen+dataLen])
		case 12: //blItem status
			replyData, err = my.serverBlStatus(newBuf[lenOfDataLen+lenOfCmd : lenOfDataLen+dataLen])
		case 20: //system report
			replyData, err = my.serverReport()
		case 21: //fileLog.TailFile
			replyData, err = my.serverTailFile(newBuf[lenOfDataLen+lenOfCmd : lenOfDataLen+dataLen])
		case 22: //fileLog.UpdateLogCheckTime
			helper.GetFileLogInstance().UpdateLogCheckTime()
		default:
			replyData = newBuf[lenOfDataLen+lenOfCmd : lenOfDataLen+dataLen]
		}
		if err != nil {
			log.Println(err)
			return
		}
		//reply
		my.serverReply(conn, newBuf, replyData)
	}
}

func (my *TcpFilter) serverTailFile(buf []byte) (result []byte, err error) {
	//db, fun byte, data []byte
	var cmd *Cmd21
	err = json.Unmarshal(buf, &cmd)
	if err != nil {
		log.Println(err)
		return
	}

	return json.Marshal(&Cmd21Response{
		TailContent: helper.GetFileLogInstance().TailFile(cmd.TailSize),
		LogMod:      helper.GetFileLogInstance().GetLogModifyTime().Unix(),
		LogSize:     helper.GetFileLogInstance().FileSize(),
	})
}

func (my *TcpFilter) serverReport() (result []byte, err error) {
	runtime.ReadMemStats(&my.mem)
	memAvailable, total := helper.GetMemInfoFromProc()

	//len(map) is not thread safe
	my.bloomFilterMutex.Lock()
	blSize := len(my.blsItems)
	my.bloomFilterMutex.Unlock()
	return json.Marshal(&Cmd20Response{
		BlAliveSize:  blSize,
		BlTestCount:  my.blTestCount,
		Goroutine:    runtime.NumGoroutine(),
		Sockets:      helper.GetSocketEstablishedCountLazy(),
		Load:         helper.GetSystemLoadFromProc(),
		MemRss:       helper.GetProgramRss(),
		MemAvailable: memAvailable,
		MemTotal:     total,
		MemSys:       my.mem.Sys,
		MemAlloc:     my.mem.Alloc,
		StartTime:    my.startTime.Unix(),
		LogMod:       helper.GetFileLogInstance().GetLogModifyTime().Unix(),
		LogCheck:     helper.GetFileLogInstance().GetLogCheckTime().Unix(),
		LogSize:      helper.GetFileLogInstance().FileSize(),
	})
}

func (my *TcpFilter) serverBl(buf []byte) (result []byte, err error) {
	//db, fun byte, data []byte
	var cmd10 *Cmd10
	err = json.Unmarshal(buf, &cmd10)
	if err != nil {
		log.Println(err)
		return
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	//fun 10=TestString 20=TestAndAddString
	blItem := my.getBlItem(cmd10)
	var list []byte
	for _, u := range cmd10.Urls {
		my.blTestCount++
		blItem.TestCount++
		var b byte
		if cmd10.Fun == 10 {
			if blItem.bl.TestString(u) {
				b = 1
			}
			list = append(result, b)
		} else {
			blItem.AddCount++
			if blItem.bl.TestAndAddString(u) {
				b = 1
			}
			list = append(result, b)
		}
	}

	return json.Marshal(list)
}

func (my *TcpFilter) serverBlStatus(buf []byte) (result []byte, err error) {
	//db, fun byte, data []byte
	var cmd *Cmd12
	err = json.Unmarshal(buf, &cmd)
	if err != nil {
		log.Println(err)
		return
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	return json.Marshal(my.blsItems[cmd.Db])
}

func (my *TcpFilter) serverBlClear(buf []byte) (result []byte, err error) {
	//db, fun byte, data []byte
	var cmd *Cmd11
	err = json.Unmarshal(buf, &cmd)
	if err != nil {
		log.Println(err)
		return
	}

	s := time.Now()
	if os.Remove(my.getBlFileName(cmd.Db)) == nil {
		log.Println("Remove File: " + my.getBlFileName(cmd.Db))
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	if blItem, ok := my.blsItems[cmd.Db]; ok {
		blItem.bl.ClearAll()

		log.Println("DEL: " + cmd.Db + " lastUse:" + blItem.LastUse.Format(time.Stamp) + " useCount:" + strconv.Itoa(blItem.TestCount) + " addCount:" + strconv.Itoa(blItem.AddCount) + " time:" + time.Since(s).String())

		blItem.bl = nil
		delete(my.blsItems, cmd.Db)
	}

	return
}

//getBl using with lock
func (my *TcpFilter) getBlItem(cmd10 *Cmd10) *BlsItem {
	blItem, ok := my.blsItems[cmd10.Db]
	if !ok {
		blItem = &BlsItem{}
		blItem.bl = bloom.NewWithEstimates(cmd10.Size, 0.003)
		f, _ := os.Open(my.getBlFileName(cmd10.Db))
		blItem.bl.ReadFrom(f)
		f.Close()

		my.blsItems[cmd10.Db] = blItem

		//log.Println("NEW: " + cmd10.Db + " size :" + strconv.Itoa(int(cmd10.Size)))
	}

	blItem.LastUse = time.Now()
	return blItem
}
