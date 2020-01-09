package scanner

import (
	"fmt"
	"github.com/google/gousb"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherbus"
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

const ScannerDispatcherKey = dispatcherbus.DISPATCHER_SCANNER_DATA

type Service struct {
	configValue atomic.Value
	scanners    map[string]*Scanner
	mtxScanners sync.Mutex

	diag Diagnostic
	Notify
	dispatcher    Dispatcher
	WS            *wsnotify.Service
	DeviceService IDeviceService
}

func NewService(c Config, d Diagnostic, dispatcher Dispatcher, ds IDeviceService) *Service {

	s := &Service{
		diag:          d,
		mtxScanners:   sync.Mutex{},
		scanners:      map[string]*Scanner{},
		dispatcher:    dispatcher,
		DeviceService: ds,
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
	err := s.createAndStartScannerDispatcher()
	if err != nil {
		s.diag.Error("Open", err)
		return err
	}
	return nil
}

func (s *Service) dispatchRecvData(data string) {
	if s.dispatcher == nil {
		err := errors.New("Please Inject dispatcherBus Service First")
		s.diag.Error("dispatchRecvData", err)
		return
	}
	err := s.dispatcher.Dispatch(ScannerDispatcherKey, data)
	if err != nil {
		s.diag.Error("dispatchRecvData", err)
	}
	return
}

func (s *Service) createAndStartScannerDispatcher() error {
	if s.dispatcher == nil {
		err := errors.New("Please Inject dispatcherBus Service First")
		s.diag.Error("createAndStartScannerDispatcher", err)
		return err
	}
	err := s.dispatcher.Create(ScannerDispatcherKey, 20)
	if err == nil || strings.HasPrefix(err.Error(), "Dispatcher Already Exist") {
		return s.dispatcher.Start(ScannerDispatcherKey)
	}
	return err
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) search() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	if s.DeviceService == nil {
		err := errors.New("Please Initial Device Service First")
		s.diag.Error("Scanner Search", err)
		return
	}

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
				newScanner := NewScanner(label, s.diag, d, s)
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
				newScanner := NewScanner(label, s.diag, d, s)
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

func (s *Service) sendScannerInfo(msg string) {
	s.WS.NotifyAll(wsnotify.WS_EVENT_SCANNER, msg)
}

func (s *Service) OnStatus(id string, status string) {
	s.diag.Debug(fmt.Sprintf("Scanner %s Status: %s\n", id, status))
	if status == SCANNER_STATUS_OFFLINE {
		s.removeScanner(id)
	}

	scannerStatus := []device.DeviceStatus{
		{
			SN:     id,
			Type:   device.BaseDeviceTypeScanner,
			Status: status,
		},
	}

	// 分发扫码枪状态
	s.dispatcher.Dispatch(dispatcherbus.DISPATCHER_DEVICE_STATUS, scannerStatus)
}

func (s *Service) OnRecv(id string, data string) {
	s.diag.Debug(fmt.Sprintf("Scanner %s Recv Raw Data: %s\n", id, data))
	if data == "" {
		return
	}

	barcodeData := ScannerRead{
		Src:     device.BaseDeviceTypeScanner,
		SN:      id,
		Barcode: data,
	}

	// 分发条码数据
	s.dispatcher.Dispatch(dispatcherbus.DISPATCHER_SCANNER_DATA, barcodeData)
}
