package scanner

import (
	"fmt"
	"github.com/karalabe/hid"
	"sync/atomic"
	"time"
)

const (
	SCANNER_OPEN_ITV       = 500 * time.Millisecond
	SCANNER_BUF_LEN        = 256
	SCANNER_STATUS_ONLINE  = "online"
	SCANNER_STATUS_OFFLINE = "offline"
)

type ScannerNotify interface {
	OnStatus(string, string)
	OnRecv(string, string)
}

type Scanner struct {
	devInfo *hid.DeviceInfo
	vendor  Vendor
	notify  ScannerNotify
	status  atomic.Value
}

func (s *Scanner) Start() {
	s.status.Store(SCANNER_STATUS_OFFLINE)
	go s.connect()
}

func (s *Scanner) Stop() error {
	return nil
}

func (s *Scanner) ID() string {
	vid, pid := s.vendor.ModelInfo()
	return fmt.Sprintf("%d.%d.%d", vid, pid, s.devInfo.Interface)
}

func (s *Scanner) Status() string {
	return s.status.Load().(string)
}

func (s *Scanner) connect() {
	buf := make([]byte, SCANNER_BUF_LEN)

	var dev *hid.Device
	for {
		d, err := s.devInfo.Open()
		if err == nil {
			// device online
			s.status.Store(SCANNER_STATUS_ONLINE)
			s.notify.OnStatus(s.ID(), SCANNER_STATUS_ONLINE)
			dev = d
			break
		} else {
			time.Sleep(SCANNER_OPEN_ITV)
			continue
		}
	}

	for {
		n, err := dev.Read(buf)
		if err != nil {
			// device offline
			s.status.Store(SCANNER_STATUS_OFFLINE)
			s.notify.OnStatus(s.ID(), SCANNER_STATUS_OFFLINE)
			break
		}

		if n > 0 {
			// receive data: string(buf[0:n])
			str := s.vendor.Parse(buf[0:n])
			s.notify.OnRecv(s.ID(), str)
		}
	}

}
