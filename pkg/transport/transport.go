package transport

import (
	"fmt"
	"github.com/tobyxdd/hysteria/pkg/conns/faketcp"
	"github.com/tobyxdd/hysteria/pkg/conns/udp"
	"github.com/tobyxdd/hysteria/pkg/conns/wechat"
	"github.com/tobyxdd/hysteria/pkg/obfs"
	"net"
	"time"
)

type Transport interface {
	QUICResolveUDPAddr(address string) (*net.UDPAddr, error)
	QUICPacketConn(proto string, server bool, laddr, raddr string, obfs obfs.Obfuscator) (net.PacketConn, error)

	LocalResolveIPAddr(address string) (*net.IPAddr, error)
	LocalResolveTCPAddr(address string) (*net.TCPAddr, error)
	LocalResolveUDPAddr(address string) (*net.UDPAddr, error)
	LocalDial(network, address string) (net.Conn, error)
	LocalDialTCP(laddr, raddr *net.TCPAddr) (*net.TCPConn, error)
	LocalListenTCP(laddr *net.TCPAddr) (*net.TCPListener, error)
	LocalListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error)
}

var DefaultTransport Transport = &defaultTransport{
	Timeout: 8 * time.Second,
}

var IPv6OnlyTransport Transport = &ipv6OnlyTransport{
	defaultTransport{
		Timeout: 8 * time.Second,
	},
}

type defaultTransport struct {
	Timeout time.Duration
}

func (t *defaultTransport) QUICResolveUDPAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}

func (t *defaultTransport) QUICPacketConn(proto string, server bool, laddr, raddr string, obfs obfs.Obfuscator) (net.PacketConn, error) {
	if len(proto) == 0 || proto == "udp" {
		var laddrU *net.UDPAddr
		if len(laddr) > 0 {
			var err error
			laddrU, err = t.QUICResolveUDPAddr(laddr)
			if err != nil {
				return nil, err
			}
		}
		conn, err := net.ListenUDP("udp", laddrU)
		if err != nil {
			return nil, err
		}
		if obfs != nil {
			oc := udp.NewObfsUDPConn(conn, obfs)
			return oc, nil
		} else {
			return conn, nil
		}
	} else if proto == "wechat-video" {
		var laddrU *net.UDPAddr
		if len(laddr) > 0 {
			var err error
			laddrU, err = t.QUICResolveUDPAddr(laddr)
			if err != nil {
				return nil, err
			}
		}
		conn, err := net.ListenUDP("udp", laddrU)
		if err != nil {
			return nil, err
		}
		if obfs != nil {
			oc := wechat.NewObfsWeChatUDPConn(conn, obfs)
			return oc, nil
		} else {
			return conn, nil
		}
	} else if proto == "faketcp" {
		var conn *faketcp.TCPConn
		var err error
		if server {
			conn, err = faketcp.Listen("tcp", laddr)
			if err != nil {
				return nil, err
			}
		} else {
			conn, err = faketcp.Dial("tcp", raddr)
			if err != nil {
				return nil, err
			}
		}
		if obfs != nil {
			oc := faketcp.NewObfsFakeTCPConn(conn, obfs)
			return oc, nil
		} else {
			return conn, nil
		}
	} else {
		return nil, fmt.Errorf("unsupported protocol: %s", proto)
	}
}

func (t *defaultTransport) LocalResolveIPAddr(address string) (*net.IPAddr, error) {
	return net.ResolveIPAddr("ip", address)
}

func (t *defaultTransport) LocalResolveTCPAddr(address string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", address)
}

func (t *defaultTransport) LocalResolveUDPAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}

func (t *defaultTransport) LocalDial(network, address string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: t.Timeout}
	return dialer.Dial(network, address)
}

func (t *defaultTransport) LocalDialTCP(laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	dialer := &net.Dialer{Timeout: t.Timeout, LocalAddr: laddr}
	conn, err := dialer.Dial("tcp", raddr.String())
	if err != nil {
		return nil, err
	}
	return conn.(*net.TCPConn), nil
}

func (t *defaultTransport) LocalListenTCP(laddr *net.TCPAddr) (*net.TCPListener, error) {
	return net.ListenTCP("tcp", laddr)
}

func (t *defaultTransport) LocalListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error) {
	return net.ListenUDP("udp", laddr)
}

type ipv6OnlyTransport struct {
	defaultTransport
}

func (t *ipv6OnlyTransport) LocalResolveIPAddr(address string) (*net.IPAddr, error) {
	return net.ResolveIPAddr("ip6", address)
}
