package device

import (
	"sync"
	"sync/atomic"
)

func CreateBaseDevice() BaseDevice {
	c := BaseDevice{
		mtxChildren: sync.Mutex{},
		children:    map[string]IDevice{},
		status:      atomic.Value{},
	}

	c.UpdateStatus(STATUS_OFFLINE)
	return c
}

type BaseDevice struct {
	status      atomic.Value
	mtxChildren sync.Mutex
	children    map[string]IDevice
}

func (s *BaseDevice) UpdateStatus(status string) {
	//fmt.Println(status)
	s.status.Store(status)
}

func (s *BaseDevice) Status() string {
	return s.status.Load().(string)
}

func (s *BaseDevice) AddChildren(sn string, device IDevice) {
	s.mtxChildren.Lock()
	defer s.mtxChildren.Unlock()

	s.children[sn] = device
}

func (s *BaseDevice) Children() map[string]IDevice {
	s.mtxChildren.Lock()
	defer s.mtxChildren.Unlock()

	return s.children
}

// 设备配置
func (s *BaseDevice) Config() interface{} {
	return nil
}

// 设备运行数据
func (s *BaseDevice) Data() interface{} {
	return nil
}
