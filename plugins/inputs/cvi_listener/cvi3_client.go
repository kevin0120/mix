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
	keep_alive_inteval = 5000

	// 下发超时(ms)
	REQUEST_TIMEOUT = 3000

	STATUS_ONLINE = "online"
	STATUS_OFFLINE = "offline"
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
	Results				ResultQueue
	mtx_serial			sync.Mutex
	Status				string
	mtx_status			sync.Mutex
	recv_flag			bool
	mtx_write			sync.Mutex
	RemoteConn			net.Conn
}

// 启动客户端
func (client *CVI3Client) Start() {

	client.Connect()

	// 订阅数据
	client.subscribe()

	// 启动心跳检测
	go client.keep_alive_check()

}

func (client *CVI3Client) Connect() {
	client.Status = STATUS_OFFLINE
	client.serial_no = 0

	fmt.Printf("CVI3:%s connecting ...\n", client.Config.SN)
	for {
		c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", client.Config.IP, client.Config.Port), 3 * time.Second)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		} else {
			client.Conn = c
			break
		}
	}

	client.Results = ResultQueue{}
	client.Results.Results = map[uint]string{}

	// 读取
	go client.Read()

	client.update_status(STATUS_ONLINE)

	// 启动心跳
	go client.keep_alive()

}

// PSet程序设定
func (client *CVI3Client) PSet(pset int, workorder_id int, screw_id string) (uint, error) {

	//time.Sleep(3 * time.Second)

	sdate, stime := cvi.GetDateTime()
	xml_pset := fmt.Sprintf(cvi.Xml_pset, sdate, stime, client.Config.SN, workorder_id, screw_id, pset)

	serial := client.get_serial()
	pset_packet := cvi.GeneratePacket(serial, cvi.Header_type_request_with_reply, xml_pset)
	fmt.Printf("%s\n", pset_packet)

	client.Results.update(serial, "")

	_, err := client.SafeWrite([]byte(pset_packet))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	return serial, err
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

		client.recv_flag = true

		msg := string(buffer[0:n])

		fmt.Printf("%s\n", msg)

		// 处理应答
		header_str := msg[0: cvi.HEADER_LEN]
		header := cvi.CVI3Header{}
		header.Deserialize(header_str)

		client.Results.update(header.MID, header_str)

	}
}

// 订阅数据
func (client *CVI3Client) subscribe() {

	sdate, stime := cvi.GetDateTime()
	xml_subscribe := fmt.Sprintf(cvi.Xml_subscribe, sdate, stime)

	serial := client.get_serial()
	subscribe_packet := cvi.GeneratePacket(serial, cvi.Header_type_request_with_reply, xml_subscribe)

	client.Results.update(serial, "")

	fmt.Printf("%s\n", subscribe_packet)
	_, err := client.SafeWrite([]byte(subscribe_packet))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

}

// 心跳
func (client *CVI3Client) keep_alive() {
	for {
		if client.Status == STATUS_OFFLINE {
			break
		}

		serial := client.get_serial()
		keep_alive_packet := cvi.GeneratePacket(serial, cvi.Header_type_request_with_reply, cvi.Xml_heart_beat)
		client.Results.update(serial, "")
		n, err := client.SafeWrite([]byte(keep_alive_packet))
		if err != nil {
			fmt.Printf("%d %s\n", n, err.Error())
			break
		}

		//fmt.Printf("n=%d\n", n)

		time.Sleep(keep_alive_inteval * time.Millisecond)
	}
}

// 心跳检测
func (client *CVI3Client) keep_alive_check() {

	for {

		for i:=0; i < 3; i++ {
			if client.recv_flag == true {
				client.update_status(STATUS_ONLINE)
				client.recv_flag = false
				time.Sleep(keep_alive_inteval * time.Millisecond)

				break
			} else {
				if i == 2 {
					client.update_status(STATUS_OFFLINE)

				}
			}

			time.Sleep(keep_alive_inteval * time.Millisecond)
		}

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

func (client *CVI3Client) update_status(status string) {
	defer client.mtx_status.Unlock()

	client.mtx_status.Lock()

	if status != client.Status {
		client.Status = status

		fmt.Printf("civ3:%s %s\n", client.Config.SN, client.Status)

		if client.Status == STATUS_OFFLINE {
			client.Conn.Close()
			client.RemoteConn.Close()

			// 断线重连
			go client.Connect()
		}
	}

}

func (client *CVI3Client) SafeWrite(buf []byte) (int, error) {
	defer client.mtx_write.Unlock()

	client.mtx_write.Lock()
	n, err := client.Conn.Write(buf)
	return n, err
}

type ResultQueue struct {
	Results		map[uint]string
	mtx			sync.Mutex
}

func (q *ResultQueue) update(serial uint, msg string) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	q.Results[serial] = msg
}

func (q *ResultQueue) remove(serial uint) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	delete(q.Results, serial)
}

func (q *ResultQueue) get(serial uint) (string) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	return q.Results[serial]
}