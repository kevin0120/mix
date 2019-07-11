package scanner

import (
	"fmt"
	"github.com/bep/debounce"
	"github.com/pkg/errors"
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
	NewReader(d *USBDevice) error

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

func NewDevice(vid, pid ID, d Diagnostic) *DeviceInfo {
	ret := &DeviceInfo{
		VendorID:  vid,
		ProductID: pid,
		diag:      d,
	}

	return ret
}

type Notify interface {
	OnStatus(string, string)
	OnRecv(string, string)
}

type Scanner struct {
	devInfo *DeviceInfo
	device  *USBDevice

	diag   Diagnostic
	notify Notify
	status atomic.Value

	debounced       func(f func())
	debounceTrigger bool
	init            bool
}

func NewScanner(vid, pid ID, d Diagnostic, dev *USBDevice) *Scanner {
	di := NewDevice(vid, pid, d)

	return &Scanner{devInfo: di, diag: d, device: dev, debounceTrigger: false, init: true}
}

func (s *Scanner) Start() {
	s.status.Store(SCANNER_STATUS_OFFLINE)
	go s.manage()
}

func (s *Scanner) Stop() error {
	return nil
}

func (s *Scanner) ID() string {
	return fmt.Sprintf("%s", s.device.String())
}

func (s *Scanner) getID() (ID, ID) {
	di := s.devInfo
	if di == nil {
		return 0, 0
	}
	return di.VendorID, di.ProductID
}

func (s *Scanner) Status() string {
	return s.status.Load().(string)
}

func (s *Scanner) open() (*USBDevice, error) {
	vid, pid := s.getID()
	if vid == 0 || pid == 0 {
		return nil, errors.New("Device Info is Empty\n")
	}

	err := s.device.SetAutoDetach(true)
	if err != nil {
		return nil, err
	}

	di := s.devInfo
	if di == nil {
		return nil, errors.New("DeviceInfo is Empty")
	}
	if err := di.updateDeviceService(); err != nil {
		return nil, err
	}
	if err := di.NewReader(s.device); err != nil {
		return nil, err
	}

	return s.device, err
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
	if err := d.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Scanner) manage() {
	for {
		_ = s.connectAndRecv()
		time.Sleep(SCANNER_OPEN_ITV)
	}
}

func (s *Scanner) connectAndRecv() error {
	//ctx := gousb.NewContext()
	//defer ctx.Close()
	d, err := s.open()
	if err == nil {
		// device online
		s.device = d
		s.status.Store(SCANNER_STATUS_ONLINE)
		s.notify.OnStatus(s.ID(), SCANNER_STATUS_ONLINE)
		s.recv()
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

func (s *Scanner) recv() {
	buf := make([]byte, SCANNER_BUF_LEN)
	di := s.devInfo
	if di == nil {
		return
	}

	strRecv := ""
	for {
		n, err := di.Read(buf)
		if err != nil {
			s.diag.Error("read failed", err)
			// device offline
			s.status.Store(SCANNER_STATUS_OFFLINE)
			s.notify.OnStatus(s.ID(), SCANNER_STATUS_OFFLINE)
			return
		}

		if s.devInfo == nil {
			s.diag.Error("Recv Fail, Plz init device info first",
				errors.Errorf("Scanner: %s ", s.device.String()))
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
					s.notify.OnRecv(s.ID(), strRecv)
					s.resetDebounce()
					strRecv = ""
				}
			})
		}
	}
}
