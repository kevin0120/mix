package cvi3

import (
	"net"
	"fmt"
	"time"
)

type CVI3Server struct {
	Port uint
	net.Listener
	Parent *CVI3
}

func (cvi3_server *CVI3Server) listen(port string) {
	for {
		l, err := net.Listen("tcp", port)
		if err == nil {
			cvi3_server.Listener = l
			break
		}

		time.Sleep(300 * time.Millisecond)
	}


	for {
		c, err := cvi3_server.Accept()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			break
		}

		sn := cvi3_server.Parent.setRemoteConn(c.RemoteAddr().String(), c)

		go cvi3_server.read(sn)
	}
}

func (cvi3_server *CVI3Server) read(sn string) {
	c := cvi3_server.Parent.clients[sn].RemoteConn
	buffer := make([]byte, 65535)
	for {

		n, err := c.Read(buffer)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			break
		}
		msg := string(buffer[0:n])
		if len(msg) < HEADER_LEN {
			continue
		}
		//fmt.Printf("%s\n", msg)

		header := CVI3Header{}
		header.Deserialize(msg[0: HEADER_LEN])
		var body string = msg[HEADER_LEN: n]
		var rest int = int(header.SIZ) - HEADER_LEN - n
		for {
			if rest <= 0 {
				break
			}
			n, err := c.Read(buffer)
			if err != nil {
				break
			}
			body += string(buffer[0:n])
			rest -= n
		}

		go cvi3_server.Parent.FUNCRecv(body)

		if header.TYP == Header_type_request_with_reply || header.TYP == Header_type_keep_alive {
			// 执行应答
			reply := CVI3Header{}
			reply.Init()
			reply.TYP = Header_type_reply
			reply.MID = header.MID
			reply_packet := reply.Serialize()

			_, err := c.Write([]byte(reply_packet))
			if err != nil {
				print("%s\n", err.Error())
				break
			}
		}
	}
}

func (cvi3_server *CVI3Server) Start(port string) error {

	// 开始tcp服务端侦听
	go cvi3_server.listen(port)

	return nil
}
