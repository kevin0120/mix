package controller

import (
	"github.com/masami10/rush/socket_listener"
	"github.com/masami10/rush/services/audi_vw"
	"net"
	"fmt"
	"strings"
	"encoding/xml"
)

const (
	MASTER_PC_PORT = 4710
	KEY_PSET = "PRG"
)

type CVI3Server struct {
	server	*socket_listener.SocketListener

	cvi3_client CVI3Client

	SN string
	Workorder_id int64
	Result_ids []int64
	Count int

}

func (cvi3_server *CVI3Server) Start(port uint) (error) {

	addr := fmt.Sprintf("tcp://:%d", port)
	cvi3_server.server = socket_listener.NewSocketListener(addr, cvi3_server)

	cvi3_server.server.Start()

	fmt.Printf("controller started on %d\n", port)

	return nil
}


func (cvi3_server *CVI3Server)Parse (buf []byte) ([]byte, error) {
	msg := string(buf)
	if strings.Contains(msg, KEY_PSET) {
		// 如果收到pset请求，调用客户端返回拧接结果
		pset := CVI3PSet{}
		err := xml.Unmarshal(buf, &pset)
		if err != nil {
			fmt.Printf("Parse pset err:%s\n", err.Error())
		}

		cvi3_server.cvi3_client.PushResult(pset)
	}

	return nil, nil
}

func (cvi3_server *CVI3Server) NewConn(c net.Conn) {
	// 接受新连接后启动客户端

	ip := strings.Split(c.RemoteAddr().String(), ":")[0]

	addr := fmt.Sprintf("tcp://%s:%d", ip, MASTER_PC_PORT)

	cvi3_server.cvi3_client = CVI3Client {
		SN: cvi3_server.SN,
		Workorder_id: cvi3_server.Workorder_id,
		Result_ids: cvi3_server.Result_ids,
		Count: cvi3_server.Count,
	}

	go cvi3_server.cvi3_client.Start(addr)
}

// 服务端读取
func (cvi3_server *CVI3Server) Read(c net.Conn) {
	defer cvi3_server.server.InterListener.RemoveConnection(c)
	defer c.Close()

	buffer := make([]byte, 65535)
	for {

		n, err := c.Read(buffer)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			break
		}
		msg := string(buffer[0:n])
		if len(msg) < audi_vw.HEADER_LEN {
			continue
		}
		fmt.Printf("%s\n", msg)

		header := audi_vw.CVI3Header{}
		header.Deserialize(msg[0: audi_vw.HEADER_LEN])

		go cvi3_server.Parse(buffer)

		if header.TYP == audi_vw.Header_type_request_with_reply || header.TYP == audi_vw.Header_type_keep_alive {
			// 执行应答
			reply_packet, _ := audi_vw.GeneratePacket(header.MID, audi_vw.Header_type_reply, audi_vw.Xml_heart_beat)

			_, err := c.Write([]byte(reply_packet))
			if err != nil {
				print("server reply err:%s\n", err.Error())
				break
			}
		}
	}
}