package proxy

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"errors"
	ssSocks "github.com/shadowsocks/go-shadowsocks2/socks"
	"github.com/sun8911879/shadowsocksR"
	"github.com/sun8911879/shadowsocksR/obfs"
	"github.com/sun8911879/shadowsocksR/protocol"
	"github.com/sun8911879/shadowsocksR/ssr"
	"github.com/sun8911879/shadowsocksR/tools/leakybuf"
	"github.com/sun8911879/shadowsocksR/tools/socks"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var (
	readTimeout = 600 * time.Second
)

// SSInfo fields that shadowsocks/shadowsocksr used only
type SSInfo struct {
	SSRInfo
	EncryptMethod   string
	EncryptPassword string
}

// SSRInfo fields that shadowsocksr used only
type SSRInfo struct {
	Obfs          string
	ObfsParam     string
	ObfsData      interface{}
	Protocol      string
	ProtocolParam string
	ProtocolData  interface{}
}

// BackendInfo all fields that a backend used
type BackendInfo struct {
	SSInfo
	Address string
	Type    string
}

func (bi *BackendInfo) Listen(clientRawAddr string) {
	listener, err := net.Listen("tcp", clientRawAddr)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		localConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go bi.Handle(localConn)
	}
}

func (bi *BackendInfo) Handle(src net.Conn) {
	defer src.Close()
	//src.SetKeepAlive(true)
	src.(*net.TCPConn).SetKeepAlive(true)

	socks.ReadAddr(src)
	rawaddr, err := ssSocks.Handshake(src)
	if err != nil {
		// UDP: keep the connection until disconnect then free the UDP socket
		if err == socks.Error(9) {
			buf := []byte{}
			// block here
			for {
				_, err := src.Read(buf)
				if err, ok := err.(net.Error); ok && err.Timeout() {
					continue
				}
				log.Println("UDP Associate End.")
				return
			}
		}
		log.Println("failed to get target address: %v", err)
		return
	}

	dst, err := bi.DialSSRConn(socks.Addr(rawaddr))
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return
		}
		//2019/01/08 17:22:59 ssr.go:97: ru-3.mitsuha-node.com:443 *errors.errorString connecting to SSR server failed :dial tcp 103.102.4.19:443: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
		//2019/01/08 17:23:01 ssr.go:97: hk-6.mitsuha-node.com:443 *errors.errorString connecting to SSR server failed :dial tcp 203.218.247.64:443: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
		//2019/01/08 17:23:02 ssr.go:97: hk-2.mitsuha-node.com:443 *errors.errorString connecting to SSR server failed :dial tcp 218.102.182.182:443: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
		//2019/01/08 17:23:03 ssr.go:97: jp-9.mitsuha-node.com:443 *errors.errorString connecting to SSR server failed :dial tcp 203.137.122.111:443: connectex: No connection could be made because the target machine actively refused it.
		//2019/01/08 17:23:03 ssr.go:97: hk-24.mitsuha-node.com:443 *errors.errorString connecting to SSR server failed :dial tcp 47.52.214.162:443: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
		//2019/01/08 17:23:03 ssr.go:97: ru-3.mitsuha-node.com:443 *errors.errorString connecting to SSR server failed :dial tcp 103.102.4.19:443: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
		//2019/01/08 17:23:03 ssr.go:97: hk-10.mitsuha-node.com:443 *errors.errorString connecting to SSR server failed :dial tcp 42.98.194.185:443: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
		//2019/01/08 17:23:04 ssr.go:97: tw-1.mitsuha-node.com:443 *errors.errorString connectin
		log.Println(bi.Address, reflect.TypeOf(err), err)
		return //ignore i/o timeout
	}
	defer dst.Close()
	//dst.(*net.TCPConn).SetKeepAlive(true)

	//n,n,err:=tcpRelay(src, dst)
	//if err != nil {
	//	if err, ok := err.(net.Error); ok && err.Timeout() {
	//		return // ignore i/o timeout
	//	}
	//	log.Println("relay error: %v", err)
	//}

	go bi.Pipe(src, dst)
	bi.Pipe(dst, src)
	//src.Close()
	//dst.Close()
}

func (bi *BackendInfo) DialSSRConn(rawaddr socks.Addr) (net.Conn, error) {
	u := &url.URL{
		Scheme: bi.Type,
		Host:   bi.Address,
	}
	v := u.Query()
	v.Set("encrypt-method", bi.EncryptMethod)
	v.Set("encrypt-key", bi.EncryptPassword)
	v.Set("obfs", bi.Obfs)
	v.Set("obfs-param", bi.ObfsParam)
	v.Set("protocol", bi.Protocol)
	v.Set("protocol-param", bi.ProtocolParam)
	u.RawQuery = v.Encode()

	ssrconn, err := NewSSRClient(u)
	if err != nil {
		return nil, fmt.Errorf("connecting to SSR server failed :%v", err)
	}

	if bi.ObfsData == nil {
		bi.ObfsData = ssrconn.IObfs.GetData()
	}
	ssrconn.IObfs.SetData(bi.ObfsData)

	if bi.ProtocolData == nil {
		bi.ProtocolData = ssrconn.IProtocol.GetData()
	}
	ssrconn.IProtocol.SetData(bi.ProtocolData)

	if _, err := ssrconn.Write(rawaddr); err != nil {
		ssrconn.Close()
		return nil, err
	}
	return ssrconn, nil
}

// relay copies between left and right bidirectionally. Returns number of
// bytes copied from right to left, from left to right, and any error occurred.
func tcpRelay(left, right net.Conn) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)

	go func() {
		n, err := io.Copy(right, left)
		right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
		left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
		ch <- res{n, err}
	}()

	n, err := io.Copy(left, right)
	right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
	left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
	rs := <-ch

	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}

// PipeThenClose copies data from src to dst, closes dst when done.
func (bi *BackendInfo) Pipe(src, dst net.Conn) error {
	buf := leakybuf.GlobalLeakyBuf.Get()
	for {
		src.SetReadDeadline(time.Now().Add(readTimeout))
		n, err := src.Read(buf)
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		if n > 0 {
			// Note: avoid overwrite err returned by Read.
			if _, err := dst.Write(buf[0:n]); err != nil {
				break
			}
		}
		if err != nil {
			// Always "use of closed network connection", but no easy way to
			// identify this specific error. So just leave the error along for now.
			// More info here: https://code.google.com/p/go/issues/detail?id=4373
			break
		}
	}
	leakybuf.GlobalLeakyBuf.Put(buf)
	dst.Close()
	return nil
}

func NewSSRClient(u *url.URL) (*shadowsocksr.SSTCPConn, error) {
	query := u.Query()
	encryptMethod := query.Get("encrypt-method")
	encryptKey := query.Get("encrypt-key")
	cipher, err := shadowsocksr.NewStreamCipher(encryptMethod, encryptKey)
	if err != nil {
		return nil, err
	}

	dialer := net.Dialer{
		//Timeout:   time.Millisecond * 700,
		DualStack: true,
	}
	conn, err := dialer.Dial("tcp", u.Host)
	if err != nil {
		return nil, err
	}

	conn.(*net.TCPConn).SetKeepAlive(true)

	ssconn := shadowsocksr.NewSSTCPConn(conn, cipher)
	if ssconn.Conn == nil || ssconn.RemoteAddr() == nil {
		return nil, errors.New("nil connection")
	}

	// should initialize obfs/protocol now
	rs := strings.Split(ssconn.RemoteAddr().String(), ":")
	port, _ := strconv.Atoi(rs[1])

	ssconn.IObfs = obfs.NewObfs(query.Get("obfs"))
	obfsServerInfo := &ssr.ServerInfoForObfs{
		Host:   rs[0],
		Port:   uint16(port),
		TcpMss: 1460,
		Param:  query.Get("obfs-param"),
	}
	ssconn.IObfs.SetServerInfo(obfsServerInfo)
	ssconn.IProtocol = protocol.NewProtocol(query.Get("protocol"))
	protocolServerInfo := &ssr.ServerInfoForObfs{
		Host:   rs[0],
		Port:   uint16(port),
		TcpMss: 1460,
		Param:  query.Get("protocol-param"),
	}
	ssconn.IProtocol.SetServerInfo(protocolServerInfo)

	return ssconn, nil
}