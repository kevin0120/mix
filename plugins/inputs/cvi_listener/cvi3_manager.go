package cvi_listener

import (
	"github.com/satori/go.uuid"
	"encoding/base64"
	"time"
	"fmt"
	"github.com/masami10/rush/utils/cvi"
	"net"
	"strings"
)

const (
	ERR_CVI3_NOT_FOUND = -1
	ERR_CVI3_OFFLINE = -2
	ERR_CVI3_REQUEST = -3
	ERR_CVI3_REPLY_TIMEOUT = -4
	ERR_DB = -5
)

type CVI3Manager struct {
	CVI3_clients   map[string]*CVI3Client
	Parent *streamSocketListener
}

func (cm *CVI3Manager) generate_id() (string) {
	u4 := uuid.NewV4()
	return base64.RawURLEncoding.EncodeToString(u4.Bytes())
}

func (cm *CVI3Manager) StartService(configs	[]*CVIConfig) {
	cm.CVI3_clients = map[string]*CVI3Client{}
	for _, cvi3 := range configs {
		client := CVI3Client{}
		client.Config = cvi3

		cm.CVI3_clients[cvi3.SN] = &client
		go client.Start()
	}
}

func (cm *CVI3Manager) SetRemoteConn(addr string, c net.Conn) (string) {
	ip := strings.Split(addr, ":")[0]
	for k,v := range cm.CVI3_clients {
		if v.Config.IP == ip {
			v.RemoteConn = c
			return k
		}
	}

	return ""
}

func (cm *CVI3Manager) PSet(sn string, pset int, workorder_id int) (int) {
	// 判断控制器是否存在
	cvi3_client, exist := cm.CVI3_clients[sn]
	if !exist {
		// SN对应控制器不存在
		return ERR_CVI3_NOT_FOUND
	}

	if cvi3_client.Status == STATUS_OFFLINE {
		// 控制器离线
		return ERR_CVI3_OFFLINE
	}

	// 设定pset并判断控制器响应
	screw_id := cm.generate_id()
	serial, err := cvi3_client.PSet(pset, workorder_id, screw_id)
	if err != nil {
		// 控制器请求失败
		return ERR_CVI3_REQUEST
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
		return ERR_CVI3_REPLY_TIMEOUT
	}

	fmt.Printf("reply_header:%s\n", header_str)
	header := cvi.CVI3Header{}
	header.Deserialize(header_str)
	if !header.Check() {
		// 控制器请求失败
		return ERR_CVI3_REQUEST
	}

	dberr := cm.Parent.CUR_DB.PreSave(sn, workorder_id, screw_id)
	if dberr != nil {
		return ERR_DB
	}

	return 0
}
