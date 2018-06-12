package rush

import (
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/storage"
)


type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	HTTPDService		*httpd.Service

	StorageService		*storage.Service

	diag        		Diagnostic
}

func NewService(c Config, d Diagnostic) *Service {
	if c.Enable{
		return &Service{
			diag: d,
		}
	}
	return nil
}

func (s *Service) Open() error {
	return nil
}


func (s *Service) Close() error {
	return nil
}