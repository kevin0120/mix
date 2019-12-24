package ts002

import (
	"sync/atomic"
)

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
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

}
