package controller

import (
	"github.com/masami10/rush/socket_writer"
	"net"
	"fmt"
	//"time"
	"sync"
	"time"
	"github.com/masami10/rush/services/audi_vw"
	"github.com/masami10/rush/utils"
	"math/rand"
)

type CVI3Client struct {
	client	*socket_writer.SocketWriter
	sequence		uint
	SN				string
	mtx_no			sync.Mutex
	mtx_serial		sync.Mutex
	buffer  		chan []byte
	Workorder_id	int64
	Result_ids		[]int64
	Count			int
}

func (cvi3_client *CVI3Client) Start(addr string) {
	cvi3_client.buffer = make(chan []byte, 1024)
	cvi3_client.mtx_serial = sync.Mutex{}
	cvi3_client.mtx_no = sync.Mutex{}

	if cvi3_client.client != nil {
		cvi3_client.client.Close()
	}

	cvi3_client.client = socket_writer.NewSocketWriter(addr, cvi3_client)

	cvi3_client.client.Connect(3 * time.Second)

	go cvi3_client.manage()

	go cvi3_client.AudoPushResult()

	//go cvi3_client.keepAlive()
}

func (cvi3_client *CVI3Client) get_sequence() uint {
	cvi3_client.mtx_serial.Lock()
	defer cvi3_client.mtx_serial.Unlock()

	if cvi3_client.sequence == 9999 {
		cvi3_client.sequence = 1
	} else {
		cvi3_client.sequence++
	}
	return cvi3_client.sequence
}

func (cvi3_client *CVI3Client) keepAlive() {

	for {

		seq := cvi3_client.get_sequence()
		keepAlivePacket, seq := audi_vw.GeneratePacket(seq, audi_vw.Header_type_keep_alive, audi_vw.Xml_heart_beat)
		cvi3_client.Write([]byte(keepAlivePacket), seq)

		<- time.After(5 * time.Second) // 周期性发送一次信号
	}
}

// 客户端读取
func (cvi3_client *CVI3Client) Read(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 65535)

	for {
		_, err := conn.Read(buffer)
		if err != nil {
			break
		}

		//msg := string(buffer[0:n])

		//fmt.Printf("%s\n", string(buffer))

	}
}

// 发送拧接结果
func (cvi3_client *CVI3Client) PushResult(pset CVI3PSet) {
	sn := pset.MSL_MSG.PID.Controller_sn
	workorder_id := pset.MSL_MSG.PID.Workorder_id
	result_id := pset.MSL_MSG.PID.Result_id
	count := pset.MSL_MSG.PID.Count

	s_dat, s_time := utils.GetDateTime()

	no := GenerateRangeNum(8, cvi3_client.mtx_no)
	var result = ""
	if no % 2 == 0 {
		result = "IO"
	} else {
		result = "NIO"
	}

	mi := (rand.Float64() * 10) + 10
	wi := (rand.Float64() * 180) + 5
	ti := rand.Float64()*100 + 100

	send_result := fmt.Sprintf(xml_result,
		sn,
		workorder_id,
		result_id,
		count,
		result,
		s_dat,
		s_time,
		mi,
		wi,
		ti)

	seq := cvi3_client.get_sequence()
	pkg, seq := audi_vw.GeneratePacket(seq, audi_vw.Header_type_request_with_reply, send_result)
	cvi3_client.Write([]byte(pkg), seq)
}

func (cvi3_client *CVI3Client) AudoPushResult() {

	pset := CVI3PSet{}
	pset.MSL_MSG.PID.Controller_sn = cvi3_client.SN
	pset.MSL_MSG.PID.Workorder_id = cvi3_client.Workorder_id

	for {
		for _, v := range cvi3_client.Result_ids {
			pset.MSL_MSG.PID.Result_id = v

			for i := 0; i < cvi3_client.Count; i++ {
				pset.MSL_MSG.PID.Count = i + 1

				cvi3_client.PushResult(pset)

			}
		}

	}


}

func (cvi3_client *CVI3Client) Write(buf []byte, seq uint) {
	cvi3_client.buffer <- buf
}

// 异步发送
func (cvi3_client *CVI3Client) manage() {

	for {
		v := <- cvi3_client.buffer
		err := cvi3_client.client.Write([]byte(v))
		if err != nil {
			fmt.Printf("控制器 %s 发送失败:%s\n", err.Error())
		}

		//fmt.Printf("控制器:%s 发送:%s\n", cvi3_client.SN, string(v))

		<- time.After(time.Duration(300 * time.Millisecond)) //300毫秒发送一次信号
	}
}