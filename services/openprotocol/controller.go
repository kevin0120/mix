package openprotocol

import (
	"fmt"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/socket_writer"
	"net"
)

type ControllerStatusType string

const (
	STATUS_ONLINE  ControllerStatusType = "online"
	STATUS_OFFLINE ControllerStatusType = "offline"
)

type Controller struct {
	w   *socket_writer.SocketWriter
	cfg controller.Config
}

func (c *Controller) Start() {

	c.w = socket_writer.NewSocketWriter(fmt.Sprintf("tcp://%s:%d", c.cfg.RemoteIP, c.cfg.Port), c)

	// 启动心跳检测
	//go c.keep_alive_check()

	c.Connect()
}

func (c *Controller) Connect() {

}

func (c *Controller) Read(conn net.Conn) {

}
