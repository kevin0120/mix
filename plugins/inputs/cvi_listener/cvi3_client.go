package cvi_listener

import (
	"net"
	"sync"
	"fmt"

	"github.com/masami10/rush/utils/cvi"
	"time"
)

const (
	// 心跳间隔(ms)
	keep_alive_inteval = 7000

	// 下发超时(ms)
	send_timeout = 3000
)

type CVIConfig struct {
	SN 		string		//控制器序列号，须和后端配置一致
	IP 		string		//控制器ip
	Port 	uint		//控制器端口
	HMI 	string		//对应HMI的url，用于数据推送
}

type CVI3Client struct {
	Config 				*CVIConfig
	Conn				net.Conn
	serial_no			uint				// 1 ~ 9999
	msg_queue			map[uint]string
	mtx_serial			sync.Mutex
	mtx_queue			sync.Mutex
}

// 启动客户端
func (client *CVI3Client) Start() {
	client.msg_queue = map[uint]string{}

	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.Config.IP, client.Config.Port))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	client.Conn = c

	// 读取
	go client.Read()

	// 启动心跳
	go client.keep_alive()

	// 订阅数据
	client.subscribe()

	//go client.PSet(1, 1)
}

// PSet程序设定
func (client *CVI3Client) PSet(pset int, workorder_id int) {
	//time.Sleep(3000 * time.Millisecond)

	sdate, stime := cvi.GetDateTime()
	xml_pset := fmt.Sprintf(cvi.Xml_pset, sdate, stime, client.Config.SN, workorder_id, pset)

	serial := client.get_serial()
	pset_packet := cvi.GeneratePacket(serial, cvi.Header_type_request_with_reply, xml_pset)
	fmt.Printf("%s\n", pset_packet)

	client.add_to_queue(serial, "")

	_, err := client.Conn.Write([]byte(pset_packet))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

// 读取
func (client *CVI3Client) Read(){
	defer client.Conn.Close()

	buffer := make([]byte, 65535)

	for {
		//msg, err := reader.ReadString('\n')
		n, err := client.Conn.Read(buffer)
		if err != nil {
			break
		}
		msg := string(buffer[0:n])

		//fmt.Printf("%s\n", msg)

		// 处理应答
		header := cvi.CVI3Header{}
		header.Deserialize(msg[0: cvi.HEADER_LEN])

		client.remove_from_queue(header.MID)

	}
}

// 订阅数据
func (client *CVI3Client) subscribe() {
	sdate, stime := cvi.GetDateTime()
	xml_subscribe := fmt.Sprintf(cvi.Xml_subscribe, sdate, stime)

	serial := client.get_serial()
	subscribe_packet := cvi.GeneratePacket(serial, cvi.Header_type_request_with_reply, xml_subscribe)

	client.add_to_queue(serial, "")

	fmt.Printf("%s\n", subscribe_packet)
	_, err := client.Conn.Write([]byte(subscribe_packet))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

}

// 心跳
func (client *CVI3Client) keep_alive() {
	for {
		serial := client.get_serial()
		keep_alive_packet := cvi.GeneratePacket(serial, cvi.Header_type_request_with_reply, cvi.Xml_heart_beat)

		client.add_to_queue(serial, "")

		_, err := client.Conn.Write([]byte(keep_alive_packet))
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		//fmt.Printf("n=%d\n", n)

		time.Sleep(keep_alive_inteval * time.Millisecond)
	}
}

func (client *CVI3Client) get_serial() (uint) {
	defer client.mtx_serial.Unlock()

	client.mtx_serial.Lock()
	if client.serial_no == 9999 {
		client.serial_no = 1
	} else {
		client.serial_no++
	}

	return client.serial_no
}

func (client *CVI3Client) add_to_queue(serial uint, msg string) {
	defer client.mtx_queue.Unlock()

	client.mtx_queue.Lock()
	client.msg_queue[serial] = msg
}

func (client *CVI3Client) remove_from_queue(serial uint) {
	defer client.mtx_queue.Unlock()

	client.mtx_queue.Lock()
	delete(client.msg_queue, serial)
}

