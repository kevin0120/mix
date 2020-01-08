package openprotocol

import (
	"fmt"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/pkg/errors"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	name        string

	storageService IStorageService
	backendService IBackendService

	devices []*TighteningController
	tightening_device.ITighteningProtocol
	vendors map[string]IOpenProtocolController
}

func NewService(c Config, d Diagnostic, vendors map[string]IOpenProtocolController, db IStorageService, backend IBackendService) *Service {

	s := &Service{
		name:           tightening_device.TIGHTENING_OPENPROTOCOL,
		diag:           d,
		devices:        []*TighteningController{},
		vendors:        vendors,
		storageService: db,
		backendService: backend,
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

func (s *Service) NewController(cfg *tightening_device.TighteningDeviceConfig, dp tightening_device.Dispatcher) (tightening_device.ITighteningController, error) {
	return s.newController(cfg, s.diag, s, dp)
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
	return s.backendService.TryCreateMaintenance(info)
}

func (s *Service) OnStatus(string, string) {
	s.diag.Error("OnStatus", errors.New("OpenProtocol Service Not Support OnStatus"))
}

func (s *Service) OnRecv(string, string) {
	s.diag.Error("OnRecv", errors.New("OpenProtocol Service Not Support OnRecv"))
}

func (s *Service) newController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service, dp Dispatcher) (tightening_device.ITighteningController, error) {

	c, exist := s.vendors[deviceConfig.Model]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Controller Model:%s Not Support", deviceConfig.Model))
	}

	controllerInstance := c.New()
	controllerInstance.initController(deviceConfig, d, service, dp)
	return controllerInstance, nil
}
