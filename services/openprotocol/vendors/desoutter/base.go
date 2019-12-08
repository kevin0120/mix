package desoutter

import "github.com/masami10/rush/socket_writer"

type BaseDesoutterController struct {
	sockClients []*socket_writer.SocketWriter //可能会有多个连接，每个工具有一个连接
}
