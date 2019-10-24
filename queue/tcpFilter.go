package queue

import (
	"encoding/binary"
	"encoding/json"
	"github.com/chenset/asuka/helper"
	"github.com/willf/bloom"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

//format
//cmd 10=bl filter, n=other, db 0~255, fun 10=TestString 20=TestAndAddString
//request: length[4],cmd[1],db[1],fun[1],[string]array[*]
//request: length[4],cmd[1],json[*]
//response:length[4],json[*]
const (
	lenOfDataLen = 2
	lenOfCmd     = 1
	lenOfDb      = 1
	lenOfFun     = 1
)

type TcpFilter struct {
	bls              map[string]*bloom.BloomFilter
	bloomFilterMutex sync.Mutex

	connPool chan net.Conn

	ServerHandleCount int
}

var tcpFilterInstanceOnce sync.Once
var tcpFilterInstance *TcpFilter

func GetTcpFilterInstance() *TcpFilter {
	tcpFilterInstanceOnce.Do(func() {
		tcpFilterInstance = &TcpFilter{connPool: make(chan net.Conn, 100)}

		// kill signal handing
		helper.ExitHandleFuncSlice = append(helper.ExitHandleFuncSlice, func() {

			tcpFilterInstance.bloomFilterMutex.Lock()
			defer tcpFilterInstance.bloomFilterMutex.Unlock()

			//save to file
			for k, v := range tcpFilterInstance.bls {
				f, _ := os.Create(tcpFilterInstance.getBlFileName(k))
				v.WriteTo(f)
				f.Close()
			}

			log.Println("save")
		})
	})
	return tcpFilterInstance
}
func (my *TcpFilter) getBlFileName(name string) string {
	return helper.Env().BloomFilterPath + "tcp_" + name + ".db"
}

//ClientBl db 0~255, fun 10=TestString 20=TestAndAddString
func (my *TcpFilter) ClientBl(db byte, fun byte, rawUrls []string) (result []byte, err error) {
	buf := helper.LeakyBuf().Get()
	defer helper.LeakyBuf().Put(buf)

	jsonBytes, err := json.Marshal(rawUrls)
	if err != nil {
		return result, err
	}

	dataLen := uint16(len(jsonBytes) + lenOfCmd + lenOfDb + lenOfFun)

	newBuf := buf
	if int(dataLen+lenOfDataLen) > len(buf) {
		newBuf = make([]byte, dataLen+lenOfDataLen)
	}

	newBuf[lenOfDataLen] = 10                   //cmd 10=bl filter
	newBuf[lenOfDataLen+lenOfCmd] = db          //Db
	newBuf[lenOfDataLen+lenOfCmd+lenOfDb] = fun //fun
	copy(newBuf[lenOfDataLen+lenOfCmd+lenOfDb+lenOfFun:], jsonBytes[:])

	binary.BigEndian.PutUint16(newBuf[:lenOfDataLen], dataLen)
	n, err := my.client(newBuf, dataLen+lenOfDataLen)
	if err != nil {
		return result, err
	}

	//log.Println(string(buf[:n]))
	err = json.Unmarshal(newBuf[lenOfDataLen:n], &result)
	return
}

//ClientOtherCmd db 0~255, fun 10=TestString 20=TestAndAddString
func (my *TcpFilter) ClientOtherCmd(cmd byte, cmdData map[string]interface{}) (bool, error) {
	buf := helper.LeakyBuf().Get()
	defer helper.LeakyBuf().Put(buf)

	jsonBytes, err := json.Marshal(cmdData)
	if err != nil {
		return false, err
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
		return false, err
	}

	//todo implement
	log.Println(newBuf[:n])
	//binary.BigEndian.Uint16(buf[:lenOfData])

	return true, nil
}

func (my *TcpFilter) GetConn() (conn net.Conn, err error) {
	// Grab a buffer if available; allocate if not.
	select {
	case conn = <-my.connPool:
		// Got one; nothing more to do.
	default:
		// None free, so allocate a new one.
		conn, err = (&net.Dialer{Timeout: time.Second * 5}).Dial("tcp", "127.0.0.1:7654")
		if err != nil {
			break
		}
		//defer conn.Close()
		conn.(*net.TCPConn).SetKeepAlive(true)
		conn.(*net.TCPConn).SetKeepAlivePeriod(time.Second * 58)
	}
	return
}

func (my *TcpFilter) PutConn(conn net.Conn) {
	// Reuse buffer if there's room.
	select {
	case my.connPool <- conn:
		// Buffer on free list; nothing more to do.
	default:
		// Free list full, just carry on.
	}
}

func (my *TcpFilter) client(buf []byte, writeLen uint16) (n int, err error) {
	conn, err := my.GetConn()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err == nil {
			my.PutConn(conn)
		}
	}()

	//write
	_, err = conn.Write(buf[:writeLen])
	if err != nil {
		return
	}

	//read
	n, err = io.ReadAtLeast(conn, buf, lenOfDataLen)
	if err != nil {
		log.Println(err)
		return
	}

	dataLen := binary.BigEndian.Uint16(buf[:lenOfDataLen])

	// read continue
	if uint16(n) < lenOfDataLen+dataLen {
		nn, err := io.ReadAtLeast(conn, buf[n:], int(lenOfDataLen+dataLen)-n)
		n += nn
		if err != nil {
			log.Println(err)
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

		//decode
		//cmd 10=bl filter, 20=other
		if newBuf[lenOfDataLen] == 10 {
			//length[4],cmd[1],db[1],fun[1],json[*]
			list, err := my.ServerBl(newBuf[lenOfDataLen+lenOfCmd], newBuf[lenOfDataLen+lenOfCmd+lenOfDb], newBuf[lenOfDataLen+lenOfCmd+lenOfDb+lenOfFun:lenOfDataLen+dataLen])
			if err != nil {
				log.Println(err)
				return
			}
			encStr, err := json.Marshal(list)
			if err != nil {
				log.Println(err)
				return
			}

			copy(newBuf[lenOfDataLen:], encStr[:])

			dataLen := len(encStr)
			binary.BigEndian.PutUint16(newBuf[:lenOfDataLen], uint16(dataLen))
			_, err = conn.Write(newBuf[:lenOfDataLen+dataLen])
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			//length[4],cmd[1],json[*]
			my.OtherCmd(newBuf[lenOfDataLen], newBuf[lenOfDataLen+lenOfCmd:lenOfDataLen+dataLen])
		}
	}
}

func (my *TcpFilter) ServerBl(db, fun byte, data []byte) (result []byte, err error) {
	my.ServerHandleCount++
	var rawUrls []string
	err = json.Unmarshal(data, &rawUrls)
	if err != nil {
		log.Println(err)
		return
	}

	my.bloomFilterMutex.Lock()
	defer my.bloomFilterMutex.Unlock()

	//fun ==
	//fun 10=TestString 20=TestAndAddString
	bl := my.getBl(strconv.Itoa(int(db)))
	for _, u := range rawUrls {
		if fun == 10 {
			var b byte
			if bl.TestString(u) {
				b = 1
			}
			result = append(result, b)

		} else if fun == 20 {
			var b byte
			if bl.TestAndAddString(u) {
				b = 1
			}
			result = append(result, b)
		}
	}

	return
}

//getBl using with lock
func (my *TcpFilter) getBl(name string) *bloom.BloomFilter {
	bl, ok := my.bls[name]
	if !ok {
		bl = bloom.NewWithEstimates(10000000, 0.003)
		f, _ := os.Open(my.getBlFileName(name))
		bl.ReadFrom(f)
		f.Close()

		if len(my.bls) == 0 {
			my.bls = make(map[string]*bloom.BloomFilter)
		}
		my.bls[name] = bl
	}
	return bl
}

func (my *TcpFilter) OtherCmd(cmd byte, data []byte) (err error) {
	var cmdMap map[string]interface{}
	err = json.Unmarshal(data, &cmdMap)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(cmd)
	//todo implement
	res, err := json.Marshal(cmdMap)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(res))
	return
}
