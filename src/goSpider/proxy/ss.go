package proxy

import (
	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
	"goSpider/helper"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type SsAddr struct {
	Enable     bool
	Interval   float64
	Type       string
	Name       string
	ServerAddr string
	ClientAddr string
	Weight     int
}

func SSLocalHandler() (ssAddr []*SsAddr) {
	for _, server := range helper.Env().SsServers {
		if !server.Enable {
			continue
		}

		if server.Obfs != "" || server.ObfsParam != "" || server.ProtocolParam != "" || server.Protocol != "" {

			ss := &SsAddr{
				Enable:     server.Enable,
				Interval:   server.Interval,
				Name:       server.Name,
				Type:       "ssr",
				ServerAddr: server.Server + ":" + server.ServerPort,
			}
			ssAddr = append(ssAddr, ss)

			go func(server helper.SsServer, ss *SsAddr) {
				bi := &BackendInfo{
					Address: server.Server + ":" + server.ServerPort,
					Type:    "ssr",
					SSInfo: SSInfo{
						EncryptMethod:   server.Method,
						EncryptPassword: server.Password,
						SSRInfo: SSRInfo{
							Protocol:      server.Protocol,
							ProtocolParam: server.ProtocolParam,
							Obfs:          server.Obfs,
							ObfsParam:     server.ObfsParam,
						},
					},
				}
				bi.Listen(ss)
			}(server, ss)
		} else {
			cipher, err := core.PickCipher(server.Method, []byte{}, server.Password)
			if err != nil {
				log.Fatal(err)
			}

			ss := &SsAddr{
				Enable:     server.Enable,
				Interval:   server.Interval,
				Name:       server.Name,
				Type:       "ss",
				ServerAddr: server.Server + ":" + server.ServerPort,
			}
			ssAddr = append(ssAddr, ss)

			go socksLocal(ss, cipher.StreamConn)
		}
	}

	return
}

// Create a SOCKS server listening on addr and proxy to server.
func socksLocal(ssAddr *SsAddr, shadow func(net.Conn) net.Conn) {
	tcpLocal(ssAddr, shadow, func(c net.Conn) (socks.Addr, error) { return socks.Handshake(c) })
}

// relay copies between left and right bidirectionally. Returns number of
// bytes copied from right to left, from left to right, and any error occurred.
func relay(left, right net.Conn) (int64, int64, error) {
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

func tcpLocal(ssAddr *SsAddr, shadow func(net.Conn) net.Conn, getAddr func(net.Conn) (socks.Addr, error)) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Println("SS failed to listen: ", err)
		os.Exit(200)
	}
	ssAddr.ClientAddr = "127.0.0.1:" + strconv.Itoa(l.Addr().(*net.TCPAddr).Port)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("failed to accept: %s", err)
			continue
		}
		go func() {
			defer c.Close()
			c.(*net.TCPConn).SetKeepAlive(true)
			tgt, err := getAddr(c)
			if err != nil {

				// UDP: keep the connection until disconnect then free the UDP socket
				if err == socks.InfoUDPAssociate {
					buf := []byte{}
					// block here
					for {
						_, err := c.Read(buf)
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

			rc, err := net.Dial("tcp", ssAddr.ServerAddr)
			if err != nil {
				//2019/01/13 22:56:01 ss.go:158: failed to connect to server %v: %v hk.......domain......72 dial tcp: lookup hk05.bilibilivpn.com: no such host
				log.Println("failed to connect to server %v: %v", ssAddr.ServerAddr, err)
				return
			}
			defer rc.Close()
			rc.(*net.TCPConn).SetKeepAlive(true)
			rc = shadow(rc)

			if _, err = rc.Write(tgt); err != nil {
				log.Println("failed to send target address: %v", err)
				return
			}

			//log.Println("proxy %s <-> %s <-> %s", c.RemoteAddr(), server, tgt)
			relay(rc, c)
			//_, _, err = relay(rc, c)
			//if err != nil {
			//	if err, ok := err.(net.Error); ok && err.Timeout() {
			//		return // ignore i/o timeout
			//	}
			//	log.Println(ssAddr.ServerAddr+"relay error: %v", err)
			//}
		}()
	}
}
