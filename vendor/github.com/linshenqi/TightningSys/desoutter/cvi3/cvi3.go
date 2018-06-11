package cvi3

import (
	"net"
	"strings"
	"time"
	"errors"
)

const (
	ERR_CVI3_NOT_FOUND = "CIV3 SN is invalid"
	ERR_CVI3_OFFLINE = "cvi3 offline"
	ERR_CVI3_REQUEST = "request to cvi3 failed"
	ERR_CVI3_REPLY_TIMEOUT = "cvi3 reply timeout"
	ERR_CVI3_REPLY = "cvi3 reply contains error"
	//ERR_DB = "cvi3 reply presave failed"
)

type ControllerStatus struct {
	SN string `json:"controller_sn"`
	Status string `json:"status"`
}

type CVI3 struct {
	Configs		  []CVI3Config
	server 		  *CVI3Server
	clients   	  map[string]*CVI3Client
	FUNCStatus 	  FUNC_STATUS
	FUNCRecv 	  FUNC_RECV
}

type FUNC_STATUS func(sn string, status string)
type FUNC_RECV func(recv string)

// 注册回调
func (cvi3 *CVI3) RegisterCallBack(func_status FUNC_STATUS, func_recv FUNC_RECV) {
	cvi3.FUNCStatus = func_status
	cvi3.FUNCRecv = func_recv
}

// 服务配置
func (cvi3 *CVI3) Config(configs []CVI3Config) error {
	cvi3.Configs = []CVI3Config{}
	cvi3.Configs = configs

	return nil
}

// 启动服务
func (cvi3 *CVI3) StartService(port string) error {
	cvi3.server = &CVI3Server{}
	cvi3.server.Parent = cvi3

	cvi3.clients = map[string]*CVI3Client{}

	// 启动服务端
	var err error
	err = cvi3.server.Start(port)
	if err != nil {
		return err
	}

	// 根据配置启动客户端
	for _, conf := range cvi3.Configs {
		client := CVI3Client{}
		client.Config = conf
		client.Parent = cvi3

		cvi3.clients[conf.SN] = &client
		go client.Start()
	}

	return nil
}

// 取得控制器状态
func (cvi3 *CVI3) GetControllersStatus(sn string) ([]ControllerStatus, error) {
	status := []ControllerStatus{}
	if sn != "" {
		c, exist := cvi3.clients[sn]
		if !exist {
			return status, errors.New("controller not found")
		} else {
			s := ControllerStatus{}
			s.SN = sn
			s.Status = c.GetStatus()
			status = append(status, s)
			return status, nil
		}
	} else {
		for k, v := range cvi3.clients {
			s := ControllerStatus{}
			s.SN = k
			s.Status = v.GetStatus()
			status = append(status, s)
		}

		return status, nil
	}
}

// 设置拧接程序
func (cvi3 *CVI3) PSet(sn string, pset int, workorder_id int, result_id int, count int) (error) {
	// 判断控制器是否存在
	cvi3_client, exist := cvi3.clients[sn]
	if !exist {
		// SN对应控制器不存在
		return errors.New(ERR_CVI3_NOT_FOUND)
	}

	if cvi3_client.GetStatus() == STATUS_OFFLINE {
		// 控制器离线
		return errors.New(ERR_CVI3_OFFLINE)
	}

	// 设定pset并判断控制器响应
	//screw_id := GenerateID()
	serial, err := cvi3_client.PSet(pset, workorder_id, result_id, count)
	if err != nil {
		// 控制器请求失败
		return errors.New(ERR_CVI3_REQUEST)
	}

	var header_str string
	for i := 0; i < 6; i++ {
		header_str = cvi3_client.Results.get(serial)
		if header_str != "" {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if header_str == "" {
		// 控制器请求失败
		return errors.New(ERR_CVI3_REPLY_TIMEOUT)
	}

	//fmt.Printf("reply_header:%s\n", header_str)
	header := CVI3Header{}
	header.Deserialize(header_str)
	if !header.Check() {
		// 控制器请求失败
		return errors.New(ERR_CVI3_REPLY)
	}

	return nil
}

func (cvi3 *CVI3) setRemoteConn(addr string, c net.Conn) (string) {
	ip := strings.Split(addr, ":")[0]
	for k,v := range cvi3.clients {
		if v.Config.IP == ip {
			v.RemoteConn = c
			return k
		}
	}

	return ""
}