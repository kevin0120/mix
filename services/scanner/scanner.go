package scanner

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/gousb"
	"github.com/pkg/errors"
)

const (
	SCANNER_OPEN_ITV       = 500 * time.Millisecond
	SCANNER_BUF_LEN        = 256
	SCANNER_STATUS_ONLINE  = "online"
	SCANNER_STATUS_OFFLINE = "offline"
)

type DeviceService interface {
	Parse([]byte) string
	//
	//isOpen() bool
	//
	//// vendorID, productID
	NewReader(d *USBDevice) error

	Read([]byte) (int, error)

	Close() error
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
}

func NewScanner(vid, pid ID, d Diagnostic) *Scanner {
	di := NewDevice(vid, pid, d)

	return &Scanner{devInfo: di, diag: d}
}

func (s *Scanner) Start() {
	s.status.Store(SCANNER_STATUS_OFFLINE)
	go s.manage()
}

func (s *Scanner) Stop() error {
	return nil
}

func (s *Scanner) ID() string {
	di := s.devInfo
	return fmt.Sprintf("VID:%d, PID:%d, interface: %s", di.VendorID, di.ProductID, di.Channel)
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

func (s *Scanner) open(ctx *gousb.Context) (*USBDevice, error) {
	vid, pid := s.getID()
	if vid == 0 || pid == 0 {
		return nil, errors.New("Device Info is Empty\n")
	}
	d, err := ctx.OpenDeviceWithVIDPID(vid, pid)
	if err != nil {
		return nil, errors.Errorf("Open Device vid:%d, pid: %d fail", vid, pid)
	}
	di := s.devInfo
	if di == nil {
		return nil, errors.New("DeviceInfo is Empty")
	}
	if err := di.NewReader(d); err != nil {
		return nil, err
	}
	return d, err
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
		err := s.connectAndRecv()
		if err != nil {
			time.Sleep(SCANNER_OPEN_ITV)
		} else {
			return
		}
	}
}

func (s *Scanner) connectAndRecv() error {
	ctx := gousb.NewContext()
	defer ctx.Close()
	d, err := s.open(ctx)
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

func (s *Scanner) recv() {
	buf := make([]byte, SCANNER_BUF_LEN)
	di := s.devInfo
	if di == nil {
		return
	}

	for {
		n, err := di.Read(buf)
		if err != nil {
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
			str := s.devInfo.Parse(buf[0:n])
			s.notify.OnRecv(s.ID(), str)
		}
	}
}
