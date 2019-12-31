package openprotocol

import (
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/pkg/errors"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
}

// TODO: 修改服务中的DISPATCH相关方法
type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	name        string

	DB    *storage.Service
	WS    *wsnotify.Service
	Aiis  *aiis.Service
	Minio *minio.Service
	Odoo  *odoo.Service

	devices []*TighteningController
	tightening_device.ITighteningProtocol
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		name:    tightening_device.TIGHTENING_OPENPROTOCOL,
		diag:    d,
		devices: []*TighteningController{},
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Name() string {
	return s.name
}

func (s *Service) SendIdentification(identification string) error {
	if s.WS == nil {
		return errors.New("Please Inject Notify Service First")
	}
	s.WS.NotifyAll(wsnotify.WS_EVENT_SCANNER, identification)
	return nil
}

func (s *Service) Parse(msg string) ([]byte, error) {
	return nil, nil
}

func (s *Service) Write(sn string, buf []byte) error {
	return nil
}

func (s *Service) CreateController(cfg *tightening_device.TighteningDeviceConfig, dp tightening_device.Dispatcher) (tightening_device.ITighteningController, error) {
	return NewController(cfg, s.diag, s, dp)
}

func (s *Service) Open() error {

	return nil
}

func (s *Service) Close() error {

	return nil
}

func (s *Service) GetDefaultMode() string {
	c := s.config()
	return c.DefaultMode
}

func (s *Service) generateIDInfo(info string) string {
	ids := ""
	for i := 0; i < MAX_IDS_NUM; i++ {
		if i == s.config().DataIndex {
			ids += fmt.Sprintf("%-25s", info)
		} else {
			ids += fmt.Sprintf("%25s", "")
		}
	}

	return ids
}

func (s *Service) TryCreateMaintenance(info ToolInfo) error {
	return s.Odoo.TryCreateMaintenance(info)
}

func (s *Service) OnStatus(string, string) {
	return
}

func (s *Service) OnRecv(string, string) {
	return
}
