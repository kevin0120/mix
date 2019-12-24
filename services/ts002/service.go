package ts002

import (
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/tightening_device"
	"sync/atomic"
)

type Service struct {
	configValue atomic.Value
	diag        Diagnostic

	IO               *io.Service
	TighteningDevice *tightening_device.Service
}

func NewService(c Config, d Diagnostic) *Service {
	ss := &Service{
		diag: d,
	}

	ss.configValue.Store(c)

	return ss
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	return nil
}

// 收到拧紧结果
func (s Service) onTighteningResult(data interface{}) {
	if data == nil {
		return
	}

	tighteningResult := data.(*tightening_device.TighteningResult)
	if tighteningResult.MeasureResult == tightening_device.RESULT_NOK {
		// 如果结果NOK，则触发报警
	}

}
