package scanner

import (
	"encoding/json"
	"fmt"
	"github.com/google/gousb"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/pkg/errors"
	"github.com/tarm/serial"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 在linux上编译windows版本
// 1. 安装mingw-w64： apt-get install mingw-w64
// 2. 添加编译参数：CC=x86_64-w64-mingw32-gcc;CGO_LDFLAGS=-L/home/linshenqi/tools/libusb -lusb-1.0;GOOS=windows;CGO_ENABLED=1

const (
	ServiceSearchItv = 2000 * time.Millisecond
)

type Diagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	configValue atomic.Value
	scanners    map[string]*Scanner
	mtxScanners sync.Mutex

	diag Diagnostic
	Notify

	WS            *wsnotify.Service
	DeviceService *device.Service
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag:        d,
		mtxScanners: sync.Mutex{},
		scanners:    map[string]*Scanner{},
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	go s.search()

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) search() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	label := s.config().EntityLabel
	var vid, pid int64

	if runtime.GOOS == "windows" {
		// COM 口 通过串口打开
	} else {
		ls := strings.Split(label, ":")
		if len(ls) != 2 {
			return
		}
		vid, _ = strconv.ParseInt(ls[0], 10, 16)
		pid, _ = strconv.ParseInt(ls[1], 10, 16)
	}

	for {
		var scanner *Scanner
		if len(s.scanners) > 0 {
			for _, ss := range s.scanners {
				scanner = ss
				break
			}
		}
		if scanner != nil && scanner.Status() == SCANNER_STATUS_ONLINE {
			time.Sleep(ServiceSearchItv)
			continue
		}
		if runtime.GOOS != "windows" {
			d, err := ctx.OpenDeviceWithVIDPID(ID(vid), ID(pid))
			if err == nil && d != nil {
				s.diag.Debug(fmt.Sprintf("Search Success: %s", label))
				newScanner := NewScanner(label, s.diag, d)
				s.addScanner(newScanner)
				s.DeviceService.AddDevice(fmt.Sprintf("%d:%d", vid, pid), newScanner)
			} else if err != nil {
				s.diag.Error("Search Fail", err)
			} else {
				s.diag.Error("Search Fail", errors.New(fmt.Sprintf("Open Fail VID:%d, PID:%d", vid, pid)))
			}
		} else {
			// windows
			c := &serial.Config{Name: label}
			d, err := serial.OpenPort(c)
			if err == nil {
				s.diag.Debug(fmt.Sprintf("Search Success: %s", label))
				newScanner := NewScanner(label, s.diag, d)
				s.addScanner(newScanner)
				s.DeviceService.AddDevice(fmt.Sprintf("%s", label), newScanner)
			} else {
				s.diag.Error("Search Fail", err)
			}
		}
		time.Sleep(ServiceSearchItv)
	}
}

func (s *Service) addScanner(scanner *Scanner) {
	defer s.mtxScanners.Unlock()
	s.mtxScanners.Lock()

	if _, ok := s.scanners[scanner.Channel()]; !ok {
		s.scanners[scanner.Channel()] = scanner
		scanner.notify = s
		scanner.Start()
	}
}

func (s *Service) removeScanner(id string) {
	defer s.mtxScanners.Unlock()
	s.mtxScanners.Lock()

	if _, ok := s.scanners[id]; ok {
		scanner := s.scanners[id]
		if err := scanner.Stop(); err == nil {
			delete(s.scanners, id)
		}
	}
}

func (s *Service) OnStatus(id string, status string) {
	s.diag.Debug(fmt.Sprintf("scanner %s status: %s\n", id, status))
	if status == SCANNER_STATUS_OFFLINE {
		s.removeScanner(id)
	}
	barcode, _ := json.Marshal(wsnotify.WSMsg{
		Type: WS_SCANNER_STATUS,
		Data: ScannerStatus{
			ID:     id,
			Status: status,
		},
	})
	s.WS.WSSendScanner(string(barcode))
}

func (s *Service) OnRecv(id string, str string) {
	s.diag.Debug(fmt.Sprintf("scanner %s recv: %s\n", id, str))
	barcode, _ := json.Marshal(wsnotify.WSMsg{
		Type: WS_SCANNER_READ,
		Data: ScannerRead{
			ID:      id,
			Barcode: str,
		},
	})

	s.WS.WSSendScanner(string(barcode))
}
