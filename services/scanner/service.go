package scanner

import (
	"fmt"
	"github.com/karalabe/hid"
	"github.com/kataras/iris/core/errors"
	"sync/atomic"
	"time"
)

// 编译windows版本
// 1. 安装mingw-w64： apt-get install mingw-w64
// 2. 添加编译参数：CC=i686-w64-mingw32-gcc;CGO_ENABLED=1;GOOS=windows;GOARCH=386

const (
	SEARCH_ITV = 500 * time.Millisecond
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	configValue atomic.Value
	scanners    map[string]*Scanner

	diag Diagnostic
	ScannerNotify
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {

	if !hid.Supported() {
		return errors.New("scanner service not supported")
	}

	go s.search()

	return nil
}

func (s *Service) Close() error {

	return nil
}

func (s *Service) search() {
	for {
		for k, v := range VENDOR_MODELS {
			vid, pid := v.ModelInfo()
			devs := hid.Enumerate(vid, pid)
			for _, dev := range devs {
				s.newDevice(k, &dev)
			}
		}

		time.Sleep(SEARCH_ITV)
	}
}

func (s *Service) newDevice(model string, dev *hid.DeviceInfo) {
	scanner := Scanner{
		devInfo: dev,
		notify:  s,
		vendor:  VENDOR_MODELS[model],
	}

	if _, ok := s.scanners[scanner.ID()]; !ok {
		s.scanners[scanner.ID()] = &scanner
		scanner.Start()
	}
}

func (s *Service) OnStatus(id string, status string) {
	s.diag.Debug(fmt.Sprintf("scanner %s status: %s\n", id, status))
}

func (s *Service) OnRecv(id string, str string) {
	s.diag.Debug(fmt.Sprintf("scanner %s recv: %s\n", id, str))
}
