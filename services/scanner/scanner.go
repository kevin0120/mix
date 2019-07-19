package scanner

import (
	"github.com/bep/debounce"
	"github.com/google/gousb"
	"github.com/pkg/errors"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	SCANNER_OPEN_ITV       = 500 * time.Millisecond
	SCANNER_BUF_LEN        = 256
	SCANNER_STATUS_ONLINE  = "online"
	SCANNER_STATUS_OFFLINE = "offline"
)

type DeviceService interface {
	Parse([]byte) (string, error)
	//
	//isOpen() bool
	//
	//// vendorID, productID
	NewReader(d USBDevice) error

	Read([]byte) (int, error)

	Close() error

	//Debounce() (time.Duration, time.Duration)
}

type DeviceInfo struct {
	VendorID  ID
	ProductID ID
	Channel   string //evdev(linux/mac）or COM(win)
	diag      Diagnostic
	DeviceService
}

func NewDevice(channel string, d Diagnostic) *DeviceInfo {
	ret := &DeviceInfo{
		Channel: channel,
		diag:    d,
	}
	if runtime.GOOS != "windows" {
		var vid, pid int64
		ls := strings.Split(channel, ":")
		if len(ls) != 2 {
			return ret
		}
		vid, _ = strconv.ParseInt(ls[0], 10, 16)
		pid, _ = strconv.ParseInt(ls[1], 10, 16)
		ret.VendorID = ID(vid)
		ret.ProductID = ID(pid)
	}

	return ret
}

type Notify interface {
	OnStatus(string, string)
	OnRecv(string, string)
}

type Scanner struct {
	devInfo *DeviceInfo
	device  USBDevice // maybe gousb, or serial

	diag   Diagnostic
	notify Notify
	status atomic.Value

	debounced       func(f func())
	debounceTrigger bool
	init            bool
}

func NewScanner(channel string, d Diagnostic, dev USBDevice) *Scanner {
	di := NewDevice(channel, d)

	return &Scanner{devInfo: di, diag: d, device: dev, debounceTrigger: false, init: true}
}

func (s *Scanner) Start() {
	s.status.Store(SCANNER_STATUS_OFFLINE)
	go s.manage()
}

func (s *Scanner) Stop() error {
	return s.close()
}

func (s *Scanner) Channel() string {
	return s.devInfo.Channel
}

func (s *Scanner) getVIDPID() (ID, ID) {
	di := s.devInfo
	if di == nil {
		return 0, 0
	}
	return di.VendorID, di.ProductID
}

func (s *Scanner) Status() string {
	return s.status.Load().(string)
}

func (s *Scanner) open() (USBDevice, error) {
	di := s.devInfo
	if di == nil {
		return nil, errors.New("DeviceInfo is Empty")
	}
	label := s.Channel()
	if label == "" {
		return nil, errors.New("Device Info Label is Empty\n")
	}
	if runtime.GOOS != "windows" {
		// 通过 gousb 创建vid, pid
		vid, pid := s.getVIDPID()
		if vid == 0 || pid == 0 {
			return nil, errors.New("Device Info is Empty\n")
		}
		d := s.device.(*gousb.Device)
		if d == nil {
			return nil, errors.New("Device is Empty\n")
		}
		err := d.SetAutoDetach(true)
		if err != nil {
			return nil, err
		}
	}
	if err := di.updateDeviceService(); err != nil {
		return nil, err
	}
	if err := di.NewReader(s.device); err != nil {
		return nil, err
	}

	return s.device, nil
}

func (s *Scanner) close() error {
	d := s.device
	if d == nil {
		return nil
	}
	di := s.devInfo
	if di != nil {
		if err := di.Close(); err != nil {
			return err
		}
	}
	s.device = nil
	s.devInfo = nil
	return nil
}

func (s *Scanner) manage() {
	for {
		if s.Status() == SCANNER_STATUS_OFFLINE {
			_ = s.connect()
		}
		if s.Status() == SCANNER_STATUS_ONLINE {
			s._recv()
		}
		time.Sleep(SCANNER_OPEN_ITV)
	}
}

func (s *Scanner) connect() error  {
	d, err := s.open()
	if err == nil {
		// device online
		s.device = d
		s.status.Store(SCANNER_STATUS_ONLINE)
		s.notify.OnStatus(s.Channel(), SCANNER_STATUS_ONLINE)
	}
	return err
}

func (s *Scanner) recv() error {
	//ctx := gousb.NewContext()
	//defer ctx.Close()
	d, err := s.open()
	if err == nil {
		// device online
		s.device = d
		s.status.Store(SCANNER_STATUS_ONLINE)
		s.notify.OnStatus(s.Channel(), SCANNER_STATUS_ONLINE)
		s._recv() //阻塞接收数据
		return nil
	} else {
		return err
	}
}

func (s *Scanner) resetDebounce() {
	if s.init {
		s.init = false
	}

	s.debounceTrigger = false
}

func (s *Scanner) triggerDebounce() {
	if !s.debounceTrigger {
		//debInit, debCommon := s.devInfo.Debounce()
		//deb := debCommon
		//if s.init {
		//	deb = debInit
		//}

		s.debounced = debounce.New(100 * time.Millisecond)
		s.debounceTrigger = true
	}
}

func (s *Scanner) _recv() {
	buf := make([]byte, SCANNER_BUF_LEN)
	di := s.devInfo
	if di == nil {
		return
	}

	strRecv := ""
	for {
		n, err := di.Read(buf)
		if err != nil {
			s.diag.Error("Read Fail", err)
			// device offline
			s.status.Store(SCANNER_STATUS_OFFLINE)
			s.notify.OnStatus(s.Channel(), SCANNER_STATUS_OFFLINE)
			return
		}

		if s.devInfo == nil {
			s.diag.Error("Recv Fail, Plz init device info first",
				errors.Errorf("Scanner: %s ", s.Channel()))
			continue
		}

		if n > 0 {
			s.triggerDebounce()
			str, err := s.devInfo.Parse(buf[0:n])
			if err == nil {
				strRecv += str
			}

			s.debounced(func() {
				if strRecv != "" {
					s.notify.OnRecv(s.Channel(), strRecv)
					s.resetDebounce()
					strRecv = ""
				}
			})
		}
	}
}
