package controller

import (
	"github.com/masami10/rush/socket_writer"
	"net"
	"fmt"
	//"time"
)



type CVI3Client struct {
	client	*socket_writer.SocketWriter
}

func (cvi3_client *CVI3Client) Start(addr string) {
	cvi3_client.client = socket_writer.NewSocketWriter(addr, cvi3_client)
}

func (cvi3_client *CVI3Client) keepAlive() {

	//for {
	//	if c.Status == STATUS_OFFLINE {
	//		break
	//	}
	//
	//	seq := c.get_sequence()
	//	keepAlivePacket, seq := GeneratePacket(seq, Header_type_keep_alive, Xml_heart_beat)
	//	c.Write([]byte(keepAlivePacket), seq)
	//
	//	<- time.After(time.Duration(c.keep_period)) // 周期性发送一次信号
	//}
}

// 客户端读取
func (cvi3_client *CVI3Client) Read(conn net.Conn){
	defer conn.Close()

	buffer := make([]byte, 65535)

	for {
		_, err := conn.Read(buffer)
		if err != nil {
			break
		}

		//msg := string(buffer[0:n])

		fmt.Printf("%s\n", string(buffer))

	}
}