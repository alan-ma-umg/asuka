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
	"sync"
	"time"
)

//format
//request: length[4],cmd[1],json[*]
//response:length[4],json[*]
const (
	lenOfDataLen = 2
	lenOfCmd     = 1
)

//Cmd10 bloomFilter
type Cmd10 struct {
	Urls []string
	Fun  byte
	Db   string
}

//Cmd10 server status report
type Cmd20 struct {
	Text string
}

type BlsItem struct {
	Bl       *bloom.BloomFilter
	LastUse  time.Time
	UseCount int
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

func GetTcpFilterInstance() *TcpFilter {
	tcpFilterInstanceOnce.Do(func() {
		//only client mode
		serverAddress := ""
		if helper.Env().BloomFilterClient != "" {
			u, err := url.Parse(helper.Env().BloomFilterClient)
			if err != nil {
				log.Println(err)
			} else {
				serverAddress = u.Host
			}
		}

		tcpFilterInstance = &TcpFilter{serverAddress: serverAddress, connPool: make(chan net.Conn, 100), startTime: time.Now()}
		//release idle bl
		go func() {
			for {
				time.Sleep(time.Minute * 33)
				tcpFilterInstance.bloomFilterMutex.Lock()

				for name, blItem := range tcpFilterInstance.blsItems {
					if time.Since(blItem.LastUse).Seconds() > 3600 {
						tcpFilterInstance.blSave(name, blItem)
						delete(tcpFilterInstance.blsItems, name)

						log.Println("del tcp bl: " + name)
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
			}

			if len(tcpFilterInstance.blsItems) > 0 {
				log.Println("save")
			}
		})
	})
	return tcpFilterInstance
}
func (my *TcpFilter) blSave(name string, blItem *BlsItem) {
	f, _ := os.Create(tcpFilterInstance.getBlFileName(name))
	blItem.Bl.WriteTo(f)
	f.Close()
}

func (my *TcpFilter) getBlFileName(name string) string {
	return helper.Env().BloomFilterPath + name + ".db"
}

//ClientOtherCmd db 0~255, fun 10=TestString 20=TestAndAddString
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
	n, err := my.client(newBuf, dataLen+lenOfDataLen)
	if err != nil {
		return res, err
	}
	return newBuf[lenOfDataLen:n], nil
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

func (my *TcpFilter) putConn(conn net.Conn) {
	// Reuse buffer if there's room.
	select {
	case my.connPool <- conn:
		// Buffer on free list; nothing more to do.
	default:
		// Free list full, just carry on.
	}
}

func (my *TcpFilter) client(buf []byte, writeLen uint16) (n int, err error) {
	conn, err := my.getConn()
	defer func() {
		if err != nil {
			name, _ := os.Hostname()
			helper.SendTextToWXDoOnceDurationHour(name + " TcpFilter connection failed: " + err.Error())
			log.Println(err)
			return
		}
		my.putConn(conn)
	}()

	if err != nil {
		return
	}

	//write
	_, err = conn.Write(buf[:writeLen])
	if err != nil {
		return
	}

	//read
	n, err = io.ReadAtLeast(conn, buf, lenOfDataLen)
	if err != nil {
		return
	}

	dataLen := binary.BigEndian.Uint16(buf[:lenOfDataLen])

	// read continue
	if uint16(n) < lenOfDataLen+dataLen {
		nn, err := io.ReadAtLeast(conn, buf[n:], int(lenOfDataLen+dataLen)-n)
		n += nn
		if err != nil {
			return n, err
		}
	}

	return
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

func (my *TcpFilter) handleServerConnection(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	buf := helper.LeakyBuf().Get()
	defer helper.LeakyBuf().Put(buf)

	for {
		n, err := io.ReadAtLeast(conn, buf, lenOfDataLen+lenOfCmd)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
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
		case 10:
			replyData, err = my.serverBl(newBuf[lenOfDataLen+lenOfCmd : lenOfDataLen+dataLen])
		case 20:
			replyData, err = my.serverReport()
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

func (my *TcpFilter) serverReply(conn net.Conn, buf, data []byte) (err error) {

	//todo check if reply data lager than buf

	copy(buf[lenOfDataLen:], data[:])
	dataLen := len(data)
	binary.BigEndian.PutUint16(buf[:lenOfDataLen], uint16(dataLen))
	_, err = conn.Write(buf[:lenOfDataLen+dataLen])
	if err != nil {
		log.Println(err)
	}

	return
}

func (my *TcpFilter) serverReport() (result []byte, err error) {

	runtime.ReadMemStats(&my.mem)

	memAvailable, total := helper.GetMemInfoFromProc()

	//len(map) is not thread safe
	my.bloomFilterMutex.Lock()
	blSize := len(my.blsItems)
	my.bloomFilterMutex.Unlock()

	return json.Marshal(map[string]interface{}{
		"pool_size":     len(my.connPool),
		"bl_alive_size": blSize,
		"bl_test_count": my.blTestCount,
		"goroutine":     runtime.NumGoroutine(),
		"sockets":       helper.GetSocketEstablishedCountLazy(),
		"load":          helper.GetSystemLoadFromProc(),
		"mem_rss":       helper.GetProgramRss(),
		"mem_available": memAvailable,
		"mem_total":     total,
		"mem_sys":       my.mem.Sys,
		"mem_alloc":     my.mem.Alloc,
		"start_time":    my.startTime.Unix(),
	})
}

func (my *TcpFilter) serverBl(buf []byte) (result []byte, err error) {
	//db, fun byte, data []byte
	var cmd10 Cmd10
	err = json.Unmarshal(buf, &cmd10)
	if err != nil {
		log.Println(err)
		return
	}

	my.blTestCount++

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	//fun 10=TestString 20=TestAndAddString
	bl := my.getBl(cmd10.Db)
	var list []byte
	for _, u := range cmd10.Urls {
		var b byte
		if cmd10.Fun == 10 {
			if bl.TestString(u) {
				b = 1
			}
			list = append(result, b)
		} else {
			if bl.TestAndAddString(u) {
				b = 1
			}
			list = append(result, b)
		}
	}

	return json.Marshal(list)
}

//getBl using with lock
func (my *TcpFilter) getBl(name string) *bloom.BloomFilter {
	blItem, ok := my.blsItems[name]
	if !ok {
		blItem = &BlsItem{}
		blItem.Bl = bloom.NewWithEstimates(10000000, 0.003) //todo 容量太小, 要视类型与情况加大. 可以考虑通过客户端传参数过来控制
		f, _ := os.Open(my.getBlFileName(name))
		blItem.Bl.ReadFrom(f)
		f.Close()

		if len(my.blsItems) == 0 {
			my.blsItems = make(map[string]*BlsItem)
		}
		my.blsItems[name] = blItem

		log.Println("new tcp bl: " + name)
	}

	blItem.LastUse = time.Now()
	blItem.UseCount++
	return blItem.Bl
}

//func (my *TcpFilter) otherCmd(data []byte) (result []byte, err error) {
//	var cmdMap map[string]interface{}
//	err = json.Unmarshal(data, &cmdMap)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	//todo implement
//	result, err = json.Marshal(cmdMap)
//	if err != nil {
//		log.Println(err)
//		//return
//	}
//	//log.Println(string(res))
//	//return
//}
