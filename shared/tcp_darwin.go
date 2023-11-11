package shared

import (
	"context"
	"net"

	"github.com/7c/go-ssclient/pfutil"
	"github.com/7c/go-ssclient/socks"
)

func natLookup(c net.Conn) (socks.Addr, error) {
	if tc, ok := c.(*net.TCPConn); ok {
		addr, err := pfutil.NatLookup(tc)
		return socks.ParseAddr(addr.String()), err
	}
	panic("not TCP connection")
}

func redirLocal(addr, server string, shadow func(net.Conn) net.Conn) {
	//func tcpLocal(addr, server string, shadow func(net.Conn) net.Conn, getAddr func(net.Conn) (socks.Addr, error)) {
	TcpLocal(context.Background(), addr, server, shadow, natLookup, nil)
}

func redir6Local(addr, server string, shadow func(net.Conn) net.Conn) {
	panic("TCP6 redirect not supported")
}
