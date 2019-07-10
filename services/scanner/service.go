package scanner

import (
	"fmt"
	"github.com/google/gousb"
	"sync"
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
	mtxScanners sync.Mutex

	diag Diagnostic
	Notify
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
	go s.search()

	return nil
}

func (s *Service) Close() error {

	return nil
}

func (s *Service) search() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	for {
		for _, v := range vendors {
			devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
				return desc.Vendor == v.VendorID && desc.Product == v.ProductID
			})

			if err == nil {
				for _, d := range devs {
					s.addScanner(NewScanner(v.VendorID, v.ProductID, s.diag, d))
				}
			} else {
				fmt.Println(err.Error())
			}
		}

		time.Sleep(SEARCH_ITV)
	}
}

func (s *Service) addScanner(scanner *Scanner) {
	defer s.mtxScanners.Unlock()
	s.mtxScanners.Lock()

	if _, ok := s.scanners[scanner.ID()]; !ok {
		s.scanners[scanner.ID()] = scanner
		scanner.notify = s
		scanner.Start()
	}
}

func (s *Service) removeScanner(id string) {
	defer s.mtxScanners.Unlock()
	s.mtxScanners.Lock()

	delete(s.scanners, id)
}

func (s *Service) OnStatus(id string, status string) {
	s.diag.Debug(fmt.Sprintf("scanner %s status: %s\n", id, status))
	if status == SCANNER_STATUS_OFFLINE {
		s.removeScanner(id)
	}
}

func (s *Service) OnRecv(id string, str string) {
	s.diag.Debug(fmt.Sprintf("scanner %s recv: %s\n", id, str))
}
