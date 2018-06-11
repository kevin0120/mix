package udp_driver

import "net"

type Protocol interface {
	Parse([] byte) error
	Read(c net.Conn)
	NewConn(c net.Conn)
}
